package fileserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func TestAssembleService(t *testing.T) {
	testCases := []struct {
		name            string
		datalayer       k8schianetv1.ChiaDataLayer
		expectedService struct {
			name      string
			namespace string
			ports     []corev1.ServicePort
		}
	}{
		{
			name: "With Default Config",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					FileserverConfig: k8schianetv1.FileserverConfig{
						Service: k8schianetv1.Service{},
					},
				},
			},
			expectedService: struct {
				name      string
				namespace string
				ports     []corev1.ServicePort
			}{
				name:      "test-datalayer-datalayer-http",
				namespace: "test-namespace",
				ports: []corev1.ServicePort{
					{
						Port:       80,
						TargetPort: intstr.FromString("http"),
						Protocol:   "TCP",
						Name:       "http",
					},
				},
			},
		},
		{
			name: "With Custom Config",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					FileserverConfig: k8schianetv1.FileserverConfig{
						Service: k8schianetv1.Service{
							ServiceType:           ptr.To(corev1.ServiceTypeNodePort),
							ExternalTrafficPolicy: ptr.To(corev1.ServiceExternalTrafficPolicyLocal),
							SessionAffinity:       ptr.To(corev1.ServiceAffinityClientIP),
							IPFamilyPolicy:        ptr.To(corev1.IPFamilyPolicySingleStack),
							IPFamilies:            ptr.To([]corev1.IPFamily{corev1.IPv4Protocol}),
							AdditionalMetadata: k8schianetv1.AdditionalMetadata{
								Labels:      map[string]string{"test": "label"},
								Annotations: map[string]string{"test": "annotation"},
							},
						},
					},
				},
			},
			expectedService: struct {
				name      string
				namespace string
				ports     []corev1.ServicePort
			}{
				name:      "test-datalayer-datalayer-http",
				namespace: "test-namespace",
				ports: []corev1.ServicePort{
					{
						Port:       80,
						TargetPort: intstr.FromString("http"),
						Protocol:   "TCP",
						Name:       "http",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			service := AssembleService(tc.datalayer)

			// Assert the results
			assert.Equal(t, tc.expectedService.name, service.Name, "Service name should match")
			assert.Equal(t, tc.expectedService.namespace, service.Namespace, "Service namespace should match")
			assert.Equal(t, tc.expectedService.ports, service.Spec.Ports, "Service ports should match")

			if tc.datalayer.Spec.FileserverConfig.Service.ServiceType != nil {
				assert.Equal(t, *tc.datalayer.Spec.FileserverConfig.Service.ServiceType, service.Spec.Type, "Service type should match")
			}
			if tc.datalayer.Spec.FileserverConfig.Service.ExternalTrafficPolicy != nil {
				assert.Equal(t, *tc.datalayer.Spec.FileserverConfig.Service.ExternalTrafficPolicy, service.Spec.ExternalTrafficPolicy, "External traffic policy should match")
			}
			if tc.datalayer.Spec.FileserverConfig.Service.SessionAffinity != nil {
				assert.Equal(t, *tc.datalayer.Spec.FileserverConfig.Service.SessionAffinity, service.Spec.SessionAffinity, "Session affinity should match")
			}
			if tc.datalayer.Spec.FileserverConfig.Service.IPFamilyPolicy != nil {
				assert.Equal(t, *tc.datalayer.Spec.FileserverConfig.Service.IPFamilyPolicy, *service.Spec.IPFamilyPolicy, "IP family policy should match")
			}
			if tc.datalayer.Spec.FileserverConfig.Service.IPFamilies != nil {
				assert.Equal(t, *tc.datalayer.Spec.FileserverConfig.Service.IPFamilies, service.Spec.IPFamilies, "IP families should match")
			}
		})
	}
}

