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

var _ = Describe("ChiaFarmer controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating ChiaFarmer", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaFarmer")
			ctx := context.Background()
			testFarmer := &apiv1.ChiaFarmer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaFarmer",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiafarmer",
					Namespace: "default",
				},
				Spec: apiv1.ChiaFarmerSpec{
					ChiaConfig: apiv1.ChiaFarmerSpecChia{
						CASecretName: "test-secret",
						FullNodePeers: &[]apiv1.Peer{
							{
								Host: "node.default.svc.cluster.local",
								Port: 58444,
							},
						},
						SecretKey: apiv1.ChiaSecretKey{
							Name: "testkeys",
							Key:  "key.txt",
						},
					},
				},
			}
			expect := &apiv1.ChiaFarmer{
				Spec: apiv1.ChiaFarmerSpec{
					ChiaConfig: apiv1.ChiaFarmerSpecChia{
						CASecretName: "test-secret",
						FullNodePeers: &[]apiv1.Peer{
							{
								Host: "node.default.svc.cluster.local",
								Port: 58444,
							},
						},
						SecretKey: apiv1.ChiaSecretKey{
							Name: "testkeys",
							Key:  "key.txt",
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

			// Create ChiaFarmer
			Expect(k8sClient.Create(ctx, testFarmer)).Should(Succeed())

			// Look up the created ChiaFarmer
			lookupKey := types.NamespacedName{Name: testFarmer.Name, Namespace: testFarmer.Namespace}
			createdChiaFarmer := &apiv1.ChiaFarmer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaFarmer)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaFarmer's spec equals the expected spec
			Expect(createdChiaFarmer.Spec).Should(Equal(expect.Spec))
		})
	})
})
