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
	"github.com/google/cel-go/interpreter"
)

const (
	ObjectVarName        = "object"
	ParamsVarName        = "params"
	PodMetaVarName       = "podMeta"
	PodSpecVarName       = "podSpec"
	AllContainersVarName = "allContainers"
	APIVersionsVarName   = "apiVersions"
	KubeVersionVarName   = "kubeVersion"
)

type activation struct {
	object        map[string]any
	podMeta       map[string]any
	podSpec       map[string]any
	allContainers []map[string]any
	params        any
	apiVersions   []string
	kubeVersion   any
}

func (a *activation) ResolveName(name string) (any, bool) {
	switch name {
	case ObjectVarName:
		return a.object, true
	case PodMetaVarName:
		return a.podMeta, true
	case PodSpecVarName:
		return a.podSpec, true
	case AllContainersVarName:
		return a.allContainers, true
	case ParamsVarName:
		return a.params, true
	case APIVersionsVarName:
		return a.apiVersions, true
	case KubeVersionVarName:
		return a.kubeVersion, true
	default:
		return nil, false
	}
}

func (a *activation) Parent() interpreter.Activation {
	return nil
}
