# 3X-UI Linux 编译与部署指南 / Linux Build & Deployment Guide

本文件收纳了在 Linux 平台上构建与部署 3X-UI 的常见流程。若需要 macOS、Windows 或其他场景，请参阅主文档中的对应链接。

---

## 中文指南

### 编译方式

#### 方法一：简单编译（动态链接）

```bash
# 1. 克隆项目
git clone https://github.com/Jiusi-pys/3x-ui.git
cd 3x-ui

# 2. 下载 Go 模块依赖
go mod download

# 3. 编译
CGO_ENABLED=1 go build -o x-ui main.go

# 4. 验证编译结果
./x-ui -v
```

#### 方法二：静态编译（生产环境推荐）

```bash
git clone https://github.com/Jiusi-pys/3x-ui.git
cd 3x-ui

go mod download

CGO_ENABLED=1 go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o x-ui main.go
ldd x-ui            # 期望输出 “not a dynamic executable”
ls -lh x-ui
file x-ui
```

#### 方法三：使用 Docker 编译（无需本地安装 Go）

```bash
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  sh -c "go mod download && CGO_ENABLED=1 go build -o x-ui main.go"
```

### 准备 Xray Core 与地理数据库

```bash
mkdir -p bin
cd bin

# 下载 Xray-core （按架构选择）
wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-64.zip     # amd64
unzip Xray-linux-64.zip && rm Xray-linux-64.zip
# 标准化命名，便于面板识别
mv xray xray-linux-amd64

# ARM64（如需）
# wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-arm64-v8a.zip
# unzip Xray-linux-arm64-v8a.zip && rm Xray-linux-arm64-v8a.zip
# mv xray xray-linux-arm64

# 下载地理数据库
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat

# 选用地区规则（可选）
wget -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat

cd ..
```

### 部署流程

#### 一键安装脚本（推荐）

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Jiusi-pys/3x-ui/main/install.sh)
```

脚本会根据架构自动抓取 Release 包，安装面板与 Xray 内核，并配置 systemd/OpenRC 服务。

#### 手动部署

###### 步骤 1：安装依赖

```bash
# Debian/Ubuntu
sudo apt update
sudo apt install -y build-essential gcc g++ make git wget unzip

# RHEL/CentOS/Fedora
sudo dnf install -y gcc gcc-c++ make git wget unzip
```

确保 Go 版本 ≥ 1.25.1：`go version`

###### 步骤 2：编译与准备文件

```bash
git clone https://github.com/Jiusi-pys/3x-ui.git ~/3x-ui
cd ~/3x-ui
go mod download
CGO_ENABLED=1 go build -ldflags "-w -s" -o build/x-ui main.go
./DockerInit.sh amd64   # 根据 uname -m 调整为 arm64、arm32 等
```

###### 步骤 3：安装到系统路径

```bash
sudo install -d /usr/local/x-ui/bin
sudo install -d /etc/x-ui
sudo install -m755 build/x-ui /usr/local/x-ui/x-ui
sudo cp -a bin/* /usr/local/x-ui/bin/
sudo ln -sf /usr/local/x-ui/x-ui /usr/bin/x-ui
sudo install -m755 x-ui.sh /usr/local/x-ui/x-ui.sh   # 可选
```

###### 步骤 4：首次运行与服务配置

```bash
/usr/local/x-ui/x-ui -v
sudo /usr/local/x-ui/x-ui   # 前台运行确认无误后 Ctrl+C 退出

sudo install -m644 x-ui.service /etc/systemd/system/x-ui.service
sudo systemctl daemon-reload
sudo systemctl enable --now x-ui
```

###### 常用维护命令

```bash
sudo systemctl restart x-ui
sudo systemctl status x-ui
sudo journalctl -u x-ui -f

###### 访问面板

在浏览器访问 `http://<主机IP>:2053/`，订阅服务默认端口为 `2096/tcp`。首次登录请立即修改默认凭据。
```

###### 防火墙设置

```bash
# UFW
sudo ufw allow 2053/tcp
sudo ufw allow 2096/tcp
sudo ufw reload

# Firewalld
sudo firewall-cmd --permanent --add-port=2053/tcp
sudo firewall-cmd --permanent --add-port=2096/tcp
sudo firewall-cmd --reload

# iptables
sudo iptables -I INPUT -p tcp --dport 2053 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 2096 -j ACCEPT
sudo iptables-save | sudo tee /etc/iptables/rules.v4

---

### 卸载 / Uninstall

- 使用 systemd 安装时：

```bash
sudo systemctl stop x-ui
sudo systemctl disable x-ui
sudo rm -f /etc/systemd/system/x-ui.service
sudo systemctl daemon-reload

# 清理程序与配置
sudo rm -rf /usr/local/x-ui /etc/x-ui /usr/bin/x-ui

# 可选：移除防火墙开放规则
# UFW
sudo ufw delete allow 2053/tcp || true
sudo ufw delete allow 2096/tcp || true
# Firewalld
sudo firewall-cmd --permanent --remove-port=2053/tcp || true
sudo firewall-cmd --permanent --remove-port=2096/tcp || true
sudo firewall-cmd --reload || true
```

- 若使用 OpenRC（Alpine 等）：

```bash
sudo rc-service x-ui stop || true
sudo rc-update del x-ui || true
sudo rm -f /etc/init.d/x-ui
sudo rm -rf /usr/local/x-ui /etc/x-ui
```
```

---

## English Guide

### Build Options

#### Method 1: Simple build (dynamic)
```bash
git clone https://github.com/Jiusi-pys/3x-ui.git
cd 3x-ui
go mod download
CGO_ENABLED=1 go build -o x-ui main.go
./x-ui -v
```

#### Method 2: Static build (recommended for production)
```bash
git clone https://github.com/Jiusi-pys/3x-ui.git
cd 3x-ui
go mod download
CGO_ENABLED=1 go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o x-ui main.go
ldd x-ui    # expect “not a dynamic executable”
```

#### Method 3: Build inside Docker
```bash
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  sh -c "go mod download && CGO_ENABLED=1 go build -o x-ui main.go"
```

### Fetch Xray Core & rule data
Same steps as above; ensure the binary is renamed to `xray-linux-<arch>` to match runtime expectations.

### Deployment

- **Script (recommended)**  
  `bash <(curl -Ls https://raw.githubusercontent.com/Jiusi-pys/3x-ui/main/install.sh)`

- **Manual install**
  1. Install prerequisites (`build-essential`, `gcc`, `git`, `wget`, `unzip`).
  2. Build into `build/x-ui` and populate `build/bin/` via `./DockerInit.sh <arch>`.
  3. Copy artifacts into `/usr/local/x-ui`, symlink `/usr/bin/x-ui`.
  4. Register systemd service: `sudo install -m644 x-ui.service /etc/systemd/system/x-ui.service`, then `systemctl enable --now x-ui`.

- **Firewall**  
  Open TCP ports `2053` (panel) and `2096` (subscription) via UFW/Firewalld/iptables as needed.

- **Uninstall**  
  Stop and disable the `x-ui` service, remove `/etc/systemd/system/x-ui.service`, and delete `/usr/local/x-ui` and `/etc/x-ui`. Revoke firewall rules if previously added.

The panel becomes available at `http://<host>:2053/` (subscription on `2096/tcp`). Change default credentials immediately after first login.
