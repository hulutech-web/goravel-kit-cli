<p align="center">
  <img src="https://github.com/hulutech-web/tinker/blob/master/images/logo.png?raw=true" width="900" />
</p>




# Goravel Kit CLI

A command-line tool to create new Goravel applications from templates.  

[goravel-kit](https://github.com/hulutech-web/goravel-kit)

# 步骤1：检查当前环境
echo "GOPATH: $(go env GOPATH)"
echo "GOBIN: $(go env GOBIN)"
echo "PATH: $PATH"

# 步骤2：安装工具
go install github.com/hulutech-web/goravel-kit-cli@latest

# 步骤3：检查安装结果
ls -la $(go env GOPATH)/bin/ | grep goravel-kit-cli

# 步骤4：验证
which goravel-kit-cli
goravel-kit-cli --version



# 重装步骤

# 删除已安装的二进制文件
rm -f $(go env GOPATH)/bin/goravel-kit-cli

# 重新安装
go install github.com/hulutech-web/goravel-kit-cli@latest