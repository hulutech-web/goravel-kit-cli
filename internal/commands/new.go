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
			Name:  "template",
			Usage: "GitHub template repository (format: owner/repo)",
			Value: "goravel/goravel", // 默认模板
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
	},
}

func createNewProject(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("project name is required\nUsage: goravel-kit-cli new <project-name>")
	}

	projectName := c.Args().First()
	templateRepo := c.String("template")
	branch := c.String("branch")
	force := c.Bool("force")
	verbose := c.Bool("verbose")

	if verbose {
		fmt.Printf("🚀 Creating project: %s\n", projectName)
		fmt.Printf("📦 Template: %s@%s\n", templateRepo, branch)
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
		fmt.Printf("📥 Downloading template...\n")
	}

	// 下载模板
	repoURL := fmt.Sprintf("https://github.com/%s.git", templateRepo)
	if err := utils.CloneRepository(repoURL, branch, tempDir); err != nil {
		return fmt.Errorf("failed to download template: %w", err)
	}

	if verbose {
		fmt.Printf("✅ Template downloaded successfully\n")
	}

	// 移除.git目录（如果存在）
	gitDir := filepath.Join(tempDir, ".git")
	if utils.DirectoryExists(gitDir) {
		if err := os.RemoveAll(gitDir); err != nil {
			return fmt.Errorf("failed to remove .git directory: %w", err)
		}
	}

	// 移除其他不必要的文件
	unnecessaryFiles := []string{".github", ".gitignore", "LICENSE", "README.md"}
	for _, file := range unnecessaryFiles {
		filePath := filepath.Join(tempDir, file)
		if utils.DirectoryExists(filePath) || utils.FileExists(filePath) {
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

	if verbose {
		fmt.Printf("🎉 Project '%s' created successfully!\n", projectName)
		fmt.Printf("\nNext steps:\n")
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
