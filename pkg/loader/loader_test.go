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
