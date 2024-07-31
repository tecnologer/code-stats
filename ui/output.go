package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type outputType string

const (
	NormalOutput  outputType = "normal"
	SuccessOutput outputType = "success"
	DebugOutput   outputType = "debug"
	WarnOutput    outputType = "warn"
	ErrorOutput   outputType = "error"
	InfoOutput    outputType = "info"
	TipOutput     outputType = "tip"
)

var isEmojiVisible = true //nolint:gochecknoglobals

func SetEmojiVisibility(isVisible bool) {
	isEmojiVisible = isVisible
}

type OutputColorFunc func(format string, a ...any)

func NewGoodOutputFunc() OutputColorFunc {
	return NewColorFunc(color.FgHiGreen, color.Bold)
}

func NewBadOutputFunc() OutputColorFunc {
	return NewColorFunc(color.FgHiRed, color.Bold)
}

func NewWarnOutputFunc() OutputColorFunc {
	return NewColorFunc(color.FgHiYellow, color.Bold)
}

func NewNormalOutputFunc() OutputColorFunc {
	return NewColorFunc(color.FgWhite)
}

func NewColorFunc(attrs ...color.Attribute) OutputColorFunc {
	return color.New(attrs...).PrintfFunc()
}

func Printf(format string, args ...any) {
	NewNormalOutputFunc()(format, args...)
}

func Warnf(format string, args ...any) {
	if currentOutputLevel < WarnLevel {
		return
	}

	format = formatOutput(WarnOutput, format)

	NewWarnOutputFunc()(format, args...)
}

func Debugf(format string, args ...any) {
	if currentOutputLevel < DebugLevel {
		return
	}

	format = formatOutput(DebugOutput, format)

	NewColorFunc(color.FgHiWhite)(format, args...)
}

func Successf(format string, args ...any) {
	format = formatOutput(SuccessOutput, format)

	NewGoodOutputFunc()(format, args...)
}

func Errorf(format string, args ...any) {
	if currentOutputLevel < ErrorLevel {
		return
	}

	format = formatOutput(ErrorOutput, format)

	NewBadOutputFunc()(format, args...)
}

func Infof(format string, args ...any) {
	if currentOutputLevel < InfoLevel {
		return
	}

	format = formatOutput(InfoOutput, format)

	NewColorFunc(color.FgHiBlue, color.Bold)(format, args...)
}

func Tipf(format string, args ...any) {
	format = formatOutput(TipOutput, format)

	NewColorFunc(color.FgHiCyan, color.Bold)(format, args...)
}

func LogError(err error) {
	if err != nil {
		Errorf("%v", err)
	}
}

func formatOutput(outputType outputType, format string) string {
	switch outputType {
	case SuccessOutput:
		return makeFormat("✅", outputType, format)
	case DebugOutput:
		return makeFormat("🪲", outputType, format)
	case WarnOutput:
		return makeFormat("⚠️ ", outputType, format)
	case ErrorOutput:
		return makeFormat("❌ ", outputType, format)
	case InfoOutput:
		return makeFormat("ℹ️ ", outputType, format)
	case TipOutput:
		return makeFormat("💡", outputType, format)
	case NormalOutput:
		fallthrough
	default:
		return format
	}
}

func makeFormat(emoji string, outputType outputType, format string) string {
	prefix := strings.ToUpper(string(outputType))

	if isEmojiVisible {
		return fmt.Sprintf("%s [%s]: %s\n", emoji, prefix, format)
	}

	return fmt.Sprintf("[%s]: %s\n", prefix, format)
}
