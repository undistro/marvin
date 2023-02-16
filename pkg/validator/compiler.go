package validator

import (
	"errors"
	"fmt"

	"github.com/google/cel-go/cel"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	k8scellib "k8s.io/apiserver/pkg/cel/library"

	"github.com/undistro/marvin/pkg/checks"
)

func Compile(check checks.Check, apiResources []*metav1.APIResourceList, kubeVersion *version.Info) (Validator, error) {
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
			return nil, fmt.Errorf("type-check error on validation %d: %s", i, issues.Err())
		}
		if ast.OutputType() != cel.BoolType {
			return nil, fmt.Errorf("cel expression must evaluate to a bool on validation %d", i)
		}
		prg, err := env.Program(ast,
			cel.EvalOptions(cel.OptOptimize),
			cel.OptimizeRegex(k8scellib.ExtensionLibRegexOptimizations...),
			cel.InterruptCheckFrequency(100),
		)
		if err != nil {
			return nil, fmt.Errorf("program construction error on validation %d: %s", i, err)
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
