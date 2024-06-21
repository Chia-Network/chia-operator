/*
Copyright 2024 Chia Network Inc.
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

var _ = Describe("ChiaCrawler controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250

		caSecret = "test-secret"
	)

	Context("When creating ChiaCrawler", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaCrawler")
			ctx := context.Background()
			testCrawler := &apiv1.ChiaCrawler{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaCrawler",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiafarmer",
					Namespace: "default",
				},
				Spec: apiv1.ChiaCrawlerSpec{
					ChiaConfig: apiv1.ChiaCrawlerSpecChia{
						CASecretName: &caSecret,
					},
				},
			}
			expect := &apiv1.ChiaCrawler{
				Spec: apiv1.ChiaCrawlerSpec{
					ChiaConfig: apiv1.ChiaCrawlerSpecChia{
						CommonSpecChia: apiv1.CommonSpecChia{
							Image: fmt.Sprintf("ghcr.io/chia-network/chia:%s", defaultChiaImageTag),
						},
						CASecretName: &caSecret,
					},
					CommonSpec: apiv1.CommonSpec{
						ImagePullPolicy: "Always",
						ChiaExporterConfig: apiv1.SpecChiaExporter{
							Enabled: true,
							Image:   fmt.Sprintf("ghcr.io/chia-network/chia-exporter:%s", defaultChiaExporterImageTag),
						},
					},
				},
			}

			// Create ChiaCrawler
			Expect(k8sClient.Create(ctx, testCrawler)).Should(Succeed())

			// Look up the created ChiaCrawler
			lookupKey := types.NamespacedName{Name: testCrawler.Name, Namespace: testCrawler.Namespace}
			createdChiaCrawler := &apiv1.ChiaCrawler{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaCrawler)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaCrawler's spec equals the expected spec
			Expect(createdChiaCrawler.Spec).Should(Equal(expect.Spec))
		})
	})
})
