// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package rpm implements rpm package importer.
package rpm

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"

	"github.com/google/hashr/core/hashr"
	"github.com/google/hashr/importers/common"

	rpmutils "github.com/sassoftware/go-rpmutils"
)

const (
	// RepoName contains the repository name.
	RepoName  = "rpm"
	chunkSize = 1024 * 1024 * 10 // 10MB
)

// Archive holds data related to rpm archive.
type Archive struct {
	filename        string
	remotePath      string
	localPath       string
	quickSha256hash string
	repoPath        string
}

func extractRPM(rpmPath, outputFolder string) error {
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		if err2 := os.MkdirAll(outputFolder, 0755); err2 != nil {
			return fmt.Errorf("error while creating target directory: %v", err2)
		}
	}

	fd, err := os.Open(rpmPath)
	if err != nil {
		return fmt.Errorf("failed to open rpm file: %v", err)
	}
	defer fd.Close()

	rpmFile, err := rpmutils.ReadRpm(fd)
	if err != nil {
		return fmt.Errorf("failed to parse rpm file: %v", err)
	}

	err = rpmFile.ExpandPayload(outputFolder)
	if err != nil {
		return fmt.Errorf("failed to extract rpm file: %v", err)
	}

	return nil
}

// Preprocess extracts the contents of a .rpm file.
func (a *Archive) Preprocess() (string, error) {
	var err error
	a.localPath, err = common.CopyToLocal(a.remotePath, a.ID())
	if err != nil {
		return "", fmt.Errorf("error while copying %s to local file system: %v", a.remotePath, err)
	}

	baseDir, _ := filepath.Split(a.localPath)
	extractionDir := filepath.Join(baseDir, "extracted")

	if err := extractRPM(a.localPath, extractionDir); err != nil {
		return "", err
	}

	return extractionDir, nil
}

// ID returns non-unique rpm Archive ID.
func (a *Archive) ID() string {
	return a.filename
}

// RepoName returns repository name.
func (a *Archive) RepoName() string {
	return RepoName
}

// RepoPath returns repository path.
func (a *Archive) RepoPath() string {
	return a.repoPath
}

// LocalPath returns local path to a rpm Archive .rpm file.
func (a *Archive) LocalPath() string {
	return a.localPath
}

// RemotePath returns non-local path to a rpm Archive .rpm file.
func (a *Archive) RemotePath() string {
	return a.remotePath
}

// Description provides additional description for a .rpm file.
func (a *Archive) Description() string {
	return ""
}

// QuickSHA256Hash calculates sha256 hash of .rpm file.
func (a *Archive) QuickSHA256Hash() (string, error) {
	// Check if the quick hash was already calculated.
	if a.quickSha256hash != "" {
		return a.quickSha256hash, nil
	}

	f, err := os.Open(a.remotePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return "", err
	}

	// Check if the file is smaller than 20MB, if so hash the whole file.
	if fileInfo.Size() < int64(chunkSize*2) {
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}
		a.quickSha256hash = fmt.Sprintf("%x", h.Sum(nil))
		return a.quickSha256hash, nil
	}

	header := make([]byte, chunkSize)
	_, err = f.Read(header)
	if err != nil {
		return "", err
	}

	footer := make([]byte, chunkSize)
	_, err = f.ReadAt(footer, fileInfo.Size()-int64(chunkSize))
	if err != nil {
		return "", err
	}

	a.quickSha256hash = fmt.Sprintf("%x", sha256.Sum256(append(header, footer...)))
	return a.quickSha256hash, nil
}

// NewRepo returns new instance of rpm repository.
func NewRepo(path string) *Repo {
	return &Repo{location: path}
}

// Repo holds data related to a rpm repository.
type Repo struct {
	location string
	files    []string
	Archives []*Archive
}

// RepoName returns repository name.
func (r *Repo) RepoName() string {
	return RepoName
}

// RepoPath returns repository path.
func (r *Repo) RepoPath() string {
	return r.location
}

// DiscoverRepo traverses the repository and looks for files that are related to rpm archives.
func (r *Repo) DiscoverRepo() ([]hashr.Source, error) {

	if err := filepath.Walk(r.location, walk(&r.files)); err != nil {
		return nil, err
	}

	for _, file := range r.files {
		_, filename := filepath.Split(file)

		if strings.HasSuffix(filename, ".rpm") {
			r.Archives = append(r.Archives, &Archive{filename: filename, remotePath: file, repoPath: r.location})
		}
	}

	var sources []hashr.Source
	for _, Archive := range r.Archives {
		sources = append(sources, Archive)
	}

	return sources, nil
}

func walk(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			glog.Errorf("Could not open %s: %v", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".rpm") {
			*files = append(*files, path)
		}

		return nil
	}
}
