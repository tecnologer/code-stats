package models

import "strings"

type StatType string

const (
	StatTypeBytes              StatType = "bytes"
	StatTypeCodeBytes          StatType = "code_bytes"
	StatTypeLines              StatType = "lines"
	StatTypeCode               StatType = "code"
	StatTypeComment            StatType = "comment"
	StatTypeBlank              StatType = "blank"
	StatTypeComplexity         StatType = "complexity"
	StatTypeCountFiles         StatType = "count_files"
	StatTypeWeightedComplexity StatType = "weighted_complexity"
)

func StatTypeFromString(s string) StatType {
	return StatType(s)
}

func (s StatType) String() string {
	return string(s)
}

func (s StatType) IsValid() bool {
	switch s {
	case StatTypeBytes,
		StatTypeCodeBytes,
		StatTypeLines,
		StatTypeCode,
		StatTypeComment,
		StatTypeBlank,
		StatTypeComplexity,
		StatTypeCountFiles,
		StatTypeWeightedComplexity:
		return true
	}

	return false
}

func AllStatsTypes() []StatType {
	return []StatType{
		StatTypeBytes,
		StatTypeCodeBytes,
		StatTypeLines,
		StatTypeCode,
		StatTypeComment,
		StatTypeBlank,
		StatTypeComplexity,
		StatTypeCountFiles,
		StatTypeWeightedComplexity,
	}
}

func AllStatTypesString() string {
	types := AllStatsTypes()

	var strTypes strings.Builder

	strTypes.WriteString(types[0].String())

	for _, t := range types[1:] {
		strTypes.WriteString(", ")
		strTypes.WriteString(t.String())
	}

	return strTypes.String()
}
