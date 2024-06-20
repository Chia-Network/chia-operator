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

var _ = Describe("ChiaHarvester controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating ChiaHarvester", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaHarvester")
			ctx := context.Background()
			testHarvester := &apiv1.ChiaHarvester{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaHarvester",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiaharvester",
					Namespace: "default",
				},
				Spec: apiv1.ChiaHarvesterSpec{
					ChiaConfig: apiv1.ChiaHarvesterSpecChia{
						CASecretName:  "test-secret",
						FarmerAddress: "farmer.default.svc.cluster.local",
					},
				},
			}
			expect := &apiv1.ChiaHarvester{
				Spec: apiv1.ChiaHarvesterSpec{
					ChiaConfig: apiv1.ChiaHarvesterSpecChia{
						CommonSpecChia: apiv1.CommonSpecChia{
							Image: fmt.Sprintf("ghcr.io/chia-network/chia:%s", defaultChiaImageTag),
						},
						CASecretName:  "test-secret",
						FarmerAddress: "farmer.default.svc.cluster.local",
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

			// Create ChiaHarvester
			Expect(k8sClient.Create(ctx, testHarvester)).Should(Succeed())

			// Look up the created ChiaHarvester
			lookupKey := types.NamespacedName{Name: testHarvester.Name, Namespace: testHarvester.Namespace}
			createdChiaHarvester := &apiv1.ChiaHarvester{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaHarvester)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaHarvester's spec equals the expected spec
			Expect(createdChiaHarvester.Spec).Should(Equal(expect.Spec))
		})
	})
})
