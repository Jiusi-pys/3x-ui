# 3X-UI macOS 部署指南 / macOS Deployment Guide

该文档汇总了在 macOS (Intel & Apple Silicon) 上构建、安装及以 LaunchDaemon 方式运行 3X-UI 的完整流程。

---

## 前置条件 / Prerequisites

- macOS 12+，支持 Intel x86_64 或 Apple Silicon。
- 安装 Xcode Command Line Tools：`xcode-select --install`
- Go 1.25.1+：`brew install go@1.25`，并确认 `go version`
- 常用工具：`brew install git wget unzip`（macOS 默认不含 wget）

## 构建步骤 / Build Steps

```bash
git clone https://github.com/Jiusi-pys/3x-ui.git ~/3x-ui
cd ~/3x-ui
go mod download
CGO_ENABLED=1 go build -ldflags "-w -s" -o build/x-ui main.go
```

### 下载 Xray Core 与规则文件

```bash
mkdir -p build/bin
cd build/bin

# Intel Mac
wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-macos-64.zip
unzip Xray-macos-64.zip && rm Xray-macos-64.zip

# Apple Silicon (如需)
# wget https://github.com/XTLS/Xray-core/releases/latest/download/Xray-macos-arm64-v8a.zip
# unzip Xray-macos-arm64-v8a.zip && rm Xray-macos-arm64-v8a.zip

wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat
wget -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat

cd ../..
```

## 安装与首次运行 / Installation & First Run

```bash
sudo install -d /usr/local/x-ui/bin
sudo install -d /etc/x-ui
sudo install -m755 build/x-ui /usr/local/x-ui/x-ui
sudo cp -a build/bin/* /usr/local/x-ui/bin/
sudo install -m755 x-ui.sh /usr/local/x-ui/x-ui.sh   # 可选 CLI

/usr/local/x-ui/x-ui -v
sudo /usr/local/x-ui/x-ui   # 前台试运行后 Ctrl+C
```

> 数据库默认位于 `/etc/x-ui/3x-ui.db`，日志位于 `/var/log`，可通过 `XUI_DB_FOLDER` / `XUI_LOG_FOLDER` 环境变量调整。

## 配置 LaunchDaemon / Configure LaunchDaemon

1. 可选：准备日志文件夹
   ```bash
   sudo mkdir -p /var/log/x-ui
   sudo touch /var/log/x-ui/x-ui.out /var/log/x-ui/x-ui.err
   sudo chmod 644 /var/log/x-ui/x-ui.*
   ```

2. 创建 `/Library/LaunchDaemons/com.x-ui.plist`
   ```bash
   sudo tee /Library/LaunchDaemons/com.x-ui.plist > /dev/null <<'EOF'
   <?xml version="1.0" encoding="UTF-8"?>
   <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
   <plist version="1.0">
   <dict>
       <key>Label</key>
       <string>com.x-ui</string>
       <key>ProgramArguments</key>
       <array>
           <string>/usr/local/x-ui/x-ui</string>
       </array>
       <key>WorkingDirectory</key>
       <string>/usr/local/x-ui</string>
       <key>RunAtLoad</key>
       <true/>
       <key>KeepAlive</key>
       <true/>
       <key>StandardOutPath</key>
       <string>/var/log/x-ui/x-ui.out</string>
       <key>StandardErrorPath</key>
       <string>/var/log/x-ui/x-ui.err</string>
   </dict>
   </plist>
   EOF
   sudo chown root:wheel /Library/LaunchDaemons/com.x-ui.plist
   sudo chmod 644 /Library/LaunchDaemons/com.x-ui.plist
   ```

3. 启动并常驻
   ```bash
   sudo launchctl bootstrap system /Library/LaunchDaemons/com.x-ui.plist
   sudo launchctl enable system/com.x-ui
   sudo launchctl kickstart -k system/com.x-ui
   ```

4. 调试与日志
   ```bash
   sudo launchctl print system/com.x-ui
   log show --style syslog --last 5m --predicate 'process == "x-ui" || process == "launchd"'
   ```

### 常见问题：Bootstrap failed: 5

- `plutil -lint /Library/LaunchDaemons/com.x-ui.plist` 检查语法。
- 确认 plist 权限 `root:wheel / 644`，以及 `/usr/local/x-ui/x-ui` 可执行。
- 若之前加载过旧配置，执行 `sudo launchctl bootout system/com.x-ui` 后重试。
- 结合 `log show` 输出分析具体报错。

## 运行面板 / Access the Panel

访问 `http://<mac-ip>:2053/`。默认凭据为 `admin/admin`（若脚本生成随机值，请查看安装输出），请在首次登录后立即修改。

若需要订阅服务，请放行 `2096/tcp`，示例：`sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /usr/local/x-ui/x-ui` 并确保对应端口开放。

## 常用维护命令

```bash
sudo launchctl kickstart -k system/com.x-ui   # 重启
sudo launchctl bootout system/com.x-ui        # 停止并卸载
sudo launchctl print system/com.x-ui          # 查看状态
```

临时前台运行：`sudo /usr/local/x-ui/x-ui`

## 故障排查

### fork/exec bin/xray-darwin-amd64: no such file or directory

Xray 可执行文件需要命名为 `xray-darwin-<arch>` 才能被面板识别。解压后执行：

```bash
cd /usr/local/x-ui/bin
sudo mv Xray xray-darwin-amd64        # Intel
# sudo mv Xray xray-darwin-arm64      # Apple Silicon
sudo chmod 755 xray-darwin-*
```
