package models

import "encoding/json"

type Stats struct {
	Name               string           `json:"Name,omitempty"`
	Bytes              int64            `json:"Bytes,omitempty"`
	CodeBytes          int64            `json:"CodeBytes,omitempty"`
	Lines              int64            `json:"Lines,omitempty"`
	Code               int64            `json:"Code,omitempty"`
	Comment            int64            `json:"Comment,omitempty"`
	Blank              int64            `json:"Blank,omitempty"`
	Complexity         int64            `json:"Complexity,omitempty"`
	Count              int64            `json:"Count,omitempty"`
	WeightedComplexity int64            `json:"WeightedComplexity,omitempty"`
	Files              []*Files         `json:"Files,omitempty"`
	LineLength         *json.RawMessage `json:"LineLength,omitempty"`
}

type Files struct{}

func (s *Stats) ValueOf(statType StatType) int64 {
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
	case StatTypeCount:
		return s.Count
	case StatTypeWeightedComplexity:
		return s.WeightedComplexity
	default:
		return 0
	}
}
