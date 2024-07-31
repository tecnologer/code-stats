package scc

import (
	"fmt"
	"os/exec"
	"strings"
)

func Run(omitDir ...string) ([]byte, error) {
	output, err := exec.Command("scc", "-f", "json", "--exclude-dir", strings.Join(omitDir, ","), "./").Output() //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to run scc: %w", err)
	}

	return output, nil
}
