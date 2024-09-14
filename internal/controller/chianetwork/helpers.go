/*
Copyright 2024 Chia Network Inc.
*/

package chianetwork

import "encoding/json"

func marshalNetworkOverride(name string, data interface{}) (string, error) {
	wrappedData := map[string]interface{}{
		name: data,
	}

	jsonData, err := json.Marshal(wrappedData)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
