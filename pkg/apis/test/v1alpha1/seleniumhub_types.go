package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SeleniumHubSpec defines the desired state of SeleniumHub
type SeleniumHubSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size   int32  `json:"size"`
	Memory string `json:"memory"`
	CPU    string `json:"cpu"`
}

// SeleniumHubStatus defines the observed state of SeleniumHub
type SeleniumHubStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SeleniumHub is the Schema for the seleniumhubs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=seleniumhubs,scope=Namespaced
type SeleniumHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SeleniumHubSpec   `json:"spec,omitempty"`
	Status SeleniumHubStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SeleniumHubList contains a list of SeleniumHub
type SeleniumHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SeleniumHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SeleniumHub{}, &SeleniumHubList{})
}
