package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
type BuildStack struct {
	RunImage string `json:"runImage,omitempty"`
	ID       string `json:"id,omitempty"`
}

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
type BuildBuilderSpec struct {
	Image string `json:"image,omitempty"`
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`
}