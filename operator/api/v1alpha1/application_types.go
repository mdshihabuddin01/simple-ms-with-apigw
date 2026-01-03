/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	acmeissuerv1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// +optional
	Image *string `json:"image,omitempty"`

	ContainerPort int32 `json:"containerPort"`

	// +optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +optional
	Ingress *IngressSpec `json:"ingress,omitempty"`

	// +optional
	EnvVars []EnvVar `json:"envVars,omitempty"`

	// +optional
	Resources *ResourcesSpec `json:"resources,omitempty"`

	// +optional
	TLS *TLS `json:"tls,omitempty"`
}

type ServiceSpec struct {
	Port int32 `json:"port"`
	// +optional
	Type        string            `json:"type,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type IngressSpec struct {
	Host        string            `json:"host,omitempty"`
	Path        string            `json:"path"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResourcesSpec struct {
	// +optional
	Requests *ResourceLimits `json:"requests,omitempty"`
	// +optional
	Limits *ResourceLimits `json:"limits,omitempty"`
}

type ResourceLimits struct {
	// +optional
	CPU string `json:"cpu,omitempty"`
	// +optional
	Memory string `json:"memory,omitempty"`
}

type TLS struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Enable bool `json:"enable"`
	//Host   string  `json:"host,omitempty"`
	Issuer *Issuer `json:"issuer,omitempty"`
}
type Issuer struct {
	ACMEIssuer acmeissuerv1.ACMEIssuer `json:"acmeIssuer,omitempty"`
}

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the Application resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Application
	// +required
	Spec ApplicationSpec `json:"spec"`

	// status defines the observed state of Application
	// +optional
	Status ApplicationStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
