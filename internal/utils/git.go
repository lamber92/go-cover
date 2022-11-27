package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// GetCurrentBranch 获取当前分支名称
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch-name. err: %v", err)
	}
	br := bufio.NewReader(bytes.NewBuffer(output))
	for {
		buff, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		return string(buff), nil
	}
	return "", fmt.Errorf("current branch-name is empty. info: %s\n", output)
}
