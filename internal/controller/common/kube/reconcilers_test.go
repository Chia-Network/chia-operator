/*
Copyright 2025 Chia Network Inc.
*/

package kube

import (
	"context"
	"testing"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	testEnv   *envtest.Environment
	k8sClient client.Client
)

func TestReconcilers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reconcilers Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

// ---------------------------------------------------------------------------
// Pure unit tests (no envtest required)
// ---------------------------------------------------------------------------

var _ = Describe("filterStaleContainers", func() {
	It("should remove containers not in the desired list", func() {
		current := []corev1.Container{
			{Name: "a", Image: "img"},
			{Name: "b", Image: "img"},
			{Name: "c", Image: "img"},
		}
		desired := []corev1.Container{
			{Name: "a", Image: "img"},
			{Name: "c", Image: "img"},
		}
		filtered, stale := filterStaleContainers(current, desired)
		Expect(stale).To(BeTrue())
		Expect(filtered).To(HaveLen(2))
		Expect(filtered[0].Name).To(Equal("a"))
		Expect(filtered[1].Name).To(Equal("c"))
	})

	It("should report no staleness when lists match", func() {
		current := []corev1.Container{
			{Name: "a", Image: "img"},
			{Name: "b", Image: "img"},
		}
		desired := []corev1.Container{
			{Name: "a", Image: "img"},
			{Name: "b", Image: "img"},
		}
		filtered, stale := filterStaleContainers(current, desired)
		Expect(stale).To(BeFalse())
		Expect(filtered).To(HaveLen(2))
	})

	It("should return empty and stale when desired is empty", func() {
		current := []corev1.Container{
			{Name: "a", Image: "img"},
			{Name: "b", Image: "img"},
		}
		filtered, stale := filterStaleContainers(current, []corev1.Container{})
		Expect(stale).To(BeTrue())
		Expect(filtered).To(BeEmpty())
	})

	It("should return empty and not stale when both are empty", func() {
		filtered, stale := filterStaleContainers([]corev1.Container{}, []corev1.Container{})
		Expect(stale).To(BeFalse())
		Expect(filtered).To(BeEmpty())
	})
})

var _ = Describe("filterStaleContainerFields", func() {
	It("should clear all optional fields present in current but absent in desired", func() {
		current := []corev1.Container{
			{
				Name:            "main",
				Image:           "img",
				LivenessProbe:   &corev1.Probe{},
				ReadinessProbe:  &corev1.Probe{},
				StartupProbe:    &corev1.Probe{},
				SecurityContext: &corev1.SecurityContext{},
				Resources: corev1.ResourceRequirements{
					Limits:   corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("100m")},
					Requests: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("50m")},
				},
			},
		}
		desired := []corev1.Container{
			{Name: "main", Image: "img"},
		}

		changed := filterStaleContainerFields(current, desired)
		Expect(changed).To(BeTrue())
		Expect(current[0].LivenessProbe).To(BeNil())
		Expect(current[0].ReadinessProbe).To(BeNil())
		Expect(current[0].StartupProbe).To(BeNil())
		Expect(current[0].SecurityContext).To(BeNil())
		Expect(current[0].Resources.Limits).To(BeNil())
		Expect(current[0].Resources.Requests).To(BeNil())
	})

	It("should return false when no differences exist", func() {
		probe := &corev1.Probe{}
		current := []corev1.Container{
			{Name: "main", Image: "img", LivenessProbe: probe},
		}
		desired := []corev1.Container{
			{Name: "main", Image: "img", LivenessProbe: probe},
		}
		changed := filterStaleContainerFields(current, desired)
		Expect(changed).To(BeFalse())
	})

	It("should skip containers not in the desired map", func() {
		current := []corev1.Container{
			{Name: "unknown", Image: "img", LivenessProbe: &corev1.Probe{}},
		}
		desired := []corev1.Container{
			{Name: "main", Image: "img"},
		}
		changed := filterStaleContainerFields(current, desired)
		Expect(changed).To(BeFalse())
		Expect(current[0].LivenessProbe).NotTo(BeNil())
	})
})

