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

package validator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/undistro/marvin/pkg/types"
)

// MatchesPodSpec returns true if any rule matches a Pod spec
func MatchesPodSpec(rules []types.ResourceRule) bool {
	for _, r := range rules {
		gr := r.ToGVR().GroupResource()
		if defaultPodSpecResources[gr] {
			return true
		}
	}
	return false
}

var defaultPodSpecResources = map[schema.GroupResource]bool{
	corev1.Resource("pods"):                   true,
	corev1.Resource("replicationcontrollers"): true,
	corev1.Resource("podtemplates"):           true,
	appsv1.Resource("replicasets"):            true,
	appsv1.Resource("deployments"):            true,
	appsv1.Resource("statefulsets"):           true,
	appsv1.Resource("daemonsets"):             true,
	batchv1.Resource("jobs"):                  true,
	batchv1.Resource("cronjobs"):              true,
}

// HasPodSpec returns true if the given object has a Pod spec
func HasPodSpec(u unstructured.Unstructured) bool {
	gk := u.GroupVersionKind().GroupKind()
	_, ok := defaultPodSpecTypes[gk]
	return ok
}

// ExtractPodSpec returns the metadata and Pod spec from the given object
func ExtractPodSpec(u unstructured.Unstructured) (*metav1.ObjectMeta, *corev1.PodSpec, error) {
	gk := u.GroupVersionKind().GroupKind()
	obj, ok := defaultPodSpecTypes[gk]
	if !ok {
		return nil, nil, fmt.Errorf("unexpected object type: %s", u.GetObjectKind().GroupVersionKind().String())
	}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), obj)
	if err != nil {
		return nil, nil, fmt.Errorf("from unstructured converter error: %s", err.Error())
	}
	switch o := obj.(type) {
	case *corev1.Pod:
		return &o.ObjectMeta, &o.Spec, nil
	case *corev1.PodTemplate:
		return extractPodSpecFromTemplate(&o.Template)
	case *corev1.ReplicationController:
		return extractPodSpecFromTemplate(o.Spec.Template)
	case *appsv1.ReplicaSet:
		return extractPodSpecFromTemplate(&o.Spec.Template)
	case *appsv1.Deployment:
		return extractPodSpecFromTemplate(&o.Spec.Template)
	case *appsv1.DaemonSet:
		return extractPodSpecFromTemplate(&o.Spec.Template)
	case *appsv1.StatefulSet:
		return extractPodSpecFromTemplate(&o.Spec.Template)
	case *batchv1.Job:
		return extractPodSpecFromTemplate(&o.Spec.Template)
	case *batchv1.CronJob:
		return extractPodSpecFromTemplate(&o.Spec.JobTemplate.Spec.Template)
	default:
		return nil, nil, fmt.Errorf("unexpected object type: %s", u.GetObjectKind().GroupVersionKind().String())
	}
}

var defaultPodSpecTypes = map[schema.GroupKind]any{
	corev1.SchemeGroupVersion.WithKind("Pod").GroupKind():                   &corev1.Pod{},
	corev1.SchemeGroupVersion.WithKind("ReplicationController").GroupKind(): &corev1.ReplicationController{},
	corev1.SchemeGroupVersion.WithKind("PodTemplate").GroupKind():           &corev1.PodTemplate{},
	appsv1.SchemeGroupVersion.WithKind("ReplicaSet").GroupKind():            &appsv1.ReplicaSet{},
	appsv1.SchemeGroupVersion.WithKind("Deployment").GroupKind():            &appsv1.Deployment{},
	appsv1.SchemeGroupVersion.WithKind("StatefulSet").GroupKind():           &appsv1.StatefulSet{},
	appsv1.SchemeGroupVersion.WithKind("DaemonSet").GroupKind():             &appsv1.DaemonSet{},
	batchv1.SchemeGroupVersion.WithKind("Job").GroupKind():                  &batchv1.Job{},
	batchv1.SchemeGroupVersion.WithKind("CronJob").GroupKind():              &batchv1.CronJob{},
}

func extractPodSpecFromTemplate(template *corev1.PodTemplateSpec) (*metav1.ObjectMeta, *corev1.PodSpec, error) {
	if template == nil {
		return nil, nil, nil
	}
	return &template.ObjectMeta, &template.Spec, nil
}

func extractAllContainers(podSpec *corev1.PodSpec) []corev1.Container {
	containers := append(podSpec.Containers, podSpec.InitContainers...)
	for _, ec := range podSpec.EphemeralContainers {
		c := corev1.Container{
			Name:                     ec.Name,
			Image:                    ec.Image,
			Command:                  ec.Command,
			Args:                     ec.Args,
			WorkingDir:               ec.WorkingDir,
			Ports:                    ec.Ports,
			EnvFrom:                  ec.EnvFrom,
			Env:                      ec.Env,
			Resources:                ec.Resources,
			VolumeMounts:             ec.VolumeMounts,
			VolumeDevices:            ec.VolumeDevices,
			LivenessProbe:            ec.LivenessProbe,
			ReadinessProbe:           ec.ReadinessProbe,
			StartupProbe:             ec.StartupProbe,
			Lifecycle:                ec.Lifecycle,
			TerminationMessagePath:   ec.TerminationMessagePath,
			TerminationMessagePolicy: ec.TerminationMessagePolicy,
			ImagePullPolicy:          ec.ImagePullPolicy,
			SecurityContext:          ec.SecurityContext,
			Stdin:                    ec.Stdin,
			StdinOnce:                ec.StdinOnce,
			TTY:                      ec.TTY,
		}
		containers = append(containers, c)
	}
	return containers
}
