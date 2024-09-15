/*
Copyright 2024 Chia Network Inc.
*/

package controller

import (
	"context"
	"time"

	"github.com/chia-network/go-chia-libs/pkg/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	apiv1 "github.com/chia-network/chia-operator/api/v1"
)

var _ = Describe("ChiaNetwork controller", func() {
	var (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250

		networkConfig = config.NetworkConfig{
			AddressPrefix:       "txch",
			DefaultFullNodePort: 58444,
		}
		networkConsts = apiv1.NetworkConstants{
			GenesisChallenge:               "fb00c54298fc1c149afbf4c8996fb2317ae41e4649b934ca495991b7852b841",
			GenesisPreFarmPoolPuzzleHash:   "asdlsakldlskalskdsasdasdsadsadsadsadsdsadsas",
			GenesisPreFarmFarmerPuzzleHash: "testestestestestestestesrestestestestestest",
		}
	)

	Context("When creating ChiaNetwork", func() {
		It("should apply to the cluster", func() {
			By("By creating a new ChiaNetwork")
			ctx := context.Background()
			testNetwork := &apiv1.ChiaNetwork{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaNetwork",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chianetwork",
					Namespace: "default",
				},
				Spec: apiv1.ChiaNetworkSpec{
					NetworkConstants: &networkConsts,
					NetworkConfig:    &networkConfig,
				},
			}
			expect := &apiv1.ChiaNetwork{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaNetwork",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-chianetwork",
					Namespace: "default",
				},
				Spec: apiv1.ChiaNetworkSpec{
					NetworkConstants: &networkConsts,
					NetworkConfig:    &networkConfig,
				},
			}

			// Create ChiaNetwork
			Expect(k8sClient.Create(ctx, testNetwork)).Should(Succeed())

			// Look up the created ChiaNetwork
			lookupKey := types.NamespacedName{Name: testNetwork.Name, Namespace: testNetwork.Namespace}
			createdChiaNetwork := &apiv1.ChiaNetwork{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdChiaNetwork)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaNetwork's spec equals the expected spec
			Expect(createdChiaNetwork.Spec).Should(Equal(expect.Spec))
		})
	})
})