// ---------------------------------------------------------------------------
// serverSideApply (existing tests)
// ---------------------------------------------------------------------------

var _ = Describe("serverSideApply", func() {
	Context("When applying a ConfigMap", func() {
		It("should create the ConfigMap if it doesn't exist", func() {
			ctx := context.Background()
			desired := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"test-key": "test-value",
				},
			}

			err := serverSideApply(ctx, k8sClient, desired, "ConfigMap", "v1")
			Expect(err).NotTo(HaveOccurred())

			var created corev1.ConfigMap
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      desired.Name,
				Namespace: desired.Namespace,
			}, &created)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.Data).To(Equal(desired.Data))
		})

		It("should update the ConfigMap if it already exists", func() {
			ctx := context.Background()
			existing := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap-update",
					Namespace: "default",
				},
				Data: map[string]string{
					"existing-key": "existing-value",
				},
			}

			err := serverSideApply(ctx, k8sClient, existing, "ConfigMap", "v1")
			Expect(err).NotTo(HaveOccurred())

			desired := existing.DeepCopy()
			desired.Data["new-key"] = "new-value"

			err = serverSideApply(ctx, k8sClient, desired, "ConfigMap", "v1")
			Expect(err).NotTo(HaveOccurred())

			var updated corev1.ConfigMap
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      desired.Name,
				Namespace: desired.Namespace,
			}, &updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Data).To(Equal(desired.Data))
		})
	})
})

// ---------------------------------------------------------------------------
// ReconcileConfigMap
// ---------------------------------------------------------------------------

var _ = Describe("ReconcileConfigMap", func() {
	It("should create a ConfigMap", func() {
		ctx := context.Background()
		desired := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-rcm-create",
				Namespace: "default",
			},
			Data: map[string]string{"key1": "val1"},
		}

		result, err := ReconcileConfigMap(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.ConfigMap
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-rcm-create", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Data).To(Equal(desired.Data))
	})

	It("should update an existing ConfigMap", func() {
		ctx := context.Background()
		desired := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-rcm-update",
				Namespace: "default",
			},
			Data: map[string]string{"key1": "val1"},
		}

		_, err := ReconcileConfigMap(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())

		desired.Data = map[string]string{"key1": "val1", "key2": "val2"}
		result, err := ReconcileConfigMap(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.ConfigMap
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-rcm-update", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Data).To(Equal(map[string]string{"key1": "val1", "key2": "val2"}))
	})
})

// ---------------------------------------------------------------------------
// ReconcileService
// ---------------------------------------------------------------------------

var _ = Describe("ReconcileService", func() {
	It("should create a Service when enabled", func() {
		ctx := context.Background()
		svcConfig := k8schianetv1.Service{Enabled: boolPtr(true)}
		desired := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-svc-create",
				Namespace: "default",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{Name: "peer", Port: 8444, Protocol: corev1.ProtocolTCP},
				},
				Selector: map[string]string{"app": "chia"},
			},
		}

		result, err := ReconcileService(ctx, k8sClient, svcConfig, desired, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.Service
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-svc-create", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Ports).To(HaveLen(1))
		Expect(fetched.Spec.Ports[0].Port).To(Equal(int32(8444)))
		Expect(fetched.Spec.Selector).To(Equal(map[string]string{"app": "chia"}))
	})

	It("should update an existing Service when enabled", func() {
		ctx := context.Background()
		svcConfig := k8schianetv1.Service{Enabled: boolPtr(true)}
		desired := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-svc-update",
				Namespace: "default",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{Name: "peer", Port: 8444, Protocol: corev1.ProtocolTCP},
				},
				Selector: map[string]string{"app": "chia"},
			},
		}

		_, err := ReconcileService(ctx, k8sClient, svcConfig, desired, false)
		Expect(err).NotTo(HaveOccurred())

		desired.Spec.Ports[0].Port = 8445
		result, err := ReconcileService(ctx, k8sClient, svcConfig, desired, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.Service
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-svc-update", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Ports[0].Port).To(Equal(int32(8445)))
	})

	It("should delete an existing Service when disabled", func() {
		ctx := context.Background()
		desired := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-svc-delete",
				Namespace: "default",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{Name: "peer", Port: 8444, Protocol: corev1.ProtocolTCP},
				},
				Selector: map[string]string{"app": "chia"},
			},
		}

		err := serverSideApply(ctx, k8sClient, &desired, "Service", "v1")
		Expect(err).NotTo(HaveOccurred())

		svcConfig := k8schianetv1.Service{Enabled: boolPtr(false)}
		result, err := ReconcileService(ctx, k8sClient, svcConfig, desired, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.Service
		err = k8sClient.Get(ctx, client.ObjectKey{Name: "test-svc-delete", Namespace: "default"}, &fetched)
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})

	It("should be a noop when disabled and no Service exists", func() {
		ctx := context.Background()
		svcConfig := k8schianetv1.Service{Enabled: boolPtr(false)}
		desired := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-svc-noop",
				Namespace: "default",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{Name: "peer", Port: 8444, Protocol: corev1.ProtocolTCP},
				},
				Selector: map[string]string{"app": "chia"},
			},
		}

		result, err := ReconcileService(ctx, k8sClient, svcConfig, desired, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.Service
		err = k8sClient.Get(ctx, client.ObjectKey{Name: "test-svc-noop", Namespace: "default"}, &fetched)
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})
})

