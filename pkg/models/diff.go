package models

//go:generate enumer -type=DifferenceType -trimprefix=Diff -transform=kebab
type DifferenceType byte

const (
	DiffNone DifferenceType = iota
	DiffPreviousDate
	DiffFirstDate
	DiffSpecificDate
)
