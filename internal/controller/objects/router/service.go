package router

import (
	kubewgv1 "github.com/kubewg-net/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ServiceFromRouterCRD(router *kubewgv1.Router, labels *client.MatchingLabels) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubewg-" + router.GetName() + "-" + string(router.GetUID()),
			Namespace: router.GetNamespace(),
			Labels:    *labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: *labels,
			Ports: []corev1.ServicePort{
				{
					Name:     "wireguard",
					Protocol: corev1.ProtocolUDP,
					Port:     51820,
				},
			},
		},
	}
}
