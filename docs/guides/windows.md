# 3X-UI Windows 编译与部署指南 / Windows Build & Deployment Guide

---

## 中文指南

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
.\x-ui.exe -v
```

### 准备 Xray Core 与规则文件

```powershell
New-Item -Path bin -ItemType Directory -Force
cd bin

Invoke-WebRequest -Uri "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-windows-64.zip" -OutFile "Xray-windows-64.zip"
Expand-Archive -Path "Xray-windows-64.zip" -DestinationPath .
Remove-Item "Xray-windows-64.zip"

Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" -OutFile "geoip.dat"
Invoke-WebRequest -Uri "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" -OutFile "geosite.dat"

# 可选：区域规则
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat" -OutFile "geoip_IR.dat"
Invoke-WebRequest -Uri "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat" -OutFile "geosite_IR.dat"
```

### 部署与服务配置

```powershell
# 1. 安装目录
New-Item -Path "C:\Program Files\x-ui" -ItemType Directory -Force

# 2. 拷贝文件
Copy-Item build\x-ui.exe "C:\Program Files\x-ui\"
Copy-Item -Path build\bin\* -Destination "C:\Program Files\x-ui\bin\" -Recurse

# 3. 使用 NSSM 创建服务
nssm install x-ui "C:\Program Files\x-ui\x-ui.exe"
nssm set x-ui AppDirectory "C:\Program Files\x-ui"
nssm set x-ui DisplayName "3X-UI Service"
nssm set x-ui Description "3X-UI Web Panel for Xray"
nssm set x-ui Start SERVICE_AUTO_START
nssm start x-ui

# 4. 打开防火墙
New-NetFirewallRule -DisplayName "3X-UI Panel" -Direction Inbound -Protocol TCP -LocalPort 2053 -Action Allow
New-NetFirewallRule -DisplayName "3X-UI Subscription" -Direction Inbound -Protocol TCP -LocalPort 2096 -Action Allow
```

首次运行后访问 `http://<主机IP>:2053/`，默认账号密码会在首次安装时随机生成（或保持 `admin/admin`，请立即修改）。

---

## English Guide

### Build

```powershell
git clone https://github.com/Jiusi-pys/3x-ui.git
cd 3x-ui
$env:CGO_ENABLED="1"
$env:GOOS="windows"
$env:GOARCH="amd64"
go mod download
go build -ldflags "-w -s" -o x-ui.exe main.go
.\x-ui.exe -v
```

### Prepare Xray core and rule files

Same as the Chinese section: create `bin\`, download the latest `Xray-windows-64.zip`, extract it, and fetch `geoip.dat` / `geosite.dat` plus optional regional datasets.

### Installation & service

1. Copy binaries into `C:\Program Files\x-ui\` with the `bin\` subfolder.
2. Use [NSSM](https://nssm.cc/) to register `x-ui.exe` as a Windows service (`nssm install x-ui ...`).
3. Open TCP ports `2053` and `2096` in Windows Firewall (`New-NetFirewallRule ...`).
4. Browse to `http://<host>:2053/` to complete setup.
