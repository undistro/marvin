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
	"github.com/google/cel-go/ext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	k8scellib "k8s.io/apiserver/pkg/cel/library"

	"github.com/undistro/marvin/pkg/types"
)

var baseEnvOptions = []cel.EnvOption{
	cel.HomogeneousAggregateLiterals(),
	cel.EagerlyValidateDeclarations(true),
	cel.DefaultUTCTimeZone(true),
	cel.CrossTypeNumericComparisons(true),
	cel.OptionalTypes(),

	ext.Strings(ext.StringsVersion(2)),
	ext.Sets(),

	k8scellib.URLs(),
	k8scellib.Regex(),
	k8scellib.Lists(),
	k8scellib.Quantity(),

	cel.Variable(ObjectVarName, cel.DynType),
	cel.Variable(APIVersionsVarName, cel.ListType(cel.StringType)),
	cel.Variable(KubeVersionVarName, cel.DynType),
}

var programOptions = []cel.ProgramOption{
	cel.EvalOptions(cel.OptOptimize),
	cel.CostLimit(1000000),
	cel.InterruptCheckFrequency(100),
}

var podSpecEnvOptions = []cel.EnvOption{
	cel.Variable(PodMetaVarName, cel.DynType),
	cel.Variable(PodSpecVarName, cel.DynType),
	cel.Variable(AllContainersVarName, cel.ListType(cel.DynType)),
}

// Compile compiles variables and expressions of the given check and returns a Validator
func Compile(check types.Check, apiResources []*metav1.APIResourceList, kubeVersion *version.Info) (Validator, error) {
	if len(check.Validations) == 0 {
		return nil, errors.New("invalid check: a check must have at least 1 validation")
	}
	env, err := newEnv(check)
	if err != nil {
		return nil, fmt.Errorf("environment construction error %s", err.Error())
	}

	variables, err := compileVariables(env, check.Variables)

	prgs, err := compileValidations(env, check.Validations)

	apiVersions := make([]string, 0, len(apiResources))
	for _, resource := range apiResources {
		apiVersions = append(apiVersions, resource.GroupVersion)
	}
	return &CELValidator{check: check, programs: prgs, apiVersions: apiVersions, kubeVersion: kubeVersion, variables: variables}, nil
}

func newEnv(check types.Check) (*cel.Env, error) {
	opts := baseEnvOptions
	if MatchesPodSpec(check.Match.Resources) {
		opts = append(opts, podSpecEnvOptions...)
	}
	if len(check.Variables) > 0 {
		opts = append(opts, cel.Variable(VariableVarName, cel.MapType(cel.StringType, cel.DynType)))
	}
	if len(check.Params) > 0 {
		opts = append(opts, cel.Variable(ParamsVarName, cel.DynType))
	}
	return cel.NewEnv(opts...)
}

func compileVariables(env *cel.Env, vars []types.Variable) ([]compiledVariables, error) {
	variables := make([]compiledVariables, 0, len(vars))
	for _, v := range vars {
		prg, err := compileExpression(env, v.Expression, cel.AnyType)
		if err != nil {
			return nil, fmt.Errorf("variables[%q].expression: %s", v.Name, err)
		}
		variables = append(variables, compiledVariables{name: v.Name, program: prg})
	}
	return variables, nil
}

func compileValidations(env *cel.Env, vals []types.Validation) ([]cel.Program, error) {
	prgs := make([]cel.Program, 0, len(vals))
	for i, v := range vals {
		prg, err := compileExpression(env, v.Expression, cel.BoolType)
		if err != nil {
			return nil, fmt.Errorf("validations[%d].expression: %s", i, err)
		}
		prgs = append(prgs, prg)
	}
	return prgs, nil
}

func compileExpression(env *cel.Env, exp string, allowedTypes ...*cel.Type) (cel.Program, error) {
	ast, issues := env.Compile(exp)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("type-check error: %s", issues.Err())
	}
	found := false
	for _, t := range allowedTypes {
		if ast.OutputType() == t || cel.AnyType == t {
			found = true
			break
		}
	}
	if !found {
		if len(allowedTypes) == 1 {
			return nil, fmt.Errorf("must evaluate to %v", allowedTypes[0].String())
		}
		return nil, fmt.Errorf("must evaluate to one of %v", allowedTypes)
	}
	prg, err := env.Program(ast, programOptions...)
	if err != nil {
		return nil, fmt.Errorf("program construction error: %s", err)
	}
	return prg, nil
}
