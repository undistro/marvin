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
