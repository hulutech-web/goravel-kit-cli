package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func CloneRepositoryWithContext(ctx context.Context, repoURL, branch, targetDir string, verbose bool) error {
	args := []string{"clone", "--progress", "--depth", "1"}

	if branch != "" && branch != "main" {
		args = append(args, "--branch", branch)
	}

	args = append(args, repoURL, targetDir)

	cmd := exec.CommandContext(ctx, "git", args...)

	if verbose {
		fmt.Printf("🔧 Running command: git %s\n", strings.Join(args, " "))
	}

	// 获取标准输出管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// 记录开始时间
	startTime := time.Now()

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start git clone: %w", err)
	}

	// 实时读取输出
	go streamOutput(stdoutPipe, "git", verbose)
	go streamOutput(stderrPipe, "git", verbose)

	// 等待命令完成
	err = cmd.Wait()
	duration := time.Since(startTime)

	if err != nil {
		// 检查超时
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("download timed out after %v", duration)
		}

		// 处理特定错误类型
		switch {
		case strings.Contains(err.Error(), "Authentication failed"),
			strings.Contains(err.Error(), "could not read Username"),
			strings.Contains(err.Error(), "Permission denied"):
			return fmt.Errorf("authentication failed after %v. Try: goravel-kit-cli new %s --ssh", duration, filepath.Base(targetDir))

		case strings.Contains(err.Error(), "Repository not found"):
			return fmt.Errorf("repository not found: %s (took %v)", repoURL, duration)

		case strings.Contains(err.Error(), "could not find remote ref"):
			return fmt.Errorf("branch '%s' not found (took %v)", branch, duration)

		case strings.Contains(err.Error(), "Host key verification failed"):
			return fmt.Errorf("SSH host key verification failed. Please check your SSH configuration")

		default:
			return fmt.Errorf("git clone failed after %v: %w", duration, err)
		}
	}

	if verbose {
		fmt.Printf("✅ Download completed in %v\n", duration)
	} else {
		fmt.Printf("✅ Download completed\n")
	}

	return nil
}

// streamOutput 实时流式输出
func streamOutput(reader io.Reader, prefix string, verbose bool) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if verbose {
			fmt.Printf("%s: %s\n", prefix, line)
		} else {
			// 非详细模式只显示进度信息
			if strings.Contains(line, "Receiving objects:") ||
				strings.Contains(line, "Resolving deltas:") ||
				strings.Contains(line, "remote:") {
				fmt.Printf("📦 %s\n", strings.TrimSpace(line))
			}
		}
	}
}

// 保持兼容性
func CloneRepository(repoURL, branch, targetDir string) error {
	return CloneRepositoryWithContext(context.Background(), repoURL, branch, targetDir, false)
}
