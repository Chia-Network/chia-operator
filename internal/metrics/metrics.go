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

	// ChiaIntroducers is a gauge metric that keeps a running total of deployed ChiaIntroducers
	ChiaIntroducers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chiaintroducer_total",
			Help: "Number of ChiaIntroducer objects controlled by this operator",
		},
	)

	// ChiaNodes is a gauge metric that keeps a running total of deployed ChiaNodes
	ChiaNodes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chianode_total",
			Help: "Number of ChiaNode objects controlled by this operator",
		},
	)

	// ChiaSeeders is a gauge metric that keeps a running total of deployed ChiaSeeders
	ChiaSeeders = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chia_operator_chiaseeder_total",
			Help: "Number of ChiaSeeder objects controlled by this operator",
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

	// OperatorErrors is a counter of the number of errors this exporter has encountered since it started
	OperatorErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "chia_operator_errors_total",
			Help: "Number of errors this exporter has encountered since it started",
		},
	)
)

func init() {
	metrics.Registry.MustRegister(
		ChiaCAs,
		ChiaFarmers,
		ChiaHarvesters,
		ChiaIntroducers,
		ChiaNodes,
		ChiaTimelords,
		ChiaWallets,
		OperatorErrors,
	)
}
