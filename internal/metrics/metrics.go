package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	// ChiaCAs is a gauge metric that keeps a running total of deployed ChiaCAs
	ChiaCAs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chiaca_total",
			Help: "Number of ChiaCA objects controlled by this operator",
		},
	)

	// ChiaFarmers is a gauge metric that keeps a running total of deployed ChiaFarmers
	ChiaFarmers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chiafarmer_total",
			Help: "Number of ChiaFarmer objects controlled by this operator",
		},
	)

	// ChiaHarvesters is a gauge metric that keeps a running total of deployed ChiaHarvesters
	ChiaHarvesters = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chiaharvester_total",
			Help: "Number of ChiaHarvester objects controlled by this operator",
		},
	)

	// ChiaNodes is a gauge metric that keeps a running total of deployed ChiaNodes
	ChiaNodes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chianode_total",
			Help: "Number of ChiaNode objects controlled by this operator",
		},
	)

	// ChiaTimelords is a gauge metric that keeps a running total of deployed ChiaTimelords
	ChiaTimelords = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chiatimelord_total",
			Help: "Number of ChiaTimelord objects controlled by this operator",
		},
	)

	// ChiaWallets is a gauge metric that keeps a running total of deployed ChiaWallets
	ChiaWallets = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chiawallet_total",
			Help: "Number of ChiaWallet objects controlled by this operator",
		},
	)
)

func init() {
	metrics.Registry.MustRegister(ChiaCAs, ChiaFarmers, ChiaHarvesters, ChiaNodes, ChiaTimelords, ChiaWallets)
}
