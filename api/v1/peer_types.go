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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PeerSpec defines the desired state of Peer
type PeerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Network is the selector for the network this peer is a part of
	Network NameSelectorSpec `json:"network"`

	// Pods is the selector for the pods that are peers in the network
	Pods metav1.LabelSelector `json:"pods"`

	// Init is the optional initial container configuration that is applied to the peer
	//+optional
	Init InitSpec `json:"init,omitempty"`

	//+kubebuilder:default={enabled: true}

	// Firewall is the optional firewall configuration that is applied to the peer
	//+optional
	Firewall FirewallSpec `json:"firewall,omitempty"`
}

// PeerStatus defines the observed state of Peer
type PeerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ready is a flag to indicate if the peer is ready
	Ready bool `json:"ready"`

	// ID is the ID of the peer
	ID string `json:"id,omitempty"`

	// Status is the status of the peer
	Status uint8 `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Peer is the Schema for the peers API
type Peer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PeerSpec   `json:"spec,omitempty"`
	Status PeerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:printcolumn:name="Network",type=string,JSONPath=`.spec.network.name`
//+kubebuilder:printcolumn:name="Firewalled",type=boolean,JSONPath=`.spec.firewall.enabled`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// PeerList contains a list of Peer
type PeerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Peer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Peer{}, &PeerList{})
}
