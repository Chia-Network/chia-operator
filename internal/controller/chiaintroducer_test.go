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

var _ = Describe("ChiaIntroducer controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250

		caSecret = "test-secret"
	)

	Context("When creating ChiaIntroducer", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaIntroducer")
			ctx := context.Background()
			testIntroducer := &apiv1.ChiaIntroducer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaIntroducer",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiafarmer",
					Namespace: "default",
				},
				Spec: apiv1.ChiaIntroducerSpec{
					ChiaConfig: apiv1.ChiaIntroducerSpecChia{
						CASecretName: &caSecret,
					},
				},
			}
			expect := &apiv1.ChiaIntroducer{
				Spec: apiv1.ChiaIntroducerSpec{
					ChiaConfig: apiv1.ChiaIntroducerSpecChia{
						CASecretName: &caSecret,
					},
					CommonSpec: apiv1.CommonSpec{
						ImagePullPolicy: "Always",
						ChiaExporterConfig: apiv1.SpecChiaExporter{
							Enabled: nil,
						},
					},
				},
			}

			// Create ChiaIntroducer
			Expect(k8sClient.Create(ctx, testIntroducer)).Should(Succeed())

			// Look up the created ChiaIntroducer
			lookupKey := types.NamespacedName{Name: testIntroducer.Name, Namespace: testIntroducer.Namespace}
			createdChiaIntroducer := &apiv1.ChiaIntroducer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaIntroducer)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaIntroducer's spec equals the expected spec
			Expect(createdChiaIntroducer.Spec).Should(Equal(expect.Spec))
		})
	})
})
