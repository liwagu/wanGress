/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WanGressSpec defines the desired state of WanGress
type WanGressSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Hosts is a list of host names that the Ingress should match.
	Hosts []string `json:"hosts,omitempty"`

	// TLS configuration for the Ingress.
	TLS []WangressTLS `json:"tls,omitempty"`

	// 添加路由规则
	Routes []Route `json:"routes,omitempty"`
}

type Route struct {
	Path     string    `json:"path"`
	Services []Service `json:"services"`
}

type Service struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

// WangressTLS defines the TLS configuration for a host.
type WangressTLS struct {
	Hosts      []string `json:"hosts"`
	SecretName string   `json:"secretName"`
}

type WanGressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// WanGress is the Schema for the wangresses API
type WanGress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WanGressSpec   `json:"spec,omitempty"`
	Status WanGressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WanGressList contains a list of WanGress
type WanGressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WanGress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WanGress{}, &WanGressList{})
}