// ---------------------------------------------------------------------------
// ReconcilePersistentVolumeClaim
// ---------------------------------------------------------------------------

var _ = Describe("ReconcilePersistentVolumeClaim", func() {
	It("should create a PVC when storage is configured", func() {
		ctx := context.Background()
		storage := &k8schianetv1.StorageConfig{
			ChiaRoot: &k8schianetv1.ChiaRootConfig{
				PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
					GenerateVolumeClaims: true,
				},
			},
		}
		desired := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pvc-create",
				Namespace: "default",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		}

		result, err := ReconcilePersistentVolumeClaim(ctx, k8sClient, storage, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.PersistentVolumeClaim
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-pvc-create", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.AccessModes).To(Equal([]corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}))
	})

	It("should be a noop when storage is not configured", func() {
		ctx := context.Background()
		desired := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pvc-noop",
				Namespace: "default",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		}

		result, err := ReconcilePersistentVolumeClaim(ctx, k8sClient, nil, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.PersistentVolumeClaim
		err = k8sClient.Get(ctx, client.ObjectKey{Name: "test-pvc-noop", Namespace: "default"}, &fetched)
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})

	It("should update PVC metadata when reconciled again", func() {
		ctx := context.Background()
		storage := &k8schianetv1.StorageConfig{
			ChiaRoot: &k8schianetv1.ChiaRootConfig{
				PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
					GenerateVolumeClaims: true,
				},
			},
		}
		desired := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pvc-update",
				Namespace: "default",
				Labels:    map[string]string{"version": "v1"},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		}

		_, err := ReconcilePersistentVolumeClaim(ctx, k8sClient, storage, desired)
		Expect(err).NotTo(HaveOccurred())

		desired.Labels["version"] = "v2"
		result, err := ReconcilePersistentVolumeClaim(ctx, k8sClient, storage, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched corev1.PersistentVolumeClaim
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-pvc-update", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Labels["version"]).To(Equal("v2"))
	})
})

// ---------------------------------------------------------------------------
// ReconcileIngress
// ---------------------------------------------------------------------------

