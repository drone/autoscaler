package engine

import (
	"reflect"
	"testing"
)

func TestFormatLabels(t *testing.T) {
	expected := map[string]string{
		"size":   "large",
		"region": "us-west-2",
	}

	result := formatLabels([]string{"size:large", "region:us-west-2"})

	if ok := reflect.DeepEqual(result, expected); !ok {
		t.Errorf("Error formatting labels, got: %+v, expected: %+v", result, expected)
	}
}
