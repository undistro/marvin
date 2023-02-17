package checks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
)

type Check struct {
	ID          string         `json:"id"`
	Match       Match          `json:"match"`
	Validations []Validation   `json:"validations"`
	Params      map[string]any `json:"params"`
	Severity    Severity       `json:"severity"`
	Message     string         `json:"message"`

	Builtin bool   `json:"builtin"`
	Path    string `json:"path"`
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

type Test struct {
	Name        string        `json:"name"`
	Input       string        `json:"input"`
	Params      any           `json:"params"`
	APIVersions []string      `json:"apiVersions"`
	KubeVersion *version.Info `json:"kubeVersion"`
	Pass        bool          `json:"pass"`
	Message     string        `json:"message"`
}
