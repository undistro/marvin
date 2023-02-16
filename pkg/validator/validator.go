package validator

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"

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
