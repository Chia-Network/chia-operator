/*
Copyright 2023 Chia Network Inc.
*/

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	apiv1 "github.com/chia-network/chia-operator/api/v1"
)

var _ = Describe("ChiaTimelord controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating ChiaTimelord", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaTimelord")
			ctx := context.Background()
			testTimelord := &apiv1.ChiaTimelord{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaTimelord",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiatimelord",
					Namespace: "default",
				},
				Spec: apiv1.ChiaTimelordSpec{
					ChiaConfig: apiv1.ChiaTimelordSpecChia{
						CASecretName: "test-secret",
						FullNodePeers: &[]apiv1.Peer{
							{
								Host: "node.default.svc.cluster.local",
								Port: 58444,
							},
						},
					},
				},
			}
			expect := &apiv1.ChiaTimelord{
				Spec: apiv1.ChiaTimelordSpec{
					ChiaConfig: apiv1.ChiaTimelordSpecChia{
						CASecretName: "test-secret",
						FullNodePeers: &[]apiv1.Peer{
							{
								Host: "node.default.svc.cluster.local",
								Port: 58444,
							},
						},
					},
					CommonSpec: apiv1.CommonSpec{
						ImagePullPolicy: "Always",
						ChiaExporterConfig: apiv1.SpecChiaExporter{
							Enabled: true,
						},
					},
				},
			}

			// Create ChiaTimelord
			Expect(k8sClient.Create(ctx, testTimelord)).Should(Succeed())

			// Look up the created ChiaTimelord
			lookupKey := types.NamespacedName{Name: testTimelord.Name, Namespace: testTimelord.Namespace}
			createdChiaTimelord := &apiv1.ChiaTimelord{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaTimelord)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaTimelord's spec equals the expected spec
			Expect(createdChiaTimelord.Spec).Should(Equal(expect.Spec))
		})
	})
})
