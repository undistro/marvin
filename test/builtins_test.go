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

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/undistro/marvin/pkg/loader"
	"github.com/undistro/marvin/pkg/validator"
)

func TestBuiltinChecks(t *testing.T) {
	checks, tests, err := loader.LoadChecksAndTests("../internal/builtins/")
	assert.NoError(t, err)
	assert.NotEmpty(t, checks)
	assert.GreaterOrEqual(t, len(checks), len(tests))
	for path, checkTests := range tests {
		t.Run(path, func(t *testing.T) {
			check, ok := checks[path]
			assert.True(t, ok)
			assert.NotNil(t, check)
			assert.NotEmpty(t, check.ID)
			v, err := validator.Compile(check, nil, nil)
			assert.NoError(t, err)
			assert.NotNil(t, v)
			for _, tt := range checkTests {
				t.Run(tt.Name, func(t *testing.T) {
					obj, err := parse(tt.Input)
					assert.NoError(t, err)
					assert.NotNil(t, obj)
					v.SetAPIVersions(tt.APIVersions)
					v.SetKubeVersion(tt.KubeVersion)
					got, msg, err := v.Validate(obj, tt.Params)
					assert.NoError(t, err)
					assert.Equal(t, tt.Pass, got)
					assert.Equal(t, tt.Message, msg)
				})
			}
		})
	}
}

func parse(i string) (unstructured.Unstructured, error) {
	var obj unstructured.Unstructured
	err := yaml.Unmarshal([]byte(i), &obj)
	return obj, err
}
