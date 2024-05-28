package router

import (
	kubewgv1 "github.com/kubewg-net/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NetworkPolicyFromRouterCRD(router *kubewgv1.Router, network *kubewgv1.Network, labels *client.MatchingLabels) *networkingv1.NetworkPolicy {
	var dnsConfig *corev1.PodDNSConfig
	if router.Spec.DNS.Nameservers != nil || router.Spec.DNS.Searches != nil || router.Spec.DNS.Options != nil {
		dnsConfig = &router.Spec.DNS
	} else {
		dnsConfig = &network.Spec.DNS
	}

	dnsEgress := []networkingv1.NetworkPolicyEgressRule{}
	if dnsConfig != nil {
		for _, nameserver := range dnsConfig.Nameservers {
			dnsEgress = append(dnsEgress, networkingv1.NetworkPolicyEgressRule{
				To: []networkingv1.NetworkPolicyPeer{
					{
						IPBlock: &networkingv1.IPBlock{
							CIDR: nameserver + "/32",
						},
					},
				},
			})
		}
	}
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubewg-" + router.GetName() + "-" + string(router.GetUID()),
			Namespace: router.GetNamespace(),
			Labels:    *labels,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: *labels,
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/managed-by": "kubewg",
								},
							},
						},
					},
				},
			},
			Egress: dnsEgress,
		},
	}
}
