package router

import (
	kubewgv1 "github.com/kubewg-net/operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeploymentFromRouterCRD(router *kubewgv1.Router, network *kubewgv1.Network, labels *client.MatchingLabels) *appsv1.Deployment {
	var dnsConfig *corev1.PodDNSConfig
	if router.Spec.DNS.Nameservers != nil || router.Spec.DNS.Searches != nil || router.Spec.DNS.Options != nil {
		dnsConfig = &router.Spec.DNS
	} else {
		dnsConfig = &network.Spec.DNS
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubewg-" + router.GetName() + "-" + string(router.GetUID()),
			Namespace: router.GetNamespace(),
			Labels:    *labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &router.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: *labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: *labels,
				},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						Sysctls: []corev1.Sysctl{
							{
								Name:  "net.ipv4.conf.all.src_valid_mark",
								Value: "1",
							},
							{
								Name:  "net.ipv4.ip_forward",
								Value: "1",
							},
						},
					},
					DNSPolicy: corev1.DNSClusterFirst,
					DNSConfig: dnsConfig,
					Containers: []corev1.Container{
						{
							Name:  "router",
							Image: router.Spec.Image,
							Env: []corev1.EnvVar{
								{
									Name:  "KUBEWG_NETWORK",
									Value: router.Spec.Network.Name,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "wireguard",
									ContainerPort: 51820,
									Protocol:      corev1.ProtocolUDP,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{
										"NET_ADMIN",
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/wireguard",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "kubewg-" + router.GetName() + "-" + string(router.GetUID()),
								},
							},
						},
					},
				},
			},
		},
	}
}
