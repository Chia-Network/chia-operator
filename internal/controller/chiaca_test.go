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

var _ = Describe("ChiaCA controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating ChiaCA", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaCA")
			ctx := context.Background()
			testCA := &apiv1.ChiaCA{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaCA",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiaca",
					Namespace: "default",
				},
				Spec: apiv1.ChiaCASpec{
					Secret: "test-secret",
				},
			}
			expect := &apiv1.ChiaCA{
				Spec: apiv1.ChiaCASpec{
					Image:  fmt.Sprintf("ghcr.io/chia-network/chia-operator/ca-gen:%s", defaultChiaCAImageTag),
					Secret: "test-secret",
				},
			}

			// Create ChiaCA
			Expect(k8sClient.Create(ctx, testCA)).Should(Succeed())

			// Look up the created ChiaCA
			lookupKey := types.NamespacedName{Name: testCA.Name, Namespace: testCA.Namespace}
			createdChiaCA := &apiv1.ChiaCA{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaCA)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaCA's spec is equal to the expected spec
			Expect(createdChiaCA.Spec).Should(Equal(expect.Spec))
		})
	})
})
