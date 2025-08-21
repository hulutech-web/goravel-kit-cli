package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/hulutech-web/goravel-kit-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

// 添加版权信息显示函数
func printWelcomeBanner(projectName string) {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	fmt.Printf("\n")
	fmt.Printf("\n")
	cyan.Println(" ██████   ██████  ██████   █████  ██    ██ ███████ ██          ██   ██ ██ ████████      ██████ ██      ██ ")
	cyan.Println("██       ██    ██ ██   ██ ██   ██ ██    ██ ██      ██          ██  ██  ██    ██        ██      ██      ██ ")
	cyan.Println("██   ███ ██    ██ ██████  ███████ ██    ██ █████   ██          █████   ██    ██        ██      ██      ██ ")
	cyan.Println("██    ██ ██    ██ ██   ██ ██   ██  ██  ██  ██      ██          ██  ██  ██    ██        ██      ██      ██ ")
	cyan.Println(" ██████   ██████  ██   ██ ██   ██   ████   ███████ ███████     ██   ██ ██    ██         ██████ ███████ ██ ")
	cyan.Println("         ")
	green.Println("                    +++++++++++++++++++🎉欢迎使用 Goravel Kit CLI 🏆+++++++++++++++++++")
	yellow.Printf("                    |·<<达州葫芦科技>>研发\n")
	yellow.Printf("                    |·作者: yuanhaozhuzhu@hotmail.com\n")
	yellow.Printf("                    |·开发时间: 2025-08-22\n")
	yellow.Printf("                    |·版本号: v1.0.0\n")
	yellow.Printf("                    |·版本说明: Goravel 项目脚手架工具\n")
	yellow.Printf("                    |·版本时间: 2025-08-22\n")
	cyan.Println("                    ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	fmt.Printf("\n")
	color.New(color.FgHiWhite, color.Bold).Printf("🚀 开始创建 Goravel 项目: %s\n", projectName)
	fmt.Printf("\n")
}

var NewCommand = &cli.Command{
	Name:      "new",
	Usage:     "Create a new Goravel application from template",
	ArgsUsage: "<project-name>",
	Action:    createNewProject,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force create project even if directory exists",
		},
		&cli.StringFlag{
			Name:  "branch",
			Usage: "Git branch to use",
			Value: "master",
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show verbose output",
		},
		&cli.BoolFlag{
			Name:  "ssh",
			Usage: "Use SSH URL instead of HTTPS",
			Value: true, // 默认启用 SSH
		},
		&cli.BoolFlag{
			Name:  "https",
			Usage: "Use HTTPS URL instead of SSH",
			Value: false,
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "Timeout for download operation",
			Value: 3 * time.Minute,
		},
		&cli.BoolFlag{
			Name:  "no-banner",
			Usage: "Don't show welcome banner",
		},
		&cli.BoolFlag{
			Name:  "gitee-only",
			Usage: "Use Gitee mirror only (skip GitHub)",
		},
		&cli.BoolFlag{
			Name:  "github-only",
			Usage: "Use GitHub only (skip Gitee fallback)",
		},
	},
}

