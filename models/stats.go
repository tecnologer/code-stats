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
