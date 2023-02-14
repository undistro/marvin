package validator

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/version"
)

type Validator interface {
	Validate(obj unstructured.Unstructured, params any) (bool, string, error)
	Matches(obj unstructured.Unstructured, resource string) bool
	SetAPIVersions(apiVersions []string)
	SetKubeVersion(v *version.Info)
}