func createNewProject(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("project name is required\nUsage: goravel-kit-cli new <project-name>")
	}

	projectName := c.Args().First()
	branch := c.String("branch")
	force := c.Bool("force")
	verbose := c.Bool("verbose")
	useSSH := c.Bool("ssh")
	useHTTPS := c.Bool("https")
	timeout := c.Duration("timeout")
	noBanner := c.Bool("no-banner")
	giteeOnly := c.Bool("gitee-only")
	githubOnly := c.Bool("github-only")

	// 处理协议选择逻辑：如果同时指定了 --https，优先使用 HTTPS
	var protocol string
	if useHTTPS {
		protocol = "https"
		useSSH = false
	} else {
		protocol = "ssh"
		useSSH = true
	}

	// 显示版权信息（除非指定不显示）
	if !noBanner {
		printWelcomeBanner(projectName)
	} else {
		color.New(color.FgHiWhite, color.Bold).Printf("🚀 Creating Goravel project: %s\n", projectName)
		fmt.Printf("\n")
	}

	// 智能选择镜像源策略
	var autoDetectedGiteeOnly bool
	var networkStatus string

	// 如果不是强制指定了镜像源，就自动检测网络
	if !giteeOnly && !githubOnly {
		color.New(color.FgHiCyan).Printf("🌐 检测网络连接...\n")

		if utils.CheckGiteeAccess() {
			networkStatus = "Gitee 访问正常"
			autoDetectedGiteeOnly = true
			// 自动启用 gitee-only 模式
			giteeOnly = true
		} else {
			networkStatus = "GitHub 访问失败，自动切换到 GitHub"
			autoDetectedGiteeOnly = true
		}
		color.New(color.FgHiCyan).Printf("   %s\n", networkStatus)
	}

	// 定义镜像源
	mirrors := []struct {
		name    string
		url     string
		sshURL  string
		enabled bool
	}{
		{
			name:    "GitHub",
			url:     "https://github.com/hulutech-web/goravel-kit.git",
			sshURL:  "git@github.com:hulutech-web/goravel-kit.git",
			enabled: !giteeOnly,
		},
		{
			name:    "Gitee",
			url:     "https://gitee.com/hulutech/goravel-kit.git",
			sshURL:  "git@gitee.com:hulutech/goravel-kit.git",
			enabled: !githubOnly,
		},
	}

	// 显示当前使用的镜像源策略
	color.New(color.FgHiBlue).Printf("📦 模板策略: ")
	switch {
	case giteeOnly && autoDetectedGiteeOnly:
		color.New(color.FgHiBlue).Printf("自动选择 Gitee 镜像 (网络检测)\n")
	case giteeOnly:
		color.New(color.FgHiBlue).Printf("强制使用 Gitee 镜像 (用户指定)\n")
	case githubOnly:
		color.New(color.FgHiBlue).Printf("强制使用 GitHub 镜像 (用户指定)\n")
	default:
		color.New(color.FgHiBlue).Printf("自动选择镜像 (GitHub → Gitee)\n")
	}

	color.New(color.FgHiBlue).Printf("🌿 分支: %s\n", branch)
	color.New(color.FgHiBlue).Printf("🔗 协议: %s (默认)\n", protocol)

	if verbose {
		color.New(color.FgHiMagenta).Printf("📡 可用镜像源:\n")
		for _, mirror := range mirrors {
			if mirror.enabled {
				var url string
				if useSSH {
					url = mirror.sshURL
				} else {
					url = mirror.url
				}
				color.New(color.FgHiMagenta).Printf("   - %s: %s\n", mirror.name, url)
			}
		}
		color.New(color.FgHiYellow).Printf("⏱️  超时时间: %v\n", timeout)
	}

	// 检查目录是否存在
	if utils.DirectoryExists(projectName) && !force {
		return fmt.Errorf("❌ 目录 '%s' 已存在。使用 --force 参数覆盖", projectName)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "goravel-kit-*")
	if err != nil {
		return fmt.Errorf("❌ 创建临时目录失败: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil && verbose {
			color.New(color.FgHiRed).Printf("⚠️  警告: 清理临时目录失败: %v\n", err)
		}
	}()

	var downloadError error
	var successMirror string
	var successRepoURL string

	// 尝试从各个镜像源下载
	for _, mirror := range mirrors {
		if !mirror.enabled {
			continue
		}

		// 选择 URL（默认使用 SSH）
		var repoURL string
		if useSSH {
			repoURL = mirror.sshURL
		} else {
			repoURL = mirror.url
		}

		color.New(color.FgHiGreen).Printf("\n📥 尝试从 %s 下载模板...\n", mirror.name)
		color.New(color.FgHiCyan).Printf("   📍 仓库: %s\n", repoURL)
		color.New(color.FgHiCyan).Printf("   🌿 分支: %s\n", branch)

		// 使用带超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		// 下载模板
		err := utils.CloneRepositoryWithContext(ctx, repoURL, branch, tempDir, verbose)
		cancel()

		if err != nil {
			downloadError = err
			color.New(color.FgHiRed).Printf("❌ %s 下载失败: %v\n", mirror.name, err)

			// 如果不是最后一个镜像源，继续尝试下一个
			if hasNextMirror(mirrors, mirror.name) {
				color.New(color.FgHiYellow).Printf("🔄 尝试下一个镜像源...\n")
				continue
			}
		} else {
			successMirror = mirror.name
			successRepoURL = repoURL
			downloadError = nil
			break
		}
	}

	// 检查下载结果
	if downloadError != nil {
		color.New(color.FgHiRed).Printf("\n❌ 所有镜像源下载均失败！\n")
		color.New(color.FgHiYellow).Printf("💡 解决方案:\n")
		color.New(color.FgHiYellow).Printf("   1. 检查网络连接\n")
		color.New(color.FgHiYellow).Printf("   2. 使用 --ssh 参数尝试 SSH 方式\n")
		color.New(color.FgHiYellow).Printf("   3. 使用 --gitee-only 强制使用 Gitee\n")
		color.New(color.FgHiYellow).Printf("   4. 使用 --github-only 强制使用 GitHub\n")
		color.New(color.FgHiYellow).Printf("   5. 使用 --verbose 查看详细错误信息\n")
		color.New(color.FgHiYellow).Printf("   6. 检查分支是否存在: %s\n", branch)
		return fmt.Errorf("所有镜像源下载失败")
	}

	color.New(color.FgHiGreen).Printf("\n✅ 成功从 %s 下载模板\n", successMirror)
	color.New(color.FgHiCyan).Printf("   📍 源仓库: %s\n", successRepoURL)
	color.New(color.FgHiCyan).Printf("   🌿 分支: %s\n", branch)
	color.New(color.FgHiGreen).Printf("🔄 处理模板文件中...\n")

	// 移除.git目录
	gitDir := filepath.Join(tempDir, ".git")
	if utils.DirectoryExists(gitDir) {
		if err := os.RemoveAll(gitDir); err != nil {
			return fmt.Errorf("❌ 移除 .git 目录失败: %w", err)
		}
		if verbose {
			color.New(color.FgHiYellow).Printf("🗑️  已移除 .git 目录\n")
		}
	}

	// 移除其他不必要的文件
	unnecessaryFiles := []string{".github", ".gitignore", "LICENSE", "README.md"}
	for _, file := range unnecessaryFiles {
		filePath := filepath.Join(tempDir, file)
		if utils.DirectoryExists(filePath) || utils.FileExists(filePath) {
			os.RemoveAll(filePath)
			if verbose {
				color.New(color.FgHiYellow).Printf("🗑️  已移除: %s\n", file)
			}
		}
	}

	// 如果目标目录已存在，先删除
	if utils.DirectoryExists(projectName) {
		if err := os.RemoveAll(projectName); err != nil {
			return fmt.Errorf("❌ 移除已存在目录失败: %w", err)
		}
		if verbose {
			color.New(color.FgHiYellow).Printf("🗑️  已移除已存在目录: %s\n", projectName)
		}
	}

	// 移动到目标位置
	if err := moveDirectoryCrossPlatform(tempDir, projectName); err != nil {
		return fmt.Errorf("❌ 创建项目失败: %w", err)
	}
	color.New(color.FgHiGreen).Printf("📁 项目结构创建完成\n")
	// 创建.env文件，通过copy .env.example得到，然后再更新
	// 创建 .env 文件，通过复制 .env.example 得到
	envExamplePath := filepath.Join(projectName, ".env.example")
	envPath := filepath.Join(projectName, ".env")
	if utils.FileExists(envExamplePath) {
		input, err := os.ReadFile(envExamplePath)
		if err != nil {
			return fmt.Errorf("❌ 读取 .env.example 文件失败: %w", err)
		}
		err = os.WriteFile(envPath, input, 0644)
		if err != nil {
			return fmt.Errorf("❌ 创建 .env 文件失败: %w", err)
		}
		if verbose {
			color.New(color.FgHiGreen).Printf("✅ 已从 .env.example 复制生成 .env 文件\n")
		}
	} else {
		if verbose {
			color.New(color.FgHiYellow).Printf("⚠️  未找到 .env.example，跳过 .env 文件创建\n")
		}
	}

	// 更新环境文件
	if err := updateEnvFile(projectName, projectName); err != nil {
		color.New(color.FgHiYellow).Printf("⚠️  警告: 更新 .env 文件失败: %v\n", err)
	} else {
		color.New(color.FgHiGreen).Printf("📝 已更新 .env 配置\n")
	}

	// 运行命令行工具，进入项目根路径，执行go run . artisan key:generate，之后再执行go run . artisan jwt:secret
	// 在项目根目录下依次执行 go run . artisan key:generate 和 go run . artisan jwt:secret
	commands := [][]string{
		{"go", "run", ".", "artisan", "key:generate"},
		{"go", "run", ".", "artisan", "jwt:secret"},
	}

	for _, cmdArgs := range commands {
		cmd := utils.NewCommandWithDir(cmdArgs[0], cmdArgs[1:], projectName)
		if verbose {
			color.New(color.FgHiCyan).Printf("🔧 执行命令: %s\n", strings.Join(cmdArgs, " "))
		}
		output, err := cmd.CombinedOutput()
		if verbose {
			fmt.Print(string(output))
		}
		if err != nil {
			color.New(color.FgHiRed).Printf("❌ 命令执行失败: %s\n", strings.Join(cmdArgs, " "))
			color.New(color.FgHiRed).Printf("   错误信息: %v\n", err)
			break
		}
	}

	color.New(color.FgHiCyan, color.Bold).Printf("\n🎉 项目 '%s' 创建成功！\n", projectName)
	color.New(color.FgHiWhite).Printf("\n📋 下一步操作:\n")
	color.New(color.FgHiGreen).Printf("   cd %s\n", projectName)
	color.New(color.FgHiGreen).Printf("   go mod tidy\n")
	color.New(color.FgHiGreen).Printf("   modify .env database configuration!\n")
	color.New(color.FgHiGreen).Printf("   air\n")
	color.New(color.FgHiYellow).Printf("\n💡 提示: 使用 --verbose 参数查看详细输出\n")

	return nil
}

