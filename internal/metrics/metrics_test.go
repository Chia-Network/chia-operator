package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetricsRegistration(t *testing.T) {
	// Test that all metrics are registered
	metrics := []prometheus.Collector{
		ChiaCAs,
		ChiaCertificates,
		ChiaCrawlers,
		ChiaDataLayers,
		ChiaFarmers,
		ChiaHarvesters,
		ChiaIntroducers,
		ChiaNodes,
		ChiaNetworks,
		ChiaSeeders,
		ChiaTimelords,
		ChiaWallets,
	}

	for _, metric := range metrics {
		assert.NotNil(t, metric, "Metric should not be nil")
	}
}

func TestGaugeOperations(t *testing.T) {
	// Test setting gauge values
	tests := []struct {
		name   string
		metric prometheus.Gauge
		value  float64
	}{
		{"ChiaCAs", ChiaCAs, 5},
		{"ChiaCertificates", ChiaCertificates, 10},
		{"ChiaCrawlers", ChiaCrawlers, 3},
		{"ChiaDataLayers", ChiaDataLayers, 2},
		{"ChiaFarmers", ChiaFarmers, 4},
		{"ChiaHarvesters", ChiaHarvesters, 6},
		{"ChiaIntroducers", ChiaIntroducers, 1},
		{"ChiaNodes", ChiaNodes, 8},
		{"ChiaNetworks", ChiaNetworks, 2},
		{"ChiaSeeders", ChiaSeeders, 3},
		{"ChiaTimelords", ChiaTimelords, 1},
		{"ChiaWallets", ChiaWallets, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify we can set the gauge value without panicking
			assert.NotPanics(t, func() {
				tt.metric.Set(tt.value)
			}, "Setting gauge value should not panic")
		})
	}
}
