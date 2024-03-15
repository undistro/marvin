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

package types

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
)

type Check struct {
	ID          string            `json:"id"`
	Match       Match             `json:"match"`
	Validations []Validation      `json:"validations"`
	Variables   []Variable        `json:"variables"`
	Params      map[string]any    `json:"params"`
	Severity    Severity          `json:"severity"`
	Message     string            `json:"message"`
	Labels      map[string]string `json:"labels,omitempty"`

	Builtin bool   `json:"builtin"`
	Path    string `json:"path,omitempty"`
}

type Match struct {
	Resources []ResourceRule `json:"resources"`
}

type ResourceRule struct {
	Group    string `json:"group,omitempty"`
	Version  string `json:"version"`
	Resource string `json:"resource"`
}

func (r *ResourceRule) ToGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: r.Group, Version: r.Version, Resource: r.Resource}
}

type Validation struct {
	Expression string `json:"expression"`
	Message    string `json:"message,omitempty"`
}

type Variable struct {
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

type Test struct {
	Name        string        `json:"name"`
	Input       string        `json:"input"`
	Params      any           `json:"params"`
	APIVersions []string      `json:"apiVersions"`
	KubeVersion *version.Info `json:"kubeVersion"`
	Pass        bool          `json:"pass"`
	Message     string        `json:"message"`
}
