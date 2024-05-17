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

The source code is available at <https://github.com/USA-RedDragon/kubewg>
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NetworkSpec defines the desired state of Network
type NetworkSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Router   RouterSpec   `json:"router,omitempty"`
	Init     InitSpec     `json:"init,omitempty"`
	Firewall FirewallSpec `json:"firewall,omitempty"`
}

// NetworkStatus defines the observed state of Network
type NetworkStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Ready  bool   `json:"ready"`
	ID     string `json:"id,omitempty"`
	Status uint8  `json:"status,omitempty"`
}

type RouterSpec struct {
	Replicas    uint32          `json:"replicas,omitempty"`
	Image       string          `json:"image,omitempty"`
	ExternalVPN ExternalVPNSpec `json:"externalVPN,omitempty"`
}

type ExternalVPNSpec struct {
	Connection  WireguardConnectionSpec  `json:"connection"`
	Credentials WireguardCredentialsSpec `json:"credentials"`
}

type WireguardCredentialsSpec struct {
	PrivateKey    string           `json:"privateKey,omitempty"`
	PeerPublicKey string           `json:"peerPublicKey,omitempty"`
	PreSharedKey  string           `json:"preSharedKey,omitempty"`
	Secret        NameSelectorSpec `json:"secret,omitempty"`
}

type WireguardConnectionSpec struct {
	Address string           `json:"address,omitempty"`
	Port    uint16           `json:"port,omitempty"`
	Secret  NameSelectorSpec `json:"secret,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Network is the Schema for the networks API
type Network struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkSpec   `json:"spec,omitempty"`
	Status NetworkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NetworkList contains a list of Network
type NetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Network `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Network{}, &NetworkList{})
}
