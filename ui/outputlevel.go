package ui

type OutputLevel byte

var currentOutputLevel = InfoLevel //nolint:gochecknoglobals

const (
	ErrorLevel OutputLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

func SetOutputLevel(level OutputLevel) {
	currentOutputLevel = level
}

func GetOutputLevel() OutputLevel {
	return currentOutputLevel
}

func (l OutputLevel) String() string {
	return [...]string{"Error", "Warn", "Info", "Debug"}[l]
}
