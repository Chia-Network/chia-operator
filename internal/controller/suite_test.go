/*
Copyright 2023 Chia Network Inc.
*/

package controller

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	apiv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/chiaca"
	"github.com/chia-network/chia-operator/internal/controller/chiafarmer"
	"github.com/chia-network/chia-operator/internal/controller/chiaharvester"
	"github.com/chia-network/chia-operator/internal/controller/chianode"
	"github.com/chia-network/chia-operator/internal/controller/chiawallet"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = apiv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&chiaca.ChiaCAReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&chiafarmer.ChiaFarmerReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&chiaharvester.ChiaHarvesterReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&chianode.ChiaNodeReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&chiawallet.ChiaWalletReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("ChiaCA controller", func() {
	var (
		caSecretName    = "test-secret"
		chiaCAName      = "test-chiaca"
		chiaCANamespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating ChiaCA Status", func() {
		It("Should update ChiaCA Status.Ready to true when deployment is created", func() {
			By("By creating a new ChiaCA")
			ctx := context.Background()
			ca := &apiv1.ChiaCA{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaCA",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      chiaCAName,
					Namespace: chiaCANamespace,
				},
				Spec: apiv1.ChiaCASpec{
					Secret: caSecretName,
				},
			}

			// Create ChiaCA
			Expect(k8sClient.Create(ctx, ca)).Should(Succeed())

			// Look up the created ChiaCA
			chiaCALookupKey := types.NamespacedName{Name: chiaCAName, Namespace: chiaCANamespace}
			createdChiaCA := &apiv1.ChiaCA{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, chiaCALookupKey, createdChiaCA)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaCA's spec.Secret is equal to the given Secret name
			Expect(createdChiaCA.Spec.Secret).Should(Equal(caSecretName))
		})
	})
})

var _ = Describe("ChiaFarmer controller", func() {
	var (
		chiaFarmerName      = "test-chiafarmer"
		chiaFarmerNamespace = "default"

		timeout       = time.Second * 10
		interval      = time.Millisecond * 250
		caSecretName  = "test-secret"
		testnet       = true
		timezone      = "UTC"
		logLevel      = "INFO"
		fullNodePeer  = "node.default.svc.cluster.local:58444"
		secretKeyName = "testkeys"
		secretKeyKey  = "key.txt"
	)

	Context("When updating ChiaFarmer Status", func() {
		It("Should update ChiaFarmer Status.Ready to true when deployment is created", func() {
			By("By creating a new ChiaFarmer")
			ctx := context.Background()
			farmer := &apiv1.ChiaFarmer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaFarmer",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      chiaFarmerName,
					Namespace: chiaFarmerNamespace,
				},
				Spec: apiv1.ChiaFarmerSpec{
					ChiaConfig: apiv1.ChiaFarmerConfigSpec{
						CommonChiaConfigSpec: apiv1.CommonChiaConfigSpec{
							CASecretName: caSecretName,
							Testnet:      &testnet,
							Timezone:     &timezone,
							LogLevel:     &logLevel,
						},
						FullNodePeer: fullNodePeer,
						SecretKeySpec: apiv1.ChiaKeysSpec{
							Name: secretKeyName,
							Key:  secretKeyKey,
						},
					},
					ChiaExporterConfig: apiv1.ChiaExporterConfigSpec{
						ServiceLabels: map[string]string{
							"key": "value",
						},
					},
				},
			}

			// Create ChiaFarmer
			Expect(k8sClient.Create(ctx, farmer)).Should(Succeed())

			// Look up the created ChiaFarmer
			cronjobLookupKey := types.NamespacedName{Name: chiaFarmerName, Namespace: chiaFarmerNamespace}
			createdChiaFarmer := &apiv1.ChiaFarmer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cronjobLookupKey, createdChiaFarmer)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaFarmer's spec.chia.timezone was set to the expected timezone
			Expect(*createdChiaFarmer.Spec.ChiaConfig.Timezone).Should(Equal(timezone))
		})
	})
})

var _ = Describe("ChiaHarvester controller", func() {
	var (
		chiaHarvesterName      = "test-chiaharvester"
		chiaHarvesterNamespace = "default"

		timeout      = time.Second * 10
		interval     = time.Millisecond * 250
		caSecretName = "test-secret"
		testnet      = true
		timezone     = "UTC"
		logLevel     = "INFO"
	)

	Context("When updating ChiaHarvester Status", func() {
		It("Should update ChiaHarvester Status.Ready to true when deployment is created", func() {
			By("By creating a new ChiaHarvester")
			ctx := context.Background()
			harvester := &apiv1.ChiaHarvester{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaHarvester",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      chiaHarvesterName,
					Namespace: chiaHarvesterNamespace,
				},
				Spec: apiv1.ChiaHarvesterSpec{
					ChiaConfig: apiv1.ChiaHarvesterConfigSpec{
						CASecretName: caSecretName,
						Testnet:      &testnet,
						Timezone:     &timezone,
						LogLevel:     &logLevel,
					},
					Storage: &apiv1.StorageConfig{
						Plots: &apiv1.PlotsConfig{
							HostPathVolume: []*apiv1.HostPathVolumeConfig{
								{
									Path: "/home/test/plot1",
								},
								{
									Path: "/home/test/plot2",
								},
							},
						},
					},
					ChiaExporterConfig: apiv1.ChiaExporterConfigSpec{
						ServiceLabels: map[string]string{
							"key": "value",
						},
					},
				},
			}

			// Create ChiaHarvester
			Expect(k8sClient.Create(ctx, harvester)).Should(Succeed())

			// Look up the created ChiaHarvester
			cronjobLookupKey := types.NamespacedName{Name: chiaHarvesterName, Namespace: chiaHarvesterNamespace}
			createdChiaHarvester := &apiv1.ChiaHarvester{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cronjobLookupKey, createdChiaHarvester)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaHarvester's spec.chia.timezone was set to the expected timezone
			Expect(*createdChiaHarvester.Spec.ChiaConfig.Timezone).Should(Equal(timezone))
		})
	})
})

