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

// NameSelectorSpec defines a name selector for a resource
type NameSelectorSpec struct {
	Name string `json:"name"`
}

// DNSSpec defines the DNS configuration for a container
type DNSSpec struct {
	Nameservers []string `json:"nameservers,omitempty"`
}

// InitSpec defines the initial container configuration
type InitSpec struct {
	Image string  `json:"image,omitempty"`
	DNS   DNSSpec `json:"dns,omitempty"`
}

// FirewallSpec defines the firewall configuration for a container
type FirewallSpec struct {
	Enabled                     bool                `json:"enabled,omitempty"`
	AllowWorkloadNetworkChanges bool                `json:"allowWorkloadNetworkChanges,omitempty"`
	Egress                      []FirewallRulesSpec `json:"egress,omitempty"`
	Ingress                     []FirewallRulesSpec `json:"ingress,omitempty"`
}

// FirewallRuleSpec defines a firewall rule
type FirewallRuleSpec struct {
	Protocol  string `json:"protocol"`
	StartPort uint16 `json:"startPort"`
	EndPort   uint16 `json:"endPort,omitempty"`
	IP        string `json:"ip,omitempty"`
	CIDR      string `json:"cidr,omitempty"`
}

// FirewallRulesSpec defines a list of firewall rules
type FirewallRulesSpec struct {
	Allow []FirewallRuleSpec `json:"allow,omitempty"`
	Block []FirewallRuleSpec `json:"block,omitempty"`
}
