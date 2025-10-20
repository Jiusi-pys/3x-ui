# 3X-UI 其他 Linux 场景指南 / Other Linux Scenarios

适用于特殊发行版（如 Alpine/OpenRC）、批量交叉编译或基于容器的构建部署。

---

## 中文指南

### 先决条件 / 工具安装

- 基本构建工具：Git、GCC/Make、wget、unzip、curl、tzdata
- 发行版示例安装命令：

```bash
# Alpine
sudo apk add --update go build-base git wget curl unzip tzdata

# Debian/Ubuntu
sudo apt update
sudo apt install -y build-essential gcc g++ make git wget curl unzip tzdata

# RHEL/CentOS/Fedora
sudo dnf install -y gcc gcc-c++ make git wget curl unzip tzdata

# 验证
git --version
gcc --version
go version
```

### 交叉编译

```bash
#!/bin/bash
platforms=(
  "linux/amd64"
  "linux/arm64"
  "linux/arm/7"
  "linux/arm/6"
  "linux/arm/5"
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

  export CGO_ENABLED=1
  export GOOS=$os
  export GOARCH=$arch
  [ -n "$arm" ] && export GOARM=$arm

  if [ "$os" = "linux" ]; then
    go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o "$output" main.go
  else
    go build -ldflags "-w -s" -o "$output" main.go
  fi
done
```

生成的产物可配合 `DockerInit.sh <arch>` 获取对应架构的 Xray 内核后再打包分发。

### Docker Buildx 多平台编译

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --tag jiusi-pys/3x-ui:latest \
  --output type=tar,dest=3x-ui-build.tar .
```

解压 `3x-ui-build.tar` 后将 `x-ui` 与 `bin/` 内容部署至目标主机即可。

### Alpine / OpenRC 提示

- 安装依赖：`apk add --update go build-base wget curl tzdata`
- 服务脚本可放置于 `/etc/init.d/x-ui` 并执行 `rc-update add x-ui`.
- 若使用一键脚本，会自动拉取仓库中的 `x-ui.rc` 并完成注册。

#### 卸载 / Uninstall（OpenRC）

```bash
sudo rc-service x-ui stop || true
sudo rc-update del x-ui || true
sudo rm -f /etc/init.d/x-ui
sudo rm -rf /usr/local/x-ui /etc/x-ui
```

---

## English Guide

### Prerequisites

- Build essentials: Git, GCC/Make, wget, unzip, curl, tzdata
- Examples per distro:

```bash
# Alpine
apk add --update go build-base git wget curl unzip tzdata

# Debian/Ubuntu
sudo apt update
sudo apt install -y build-essential gcc g++ make git wget curl unzip tzdata

# RHEL/CentOS/Fedora
sudo dnf install -y gcc gcc-c++ make git wget curl unzip tzdata

# Verify
git --version
gcc --version
go version
```

### Cross-compilation script

Use the same script as above to iterate over GOOS/GOARCH combinations. Static linking (`-extldflags '-static'`) is recommended for Linux targets.

### Docker Buildx

```
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --tag jiusi-pys/3x-ui:latest \
  .
```

Combine the resulting artifacts with architecture-specific Xray binaries (run `./DockerInit.sh <arch>` inside the container or locally).

### Alpine/OpenRC notes

- Packages: `apk add go build-base wget curl tzdata`
- Service: install `x-ui.rc` to `/etc/init.d/x-ui`, then `rc-update add x-ui && rc-service x-ui start`.
- Firewall: use `iptables` or `nftables` to allow TCP 2053/2096 as needed.

#### Uninstall (OpenRC)

```bash
sudo rc-service x-ui stop || true
sudo rc-update del x-ui || true
sudo rm -f /etc/init.d/x-ui
sudo rm -rf /usr/local/x-ui /etc/x-ui
```
