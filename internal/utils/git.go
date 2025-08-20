package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func CloneRepository(repoURL, branch, targetDir string) error {
	// 检查git是否安装
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not in PATH: %w", err)
	}

	// 使用浅克隆加快下载速度
	args := []string{"clone", "--depth", "1"}

	// 如果指定了分支，添加分支参数
	if branch != "" {
		args = append(args, "--branch", branch)
	}

	args = append(args, repoURL, targetDir)

	cmd := exec.Command("git", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		errorMsg := string(output)
		if strings.Contains(errorMsg, "could not read Username") {
			return fmt.Errorf("authentication failed. Try using SSH mode with --ssh flag")
		}
		if strings.Contains(errorMsg, "Repository not found") {
			return fmt.Errorf("repository not found: %s", repoURL)
		}
		if strings.Contains(errorMsg, "could not find remote ref") {
			return fmt.Errorf("branch not found: %s", branch)
		}
		return fmt.Errorf("git clone failed: %s, %w", errorMsg, err)
	}

	return nil
}
