package engine

import (
	"strings"
)

// formatLabels takes a comma separated list of agent labels
// and formats it into map[string]string for agent label matching
func formatLabels(labels []string) map[string]string {
	labelsMap := make(map[string]string)

	for _, label := range labels {
		keyVal := strings.Split(label, ":")

		labelsMap[keyVal[0]] = keyVal[1]
	}

	return labelsMap
}
