/*
Copyright 2025 Chia Network Inc.
*/

package kube

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
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

			// Verify the ConfigMap was created
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

			// Create the initial ConfigMap
			err := serverSideApply(ctx, k8sClient, existing, "ConfigMap", "v1")
			Expect(err).NotTo(HaveOccurred())

			// Update the ConfigMap
			desired := existing.DeepCopy()
			desired.Data["new-key"] = "new-value"

			err = serverSideApply(ctx, k8sClient, desired, "ConfigMap", "v1")
			Expect(err).NotTo(HaveOccurred())

			// Verify the ConfigMap was updated
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
