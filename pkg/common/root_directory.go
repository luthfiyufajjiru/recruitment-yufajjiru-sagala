package common

import (
	"os/exec"
	"strings"
)

func RootDirectory() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, _ := cmd.Output()
	outStr := string(output)
	return strings.Replace(outStr, "\n", "", -1)
}
