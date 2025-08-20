package commands

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hulutech-web/goravel-kit-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

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
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "Timeout for download operation",
			Value: 5 * time.Minute,
		},
	},
}

func display_banner() {
	fmt.Println(" ██████   ██████  ██████   █████  ██    ██ ███████ ██          ██   ██ ██ ████████      ██████ ██      ██ ")
	fmt.Println("██       ██    ██ ██   ██ ██   ██ ██    ██ ██      ██          ██  ██  ██    ██        ██      ██      ██ ")
	fmt.Println("██   ███ ██    ██ ██████  ███████ ██    ██ █████   ██          █████   ██    ██        ██      ██      ██ ")
	fmt.Println("██    ██ ██    ██ ██   ██ ██   ██  ██  ██  ██      ██          ██  ██  ██    ██        ██      ██      ██ ")
	fmt.Println(" ██████   ██████  ██   ██ ██   ██   ████   ███████ ███████     ██   ██ ██    ██         ██████ ███████ ██ ")
	fmt.Println("                                                                                                          ")
	fmt.Println("                                                                                                          ")
}

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

func createNewProject(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("project name is required\nUsage: goravel-kit-cli new <project-name>")
	}

	projectName := c.Args().First()
	branch := c.String("branch")
	force := c.Bool("force")
	verbose := c.Bool("verbose")
	useSSH := c.Bool("ssh")
	timeout := c.Duration("timeout")

	// 设置固定的模板仓库
	var repoURL string
	if useSSH {
		repoURL = "git@github.com:hulutech-web/goravel-kit.git"
	} else {
		repoURL = "https://github.com/hulutech-web/goravel-kit.git"
	}

	printWelcomeBanner(projectName)
	fmt.Printf("📦 Template: hulutech-web/goravel-kit@%s\n", branch)

	if verbose {
		fmt.Printf("🔗 URL: %s\n", repoURL)
		fmt.Printf("⏱️  Timeout: %v\n", timeout)
	}

	// 检查目录是否存在
	if utils.DirectoryExists(projectName) && !force {
		return fmt.Errorf("❌ Directory '%s' already exists. Use --force to overwrite", projectName)
	}

	// 检查网络连接
	if verbose {
		fmt.Printf("🌐 Checking network connection...\n")
	}
	if !utils.CheckGitHubAccess() {
		return fmt.Errorf("❌ Cannot access GitHub. Please check your network connection")
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "goravel-kit-*")
	if err != nil {
		return fmt.Errorf("❌ Failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil && verbose {
			fmt.Printf("⚠️  Warning: failed to clean temp directory: %v\n", err)
		}
	}()

	fmt.Printf("📥 Downloading template...\n")
	fmt.Printf("   This may take a few moments depending on your network speed.\n")

	// 使用带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 下载模板
	if err := utils.CloneRepositoryWithContext(ctx, repoURL, branch, tempDir, verbose); err != nil {
		fmt.Printf("❌ Download failed: %v\n", err)
		fmt.Printf("💡 Tips:\n")
		fmt.Printf("   - Try using --ssh flag if you have SSH keys configured\n")
		fmt.Printf("   - Check your internet connection\n")
		fmt.Printf("   - Use --verbose for more details\n")
		return fmt.Errorf("failed to download template")
	}

	fmt.Printf("✅ Template downloaded successfully\n")
	fmt.Printf("🔄 Processing template files...\n")

	// 移除.git目录
	gitDir := filepath.Join(tempDir, ".git")
	if utils.DirectoryExists(gitDir) {
		if err := os.RemoveAll(gitDir); err != nil {
			return fmt.Errorf("❌ Failed to remove .git directory: %w", err)
		}
		fmt.Printf("🗑️  Removed .git directory\n")
	}

	// 移除其他不必要的文件
	unnecessaryFiles := []string{".github", ".gitignore", "LICENSE", "README.md"}
	for _, file := range unnecessaryFiles {
		filePath := filepath.Join(tempDir, file)
		if utils.DirectoryExists(filePath) || utils.FileExists(filePath) {
			os.RemoveAll(filePath)
			if verbose {
				fmt.Printf("🗑️  Removed: %s\n", file)
			}
		}
	}

	// 如果目标目录已存在，先删除
	if utils.DirectoryExists(projectName) {
		if err := os.RemoveAll(projectName); err != nil {
			return fmt.Errorf("❌ Failed to remove existing directory: %w", err)
		}
		fmt.Printf("🗑️  Removed existing directory: %s\n", projectName)
	}

	// 移动到目标位置
	if err := utils.MoveDirectory(tempDir, projectName); err != nil {
		return fmt.Errorf("❌ Failed to create project: %w", err)
	}
	fmt.Printf("📁 Project structure created\n")

	// 更新项目中的模块名称
	if err := updateModuleName(projectName, projectName); err != nil {
		fmt.Printf("⚠️  Warning: failed to update module name: %v\n", err)
	} else {
		fmt.Printf("📝 Updated go.mod module name\n")
	}

	// 更新环境文件
	if err := updateEnvFile(projectName, projectName); err != nil {
		fmt.Printf("⚠️  Warning: failed to update .env file: %v\n", err)
	} else {
		fmt.Printf("📝 Updated .env configuration\n")
	}

	fmt.Printf("\n🎉 Project '%s' created successfully!\n", projectName)
	fmt.Printf("\n📋 Next steps:\n")
	fmt.Printf("   cd %s\n", projectName)
	fmt.Printf("   go mod tidy\n")
	fmt.Printf("   go run .\n")
	fmt.Printf("\n💡 Tip: Run with --verbose for detailed output\n")

	return nil
}

func updateModuleName(projectDir, moduleName string) error {
	goModPath := filepath.Join(projectDir, "go.mod")
	if !utils.FileExists(goModPath) {
		return nil
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "module ") {
		lines[0] = "module " + moduleName
	}

	newContent := strings.Join(lines, "\n")
	return os.WriteFile(goModPath, []byte(newContent), 0644)
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
