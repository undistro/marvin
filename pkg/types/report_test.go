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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/version"
)

var (
	kubeVersion = &version.Info{
		Major:      "1",
		Minor:      "25",
		GitVersion: "v1.25.3",
	}
	check = Check{
		ID:       "foo",
		Severity: SeverityHigh,
		Message:  "bar",
		Builtin:  true,
		Path:     "path.yml",
	}
)

func TestReport(t *testing.T) {
	rep := NewReport(kubeVersion)
	assert.NotNil(t, rep)
	assert.NotNil(t, rep.KubeVersion)
	assert.Equal(t, "25", rep.KubeVersion.Minor)
	assert.Len(t, rep.Checks, 0)

	cr := NewCheckResult(check)
	assert.NotNil(t, cr)
	assert.Equal(t, "foo", cr.ID)
	assert.Equal(t, "bar", cr.Message)
	assert.Equal(t, SeverityHigh, cr.Severity)
	assert.True(t, cr.Builtin)
	assert.Equal(t, "path.yml", cr.Path)
	assert.Len(t, cr.Failed, 0)
	assert.Len(t, cr.Passed, 0)
	assert.Len(t, cr.Skipped, 0)
	assert.Len(t, cr.Errors, 0)

	cr.AddSkipped(obj("apps/v1", "Deployment", "ns", "skipped-deploy-1"))
	cr.AddSkipped(obj("apps/v1", "Deployment", "ns", "skipped-deploy-2"))
	cr.AddSkipped(obj("v1", "Pod", "", "skipped-pod"))

	assert.Len(t, cr.Skipped, 2)
	assert.Len(t, cr.Skipped["apps/v1/Deployment"], 2)
	assert.Len(t, cr.Skipped["v1/Pod"], 1)
	assert.Len(t, cr.Failed, 0)
	assert.Len(t, cr.Passed, 0)
	assert.Len(t, cr.Errors, 0)

	cr.UpdateStatus()

	assert.Equal(t, StatusSkipped, cr.Status)

	cr.AddPassed(obj("v1", "Pod", "ns", "passed-1"))
	cr.AddPassed(obj("v1", "Pod", "", "passed-2"))

	assert.Len(t, cr.Skipped, 2)
	assert.Len(t, cr.Passed, 1)
	assert.Len(t, cr.Passed["v1/Pod"], 2)
	assert.Equal(t, cr.Passed["v1/Pod"][0], "ns/passed-1")
	assert.Equal(t, cr.Passed["v1/Pod"][1], "passed-2")
	assert.Len(t, cr.Failed, 0)
	assert.Len(t, cr.Errors, 0)

	cr.UpdateStatus()

	assert.Equal(t, StatusPassed, cr.Status)

	cr.AddFailed(obj("batch/v1", "CronJob", "ns", "failed"))

	assert.Len(t, cr.Skipped, 2)
	assert.Len(t, cr.Passed, 1)
	assert.Len(t, cr.Failed, 1)
	assert.Len(t, cr.Failed["batch/v1/CronJob"], 1)
	assert.Len(t, cr.Errors, 0)

	cr.UpdateStatus()

	assert.Equal(t, StatusFailed, cr.Status)

	cr.AddError(errors.New("list deployments error"))
	assert.Len(t, cr.Skipped, 2)
	assert.Len(t, cr.Passed, 1)
	assert.Len(t, cr.Failed, 1)
	assert.Len(t, cr.Errors, 1)

	cr.UpdateStatus()

	assert.Equal(t, StatusError, cr.Status)

	rep.Add(cr)

	assert.Len(t, rep.Checks, 1)
}

func obj(apiVersion, kind, ns, name string) unstructured.Unstructured {
	return unstructured.Unstructured{Object: map[string]any{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata":   map[string]any{"namespace": ns, "name": name},
	}}
}
