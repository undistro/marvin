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

package types

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/version"
)

type Report struct {
	KubeVersion *version.Info     `json:"kubeVersion"`
	Checks      []*CheckResult    `json:"checks"`
	GVRs        map[string]string `json:"gvrs,omitempty"`
}

func NewReport(kubeVersion *version.Info) *Report {
	return &Report{KubeVersion: kubeVersion}
}

func (r *Report) Add(cr *CheckResult) {
	r.Checks = append(r.Checks, cr)
}

func (r *Report) HasError() bool {
	for _, check := range r.Checks {
		if len(check.Errors) > 0 {
			return true
		}
	}
	return false
}

type CheckResult struct {
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Builtin  bool     `json:"builtin"`
	Path     string   `json:"path"`

	Status  CheckStatus         `json:"status"`
	Failed  map[string][]string `json:"failed"`
	Passed  map[string][]string `json:"passed"`
	Skipped map[string][]string `json:"skipped"`
	Errors  []string            `json:"errors"`

	TotalFailed  int `json:"totalFailed"`
	TotalPassed  int `json:"totalPassed"`
	TotalSkipped int `json:"totalSkipped"`
}

func NewCheckResult(check Check) *CheckResult {
	return &CheckResult{
		ID:       check.ID,
		Message:  check.Message,
		Severity: check.Severity,
		Builtin:  check.Builtin,
		Path:     check.Path,

		Failed:  map[string][]string{},
		Passed:  map[string][]string{},
		Skipped: map[string][]string{},
		Errors:  []string{},
	}
}

func (r *CheckResult) AddFailed(obj unstructured.Unstructured) {
	addResource(obj, r.Failed)
	r.TotalFailed++
}

func (r *CheckResult) AddPassed(obj unstructured.Unstructured) {
	addResource(obj, r.Passed)
	r.TotalPassed++
}

func (r *CheckResult) AddSkipped(obj unstructured.Unstructured) {
	addResource(obj, r.Skipped)
	r.TotalSkipped++
}

func (r *CheckResult) AddError(err error) {
	r.Errors = append(r.Errors, err.Error())
}

func (r *CheckResult) AddErrors(errs ...error) {
	for _, err := range errs {
		r.AddError(err)
	}
}

func (r *CheckResult) UpdateStatus() {
	if len(r.Errors) > 0 {
		r.Status = StatusError
		return
	}
	if len(r.Failed) > 0 {
		r.Status = StatusFailed
		return
	}
	if len(r.Passed) == 0 && len(r.Skipped) > 0 {
		r.Status = StatusSkipped
		return
	}
	r.Status = StatusPassed
	return
}

// GVK returns the GroupVersionKind string of the given resource
func GVK(obj unstructured.Unstructured) string {
	gvk := obj.GroupVersionKind()
	return fmt.Sprintf("%s/%s", gvk.GroupVersion().String(), gvk.Kind)
}

// NamespacedName returns the namespaced name string of the given resource
func NamespacedName(obj unstructured.Unstructured) string {
	if len(obj.GetNamespace()) > 0 {
		return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
	}
	return obj.GetName()
}

func addResource(obj unstructured.Unstructured, m map[string][]string) {
	k := GVK(obj)
	v := NamespacedName(obj)
	if _, ok := m[k]; ok {
		m[k] = append(m[k], v)
	} else {
		m[k] = []string{v}
	}
}
