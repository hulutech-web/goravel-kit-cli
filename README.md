<p align="center">
  <img src="https://github.com/hulutech-web/goravel-kit-cli/blob/master/image/logo.png?raw=true" width="900" />
</p>




# Goravel Kit CLI

A command-line tool to create new Goravel applications from templates.  

[goravel-kit](https://github.com/hulutech-web/goravel-kit)

# 步骤1：检查当前环境
``bash
echo "GOPATH: $(go env GOPATH)"
echo "GOBIN: $(go env GOBIN)"
echo "PATH: $PATH"
``
# 查看当前 PATH
``bash
echo $PATH
``
# 检查是否包含 Go 的 bin 目录
``bash
which go
``
# 添加 Go bin 目录到 PATH（如果不在）
``bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
``
# 步骤2：安装工具
``bash
go install github.com/hulutech-web/goravel-kit-cli@latest
``
# 步骤3：检查安装结果
``bash
ls -la $(go env GOPATH)/bin/ | grep goravel-kit-cli
``
# 步骤4：验证
``bash
which goravel-kit-cli
goravel-kit-cli --version
``


# 重装步骤

# 删除已安装的二进制文件
``bash
rm -f $(go env GOPATH)/bin/goravel-kit-cli
``
# 重新安装
``bash
go install github.com/hulutech-web/goravel-kit-cli@latest
``
