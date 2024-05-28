package router

import (
	kubewgv1 "github.com/kubewg-net/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SecretFromRouterCRD(router *kubewgv1.Router, network *kubewgv1.Network, labels *client.MatchingLabels) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubewg-" + router.GetName() + "-" + string(router.GetUID()),
			Namespace: router.GetNamespace(),
			Labels:    *labels,
		},
	}
}
