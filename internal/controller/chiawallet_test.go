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

var _ = Describe("ChiaWallet controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating ChiaWallet", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaWallet")
			ctx := context.Background()
			testWallet := &apiv1.ChiaWallet{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaWallet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiawallet",
					Namespace: "default",
				},
				Spec: apiv1.ChiaWalletSpec{
					ChiaConfig: apiv1.ChiaWalletSpecChia{
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
			expect := &apiv1.ChiaWallet{
				Spec: apiv1.ChiaWalletSpec{
					ChiaConfig: apiv1.ChiaWalletSpecChia{
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
							Enabled: nil,
						},
					},
				},
			}

			// Create ChiaWallet
			Expect(k8sClient.Create(ctx, testWallet)).Should(Succeed())

			// Look up the created ChiaWallet
			lookupKey := types.NamespacedName{Name: testWallet.Name, Namespace: testWallet.Namespace}
			createdChiaWallet := &apiv1.ChiaWallet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaWallet)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaWallet's spec equals the expected spec
			Expect(createdChiaWallet.Spec).Should(Equal(expect.Spec))
		})
	})
})
