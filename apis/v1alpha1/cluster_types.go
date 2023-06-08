package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterSpec struct {
	// AccessObjectRefs represents references to objects providing access info to the cluster.
	// It could be a kubeconf stored in a secret
	AccessObjectRefs []AccessObjectRef `json:"accessObjectRef,omitempty"`

	// HealthProbe is used to coordinate the heartbeat time of to check the healthiness of the cluster.
	HealthProbe HealthProbe `json:"healthProbe"`

	// Taints is a property of cluster that allow the cluster to be repelled when scheduling.
	// +optional
	Taints []Taint `json:"taints,omitempty"`
}

type HealthProbe struct {
	// HeartbeatIntervalSeconds is the interval of the cluster's heartbeat to check the
	// availability of the cluster.
	HeartbeatIntervalSeconds int32 `json:"heatbeatIntervalSeconds"`
}

type AccessObjectRef struct {
	// Type is type of the access info. If the type is KUBECONFIG, the realted object
	// should be a secret containing kubeconfig key.
	Type string `json:"type"`

	// Group is the API Group of the Kubernetes resource,
	// empty string indicates it is in core group.
	// +optional
	Group string `json:"group"`

	// Resource is the resource name of the Kubernetes resource.
	// +kubebuilder:validation:Required
	// +required
	Resource string `json:"resource"`

	// Name is the name of the Kubernetes resource.
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// Name is the namespace of the Kubernetes resource, empty string indicates
	// it is a cluster scoped resource.
	// +optional
	Namespace string `json:"namespace"`
}

// The managed cluster this Taint is attached to has the "effect" on
// any placement that does not tolerate the Taint.
type Taint struct {
	// Key is the taint key applied to a cluster. e.g. bar or foo.example.com/bar.
	// The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:MaxLength=316
	// +required
	Key string `json:"key"`
	// Value is the taint value corresponding to the taint key.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Value string `json:"value,omitempty"`
	// Effect indicates the effect of the taint
	// Valid effects are NoSelect, PreferNoSelect and NoSelectIfNew.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=NoSelect;PreferNoSelect;NoSelectIfNew
	// +required
	Effect TaintEffect `json:"effect"`
	// TimeAdded represents the time at which the taint was added.
	// +nullable
	// +required
	TimeAdded metav1.Time `json:"timeAdded"`
}

type TaintEffect string

const (
	// TaintEffectNoSelect means not allowed to select the cluster unless tolerating the taint.
	// The cluster will be removed from the scheduler decisions if the scheduler has already selected
	// this cluster.
	TaintEffectNoSelect TaintEffect = "NoSelect"
	// TaintEffectPreferNoSelect means the scheduler tries not to select the cluster, rather than prohibiting
	// from selecting the cluster entirely.
	TaintEffectPreferNoSelect TaintEffect = "PreferNoSelect"
	// TaintEffectNoSelectIfNew means scheduler are not allowed to select the cluster unless
	// 1) they tolerate the taint;
	// 2) they have already had the cluster in their scheduler decisions;
	TaintEffectNoSelectIfNew TaintEffect = "NoSelectIfNew"
)

type ClusterStatus struct {
	// Conditions contains the different condition statuses for this cluster.
	Conditions []metav1.Condition `json:"conditions"`

	// Version represents the kubernetes version of the cluster.
	Version ClusterVersion `json:"version,omitempty"`

	// Resource represents the resource of the cluster.
	Resources Resources `json:"resources,omitempty"`

	// Properties represents properties of collected from the cluster,
	// for example a unique cluster identifier (id.k8s.io).
	// The set of properties is not uniform across a fleet, some properties can be
	// vendor or version specific and may not be included from all clusters.
	// +optional
	Properties []Property `json:"properties,omitempty"`
}

// ManagedClusterVersion represents version information about the cluster.
type ClusterVersion struct {
	// Kubernetes is the kubernetes version of managed cluster.
	// +optional
	Kubernetes string `json:"kubernetes,omitempty"`
}

type Resources struct {
	// Capacity represents the total resource capacity from all nodeStatuses
	// on the cluster.
	Capacity ResourceList `json:"capacity,omitempty"`

	// Allocatable represents the total allocatable resources on the cluster.
	Allocatable ResourceList `json:"allocatable,omitempty"`
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

const (
	// ResourceCPU defines the number of CPUs in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// ResourceMemory defines the amount of memory in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
)

// ResourceList defines a map for the quantity of different resources, the definition
// matches the ResourceList defined in k8s.io/api/core/v1.
type ResourceList map[ResourceName]resource.Quantity

// Property represents a Property collected from a cluster.
type Property struct {
	// Name is the name of a propertie resource on cluster. It's a well known
	// or customized name to identify the propertie.
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// Value is a property-dependent string
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:MinLength=1
	Value string `json:"value,omitempty"`
}

const (
	// ClusterConditionJoined means the cluster has successfully joined the control.
	ClusterConditionJoined string = "Joined"
	// Available means the cluster is available.
	ClusterConditionAvailable string = "Available"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Cluster is the Schema for the cluster inventory API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec defines the spec of a cluster.
	// +optional
	Spec ClusterSpec `json:"spec,omitempty"`
	// status defines the status of a cluster.
	Status ClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Clusters
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of clusters.
	// +listType=set
	Items []Cluster `json:"items"`
}
