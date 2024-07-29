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

var _ = Describe("ChiaSeeder controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250

		domainName = "seeder.example.com."
		nameserver = "example.com."
		caSecret   = "test-secret"
	)

	Context("When creating ChiaSeeder", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaSeeder")
			ctx := context.Background()
			testSeeder := &apiv1.ChiaSeeder{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaSeeder",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chiafarmer",
					Namespace: "default",
				},
				Spec: apiv1.ChiaSeederSpec{
					ChiaConfig: apiv1.ChiaSeederSpecChia{
						CASecretName: &caSecret,
						DomainName:   domainName,
						Nameserver:   nameserver,
					},
				},
			}
			expect := &apiv1.ChiaSeeder{
				Spec: apiv1.ChiaSeederSpec{
					ChiaConfig: apiv1.ChiaSeederSpecChia{
						CommonSpecChia: apiv1.CommonSpecChia{
							Image: fmt.Sprintf("ghcr.io/chia-network/chia:%s", defaultChiaImageTag),
						},
						CASecretName: &caSecret,
						DomainName:   domainName,
						Nameserver:   nameserver,
					},
					ChiaHealthcheckConfig: apiv1.SpecChiaHealthcheck{
						Enabled:     false,
						Image:       fmt.Sprintf("ghcr.io/chia-network/chia-healthcheck:%s", defaultChiaHealthcheckImageTag),
						DNSHostname: nil,
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

			// Create ChiaSeeder
			Expect(k8sClient.Create(ctx, testSeeder)).Should(Succeed())

			// Look up the created ChiaSeeder
			lookupKey := types.NamespacedName{Name: testSeeder.Name, Namespace: testSeeder.Namespace}
			createdChiaSeeder := &apiv1.ChiaSeeder{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaSeeder)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaSeeder's spec equals the expected spec
			Expect(createdChiaSeeder.Spec).Should(Equal(expect.Spec))
		})
	})
})
