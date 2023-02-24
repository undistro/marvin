// Copyright 2023 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadChecks(t *testing.T) {
	emptyDir := t.TempDir()
	assert.NotEmpty(t, emptyDir)

	tests := []struct {
		name    string
		path    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "recursive",
			path:    "testdata/checks/",
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "depth 1",
			path:    "testdata/checks/workloads/",
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "valid file",
			path:    "testdata/checks/workloads/replicas.yaml",
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "invalid file",
			path:    "testdata/checks/workloads/unsupported.txt",
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "not found",
			path:    "testdata/checks/notfound",
			wantErr: true,
		},
		{
			name:    "empty",
			path:    emptyDir,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "invalid YAML",
			path:    "testdata/invalid/",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadChecks(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadChecks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("LoadChecks() got = %v items, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestLoadChecksAndTests(t *testing.T) {
	emptyDir := t.TempDir()
	assert.NotEmpty(t, emptyDir)

	tests := []struct {
		name          string
		path          string
		wantChecksLen int
		wantTestsLen  int
		wantErr       bool
	}{
		{
			name:          "recursive",
			path:          "testdata/checks/",
			wantChecksLen: 2,
			wantTestsLen:  1,
			wantErr:       false,
		},
		{
			name:          "depth 1",
			path:          "testdata/checks/workloads/",
			wantChecksLen: 1,
			wantTestsLen:  1,
			wantErr:       false,
		},
		{
			name:          "valid file",
			path:          "testdata/checks/workloads/replicas.yaml",
			wantChecksLen: 1,
			wantErr:       false,
		},
		{
			name:          "invalid file",
			path:          "testdata/checks/workloads/unsupported.txt",
			wantChecksLen: 0,
			wantErr:       false,
		},
		{
			name:    "not found",
			path:    "testdata/checks/notfound",
			wantErr: true,
		},
		{
			name:          "empty",
			path:          emptyDir,
			wantErr:       false,
			wantChecksLen: 0,
		},
		{
			name:    "invalid YAML",
			path:    "testdata/invalid/",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChecks, gotTests, err := LoadChecksAndTests(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadChecksAndTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotChecks) != tt.wantChecksLen {
				t.Errorf("LoadChecksAndTests() got = %v items, want %v", len(gotChecks), tt.wantChecksLen)
			}
			if len(gotTests) != tt.wantTestsLen {
				t.Errorf("LoadChecksAndTests() got = %v items, want %v", len(gotTests), tt.wantTestsLen)
			}
		})
	}
}
