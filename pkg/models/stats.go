package models

import (
	"slices"
	"strings"
)

type Stats struct {
	Name               string  `json:"Name,omitempty"`
	Bytes              float64 `json:"Bytes,omitempty"`
	CodeBytes          float64 `json:"CodeBytes,omitempty"`
	Lines              float64 `json:"Lines,omitempty"`
	Code               float64 `json:"Code,omitempty"`
	Comment            float64 `json:"Comment,omitempty"`
	Blank              float64 `json:"Blank,omitempty"`
	Complexity         float64 `json:"Complexity,omitempty"`
	Count              float64 `json:"Count,omitempty"`
	WeightedComplexity float64 `json:"WeightedComplexity,omitempty"`
}

func (s *Stats) ValueOf(statType StatType) float64 { //nolint:cyclop
	switch statType {
	case StatTypeBytes:
		return s.Bytes
	case StatTypeCodeBytes:
		return s.CodeBytes
	case StatTypeLines:
		return s.Lines
	case StatTypeCode:
		return s.Code
	case StatTypeComment:
		return s.Comment
	case StatTypeBlank:
		return s.Blank
	case StatTypeComplexity:
		return s.Complexity
	case StatTypeCountFiles:
		return s.Count
	case StatTypeWeightedComplexity:
		return s.WeightedComplexity
	default:
		return 0
	}
}

func (s *Stats) IsInLanguageList(list []string) bool {
	return slices.ContainsFunc(list, s.EqualsName)
}

func (s *Stats) EqualsName(name string) bool {
	return strings.EqualFold(s.Name, name)
}
