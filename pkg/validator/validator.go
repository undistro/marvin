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

package validator

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	marvin "github.com/undistro/marvin/pkg/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/cel/lazy"
)

// CELValidator is a Validator that performs CEL expressions
type CELValidator struct {
	check       marvin.Check
	programs    []cel.Program
	apiVersions []string
	kubeVersion *version.Info
	variables   []compiledVariable
}

type compiledVariable struct {
	name    string
	program cel.Program
}

func (r *CELValidator) SetAPIVersions(apiVersions []string) {
	r.apiVersions = apiVersions
}

func (r *CELValidator) SetKubeVersion(v *version.Info) {
	r.kubeVersion = v
}

func (r *CELValidator) Validate(obj unstructured.Unstructured, params any) (bool, string, error) {
	if params == nil {
		params = r.check.Params
	}
	input := &activation{object: obj.UnstructuredContent(), apiVersions: r.apiVersions, params: params, variables: make(map[string]any)}
	if err := r.setPodSpecParams(obj, input); err != nil {
		return false, "", err
	}
	lazyMap := lazy.NewMapValue(types.MapType)
	for _, v := range r.variables {
		lazyMap.Append(v.name, callback(v, input))
	}
	input.variables = lazyMap
	for i, prg := range r.programs {
		out, _, err := prg.Eval(input)
		if err != nil {
			return false, "", fmt.Errorf("evaluate error: %s", err)
		}
		if out != types.True {
			return false, r.check.Validations[i].Message, nil
		}
	}
	return true, "", nil
}

func callback(v compiledVariable, activation any) lazy.GetFieldFunc {
	return func(_ *lazy.MapValue) ref.Val {
		val, _, err := v.program.Eval(activation)
		if err != nil {
			return types.NewErr("variable %q fails to evaluate: %v", v.name, err)
		}
		return val
	}
}

func (r *CELValidator) setPodSpecParams(obj unstructured.Unstructured, input *activation) error {
	if !HasPodSpec(obj) {
		return nil
	}
	meta, spec, err := ExtractPodSpec(obj)
	if err != nil {
		return fmt.Errorf("pod spec extract error: %s", err)
	}
	podSpec, err := runtime.DefaultUnstructuredConverter.ToUnstructured(spec)
	if err != nil {
		return fmt.Errorf("podSpec to unstructured converter error: %s", err.Error())
	}
	podMeta, err := runtime.DefaultUnstructuredConverter.ToUnstructured(meta)
	if err != nil {
		return fmt.Errorf("podMeta to unstructured converter error: %s", err.Error())
	}
	input.podSpec = podSpec
	input.podMeta = podMeta
	for _, container := range extractAllContainers(spec) {
		c, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&container)
		if err != nil {
			return fmt.Errorf("container to unstructured converter error: %s", err.Error())
		}
		input.allContainers = append(input.allContainers, c)
	}
	return nil
}
