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

package version

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		version string
		commit  string
		date    string
	}
	tests := []struct {
		name string
		args *args
		want Info
	}{
		{
			name: "default",
			want: Info{
				Version: "dev",
				Commit:  "",
				Date:    "",
				Major:   0,
				Minor:   0,
			},
		},
		{
			name: "version",
			args: &args{
				version: "0.1.0",
				commit:  "commit",
				date:    "date",
			},
			want: Info{
				Version: "0.1.0",
				Major:   0,
				Minor:   1,
				Commit:  "commit",
				Date:    "date",
			},
		},
		{
			name: "prefixed",
			args: &args{
				version: "v0.1.0",
				commit:  "commit",
				date:    "date",
			},
			want: Info{
				Version: "v0.1.0",
				Major:   0,
				Minor:   1,
				Commit:  "commit",
				Date:    "date",
			},
		},
		{
			name: "pre-release",
			args: &args{
				version: "v0.1.1-next",
				commit:  "commit",
				date:    "date",
			},
			want: Info{
				Version: "v0.1.1-next",
				Major:   0,
				Minor:   1,
				Commit:  "commit",
				Date:    "date",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args != nil {
				version = tt.args.version
				commit = tt.args.commit
				date = tt.args.date
			}
			got := Get()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
			if got.String() != tt.want.Version {
				t.Errorf("String() = %v, want %v", got.String(), tt.want.Version)
			}
		})
	}
}
