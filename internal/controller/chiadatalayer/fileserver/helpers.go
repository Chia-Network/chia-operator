package fileserver

import (
	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// ShouldAddFileserverContainer determines if a fileserver container should be added based on its configuration settings.
func ShouldAddFileserverContainer(fileserver k8schianetv1.FileserverConfig) bool {
	return fileserver.Enabled != nil && *fileserver.Enabled
}

// ShouldMakeFileserverService determines whether the fileserver service should be created based on configuration flags.
// It returns true if both FileserverConfig.Enabled and FileserverConfig.Service.Enabled are set to true.
func ShouldMakeFileserverService(fileserver k8schianetv1.FileserverConfig) bool {
	return fileserver.Enabled != nil && *fileserver.Enabled && fileserver.Service.Enabled != nil && *fileserver.Service.Enabled
}

func getChiaContainerEnv(mountPath string) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "service",
			Value: "data_layer_http",
		},
		{
			Name:  "keys",
			Value: "none",
		},
		{
			Name:  "chia.data_layer.server_files_location",
			Value: mountPath,
		},
		{
			Name:  "chia.daemon_port",
			Value: "55401", // Avoids port conflict with the main chia container
		},
	}
}