var _ = Describe("ReconcileIngress", func() {
	pathTypePrefix := networkingv1.PathTypePrefix

	makeIngress := func(name string) networkingv1.Ingress {
		return networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     "/",
										PathType: &pathTypePrefix,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "backend-svc",
												Port: networkingv1.ServiceBackendPort{Number: 8080},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	It("should create an Ingress when enabled", func() {
		ctx := context.Background()
		ingressConfig := k8schianetv1.IngressConfig{Enabled: boolPtr(true)}
		desired := makeIngress("test-ingress-create")

		result, err := ReconcileIngress(ctx, k8sClient, ingressConfig, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched networkingv1.Ingress
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-ingress-create", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Rules).To(HaveLen(1))
		Expect(fetched.Spec.Rules[0].Host).To(Equal("example.com"))
	})

	It("should update an existing Ingress when enabled", func() {
		ctx := context.Background()
		ingressConfig := k8schianetv1.IngressConfig{Enabled: boolPtr(true)}
		desired := makeIngress("test-ingress-update")

		_, err := ReconcileIngress(ctx, k8sClient, ingressConfig, desired)
		Expect(err).NotTo(HaveOccurred())

		desired.Spec.Rules[0].Host = "updated.example.com"
		result, err := ReconcileIngress(ctx, k8sClient, ingressConfig, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched networkingv1.Ingress
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-ingress-update", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Rules[0].Host).To(Equal("updated.example.com"))
	})

	It("should delete an existing Ingress when disabled", func() {
		ctx := context.Background()
		desired := makeIngress("test-ingress-delete")

		err := serverSideApply(ctx, k8sClient, &desired, "Ingress", "networking.k8s.io/v1")
		Expect(err).NotTo(HaveOccurred())

		ingressConfig := k8schianetv1.IngressConfig{Enabled: boolPtr(false)}
		result, err := ReconcileIngress(ctx, k8sClient, ingressConfig, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched networkingv1.Ingress
		err = k8sClient.Get(ctx, client.ObjectKey{Name: "test-ingress-delete", Namespace: "default"}, &fetched)
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})

	It("should be a noop when disabled and no Ingress exists", func() {
		ctx := context.Background()
		ingressConfig := k8schianetv1.IngressConfig{Enabled: boolPtr(false)}
		desired := makeIngress("test-ingress-noop")

		result, err := ReconcileIngress(ctx, k8sClient, ingressConfig, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched networkingv1.Ingress
		err = k8sClient.Get(ctx, client.ObjectKey{Name: "test-ingress-noop", Namespace: "default"}, &fetched)
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})
})

// ---------------------------------------------------------------------------
// ReconcileDeployment
// ---------------------------------------------------------------------------

var _ = Describe("ReconcileDeployment", func() {
	makeDeployment := func(name string, labels map[string]string, containers []corev1.Container) appsv1.Deployment {
		return appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: containers,
					},
				},
			},
		}
	}

	It("should create a new Deployment", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-deploy-create"}
		desired := makeDeployment("test-deploy-create", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileDeployment(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.Deployment
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-deploy-create", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(fetched.Spec.Template.Spec.Containers[0].Image).To(Equal("ghcr.io/chia-network/chia:latest"))
	})

	It("should update an existing Deployment with same selectors", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-deploy-update"}
		desired := makeDeployment("test-deploy-update", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:v1"},
		})

		_, err := ReconcileDeployment(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())

		desired.Spec.Template.Spec.Containers[0].Image = "ghcr.io/chia-network/chia:v2"
		result, err := ReconcileDeployment(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.Deployment
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-deploy-update", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers[0].Image).To(Equal("ghcr.io/chia-network/chia:v2"))
	})

	It("should recreate a Deployment when selector labels change", func() {
		ctx := context.Background()
		oldLabels := map[string]string{"app": "chia-deploy-old"}
		desired := makeDeployment("test-deploy-recreate", oldLabels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		_, err := ReconcileDeployment(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())

		newLabels := map[string]string{"app": "chia-deploy-new"}
		desired = makeDeployment("test-deploy-recreate", newLabels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileDeployment(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.Deployment
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-deploy-recreate", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Selector.MatchLabels).To(Equal(newLabels))
		Expect(fetched.Spec.Template.Labels).To(Equal(newLabels))
	})

	It("should remove stale containers from a Deployment", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-deploy-stale-ctr"}
		initial := makeDeployment("test-deploy-stale-ctr", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
			{Name: "sidecar", Image: "sidecar:latest"},
		})

		err := serverSideApply(ctx, k8sClient, &initial, "Deployment", "apps/v1")
		Expect(err).NotTo(HaveOccurred())

		desired := makeDeployment("test-deploy-stale-ctr", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileDeployment(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.Deployment
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-deploy-stale-ctr", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(fetched.Spec.Template.Spec.Containers[0].Name).To(Equal("chia"))
	})

	It("should remove stale fields from Deployment containers", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-deploy-stale-fields"}
		initial := makeDeployment("test-deploy-stale-fields", labels, []corev1.Container{
			{
				Name:  "chia",
				Image: "ghcr.io/chia-network/chia:latest",
				LivenessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						Exec: &corev1.ExecAction{Command: []string{"true"}},
					},
				},
			},
		})

		err := serverSideApply(ctx, k8sClient, &initial, "Deployment", "apps/v1")
		Expect(err).NotTo(HaveOccurred())

		desired := makeDeployment("test-deploy-stale-fields", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileDeployment(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.Deployment
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-deploy-stale-fields", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
	})
})

// ---------------------------------------------------------------------------
// ReconcileStatefulset
// ---------------------------------------------------------------------------

var _ = Describe("ReconcileStatefulset", func() {
	makeStatefulSet := func(name string, labels map[string]string, containers []corev1.Container) appsv1.StatefulSet {
		return appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: appsv1.StatefulSetSpec{
				ServiceName: "test-headless",
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: containers,
					},
				},
			},
		}
	}

	It("should create a new StatefulSet", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-sts-create"}
		desired := makeStatefulSet("test-sts-create", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileStatefulset(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.StatefulSet
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-sts-create", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(fetched.Spec.Template.Spec.Containers[0].Image).To(Equal("ghcr.io/chia-network/chia:latest"))
	})

	It("should update an existing StatefulSet with same selectors", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-sts-update"}
		desired := makeStatefulSet("test-sts-update", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:v1"},
		})

		_, err := ReconcileStatefulset(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())

		desired.Spec.Template.Spec.Containers[0].Image = "ghcr.io/chia-network/chia:v2"
		result, err := ReconcileStatefulset(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.StatefulSet
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-sts-update", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers[0].Image).To(Equal("ghcr.io/chia-network/chia:v2"))
	})

	It("should recreate a StatefulSet when selector labels change", func() {
		ctx := context.Background()
		oldLabels := map[string]string{"app": "chia-sts-old"}
		desired := makeStatefulSet("test-sts-recreate", oldLabels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		_, err := ReconcileStatefulset(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())

		newLabels := map[string]string{"app": "chia-sts-new"}
		desired = makeStatefulSet("test-sts-recreate", newLabels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileStatefulset(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.StatefulSet
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-sts-recreate", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Selector.MatchLabels).To(Equal(newLabels))
		Expect(fetched.Spec.Template.Labels).To(Equal(newLabels))
	})

	It("should remove stale containers from a StatefulSet", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-sts-stale-ctr"}
		initial := makeStatefulSet("test-sts-stale-ctr", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
			{Name: "sidecar", Image: "sidecar:latest"},
		})

		err := serverSideApply(ctx, k8sClient, &initial, "StatefulSet", "apps/v1")
		Expect(err).NotTo(HaveOccurred())

		desired := makeStatefulSet("test-sts-stale-ctr", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileStatefulset(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.StatefulSet
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-sts-stale-ctr", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(fetched.Spec.Template.Spec.Containers[0].Name).To(Equal("chia"))
	})

	It("should remove stale fields from StatefulSet containers", func() {
		ctx := context.Background()
		labels := map[string]string{"app": "chia-sts-stale-fields"}
		initial := makeStatefulSet("test-sts-stale-fields", labels, []corev1.Container{
			{
				Name:  "chia",
				Image: "ghcr.io/chia-network/chia:latest",
				LivenessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						Exec: &corev1.ExecAction{Command: []string{"true"}},
					},
				},
			},
		})

		err := serverSideApply(ctx, k8sClient, &initial, "StatefulSet", "apps/v1")
		Expect(err).NotTo(HaveOccurred())

		desired := makeStatefulSet("test-sts-stale-fields", labels, []corev1.Container{
			{Name: "chia", Image: "ghcr.io/chia-network/chia:latest"},
		})

		result, err := ReconcileStatefulset(ctx, k8sClient, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var fetched appsv1.StatefulSet
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "test-sts-stale-fields", Namespace: "default"}, &fetched)).To(Succeed())
		Expect(fetched.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
	})
})
