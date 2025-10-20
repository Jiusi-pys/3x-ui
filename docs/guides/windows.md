# 3X-UI Windows 编译与部署指南 / Windows Build & Deployment Guide

---

## 中文指南

### 先决条件 / 工具安装

- Git（用于克隆仓库）
- GCC（用于 CGO，推荐 MSYS2 提供的 mingw-w64 工具链）
- NSSM（将面板注册为 Windows 服务）

使用 winget 快速安装（管理员 PowerShell）：

```powershell
# 安装 Git / MSYS2 / NSSM
winget install --id Git.Git -e
winget install --id MSYS2.MSYS2 -e
winget install --id NSSM.NSSM -e

# 在 MSYS2 里安装 GCC（打开“MSYS2 UCRT64”或“MSYS2 MINGW64”终端）
# 选择其一：
# UCRT64（推荐）：
pacman -S --needed --noconfirm mingw-w64-ucrt-x86_64-gcc
# 或 MINGW64：
pacman -S --needed --noconfirm mingw-w64-x86_64-gcc

# 将对应二进制目录加入 PATH（按你实际安装的环境选择其一）
# UCRT64：
setx PATH "$env:PATH;C:\\msys64\\ucrt64\\bin"
# 或 MINGW64：
setx PATH "$env:PATH;C:\\msys64\\mingw64\\bin"

# 验证
git --version
gcc --version
nssm --version
```

### 编译步骤

```powershell
# 1. 克隆项目
git clone https://github.com/Jiusi-pys/3x-ui.git
cd 3x-ui

# 2. 设置环境变量
$env:CGO_ENABLED="1"
$env:GOOS="windows"
$env:GOARCH="amd64"    # x86_64 平台，如需 32 位可改为 386

# 3. 下载依赖并编译
go mod download
go build -ldflags "-w -s" -o x-ui.exe main.go

# 4. 验证
./x-ui.exe -v
```

### 准备 Xray Core 与规则文件

```powershell
# 在项目根目录下执行
New-Item -Path bin -ItemType Directory -Force
cd bin

Invoke-WebRequest -Uri "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-windows-64.zip" -OutFile "Xray-windows-64.zip"
Expand-Archive -Path "Xray-windows-64.zip" -DestinationPath .
Remove-Item "Xray-windows-64.zip"; cd ..

# 标准化命名，便于面板识别（兼容不同解压结构）
# 若 xray.exe 在当前目录则重命名；若在子目录（Xray-windows-64）则移动并重命名
Rename-Item -Path ".\bin\xray.exe" -NewName "xray-windows-amd64.exe" -ErrorAction SilentlyContinue
Move-Item -Path ".\bin\Xray-windows-64\xray.exe" -Destination ".\bin\xray-windows-amd64.exe" -ErrorAction SilentlyContinue

Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" -OutFile ".\\bin\\geoip.dat"
Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" -OutFile ".\\bin\\geosite.dat"

# 可选：区域规则（伊朗）
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat" -OutFile ".\\bin\\geoip_IR.dat"
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat" -OutFile ".\\bin\\geosite_IR.dat"

# 返回项目根目录
cd ..
```

### 部署与服务配置

```powershell
# 1. 安装目录
New-Item -Path "C:\Program Files\x-ui" -ItemType Directory -Force
New-Item -Path "C:\Program Files\x-ui\bin" -ItemType Directory -Force

# 2. 拷贝文件（从项目根目录）
Copy-Item .\x-ui.exe "C:\Program Files\x-ui\"
Copy-Item -Path .\bin\* -Destination "C:\Program Files\x-ui\bin\" -Recurse -Force

# 3. 使用 NSSM 创建服务（需提前安装 NSSM: https://nssm.cc/）
nssm install x-ui "C:\Program Files\x-ui\x-ui.exe"
nssm set x-ui AppDirectory "C:\Program Files\x-ui"
nssm set x-ui DisplayName "3X-UI Service"
nssm set x-ui Description "3X-UI Web Panel for Xray"
nssm set x-ui Start SERVICE_AUTO_START
nssm start x-ui

# 4. 打开防火墙（如有自定义端口，请同步调整）
New-NetFirewallRule -DisplayName "3X-UI Panel" -Direction Inbound -Protocol TCP -LocalPort 2053 -Action Allow
New-NetFirewallRule -DisplayName "3X-UI Subscription" -Direction Inbound -Protocol TCP -LocalPort 2096 -Action Allow
```

