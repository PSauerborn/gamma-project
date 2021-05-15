package utils

import "github.com/PSauerborn/gamma-project/internal/pkg/utils"

// function to create a new config map
func NewConfigMap() *utils.ConfigMap {
	return &utils.ConfigMap{
		ValueMaps: map[string]string{},
	}
}

// function to create new config map with a
// collection of defaults provided in a
// map instance
func NewConfigMapWithValues(defaults map[string]string) *utils.ConfigMap {
	return &utils.ConfigMap{
		ValueMaps: defaults,
	}
}
