// Copyright 2024 Undistro Authors
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

package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/undistro/marvin/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
)

func TestCompile(t *testing.T) {
	var apiResources []*metav1.APIResourceList
	kubeVersion := &version.Info{Major: "1", Minor: "29", GitVersion: "v1.29.2"}
	podsMatch := types.Match{Resources: []types.ResourceRule{{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}}}

	tests := []struct {
		check   types.Check
		wantErr assert.ErrorAssertionFunc
	}{
		{
			check: types.Check{
				ID:    "ok",
				Match: podsMatch,
				Validations: []types.Validation{{
					Expression: `variables.isWindows || allContainers.size() > 0`,
				}},
				Variables: []types.Variable{{
					Name:       "isWindows",
					Expression: `podSpec.?os.?name.orValue("") == "windows"`,
				}},
			},
			wantErr: assert.NoError,
		},
		{
			check: types.Check{
				ID:    "validation error",
				Match: podsMatch,
				Validations: []types.Validation{{
					Expression: `allContainers.sizeX() > 0`,
				}},
			},
			wantErr: assert.Error,
		},
		{
			check: types.Check{
				ID:    "variable error",
				Match: podsMatch,
				Validations: []types.Validation{{
					Expression: `variables.isWindows || allContainers.size() > 0`,
				}},
				Variables: []types.Variable{{
					Name:       "isWindows",
					Expression: `foo`,
				}},
			},
			wantErr: assert.Error,
		},
		{
			check: types.Check{
				ID: "no workload",
				Match: types.Match{Resources: []types.ResourceRule{{
					Group:    "",
					Version:  "v1",
					Resource: "configmaps",
				}}},
				Validations: []types.Validation{{
					Expression: `allContainers.size() > 0`,
				}},
			},
			wantErr: assert.Error,
		},
		{
			check: types.Check{
				ID:          "no validations",
				Match:       podsMatch,
				Validations: nil,
				Variables: []types.Variable{{
					Name:       "isWindows",
					Expression: `podSpec.?os.?name.orValue("") == "windows"`,
				}},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.check.ID, func(t *testing.T) {
			_, err := Compile(tt.check, apiResources, kubeVersion, 1000000)
			if !tt.wantErr(t, err, fmt.Sprintf("Compile(%v, %v, %v)", tt.check, apiResources, kubeVersion)) {
				return
			}
		})
	}
}
