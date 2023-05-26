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
	"errors"
	"fmt"

	"github.com/google/cel-go/cel"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	k8scellib "k8s.io/apiserver/pkg/cel/library"

	"github.com/undistro/marvin/pkg/types"
)

// Compile compiles the expressions of the given check and returns a Validator
func Compile(check types.Check, apiResources []*metav1.APIResourceList, kubeVersion *version.Info) (Validator, error) {
	if len(check.Validations) == 0 {
		return nil, errors.New("invalid check: a check must have at least 1 validation")
	}
	hasPodSpec := MatchesPodSpec(check.Match.Resources)
	env, err := newEnv(hasPodSpec)
	if err != nil {
		return nil, fmt.Errorf("environment construction error %s", err.Error())
	}
	prgs := make([]cel.Program, 0, len(check.Validations))
	for i, v := range check.Validations {
		ast, issues := env.Compile(v.Expression)
		if issues != nil && issues.Err() != nil {
			return nil, fmt.Errorf("validation[%d].expression: type-check error: %s", i, issues.Err())
		}
		if ast.OutputType() != cel.BoolType {
			return nil, fmt.Errorf("validation[%d].expression: cel expression must evaluate to a bool", i)
		}
		prg, err := env.Program(ast,
			cel.EvalOptions(cel.OptOptimize),
			cel.OptimizeRegex(k8scellib.ExtensionLibRegexOptimizations...),
			cel.InterruptCheckFrequency(100),
		)
		if err != nil {
			return nil, fmt.Errorf("validation[%d].expression: program construction error: %s", i, err)
		}
		prgs = append(prgs, prg)
	}
	apiVersions := make([]string, 0, len(apiResources))
	for _, resource := range apiResources {
		apiVersions = append(apiVersions, resource.GroupVersion)
	}
	return &CELValidator{check: check, programs: prgs, hasPodSpec: hasPodSpec, apiVersions: apiVersions, kubeVersion: kubeVersion}, nil
}

func newEnv(podSpec bool) (*cel.Env, error) {
	opts := []cel.EnvOption{
		cel.HomogeneousAggregateLiterals(),
		cel.EagerlyValidateDeclarations(true),
		cel.DefaultUTCTimeZone(true),
		cel.Variable(ObjectVarName, cel.DynType),
		cel.Variable(ParamsVarName, cel.DynType),
		cel.Variable(APIVersionsVarName, cel.ListType(cel.StringType)),
		cel.Variable(KubeVersionVarName, cel.DynType),
	}
	opts = append(opts, k8scellib.ExtensionLibs...)
	if podSpec {
		opts = append(opts,
			cel.Variable(PodMetaVarName, cel.DynType),
			cel.Variable(PodSpecVarName, cel.DynType),
			cel.Variable(AllContainersVarName, cel.ListType(cel.DynType)),
		)
	}
	return cel.NewEnv(opts...)
}