首次运行后访问 `http://<主机IP>:2053/`，默认账号密码会在首次安装时随机生成（或保持 `admin/admin`，请立即修改）。

### 卸载 / Uninstall

```powershell
# 停止并移除服务
nssm stop x-ui
nssm remove x-ui confirm

# 删除程序与配置
Remove-Item -Path "C:\Program Files\x-ui" -Recurse -Force

# 可选：移除防火墙规则（如按本指南创建）
Get-NetFirewallRule -DisplayName "3X-UI Panel" -ErrorAction SilentlyContinue | Remove-NetFirewallRule
Get-NetFirewallRule -DisplayName "3X-UI Subscription" -ErrorAction SilentlyContinue | Remove-NetFirewallRule
```

---

## English Guide

### Prerequisites

- Git for Windows
- GCC for CGO (via MSYS2 mingw-w64 toolchain)
- NSSM (Windows service helper)

Quick install with winget (Admin PowerShell):

```powershell
winget install --id Git.Git -e
winget install --id MSYS2.MSYS2 -e
winget install --id NSSM.NSSM -e

# In MSYS2 (UCRT64 or MINGW64 shell), install GCC:
pacman -S --needed --noconfirm mingw-w64-ucrt-x86_64-gcc    # UCRT64
# or
pacman -S --needed --noconfirm mingw-w64-x86_64-gcc          # MINGW64

# Add to PATH (pick the one you actually use):
setx PATH "$env:PATH;C:\\msys64\\ucrt64\\bin"
# or
setx PATH "$env:PATH;C:\\msys64\\mingw64\\bin"

# Verify
git --version
gcc --version
nssm --version
```

### Build

```powershell
git clone https://github.com/Jiusi-pys/3x-ui.git
cd 3x-ui
$env:CGO_ENABLED="1"
$env:GOOS="windows"
$env:GOARCH="amd64"
go mod download
go build -ldflags "-w -s" -o x-ui.exe main.go
./x-ui.exe -v
```

### Prepare Xray core and rule files

1) From the repo root, create `bin\`, download `Xray-windows-64.zip`, extract it, and fetch `geoip.dat` / `geosite.dat` (plus optional regional datasets).
2) Normalize the Xray binary name so the panel can detect it:

```powershell
Rename-Item -Path ".\bin\xray.exe" -NewName "xray-windows-amd64.exe" -ErrorAction SilentlyContinue
Move-Item -Path ".\bin\Xray-windows-64\xray.exe" -Destination ".\bin\xray-windows-amd64.exe" -ErrorAction SilentlyContinue
```

### Installation & service

1. Copy binaries into `C:\Program Files\x-ui\` with the `bin\` subfolder.
2. Use NSSM to register `x-ui.exe` as a Windows service (`nssm install x-ui ...`).
3. Open TCP ports `2053` and `2096` in Windows Firewall (`New-NetFirewallRule ...`).
4. Browse to `http://<host>:2053/` to complete setup.

### Uninstall

```powershell
nssm stop x-ui
nssm remove x-ui confirm
Remove-Item -Path "C:\Program Files\x-ui" -Recurse -Force
Get-NetFirewallRule -DisplayName "3X-UI Panel" -ErrorAction SilentlyContinue | Remove-NetFirewallRule
Get-NetFirewallRule -DisplayName "3X-UI Subscription" -ErrorAction SilentlyContinue | Remove-NetFirewallRule
```
