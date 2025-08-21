package utils

import (
	"os/exec"
)

// NewCommandWithDir 创建一个在指定目录下执行的命令
// name: 命令名称
// args: 命令参数
// dir: 工作目录
func NewCommandWithDir(name string, args []string, dir string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd
}
