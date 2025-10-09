# 3X-UI 编译部署指南 / Compilation and Deployment Guide

[中文](#中文指南) | [English](#english-guide)

---

## 中文指南

本指南将帮助您在 Linux、Windows 和 macOS 平台上编译和部署 3X-UI 项目。

### 目录
1. [系统要求](#系统要求)
2. [依赖安装](#依赖安装)
3. [Linux 编译](#linux-编译)
4. [Windows 编译](#windows-编译)
5. [macOS 编译](#macos-编译)
6. [交叉编译](#交叉编译)
7. [部署说明](#部署说明)
8. [常见问题](#常见问题)

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

#### 方法一：简单编译（动态链接）

```bash
# 1. 克隆项目
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. 下载 Go 模块依赖
go mod download

# 3. 编译
CGO_ENABLED=1 go build -o x-ui main.go

# 4. 验证编译结果
./x-ui -v
```

#### 方法二：静态编译（推荐用于生产环境）

```bash
# 1. 克隆项目
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. 下载依赖
go mod download

# 3. 静态编译
CGO_ENABLED=1 go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o x-ui main.go

# 4. 验证是静态链接
ldd x-ui
# 应该显示 "not a dynamic executable" 或无输出

# 5. 检查文件大小和类型
ls -lh x-ui
file x-ui
```

#### 方法三：使用 Docker 编译（无需本地安装 Go）

```bash
# 使用官方 Go 镜像编译
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  sh -c "go mod download && CGO_ENABLED=1 go build -o x-ui main.go"
```

#### 下载 Xray-core 和地理数据库

```bash
# 创建 bin 目录
mkdir -p bin
cd bin

# 下载 Xray-core (根据架构选择)
# amd64
wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-64.zip
unzip Xray-linux-64.zip
rm Xray-linux-64.zip

# arm64
# wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-arm64-v8a.zip
# unzip Xray-linux-arm64-v8a.zip

# 下载地理数据库
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat

# 下载区域特定的地理数据库（可选）
wget -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat

cd ..
```

---

### Windows 编译

#### 方法一：使用 PowerShell 编译

```powershell
# 1. 克隆项目
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. 设置环境变量
$env:CGO_ENABLED="1"
$env:GOOS="windows"
$env:GOARCH="amd64"

# 3. 下载依赖
go mod download

# 4. 编译
go build -ldflags "-w -s" -o x-ui.exe main.go

# 5. 验证
.\x-ui.exe -v
```

#### 下载 Xray-core 和地理数据库

```powershell
# 创建 bin 目录
New-Item -Path bin -ItemType Directory -Force
cd bin

# 下载 Xray-core
Invoke-WebRequest -Uri "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-windows-64.zip" -OutFile "Xray-windows-64.zip"
Expand-Archive -Path "Xray-windows-64.zip" -DestinationPath .
Remove-Item "Xray-windows-64.zip"

# 下载地理数据库
Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" -OutFile "geoip.dat"
Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" -OutFile "geosite.dat"

# 下载区域特定数据库（可选）
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat" -OutFile "geoip_IR.dat"
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat" -OutFile "geosite_IR.dat"
Invoke-WebRequest -Uri "https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat" -OutFile "geoip_RU.dat"
Invoke-WebRequest -Uri "https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat" -OutFile "geosite_RU.dat"

cd ..
```

#### 创建 Windows 服务（可选）

```powershell
# 使用 NSSM 创建 Windows 服务
# 1. 下载 NSSM: https://nssm.cc/download
# 2. 安装服务
nssm install x-ui "C:\path\to\x-ui.exe"
nssm set x-ui AppDirectory "C:\path\to\"
nssm start x-ui
```

---

### macOS 编译

#### 编译步骤

```bash
# 1. 克隆项目
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. 设置环境变量
export CGO_ENABLED=1
export GOOS=darwin
export GOARCH=amd64  # 或 arm64 用于 Apple Silicon

# 3. 下载依赖
go mod download

# 4. 编译
go build -ldflags "-w -s" -o x-ui main.go

# 5. 赋予执行权限
chmod +x x-ui

# 6. 验证
./x-ui -v
```

#### 下载 Xray-core 和地理数据库

```bash
# 创建 bin 目录
mkdir -p bin
cd bin

# 下载 Xray-core
# Intel Mac (amd64)
wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-macos-64.zip
unzip Xray-macos-64.zip
rm Xray-macos-64.zip

# Apple Silicon (arm64)
# wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-macos-arm64-v8a.zip
# unzip Xray-macos-arm64-v8a.zip

# 下载地理数据库
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat
wget -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat

cd ..
```

#### 创建 macOS LaunchDaemon（可选）

```bash
# 创建 plist 文件
sudo tee /Library/LaunchDaemons/com.x-ui.plist > /dev/null <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.x-ui</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/x-ui</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/var/log/x-ui.err</string>
    <key>StandardOutPath</key>
    <string>/var/log/x-ui.out</string>
</dict>
</plist>
EOF

# 加载服务
sudo launchctl load /Library/LaunchDaemons/com.x-ui.plist
```

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

#### Linux 部署

##### 1. 使用安装脚本（推荐）

```bash
# 一键安装
bash <(curl -Ls https://raw.githubusercontent.com/mhsanaei/3x-ui/main/install.sh)
```

##### 2. 手动部署

```bash
# 创建目录结构
sudo mkdir -p /usr/local/x-ui/bin

# 复制文件
sudo cp x-ui /usr/local/x-ui/
sudo cp -r bin/* /usr/local/x-ui/bin/
sudo chmod +x /usr/local/x-ui/x-ui
sudo chmod +x /usr/local/x-ui/bin/xray

# 创建 systemd 服务
sudo tee /etc/systemd/system/x-ui.service > /dev/null <<EOF
[Unit]
Description=3x-ui Service
After=network.target nss-lookup.target

[Service]
User=root
WorkingDirectory=/usr/local/x-ui
ExecStart=/usr/local/x-ui/x-ui
Restart=on-failure
RestartPreventExitStatus=23
LimitNPROC=10000
LimitNOFILE=1000000

[Install]
WantedBy=multi-user.target
EOF

# 重载 systemd 并启动服务
sudo systemctl daemon-reload
sudo systemctl enable x-ui
sudo systemctl start x-ui

# 查看服务状态
sudo systemctl status x-ui

# 查看日志
sudo journalctl -u x-ui -f
```

##### 3. 配置防火墙

```bash
# UFW (Ubuntu/Debian)
sudo ufw allow 54321/tcp  # 默认面板端口
sudo ufw allow 443/tcp    # HTTPS
sudo ufw reload

# Firewalld (RHEL/CentOS)
sudo firewall-cmd --permanent --add-port=54321/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --reload

# iptables
sudo iptables -I INPUT -p tcp --dport 54321 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 443 -j ACCEPT
sudo iptables-save > /etc/iptables/rules.v4
```

#### Windows 部署

```powershell
# 1. 创建安装目录
New-Item -Path "C:\Program Files\x-ui" -ItemType Directory -Force

# 2. 复制文件
Copy-Item x-ui.exe "C:\Program Files\x-ui\"
Copy-Item -Path bin\* -Destination "C:\Program Files\x-ui\bin\" -Recurse

# 3. 创建 Windows 服务（使用 NSSM）
nssm install x-ui "C:\Program Files\x-ui\x-ui.exe"
nssm set x-ui AppDirectory "C:\Program Files\x-ui"
nssm set x-ui DisplayName "3X-UI Service"
nssm set x-ui Description "3X-UI Web Panel for Xray"
nssm set x-ui Start SERVICE_AUTO_START

# 4. 启动服务
nssm start x-ui

# 5. 配置防火墙
New-NetFirewallRule -DisplayName "3X-UI Panel" -Direction Inbound -Protocol TCP -LocalPort 54321 -Action Allow
New-NetFirewallRule -DisplayName "3X-UI HTTPS" -Direction Inbound -Protocol TCP -LocalPort 443 -Action Allow
```

#### macOS 部署

```bash
# 1. 安装到系统目录
sudo mkdir -p /usr/local/x-ui/bin
sudo cp x-ui /usr/local/x-ui/
sudo cp -r bin/* /usr/local/x-ui/bin/
sudo chmod +x /usr/local/x-ui/x-ui
sudo chmod +x /usr/local/x-ui/bin/xray

# 2. 创建 LaunchDaemon（见上文）

# 3. 或者使用前台运行
cd /usr/local/x-ui
./x-ui
```

#### 访问面板

编译部署完成后，通过浏览器访问：

```
http://服务器IP:54321/
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
sudo lsof -i :54321
# 或
sudo netstat -tulpn | grep 54321

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

#### Method 1: Simple Build (Dynamic Linking)

```bash
# 1. Clone the project
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. Download Go module dependencies
go mod download

# 3. Build
CGO_ENABLED=1 go build -o x-ui main.go

# 4. Verify build
./x-ui -v
```

#### Method 2: Static Build (Recommended for Production)

```bash
# 1. Clone the project
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. Download dependencies
go mod download

# 3. Static build
CGO_ENABLED=1 go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o x-ui main.go

# 4. Verify static linking
ldd x-ui
# Should show "not a dynamic executable" or no output

# 5. Check file size and type
ls -lh x-ui
file x-ui
```

#### Method 3: Build with Docker (No Local Go Installation Required)

```bash
# Use official Go image for building
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  sh -c "go mod download && CGO_ENABLED=1 go build -o x-ui main.go"
```

#### Download Xray-core and Geodata

```bash
# Create bin directory
mkdir -p bin
cd bin

# Download Xray-core (choose based on architecture)
# amd64
wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-64.zip
unzip Xray-linux-64.zip
rm Xray-linux-64.zip

# arm64
# wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-arm64-v8a.zip
# unzip Xray-linux-arm64-v8a.zip

# Download geodata
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat

# Download region-specific geodata (optional)
wget -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat

cd ..
```

---

### Windows Compilation (EN)

#### Method 1: Build with PowerShell

```powershell
# 1. Clone the project
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. Set environment variables
$env:CGO_ENABLED="1"
$env:GOOS="windows"
$env:GOARCH="amd64"

# 3. Download dependencies
go mod download

# 4. Build
go build -ldflags "-w -s" -o x-ui.exe main.go

# 5. Verify
.\x-ui.exe -v
```

#### Download Xray-core and Geodata

```powershell
# Create bin directory
New-Item -Path bin -ItemType Directory -Force
cd bin

# Download Xray-core
Invoke-WebRequest -Uri "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-windows-64.zip" -OutFile "Xray-windows-64.zip"
Expand-Archive -Path "Xray-windows-64.zip" -DestinationPath .
Remove-Item "Xray-windows-64.zip"

# Download geodata
Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" -OutFile "geoip.dat"
Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" -OutFile "geosite.dat"

# Download region-specific data (optional)
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat" -OutFile "geoip_IR.dat"
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat" -OutFile "geosite_IR.dat"
Invoke-WebRequest -Uri "https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat" -OutFile "geoip_RU.dat"
Invoke-WebRequest -Uri "https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat" -OutFile "geosite_RU.dat"

cd ..
```

#### Create Windows Service (Optional)

```powershell
# Create Windows service using NSSM
# 1. Download NSSM: https://nssm.cc/download
# 2. Install service
nssm install x-ui "C:\path\to\x-ui.exe"
nssm set x-ui AppDirectory "C:\path\to\"
nssm start x-ui
```

---

### macOS Compilation (EN)

#### Build Steps

```bash
# 1. Clone the project
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 2. Set environment variables
export CGO_ENABLED=1
export GOOS=darwin
export GOARCH=amd64  # or arm64 for Apple Silicon

# 3. Download dependencies
go mod download

# 4. Build
go build -ldflags "-w -s" -o x-ui main.go

# 5. Grant execute permission
chmod +x x-ui

# 6. Verify
./x-ui -v
```

#### Download Xray-core and Geodata

```bash
# Create bin directory
mkdir -p bin
cd bin

# Download Xray-core
# Intel Mac (amd64)
wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-macos-64.zip
unzip Xray-macos-64.zip
rm Xray-macos-64.zip

# Apple Silicon (arm64)
# wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-macos-arm64-v8a.zip
# unzip Xray-macos-arm64-v8a.zip

# Download geodata
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat
wget -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat

cd ..
```

---

### Cross-Compilation (EN)

#### Cross-Compile on Linux to Multiple Architectures

```bash
#!/bin/bash
# build-all.sh - Cross-compilation script

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

    echo "Building ${os}/${arch}${arm:+v$arm}..."

    export CGO_ENABLED=1
    export GOOS=$os
    export GOARCH=$arch
    [ -n "$arm" ] && export GOARM=$arm

    # Set cross-compilation toolchain (for Linux static builds)
    if [ "$os" = "linux" ]; then
        go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o "$output" main.go
    else
        go build -ldflags "-w -s" -o "$output" main.go
    fi

    [ $? -eq 0 ] && echo "✓ $output build successful" || echo "✗ $output build failed"
done
```

---

### Deployment (EN)

#### Linux Deployment

##### Using Installation Script (Recommended)

```bash
# One-line installation
bash <(curl -Ls https://raw.githubusercontent.com/mhsanaei/3x-ui/main/install.sh)
```

##### Manual Deployment

```bash
# Create directory structure
sudo mkdir -p /usr/local/x-ui/bin

# Copy files
sudo cp x-ui /usr/local/x-ui/
sudo cp -r bin/* /usr/local/x-ui/bin/
sudo chmod +x /usr/local/x-ui/x-ui
sudo chmod +x /usr/local/x-ui/bin/xray

# Create systemd service
sudo tee /etc/systemd/system/x-ui.service > /dev/null <<EOF
[Unit]
Description=3x-ui Service
After=network.target nss-lookup.target

[Service]
User=root
WorkingDirectory=/usr/local/x-ui
ExecStart=/usr/local/x-ui/x-ui
Restart=on-failure
RestartPreventExitStatus=23
LimitNPROC=10000
LimitNOFILE=1000000

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and start service
sudo systemctl daemon-reload
sudo systemctl enable x-ui
sudo systemctl start x-ui

# Check service status
sudo systemctl status x-ui

# View logs
sudo journalctl -u x-ui -f
```

#### Access the Panel

After deployment, access the panel via browser:

```
http://server-ip:54321/
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
