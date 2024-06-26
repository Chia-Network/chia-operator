package kube

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AssembleCommonServiceInputs contains configuration inputs to the AssembleCommonService function
type AssembleCommonServiceInputs struct {
	Name           string
	Namespace      string
	Labels         map[string]string
	Annotations    map[string]string
	OwnerReference []metav1.OwnerReference
	ServiceType    corev1.ServiceType
	Ports          []corev1.ServicePort
	SelectorLabels map[string]string
}

// AssembleCommonService accepts some values and outputs a kubernetes Service definition in a standard way
func AssembleCommonService(input AssembleCommonServiceInputs) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            input.Name,
			Namespace:       input.Namespace,
			Labels:          input.Labels,
			Annotations:     input.Annotations,
			OwnerReferences: input.OwnerReference,
		},
		Spec: corev1.ServiceSpec{
			Type:     input.ServiceType,
			Ports:    input.Ports,
			Selector: input.SelectorLabels,
		},
	}
}
