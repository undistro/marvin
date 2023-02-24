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

package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestResourceRule_ToGVR(t *testing.T) {
	type fields struct {
		Group    string
		Version  string
		Resource string
	}
	tests := []struct {
		name   string
		fields fields
		want   schema.GroupVersionResource
	}{
		{
			name: "services",
			fields: fields{
				Group:    "",
				Version:  "v1",
				Resource: "services",
			},
			want: corev1.SchemeGroupVersion.WithResource("services"),
		},
		{
			name: "deployments",
			fields: fields{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			},
			want: appsv1.SchemeGroupVersion.WithResource("deployments"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceRule{
				Group:    tt.fields.Group,
				Version:  tt.fields.Version,
				Resource: tt.fields.Resource,
			}
			assert.Equalf(t, tt.want, r.ToGVR(), "ToGVR()")
		})
	}
}
