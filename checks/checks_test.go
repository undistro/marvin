package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/undistro/marvin/pkg/loader"
	"github.com/undistro/marvin/pkg/validator"
)

func TestChecks(t *testing.T) {
	checks, tests, err := loader.LoadChecksAndTests(".")
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
