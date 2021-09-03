package servicebinding

import (
	corev1 "k8s.io/api/core/v1"
)

type ServiceBinding struct {
	Name string
	SecretRef *corev1.LocalObjectReference
}

func (s *ServiceBinding) ServiceName() string {
	return s.Name
}

type V1Alpha1ServiceBinding struct {
	Name string
	SecretRef *corev1.LocalObjectReference
	MetadataRef *corev1.LocalObjectReference
}

func (v *V1Alpha1ServiceBinding) ServiceName() string {
	return v.Name
}
