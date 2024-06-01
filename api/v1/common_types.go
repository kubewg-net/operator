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

// NameSelectorSpec defines a name selector for a resource
type NameSelectorSpec struct {
	// Name is the name of the resource
	Name string `json:"name"`
}

// InitSpec defines the initial container configuration
type InitSpec struct {
	//+kubebuilder:default="ghcr.io/kubewg-net/container:v0.0.0@sha256:c6f1c2fa01fe79caaa8de7f546e53e089091710bc0dab3fc676a4206424c8040"

	// Image is the container image
	// If not specified, the default image of ghcr.io/kubewg-net/container:v0.0.0@sha256:c6f1c2fa01fe79caaa8de7f546e53e089091710bc0dab3fc676a4206424c8040 is used
	//+optional
	Image string `json:"image,omitempty"`
}

// FirewallSpec defines the firewall configuration for a container
type FirewallSpec struct {
	//+kubebuilder:default=true

	// Enabled is a flag to enable the firewall.
	// The default firewall configuration is to block all non-VPN traffic, aka a kill switch.
	//+optional
	Enabled bool `json:"enabled"`

	// AllowWorkloadNetworkChanges is a flag to allow pods that could potentially make changes to the workload network
	// This is disabled by default and will reject any containers with the NET_RAW or NET_ADMIN capabilities as
	// these capabilities can be used to make changes to the network. Enabling this flag will allow containers with
	// these capabilities to be deployed.
	//+optional
	AllowWorkloadNetworkChanges bool `json:"allowWorkloadNetworkChanges,omitempty"`

	// Egress is a list of egress firewall rules
	// These rules are applied to traffic leaving the container
	// The default egress rules are to block all RFC1918 IPs and allow all other traffic
	//+optional
	Egress []FirewallRulesSpec `json:"egress,omitempty"`

	// Ingress is a list of ingress firewall rules
	// These rules are applied to traffic entering the container
	// The default ingress rules are to block all traffic
	//+optional
	Ingress []FirewallRulesSpec `json:"ingress,omitempty"`
}

//+kubebuilder:validation:Enum=UDP;TCP;ICMP;ALL

// Protocol defines a network protocol
type Protocol string

// FirewallRuleSpec defines a firewall rule
type FirewallRuleSpec struct {
	//+kubebuilder:default=ALL

	// Protocol is the network protocol
	// If not specified, the default protocol of ALL is used
	//+optional
	Protocol Protocol `json:"protocol"`

	//+kubebuilder:validation:Minimum=1024
	//+kubebuilder:validation:Maximum=65535

	// StartPort is the start port for a range of ports
	// If the end port is not specified, the default end port is the same as the start port
	StartPort uint16 `json:"startPort"`

	//+kubebuilder:validation:Minimum=1024
	//+kubebuilder:validation:Maximum=65535

	// EndPort is the end port for a range of ports
	// If not specified, the default end port is the same as the start port
	//+optional
	EndPort uint16 `json:"endPort,omitempty"`

	// IP is the IP address of the subject
	// Either an IP or CIDR must be specified
	//+optional
	IP string `json:"ip,omitempty"`

	// CIDR is the CIDR block of the subject
	// Either an IP or CIDR must be specified
	//+optional
	CIDR string `json:"cidr,omitempty"`
}

// FirewallRulesSpec defines a list of firewall rules
type FirewallRulesSpec struct {
	// Allow is a list of firewall rules to allow traffic
	//+optional
	Allow []FirewallRuleSpec `json:"allow,omitempty"`

	// Block is a list of firewall rules to block traffic
	//+optional
	Block []FirewallRuleSpec `json:"block,omitempty"`
}
