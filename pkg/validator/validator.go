package validator

import (
	"errors"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	k8scellib "k8s.io/apiserver/pkg/cel/library"

	"github.com/undistro/marvin/pkg/checks"
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

type CELValidator struct {
	check       checks.Check
	programs    []cel.Program
	hasPodSpec  bool
	apiVersions []string
	kubeVersion *version.Info
}

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
	input := &activation{object: obj.UnstructuredContent(), apiVersions: r.apiVersions, params: params}
	if err := r.setPodSpecParams(obj, input); err != nil {
		return false, "", err
	}
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

func (r *CELValidator) setPodSpecParams(obj unstructured.Unstructured, input *activation) error {
	if !r.hasPodSpec || !HasPodSpec(obj) {
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

func newEnv(podSpec bool) (*cel.Env, error) {
	var opts []cel.EnvOption
	opts = append(opts, cel.HomogeneousAggregateLiterals())
	opts = append(opts, cel.EagerlyValidateDeclarations(true), cel.DefaultUTCTimeZone(true))
	opts = append(opts, k8scellib.ExtensionLibs...)
	opts = append(opts, cel.Variable(ObjectVarName, cel.DynType))
	opts = append(opts, cel.Variable(ParamsVarName, cel.DynType))
	opts = append(opts, cel.Variable(APIVersionsVarName, cel.ListType(cel.StringType)))
	opts = append(opts, cel.Variable(KubeVersionVarName, cel.DynType))
	if podSpec {
		opts = append(opts, cel.Variable(PodMetaVarName, cel.DynType))
		opts = append(opts, cel.Variable(PodSpecVarName, cel.DynType))
		opts = append(opts, cel.Variable(AllContainersVarName, cel.ListType(cel.DynType)))
	}
	return cel.NewEnv(opts...)
}
