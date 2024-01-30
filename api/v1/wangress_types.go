package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WangressControllerSpec defines the desired state of WangressController
type WangressControllerSpec struct {
	// todo Define specifications for routing, service discovery, etc.
	// Example: RoutingRules, ServiceDiscoveryConfig, SecurityConfig, etc.
}

// WangressControllerStatus defines the observed state of WangressController
type WangressControllerStatus struct {
	// todo Fields to represent the operational status of the IngressController
	// Example: AvailableEndpoints, SecurityStatus, etc.
}

// WangressController is the Schema for the Wangress Controllers API
type WangressController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WangressControllerSpec   `json:"spec,omitempty"`
	Status WangressControllerStatus `json:"status,omitempty"`
}

// Add schema information here
// ...

func init() {
	// Register the WangressController CRD
	// ...
}
