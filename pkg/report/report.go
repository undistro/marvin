package report

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/version"

	"github.com/undistro/marvin/pkg/checks"
)

type Report struct {
	KubeVersion *version.Info  `json:"kubeVersion"`
	Checks      []*CheckResult `json:"checks"`
}

func New(kubeVersion *version.Info) *Report {
	return &Report{KubeVersion: kubeVersion}
}

func (r *Report) Add(cr *CheckResult) {
	r.Checks = append(r.Checks, cr)
}

type CheckResult struct {
	ID       string          `json:"id"`
	Message  string          `json:"message"`
	Severity checks.Severity `json:"severity"`
	Builtin  bool            `json:"builtin"`
	Path     string          `json:"path"`

	Passed  map[string][]string `json:"passed"`
	Failed  map[string][]string `json:"failed"`
	Skipped map[string][]string `json:"skipped"`
	Errors  []string            `json:"errors"`
}

func NewCheckResult(c checks.Check) *CheckResult {
	return &CheckResult{
		ID:       c.ID,
		Message:  c.Message,
		Severity: c.Severity,
		Builtin:  c.Builtin,
		Path:     c.Path,

		Passed:  map[string][]string{},
		Failed:  map[string][]string{},
		Skipped: map[string][]string{},
		Errors:  []string{},
	}
}

func (r *CheckResult) AddPassed(obj unstructured.Unstructured) {
	k := key(obj)
	v := value(obj)
	if _, ok := r.Passed[k]; ok {
		r.Passed[k] = append(r.Passed[k], v)
	} else {
		r.Passed[k] = []string{v}
	}
}

func (r *CheckResult) AddFailed(obj unstructured.Unstructured) {
	k := key(obj)
	v := value(obj)
	if _, ok := r.Failed[k]; ok {
		r.Failed[k] = append(r.Failed[k], v)
	} else {
		r.Failed[k] = []string{v}
	}
}

func (r *CheckResult) AddSkipped(obj unstructured.Unstructured) {
	k := key(obj)
	v := value(obj)
	if _, ok := r.Skipped[k]; ok {
		r.Skipped[k] = append(r.Skipped[k], v)
	} else {
		r.Skipped[k] = []string{v}
	}
}

func (r *CheckResult) AddError(err error) {
	r.Errors = append(r.Errors, err.Error())
}

func key(obj unstructured.Unstructured) string {
	gvk := obj.GroupVersionKind()
	return fmt.Sprintf("%s/%s", gvk.GroupVersion().String(), gvk.Kind)
}

func value(obj unstructured.Unstructured) string {
	if len(obj.GetNamespace()) > 0 {
		return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
	}
	return obj.GetName()
}
