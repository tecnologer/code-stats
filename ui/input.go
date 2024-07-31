package ui

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
)

type InputStringColorFunc func(format string, a ...any) string

func NewPromptOutputFunc() OutputColorFunc {
	return NewColorFunc(color.FgHiCyan, color.Bold)
}

func NewStringInputFunc() InputStringColorFunc {
	return ReadStringInput
}

func ReadStringInput(format string, attrs ...any) string {
	warn := NewWarnOutputFunc()
	promptOut := NewPromptOutputFunc()

	var input string

	r := bufio.NewReader(os.Stdin)

	for {
		promptOut(format, attrs...)

		input, _ = r.ReadString('\n')
		if input != "" {
			break
		}

		warn("please enter attrs value, and try again\n")
	}

	return strings.TrimSpace(input)
}