var _ = Describe("ChiaNode controller", func() {
	var (
		chiaNodeName      = "test-chianode"
		chiaNodeNamespace = "default"

		timeout         = time.Second * 10
		interval        = time.Millisecond * 250
		caSecretName    = "test-secret"
		testnet         = true
		timezone        = "UTC"
		logLevel        = "INFO"
		storageClass    = ""
		resourceRequest = "250Gi"
	)

	Context("When updating ChiaNode Status", func() {
		It("Should update ChiaNode Status.Ready to true when deployment is created", func() {
			By("By creating a new ChiaNode")
			ctx := context.Background()
			node := &apiv1.ChiaNode{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaNode",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      chiaNodeName,
					Namespace: chiaNodeNamespace,
				},
				Spec: apiv1.ChiaNodeSpec{
					ChiaConfig: apiv1.ChiaNodeConfigSpec{
						CommonChiaConfigSpec: apiv1.CommonChiaConfigSpec{
							CASecretName: caSecretName,
							Testnet:      &testnet,
							Timezone:     &timezone,
							LogLevel:     &logLevel,
						},
					},
					Storage: &apiv1.StorageConfig{
						ChiaRoot: &apiv1.ChiaRootConfig{
							PersistentVolumeClaim: &apiv1.PersistentVolumeClaimConfig{
								StorageClass:    storageClass,
								ResourceRequest: resourceRequest,
							},
						},
					},
					ChiaExporterConfig: apiv1.ChiaExporterConfigSpec{
						Enabled: true,
						ServiceLabels: map[string]string{
							"key": "value",
						},
					},
				},
			}

			// Create ChiaNode
			Expect(k8sClient.Create(ctx, node)).Should(Succeed())

			// Look up the created ChiaNode
			cronjobLookupKey := types.NamespacedName{Name: chiaNodeName, Namespace: chiaNodeNamespace}
			createdChiaNode := &apiv1.ChiaNode{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cronjobLookupKey, createdChiaNode)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaNode's spec.chia.timezone was set to the expected timezone
			Expect(*createdChiaNode.Spec.ChiaConfig.Timezone).Should(Equal(timezone))
			Expect(createdChiaNode.Spec.ChiaConfig.Image).Should(Equal("ghcr.io/chia-network/chia:latest"))
			Expect(createdChiaNode.Spec.ChiaExporterConfig.Image).Should(Equal(consts.DefaultChiaExporterImage))
		})
	})
})

var _ = Describe("ChiaWallet controller", func() {
	var (
		chiaWalletName      = "test-chiawallet"
		chiaWalletNamespace = "default"

		timeout       = time.Second * 10
		interval      = time.Millisecond * 250
		caSecretName  = "test-secret"
		testnet       = true
		timezone      = "UTC"
		logLevel      = "INFO"
		fullNodePeer  = "node.default.svc.cluster.local:58444"
		secretKeyName = "testkeys"
		secretKeyKey  = "key.txt"
	)

	Context("When updating ChiaWallet Status", func() {
		It("Should update ChiaWallet Status.Ready to true when deployment is created", func() {
			By("By creating a new ChiaWallet")
			ctx := context.Background()
			wallet := &apiv1.ChiaWallet{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaWallet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      chiaWalletName,
					Namespace: chiaWalletNamespace,
				},
				Spec: apiv1.ChiaWalletSpec{
					ChiaConfig: apiv1.ChiaWalletConfigSpec{
						CASecretName: caSecretName,
						Testnet:      &testnet,
						Timezone:     &timezone,
						LogLevel:     &logLevel,
						FullNodePeer: fullNodePeer,
						SecretKeySpec: apiv1.ChiaKeysSpec{
							Name: secretKeyName,
							Key:  secretKeyKey,
						},
					},
					ChiaExporterConfig: apiv1.ChiaExporterConfigSpec{
						ServiceLabels: map[string]string{
							"key": "value",
						},
					},
				},
			}

			// Create ChiaWallet
			Expect(k8sClient.Create(ctx, wallet)).Should(Succeed())

			// Look up the created ChiaWallet
			cronjobLookupKey := types.NamespacedName{Name: chiaWalletName, Namespace: chiaWalletNamespace}
			createdChiaWallet := &apiv1.ChiaWallet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cronjobLookupKey, createdChiaWallet)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Ensure the ChiaWallet's spec.chia.fullNodePeer was set to the expected fullNodePeer
			Expect(createdChiaWallet.Spec.ChiaConfig.FullNodePeer).Should(Equal(fullNodePeer))
		})
	})
})
