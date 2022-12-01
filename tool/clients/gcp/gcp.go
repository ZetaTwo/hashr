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

package gcp

import (
	"context"

	"cloud.google.com/go/spanner"
)

const (
	Name = "GCP"
)

// Name returns exporter name.
func (e *Client) ClientName() string {
	return Name
}

// NewClient returns new instance of a GCP storage client.
func NewClient(ctx context.Context, spannerDBPath *string) (*Client, error) {
	spannerClient, err := spanner.NewClient(ctx, *spannerDBPath)
	if err != nil {
		return nil, err
	}
	return &Client{spannerClient: spannerClient}, nil
}

func (e *Client) LookupSamples(hashes []string) ([]bool, error) {
	// TODO: actually lookup data in database
	res := make([]bool, len(hashes))
	for i, _ := range hashes {
		res[i] = false
	}
	return res, nil
}

type Client struct {
	spannerClient *spanner.Client
}
