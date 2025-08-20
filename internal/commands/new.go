package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			Value: "main", // 默认分支
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show verbose output",
		},
		&cli.BoolFlag{
			Name:  "ssh",
			Usage: "Use SSH URL instead of HTTPS",
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

	// 设置固定的模板仓库
	templateRepo := "hulutech-web/goravel-kit"
	var repoURL string
	if useSSH {
		repoURL = "git@github.com:hulutech-web/goravel-kit.git"
	} else {
		repoURL = "https://github.com/hulutech-web/goravel-kit.git"
	}

	if verbose {
		fmt.Printf("🚀 Creating project: %s\n", projectName)
		fmt.Printf("📦 Template: %s@%s\n", templateRepo, branch)
		fmt.Printf("🔗 URL: %s\n", repoURL)
	} else {
		fmt.Printf("Creating project: %s\n", projectName)
	}

	// 检查目录是否存在
	if utils.DirectoryExists(projectName) && !force {
		return fmt.Errorf("directory '%s' already exists. Use --force to overwrite", projectName)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "goravel-kit-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	if verbose {
		fmt.Printf("📥 Downloading template from %s...\n", repoURL)
	}

	// 下载模板
	if err := utils.CloneRepository(repoURL, branch, tempDir); err != nil {
		return fmt.Errorf("failed to download template: %w", err)
	}

	if verbose {
		fmt.Printf("✅ Template downloaded successfully\n")
		fmt.Printf("🔄 Processing template files...\n")
	}

	// 移除.git目录（如果存在）
	gitDir := filepath.Join(tempDir, ".git")
	if utils.DirectoryExists(gitDir) {
		if err := os.RemoveAll(gitDir); err != nil {
			return fmt.Errorf("failed to remove .git directory: %w", err)
		}
	}

	// 移除其他不必要的文件（可选）
	unnecessaryFiles := []string{".github", ".gitignore", "LICENSE", "README.md"}
	for _, file := range unnecessaryFiles {
		filePath := filepath.Join(tempDir, file)
		if utils.DirectoryExists(filePath) || utils.FileExists(filePath) {
			if verbose {
				fmt.Printf("🗑️  Removing: %s\n", file)
			}
			os.RemoveAll(filePath)
		}
	}

	// 重命名并移动到目标位置
	if err := utils.MoveDirectory(tempDir, projectName); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	// 更新项目中的模块名称
	if err := updateModuleName(projectName, projectName); err != nil {
		if verbose {
			fmt.Printf("⚠️  Warning: failed to update module name: %v\n", err)
		}
	}

	// 更新其他可能需要修改的文件
	if err := updateProjectFiles(projectName, projectName); err != nil {
		if verbose {
			fmt.Printf("⚠️  Warning: failed to update project files: %v\n", err)
		}
	}

	if verbose {
		fmt.Printf("🎉 Project '%s' created successfully!\n", projectName)
		fmt.Printf("\n📋 Next steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  go mod tidy\n")
		fmt.Printf("  go run .\n")
	} else {
		fmt.Printf("Project '%s' created successfully!\n", projectName)
	}

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

func updateProjectFiles(projectDir, projectName string) error {
	// 这里可以添加其他需要更新的文件
	// 例如：配置文件、环境文件等

	// 示例：更新 .env 文件中的 APP_NAME
	envPath := filepath.Join(projectDir, ".env")
	if utils.FileExists(envPath) {
		content, err := os.ReadFile(envPath)
		if err != nil {
			return err
		}

		envContent := string(content)
		envContent = strings.Replace(envContent, "APP_NAME=Goravel", "APP_NAME="+projectName, 1)

		return os.WriteFile(envPath, []byte(envContent), 0644)
	}

	return nil
}
