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
			_, err := Compile(tt.check, apiResources, kubeVersion)
			if !tt.wantErr(t, err, fmt.Sprintf("Compile(%v, %v, %v)", tt.check, apiResources, kubeVersion)) {
				return
			}
		})
	}
}
