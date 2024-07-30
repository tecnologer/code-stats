package command

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunSCC(omitDir ...string) ([]byte, error) {
	output, err := exec.Command("scc", "-f", "json", "--exclude-dir", strings.Join(omitDir, ","), "./").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run scc: %w", err)
	}

	return output, nil
}
