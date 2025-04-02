/*
Copyright 2024 Chia Network Inc.
*/

package chianetwork

import (
	"fmt"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func assembleConfigMap(network k8schianetv1.ChiaNetwork) (corev1.ConfigMap, error) {
	data, err := assembleConfigMapData(network)
	if err != nil {
		return corev1.ConfigMap{}, err
	}

	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      network.Name,
			Namespace: network.Namespace,
		},
		Data: data,
	}, nil
}

func assembleConfigMapData(network k8schianetv1.ChiaNetwork) (map[string]string, error) {
	var data = make(map[string]string)

	// network env var
	if network.Spec.NetworkName != nil && *network.Spec.NetworkName != "" {
		data["network"] = *network.Spec.NetworkName
	} else {
		data["network"] = network.Name
	}

	// network_port env var
	if network.Spec.NetworkPort != nil && *network.Spec.NetworkPort != 0 {
		data["network_port"] = strconv.Itoa(int(*network.Spec.NetworkPort))
	}

	// introducer_address env var
	if network.Spec.IntroducerAddress != nil && *network.Spec.IntroducerAddress != "" {
		data["introducer_address"] = *network.Spec.IntroducerAddress
	}

	// dns_introducer_address env var
	if network.Spec.DNSIntroducerAddress != nil && *network.Spec.DNSIntroducerAddress != "" {
		data["dns_introducer_address"] = *network.Spec.DNSIntroducerAddress
	}

	// chia.network_overrides.constants env var
	if network.Spec.NetworkConstants != nil {
		networkConstants, err := marshalNetworkOverride(data["network"], *network.Spec.NetworkConstants)
		if err != nil {
			return nil, fmt.Errorf("error marshaling network constants: %v", err)
		}
		data["chia.network_overrides.constants"] = networkConstants
	}

	// chia.network_overrides.config env var
	if network.Spec.NetworkConfig != nil {
		networkConfig, err := marshalNetworkOverride(data["network"], *network.Spec.NetworkConfig)
		if err != nil {
			return nil, fmt.Errorf("error marshaling network config: %v", err)
		}
		data["chia.network_overrides.config"] = networkConfig
	}

	return data, nil
}
