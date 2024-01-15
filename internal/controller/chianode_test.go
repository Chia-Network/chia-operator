/*
Copyright 2023 Chia Network Inc.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	apiv1 "github.com/chia-network/chia-operator/api/v1"
)

var _ = Describe("ChiaNode controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating ChiaNode", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaNode")
			ctx := context.Background()
			testNode := &apiv1.ChiaNode{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaNode",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chianode",
					Namespace: "default",
				},
				Spec: apiv1.ChiaNodeSpec{
					ChiaConfig: apiv1.ChiaNodeSpecChia{
						CommonSpecChia: apiv1.CommonSpecChia{
							CASecretName: "test-secret",
						},
					},
				},
			}
			expect := &apiv1.ChiaNode{
				Spec: apiv1.ChiaNodeSpec{
					Replicas: 1,
					ChiaConfig: apiv1.ChiaNodeSpecChia{
						CommonSpecChia: apiv1.CommonSpecChia{
							Image:        fmt.Sprintf("ghcr.io/chia-network/chia:%s", defaultChiaImageTag),
							CASecretName: "test-secret",
						},
					},
					CommonSpec: apiv1.CommonSpec{
						ServiceType:     "ClusterIP",
						ImagePullPolicy: "Always",
						ChiaExporterConfig: apiv1.SpecChiaExporter{
							Enabled: true,
							Image:   fmt.Sprintf("ghcr.io/chia-network/chia-exporter:%s", defaultChiaExporterImageTag),
						},
					},
				},
			}

			// Create ChiaNode
			Expect(k8sClient.Create(ctx, testNode)).Should(Succeed())

			// Look up the created ChiaNode
			lookupKey := types.NamespacedName{Name: testNode.Name, Namespace: testNode.Namespace}
			createdChiaNode := &apiv1.ChiaNode{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaNode)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaNode's spec equals the expected spec
			Expect(createdChiaNode.Spec).Should(Equal(expect.Spec))
		})
	})
})
