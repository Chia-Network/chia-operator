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
	)

	Context("When creating ChiaSeeder", func() {
		It("should update its Spec with API defaults", func() {
			By("By creating a new ChiaSeeder")
			ctx := context.Background()
			testFarmer := &apiv1.ChiaSeeder{
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
						CommonSpecChia: apiv1.CommonSpecChia{
							CASecretName: "test-secret",
						},
						BootstrapPeer: "node.default.svc.cluster.local:58444",
						MinimumHeight: uint64(100),
						DomainName:    "seeder.example.com",
						Nameserver:    "example.com",
						Rname:         "admin.example.com",
					},
				},
			}
			expect := &apiv1.ChiaSeeder{
				Spec: apiv1.ChiaSeederSpec{
					ChiaConfig: apiv1.ChiaSeederSpecChia{
						CommonSpecChia: apiv1.CommonSpecChia{
							Image:        fmt.Sprintf("ghcr.io/chia-network/chia:%s", defaultChiaImageTag),
							CASecretName: "test-secret",
						},
						BootstrapPeer: "node.default.svc.cluster.local:58444",
						MinimumHeight: uint64(100),
						DomainName:    "seeder.example.com",
						Nameserver:    "example.com",
						Rname:         "admin.example.com",
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

			// Create ChiaSeeder
			Expect(k8sClient.Create(ctx, testFarmer)).Should(Succeed())

			// Look up the created ChiaSeeder
			lookupKey := types.NamespacedName{Name: testFarmer.Name, Namespace: testFarmer.Namespace}
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
