package utils

import (
	"fmt"
	"os/exec"
)

func CloneRepository(repoURL, branch, targetDir string) error {
	// 使用浅克隆加快下载速度
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", branch, repoURL, targetDir)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %s, %w", string(output), err)
	}

	return nil
}
