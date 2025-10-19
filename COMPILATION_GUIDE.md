# 3X-UI 编译部署指南 / Compilation and Deployment Guide

[中文](#中文指南) | [English](#english-guide)

---

## 中文指南

本指南提供跨平台安装的总览说明，并链接至拆分出的各平台文档。

### 目录
1. [平台导航](#平台导航)
2. [系统要求](#系统要求)
3. [依赖安装](#依赖安装)
4. [常见问题](#常见问题)
5. [English Guide](#english-guide)

### 平台导航

- [Linux 编译与部署](./docs/guides/linux.md)
- [其他 Linux 场景（交叉编译、Docker、Alpine 等）](./docs/guides/linux_other.md)
- [Windows 编译部署](./docs/guides/windows.md)
- [macOS 部署指南](./docs/guides/macos.md)

如需快速安装，可直接运行：

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Jiusi-pys/3x-ui/main/install.sh)
```

脚本会自动下载 Release 包、部署 Xray 与面板，并注册系统服务。

---

### 系统要求

#### 最低配置
- **CPU**: 1核心
- **内存**: 512 MB RAM
- **磁盘**: 100 MB 可用空间
- **网络**: 互联网连接（用于下载依赖）

#### 推荐配置
- **CPU**: 2核心及以上
- **内存**: 1 GB RAM 及以上
- **磁盘**: 500 MB 可用空间
- **网络**: 稳定的互联网连接

---

### 依赖安装

#### Go 环境要求
本项目需要 **Go 1.25.1** 或更高版本。

##### Linux 安装 Go
```bash
# 下载 Go
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz

# 解压到 /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

##### Windows 安装 Go
1. 下载安装包：https://go.dev/dl/go1.25.1.windows-amd64.msi
2. 双击运行安装程序
3. 重启终端并验证：
```powershell
go version
```

##### macOS 安装 Go
```bash
# 使用 Homebrew
brew install go@1.25

# 或者下载安装包
# https://go.dev/dl/go1.25.1.darwin-amd64.pkg

# 验证安装
go version
```

#### 其他依赖

##### Linux 依赖
```bash
# Debian/Ubuntu
sudo apt-get update
sudo apt-get install -y gcc g++ make git unzip wget

# RHEL/CentOS/Fedora
sudo yum install -y gcc gcc-c++ make git unzip wget

# Arch Linux
sudo pacman -S gcc make git unzip wget
```

##### Windows 依赖
- **Git**: https://git-scm.com/download/win
- **MinGW-w64** (GCC编译器): https://www.mingw-w64.org/downloads/
- 或者安装 **Visual Studio Build Tools**

##### macOS 依赖
```bash
# 安装 Xcode Command Line Tools
xcode-select --install

# 或使用 Homebrew
brew install git wget
```

---

### Linux 编译

完整的构建、部署与维护流程请参阅 [docs/guides/linux.md](./docs/guides/linux.md)。

---

### Windows 编译

完整编译与服务配置指南请参阅 [docs/guides/windows.md](./docs/guides/windows.md)。

---

### macOS 编译

详见 [docs/guides/macos.md](./docs/guides/macos.md)。

---

### 交叉编译

#### Linux 交叉编译到多架构

```bash
#!/bin/bash
# build-all.sh - 交叉编译脚本

platforms=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm/7"    # ARMv7
    "linux/arm/6"    # ARMv6
    "linux/arm/5"    # ARMv5
    "linux/386"
    "linux/s390x"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r os arch arm <<< "$platform"

    output="x-ui-${os}-${arch}"
    [ -n "$arm" ] && output="${output}v${arm}"
    [ "$os" = "windows" ] && output="${output}.exe"

    echo "编译 ${os}/${arch}${arm:+v$arm}..."

    export CGO_ENABLED=1
    export GOOS=$os
    export GOARCH=$arch
    [ -n "$arm" ] && export GOARM=$arm

    # 设置交叉编译工具链（Linux 静态编译需要）
    if [ "$os" = "linux" ]; then
        go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o "$output" main.go
    else
        go build -ldflags "-w -s" -o "$output" main.go
    fi

    [ $? -eq 0 ] && echo "✓ $output 编译成功" || echo "✗ $output 编译失败"
done
```

#### 使用 Docker 交叉编译（推荐）

```bash
# 使用 Docker Buildx 进行多架构编译
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --tag 3x-ui:latest \
  .
```

---

### 部署说明

- Linux 安装/手动部署：见 [docs/guides/linux.md](./docs/guides/linux.md)
- Windows 安装：见 [docs/guides/windows.md](./docs/guides/windows.md)
- macOS 安装：见 [docs/guides/macos.md](./docs/guides/macos.md)
- 交叉编译、Alpine/OpenRC、Docker：见 [docs/guides/linux_other.md](./docs/guides/linux_other.md)

#### 访问面板

编译部署完成后，通过浏览器访问：

```
http://服务器IP:2053/
```

默认登录凭据：
- **用户名**: `admin`
- **密码**: `admin`

⚠️ **重要安全提示**：首次登录后请立即修改默认密码！

---

### 常见问题

#### 1. 编译错误：`gcc: command not found`

**解决方案**：安装 GCC 编译器

```bash
# Debian/Ubuntu
sudo apt-get install gcc

# RHEL/CentOS
sudo yum install gcc

# macOS
xcode-select --install

# Windows
# 安装 MinGW-w64 或 Visual Studio Build Tools
```

#### 2. 编译错误：`package embed is not in GOROOT`

**解决方案**：升级 Go 版本到 1.16 或更高

```bash
go version  # 检查当前版本
# 访问 https://go.dev/dl/ 下载最新版本
```

#### 3. 运行错误：`cannot open shared object file`

**解决方案**：使用静态编译或安装缺失的库

```bash
# 静态编译
CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags '-static'" -o x-ui main.go

# 或安装缺失的库
ldd x-ui  # 查看缺失的库
sudo apt-get install <缺失的库>
```

#### 4. 权限错误：`permission denied`

**解决方案**：赋予执行权限

```bash
chmod +x x-ui
# 或使用 sudo
sudo ./x-ui
```

#### 5. 端口被占用：`address already in use`

**解决方案**：更改默认端口或终止占用进程

```bash
# 查找占用端口的进程
sudo lsof -i :2053
# 或
sudo netstat -tulpn | grep 2053

# 终止进程
sudo kill -9 <PID>

# 或修改配置文件中的端口
```

#### 6. Xray 启动失败

**解决方案**：检查 Xray 配置和日志

```bash
# 查看 x-ui 日志
sudo journalctl -u x-ui -n 100

# 手动测试 Xray
./bin/xray -test -config /etc/x-ui/xray.json

# 检查 Xray 版本兼容性
./bin/xray -version
```

#### 7. 数据库初始化错误

**解决方案**：删除数据库文件重新初始化

```bash
# Linux/macOS
sudo rm /etc/x-ui/x-ui.db
sudo systemctl restart x-ui

# Windows
Remove-Item "C:\Program Files\x-ui\x-ui.db"
Restart-Service x-ui
```

#### 8. macOS 上提示"无法验证开发者"

**解决方案**：移除隔离属性

```bash
sudo xattr -r -d com.apple.quarantine /usr/local/x-ui/x-ui
```

#### 9. 交叉编译 CGO 错误

**解决方案**：安装交叉编译工具链

```bash
# Linux 编译到 Windows
sudo apt-get install gcc-mingw-w64

# 设置交叉编译器
export CC=x86_64-w64-mingw32-gcc
export CXX=x-86_64-w64-mingw32-g++
export CGO_ENABLED=1
export GOOS=windows
export GOARCH=amd64

go build -o x-ui.exe main.go
```

#### 10. Docker 编译权限问题

**解决方案**：使用用户映射

```bash
docker run --rm \
  --user $(id -u):$(id -g) \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  sh -c "go mod download && go build -o x-ui main.go"
```

---

## English Guide

### Table of Contents
1. [System Requirements](#system-requirements-en)
2. [Dependency Installation](#dependency-installation-en)
3. [Linux Compilation](#linux-compilation-en)
4. [Windows Compilation](#windows-compilation-en)
5. [macOS Compilation](#macos-compilation-en)
6. [Cross-Compilation](#cross-compilation-en)
7. [Deployment](#deployment-en)
8. [Troubleshooting](#troubleshooting-en)

---

### System Requirements (EN)

#### Minimum Requirements
- **CPU**: 1 core
- **RAM**: 512 MB
- **Disk**: 100 MB available space
- **Network**: Internet connection (for downloading dependencies)

#### Recommended Requirements
- **CPU**: 2 cores or more
- **RAM**: 1 GB or more
- **Disk**: 500 MB available space
- **Network**: Stable internet connection

---

### Dependency Installation (EN)

#### Go Environment
This project requires **Go 1.25.1** or higher.

##### Install Go on Linux
```bash
# Download Go
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz

# Extract to /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz

# Configure environment variables
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

##### Install Go on Windows
1. Download installer: https://go.dev/dl/go1.25.1.windows-amd64.msi
2. Run the installer
3. Restart terminal and verify:
```powershell
go version
```

##### Install Go on macOS
```bash
# Using Homebrew
brew install go@1.25

# Or download installer
# https://go.dev/dl/go1.25.1.darwin-amd64.pkg

# Verify installation
go version
```

#### Additional Dependencies

##### Linux Dependencies
```bash
# Debian/Ubuntu
sudo apt-get update
sudo apt-get install -y gcc g++ make git unzip wget

# RHEL/CentOS/Fedora
sudo yum install -y gcc gcc-c++ make git unzip wget

# Arch Linux
sudo pacman -S gcc make git unzip wget
```

##### Windows Dependencies
- **Git**: https://git-scm.com/download/win
- **MinGW-w64** (GCC compiler): https://www.mingw-w64.org/downloads/
- Or install **Visual Studio Build Tools**

##### macOS Dependencies
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Or use Homebrew
brew install git wget
```

---

### Linux Compilation (EN)

See the dedicated [Linux guide](./docs/guides/linux.md) for build and deployment steps.

---

### Windows Compilation (EN)

Refer to [docs/guides/windows.md](./docs/guides/windows.md) for build, service, and firewall instructions.

---

### macOS Compilation (EN)

Use [docs/guides/macos.md](./docs/guides/macos.md) for build, LaunchDaemon bootstrap, and troubleshooting.

---

### Cross-Compilation (EN)

Docker Buildx, multi-arch builds, and Alpine/OpenRC notes are documented in [docs/guides/linux_other.md](./docs/guides/linux_other.md).

---

### Deployment (EN)

- Linux deployment instructions: [docs/guides/linux.md](./docs/guides/linux.md)
- Windows deployment: [docs/guides/windows.md](./docs/guides/windows.md)
- macOS deployment: [docs/guides/macos.md](./docs/guides/macos.md)
- Cross-platform builds & Alpine/OpenRC: [docs/guides/linux_other.md](./docs/guides/linux_other.md)

#### Access the Panel

After deployment, access the panel via browser:

```
http://server-ip:2053/
```

Default credentials:
- **Username**: `admin`
- **Password**: `admin`

⚠️ **Important Security Notice**: Change the default password immediately after first login!

---

### Troubleshooting (EN)

For detailed troubleshooting, please refer to the Chinese guide above or visit the project's GitHub Issues page.

Common issues include:
- GCC compiler not found
- Go version too old
- Shared library errors
- Permission denied
- Port already in use
- Xray startup failures
- Database initialization errors
- Cross-compilation CGO errors

---

## 新增的出站规则模块 / New Outbound Rules Module

### 功能说明 / Features

本次更新添加了完整的出站规则管理模块，参照入站规则的实现方式：

This update adds a complete outbound rules management module, following the inbound rules implementation pattern:

#### 数据库模型 / Database Model
- 新增 `Outbound` 表用于存储出站配置
- 支持多种协议：freedom, blackhole, vmess, vless, trojan, shadowsocks, socks, http, dns
- Added `Outbound` table for storing outbound configurations
- Supports multiple protocols: freedom, blackhole, vmess, vless, trojan, shadowsocks, socks, http, dns

#### API 端点 / API Endpoints
- `GET /panel/api/outbounds/list` - 获取出站列表 / Get outbound list
- `GET /panel/api/outbounds/get/:id` - 获取单个出站 / Get single outbound
- `GET /panel/api/outbounds/tags` - 获取出站标签 / Get outbound tags
- `POST /panel/api/outbounds/add` - 添加出站 / Add outbound
- `POST /panel/api/outbounds/update/:id` - 更新出站 / Update outbound
- `POST /panel/api/outbounds/del/:id` - 删除出站 / Delete outbound

#### 自动集成 / Automatic Integration
- 数据库中的出站配置自动与 Xray 配置模板合并
- 修改出站后自动标记需要重启 Xray
- Database outbounds automatically merge with Xray config template
- Automatically marks Xray for restart after outbound modifications

### 使用方法 / Usage

1. 访问面板 `/panel/outbounds` 页面
2. 添加、编辑或删除出站配置
3. 配置将自动应用到 Xray（需重启）

1. Access the panel at `/panel/outbounds`
2. Add, edit, or delete outbound configurations
3. Configurations will be automatically applied to Xray (restart required)

---

## 版权信息 / Copyright

本项目基于 GPL-3.0 许可证开源。详见 LICENSE 文件。

This project is open-sourced under the GPL-3.0 License. See LICENSE file for details.

## 支持 / Support

- **GitHub Issues**: https://github.com/MHSanaei/3x-ui/issues
- **Wiki**: https://github.com/MHSanaei/3x-ui/wiki
- **Telegram**: (请查看项目 README / See project README)

---

**文档版本 / Document Version**: 1.0
**最后更新 / Last Updated**: 2025-10-10
**项目版本 / Project Version**: Based on 3x-ui v2.x
