package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/undistro/marvin/pkg/checks"
)

func TestMatchesPodSpec(t *testing.T) {
	tests := []struct {
		name  string
		rules []checks.ResourceRule
		want  bool
	}{
		{
			name: "deployments",
			rules: []checks.ResourceRule{{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}},
			want: true,
		},
		{
			name: "pods and services",
			rules: []checks.ResourceRule{
				{
					Group:    "",
					Version:  "v1",
					Resource: "pods",
				},
				{
					Group:    "",
					Version:  "v1",
					Resource: "services",
				},
			},
			want: true,
		},
		{
			name: "services and cronjobs",
			rules: []checks.ResourceRule{
				{
					Group:    "",
					Version:  "v1",
					Resource: "services",
				},
				{
					Group:    "batch",
					Version:  "v1",
					Resource: "cronjobs",
				},
			},
			want: true,
		},
		{
			name: "services",
			rules: []checks.ResourceRule{{
				Group:    "",
				Version:  "v1",
				Resource: "services",
			}},
			want: false,
		},
		{
			name:  "empty",
			rules: []checks.ResourceRule{},
			want:  false,
		},
		{
			name:  "nil",
			rules: nil,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchesPodSpec(tt.rules); got != tt.want {
				t.Errorf("MatchesPodSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultExtractPodSpec(t *testing.T) {
	metadata := map[string]any{
		"name": "foo-pod",
	}
	expectedMeta := &v1.ObjectMeta{
		Name: "foo-pod",
	}
	spec := map[string]any{
		"containers": []map[string]any{{
			"name": "foo-container",
		}},
	}
	expectedSpec := &corev1.PodSpec{
		Containers: []corev1.Container{{Name: "foo-container"}},
	}
	objects := []map[string]any{
		{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata":   metadata,
			"spec":       spec,
		},
		{
			"apiVersion": "v1",
			"kind":       "PodTemplate",
			"metadata":   map[string]any{"name": "foo-template"},
			"template": map[string]any{
				"metadata": metadata,
				"spec":     spec,
			},
		},
		{
			"apiVersion": "v1",
			"kind":       "ReplicationController",
			"metadata":   map[string]any{"name": "foo-rc"},
			"spec": map[string]any{
				"template": map[string]any{
					"metadata": metadata,
					"spec":     spec,
				},
			},
		},
		{
			"apiVersion": "apps/v1",
			"kind":       "ReplicaSet",
			"metadata":   map[string]any{"name": "foo-rs"},
			"spec": map[string]any{
				"template": map[string]any{
					"metadata": metadata,
					"spec":     spec,
				},
			},
		},
		{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "foo-deployment"},
			"spec": map[string]any{
				"template": map[string]any{
					"metadata": metadata,
					"spec":     spec,
				},
			},
		},
		{
			"apiVersion": "apps/v1",
			"kind":       "StatefulSet",
			"metadata":   map[string]any{"name": "foo-ss"},
			"spec": map[string]any{
				"template": map[string]any{
					"metadata": metadata,
					"spec":     spec,
				},
			},
		},
		{
			"apiVersion": "apps/v1",
			"kind":       "DaemonSet",
			"metadata":   map[string]any{"name": "foo-ds"},
			"spec": map[string]any{
				"template": map[string]any{
					"metadata": metadata,
					"spec":     spec,
				},
			},
		},
		{
			"apiVersion": "batch/v1",
			"kind":       "Job",
			"metadata":   map[string]any{"name": "foo-job"},
			"spec": map[string]any{
				"template": map[string]any{
					"metadata": metadata,
					"spec":     spec,
				},
			},
		},
		{
			"apiVersion": "batch/v1",
			"kind":       "CronJob",
			"metadata":   map[string]any{"name": "foo-cronjob"},
			"spec": map[string]any{
				"jobTemplate": map[string]any{
					"spec": map[string]any{
						"template": map[string]any{
							"metadata": metadata,
							"spec":     spec,
						},
					},
				},
			},
		},
	}
	for _, obj := range objects {
		u := unstructured.Unstructured{Object: obj}
		name := u.GetName()
		actualMetadata, actualSpec, err := ExtractPodSpec(u)
		assert.NoError(t, err, name)
		assert.Equal(t, expectedMeta, actualMetadata, "%s: Metadata mismatch", name)
		assert.Equal(t, expectedSpec, actualSpec, "%s: PodSpec mismatch", name)
	}

	var service = map[string]any{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]any{
			"name": "foo-svc",
		},
	}
	_, _, err := ExtractPodSpec(unstructured.Unstructured{Object: service})
	assert.Error(t, err, "service should not have an extractable pod spec")
}
