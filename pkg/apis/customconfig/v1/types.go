package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CustomConfig struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	Spec CustomConfigSpec `json:"spec"`
}

// CustomConfigSpec is the spec for a CustomConfig resource
type CustomConfigSpec struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	ConfigmapName string `json:"configmapName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CustomConfigList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []CustomConfig `json:"items"`
}
