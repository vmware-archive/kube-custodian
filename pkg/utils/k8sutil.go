// Package utils implements common utils used by other pkg/'s
package utils

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

// LabelsSubSet checks that subSet label keys are all present in Labels
func LabelsSubSet(Labels map[string]string, subSet []string) bool {
	labelKeys := []string{}
	for k := range Labels {
		labelKeys = append(labelKeys, k)
	}
	labelSet := sets.NewString(labelKeys...)

	if labelSet.HasAll(subSet...) {
		return true
	}
	return false
}
