# Goravel Kit CLI

A command-line tool to create new Goravel applications from templates.


# 步骤1：检查当前环境
echo "GOPATH: $(go env GOPATH)"
echo "GOBIN: $(go env GOBIN)"
echo "PATH: $PATH"

# 步骤2：安装工具
go install github.com/hulutech-web/goravel-kit-cli/cmd/goravel-kit-cli@latest

# 步骤3：检查安装结果
ls -la /Users/yuanhao/go/bin/ | grep goravel

# 步骤4：更新配置文件
echo 'export PATH=$PATH:/Users/yuanhao/go/bin' >> ~/.zshrc
source ~/.zshrc

# 步骤5：验证
which goravel-kit-cli
goravel-kit-cli --version