// hasNextMirror 检查是否还有下一个可用的镜像源
func hasNextMirror(mirrors []struct {
	name    string
	url     string
	sshURL  string
	enabled bool
}, currentMirror string) bool {
	foundCurrent := false
	for _, mirror := range mirrors {
		if !mirror.enabled {
			continue
		}
		if foundCurrent {
			return true
		}
		if mirror.name == currentMirror {
			foundCurrent = true
		}
	}
	return false
}

func updateEnvFile(projectDir, projectName string) error {
	envPath := filepath.Join(projectDir, ".env")
	if !utils.FileExists(envPath) {
		return nil
	}

	content, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}

	envContent := string(content)
	envContent = strings.Replace(envContent, "APP_NAME=Goravel", "APP_NAME="+projectName, 1)
	envContent = strings.Replace(envContent, "APP_URL=http://localhost", "APP_URL=http://localhost:3000", 1)

	return os.WriteFile(envPath, []byte(envContent), 0644)
}

// moveDirectoryCrossPlatform 跨平台的目录移动函数
func moveDirectoryCrossPlatform(source, destination string) error {
	// 尝试直接重命名（同磁盘分区时有效）
	err := os.Rename(source, destination)
	if err == nil {
		return nil
	}

	// 如果重命名失败（可能是因为跨磁盘），使用复制+删除的方式
	color.New(color.FgHiYellow).Printf("⚠️  跨磁盘操作，使用复制方式移动文件...\n")

	// 创建目标目录
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 复制所有文件和子目录
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destination, relPath)

		if info.IsDir() {
			// 创建目录
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// 复制文件
			return copyFile(path, destPath)
		}
	})

	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	// 删除源目录
	if err := os.RemoveAll(source); err != nil {
		return fmt.Errorf("清理源目录失败: %w", err)
	}

	return nil
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 复制文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}