func TestAssembleContainer(t *testing.T) {
	testCases := []struct {
		name              string
		datalayer         k8schianetv1.ChiaDataLayer
		expectedContainer struct {
			name            string
			image           string
			imagePullPolicy corev1.PullPolicy
			ports           []corev1.ContainerPort
		}
	}{
		{
			name: "With Default Config",
			datalayer: k8schianetv1.ChiaDataLayer{
				Spec: k8schianetv1.ChiaDataLayerSpec{
					FileserverConfig: k8schianetv1.FileserverConfig{
						Enabled: boolPtr(true),
						Image:   stringPtr("chia-network/chia:latest"),
					},
				},
			},
			expectedContainer: struct {
				name            string
				image           string
				imagePullPolicy corev1.PullPolicy
				ports           []corev1.ContainerPort
			}{
				name:  "fileserver",
				image: "chia-network/chia:latest",
				ports: []corev1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8575,
						Protocol:      corev1.ProtocolTCP,
					},
				},
			},
		},
		{
			name: "With Custom Config",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					CommonSpec: k8schianetv1.CommonSpec{
						ImagePullPolicy: corev1.PullIfNotPresent,
					},
					FileserverConfig: k8schianetv1.FileserverConfig{
						Image:         ptr.To("custom-image:tag"),
						ContainerPort: ptr.To(8080),
					},
				},
			},
			expectedContainer: struct {
				name            string
				image           string
				imagePullPolicy corev1.PullPolicy
				ports           []corev1.ContainerPort
			}{
				name:            "fileserver",
				image:           "custom-image:tag",
				imagePullPolicy: corev1.PullIfNotPresent,
				ports: []corev1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
						Protocol:      "TCP",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			container := AssembleContainer(tc.datalayer)

			// Assert the results
			assert.Equal(t, tc.expectedContainer.name, container.Name, "Container name should match")
			assert.Equal(t, tc.expectedContainer.image, container.Image, "Container image should match")
			assert.Equal(t, tc.expectedContainer.imagePullPolicy, container.ImagePullPolicy, "Image pull policy should match")
			assert.Equal(t, tc.expectedContainer.ports, container.Ports, "Container ports should match")
		})
	}
}

func TestAssembleIngress(t *testing.T) {
	// Test cases
	testCases := []struct {
		name            string
		datalayer       k8schianetv1.ChiaDataLayer
		expectedIngress struct {
			name      string
			namespace string
		}
	}{
		{
			name: "With Default Config",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					FileserverConfig: k8schianetv1.FileserverConfig{
						Ingress: k8schianetv1.IngressConfig{},
					},
				},
			},
			expectedIngress: struct {
				name      string
				namespace string
			}{
				name:      "test-datalayer-datalayer-http",
				namespace: "test-namespace",
			},
		},
		{
			name: "With Custom Config",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					FileserverConfig: k8schianetv1.FileserverConfig{
						Ingress: k8schianetv1.IngressConfig{
							Enabled:          ptr.To(true),
							IngressClassName: ptr.To("nginx"),
							Host:             ptr.To("test.example.com"),
							AdditionalMetadata: k8schianetv1.AdditionalMetadata{
								Labels:      map[string]string{"test": "label"},
								Annotations: map[string]string{"test": "annotation"},
							},
						},
					},
				},
			},
			expectedIngress: struct {
				name      string
				namespace string
			}{
				name:      "test-datalayer-datalayer-http",
				namespace: "test-namespace",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			ingress := AssembleIngress(tc.datalayer)

			// Assert the results
			assert.Equal(t, tc.expectedIngress.name, ingress.Name, "Ingress name should match")
			assert.Equal(t, tc.expectedIngress.namespace, ingress.Namespace, "Ingress namespace should match")

			if tc.datalayer.Spec.FileserverConfig.Ingress.Enabled != nil && *tc.datalayer.Spec.FileserverConfig.Ingress.Enabled {
				assert.NotNil(t, ingress.Spec.IngressClassName, "Ingress class name should be set")
				assert.NotNil(t, ingress.Spec.Rules, "Ingress rules should be set")
			}
		})
	}
}
