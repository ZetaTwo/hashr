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

package common

// Client represents a client to a data source where hashes can be looked up
type Client interface {
	// LookupSamples returns whether or not each hash in a list exist in the database.
	LookupSamples([]string) ([]bool, error)
	// CLientName() returns the client name.
	ClientName() string
}
