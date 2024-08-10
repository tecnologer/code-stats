package models //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatTypeIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statType StatType
		expected bool
	}{
		{StatTypeBytes, true},
		{StatTypeCodeBytes, true},
		{StatTypeLines, true},
		{StatTypeCode, true},
		{StatTypeComment, true},
		{StatTypeBlank, true},
		{StatTypeComplexity, true},
		{StatTypeCountFiles, true},
		{StatTypeWeightedComplexity, true},
		{StatTypeFromString("unknown"), false},
		{StatTypeFromString("another_invalid"), false},
	}

	for _, test := range tests {
		t.Run(test.statType.String(), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.statType.IsValid())
		})
	}
}

func TestStatTypeTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statType StatType
		expected string
	}{
		{StatTypeBytes, "Bytes"},
		{StatTypeCodeBytes, "Code Bytes"},
		{StatTypeLines, "Lines"},
		{StatTypeCode, "Code"},
		{StatTypeComment, "Comment"},
		{StatTypeBlank, "Blank"},
		{StatTypeComplexity, "Complexity"},
		{StatTypeCountFiles, "Count Files"},
		{StatTypeWeightedComplexity, "Weighted Complexity"},
	}

	for _, test := range tests {
		t.Run(test.statType.String(), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.statType.Title())
		})
	}
}

func TestStatTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statType StatType
		expected string
	}{
		{StatTypeBytes, "bytes"},
		{StatTypeCodeBytes, "code_bytes"},
		{StatTypeLines, "lines"},
		{StatTypeCode, "code"},
		{StatTypeComment, "comment"},
		{StatTypeBlank, "blank"},
		{StatTypeComplexity, "complexity"},
		{StatTypeCountFiles, "count_files"},
		{StatTypeWeightedComplexity, "weighted_complexity"},
	}

	for _, test := range tests {
		t.Run(test.statType.String(), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.statType.String())
		})
	}
}
