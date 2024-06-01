/*
SPDX-License-Identifier: AGPL-3.0-or-later
KubeWG - Wireguard in your Kubernetes cluster
Copyright (C) 2024 Jacob McSwain

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

The source code is available at <https://github.com/kubewg-net/operator>
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RouterSpec defines the desired state of Router
type RouterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:default=0

	// Replicas is the number of router replicas
	// This defaults to 0, the same as disabling the router
	//+optional
	Replicas int32 `json:"replicas,omitempty"`

	// Network is the selector for the network this peer is a part of
	Network NameSelectorSpec `json:"network"`

	//+kubebuilder:default="ghcr.io/kubewg-net/container:v0.0.0@sha256:c6f1c2fa01fe79caaa8de7f546e53e089091710bc0dab3fc676a4206424c8040"

	// Image is the container image for the router
	// This defaults to ghcr.io/kubewg-net/container:v0.0.0@sha256:c6f1c2fa01fe79caaa8de7f546e53e089091710bc0dab3fc676a4206424c8040
	Image string `json:"image,omitempty"`

	// DNS is the optional DNS configuration
	// This overrides the default DNS configuration from the Network
	//+optional
	DNS corev1.PodDNSConfig `json:"dns"`

	// ExternalVPN is the optional external VPN configuration
	// If specified, the router will route traffic through the external VPN
	// Paired with enabling the firewall, this can be used to create a VPN kill-switched
	// connection to an external VPN provider from all pods in the network
	//+optional
	ExternalVPN ExternalVPNSpec `json:"externalVPN,omitempty"`

	//+kubebuilder:default={enabled: true}

	// Firewall is the optional firewall configuration that is applied to the peer
	//+optional
	Firewall FirewallSpec `json:"firewall,omitempty"`
}

// RouterStatus defines the observed state of Router
type RouterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ready is a flag to indicate if the network is ready
	Ready bool `json:"ready"`

	// ID is the ID of the network
	ID string `json:"id,omitempty"`

	// Status is the status of the network
	Status uint8 `json:"status,omitempty"`

	// Replicas is the number of router replicas
	Replicas int32 `json:"replicas"`

	// Selector is the selector for scaling the router pods
	Selector string `json:"selector"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Network",type=string,JSONPath=`.spec.network.name`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector

// Router is the Schema for the routers API
type Router struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouterSpec   `json:"spec,omitempty"`
	Status RouterStatus `json:"status,omitempty"`
}

// ExternalVPNSpec defines the an external VPN connection
type ExternalVPNSpec struct {
	// Connection is the Wireguard connection configuration
	Connection WireguardConnectionSpec `json:"connection"`

	// Credentials are the external VPN Wireguard credentials
	Credentials WireguardCredentialsSpec `json:"credentials"`
}

//+kubebuilder:validation:MinLength=44
//+kubebuilder:validation:MaxLength=44

// WireguardKey is a 44-character base64-encoded Wireguard key
type WireguardKey string

// WireguardCredentialsSpec defines a set of Wireguard credentials
type WireguardCredentialsSpec struct {

	// PrivateKey is the 44-character private key for the Wireguard client in base64 format
	PrivateKey WireguardKey `json:"privateKey,omitempty"`

	// PeerPublicKey is the 44-character public key for the peer in base64 format
	PeerPublicKey WireguardKey `json:"peerPublicKey,omitempty"`

	// PreSharedKey is the optional pre-shared key for the Wireguard connection
	//+optional
	PreSharedKey string `json:"preSharedKey,omitempty"`

	// Secret is the name of the secret containing the Wireguard credentials in the keys "privateKey", "peerPublicKey", and "preSharedKey"
	//+optional
	Secret NameSelectorSpec `json:"secret,omitempty"`
}

// WireguardConnectionSpec defines a Wireguard connection
type WireguardConnectionSpec struct {
	// Address is the IP address or hostname of the Wireguard server
	//+optional
	Address string `json:"address,omitempty"`

	// Port is the port of the Wireguard server
	//+optional
	Port uint16 `json:"port,omitempty"`

	// Secret is the selector for the secret containing the Wireguard connection configuration in the keys "address" and "port"
	//+optional
	Secret NameSelectorSpec `json:"secret,omitempty"`
}

//+kubebuilder:object:root=true

// RouterList contains a list of Router
type RouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Router `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Router{}, &RouterList{})
}
