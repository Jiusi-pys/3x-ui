# 3X-UI 项目完整分析文档

## 项目概述

**3X-UI** 是一个基于 Xray-core 的高级开源 Web 管理面板，为各种 VPN 和代理协议提供用户友好的配置和监控界面。作为原始 X-UI 项目的增强分支，3X-UI 提供了改进的稳定性、更广泛的协议支持和额外功能。

### 核心信息
- **语言**: Go (后端) + HTML/JavaScript (前端)
- **框架**: Gin (Web框架), GORM (ORM), Vue.js (前端)
- **数据库**: SQLite
- **协议支持**: VMess, VLESS, Trojan, Shadowsocks, HTTP, Mixed, WireGuard

## 项目架构分析

### 1. 目录结构
```
3x-ui/
├── main.go                 # 应用程序入口点
├── go.mod/go.sum           # Go 模块依赖
├── config/                 # 配置文件目录
├── database/               # 数据库相关
│   └── model/             # 数据模型定义
├── web/                   # Web 服务器实现
│   ├── controller/        # 控制器层
│   ├── service/          # 业务逻辑层
│   ├── middleware/       # 中间件
│   ├── html/            # HTML 模板
│   ├── assets/          # 静态资源
│   ├── locale/          # 国际化支持
│   └── job/             # 后台任务
├── util/                 # 工具函数
├── logger/              # 日志系统
├── sub/                 # 订阅服务
├── xray/                # Xray 集成
└── windows_files/       # Windows 特定文件
```

### 2. 核心模块功能分析

#### 2.1 主程序入口 (main.go)
- **功能**: 应用程序启动点，处理命令行参数
- **关键职责**:
  - 初始化数据库连接
  - 启动 Web 服务器和订阅服务器
  - 处理系统信号（SIGHUP 重启、SIGTERM 关闭）
  - 提供命令行工具功能（设置修改、数据库迁移）

#### 2.2 Web 服务器模块 (web/)
- **web/web.go**: 主要 Web 服务器实现
  - HTTP/HTTPS 服务配置
  - 路由初始化和中间件注册
  - 模板引擎配置
  - 静态资源服务
  - 后台任务调度

- **控制器层 (web/controller/)**:
  - `api.go`: API 路由控制器，处理认证和路由分发
  - `inbound.go`: 入站规则管理 (CRUD 操作)
  - `index.go`: 主页和系统状态
  - `setting.go`: 系统设置管理
  - `server.go`: 服务器信息和操作
  - `xray_setting.go`: Xray 配置管理
  - `xui.go`: 主面板控制器

- **服务层 (web/service/)**:
  - `inbound.go`: 入站规则业务逻辑
  - `outbound.go`: 出站规则业务逻辑
  - `user.go`: 用户管理
  - `setting.go`: 系统设置
  - `xray.go`: Xray 核心集成
  - `tgbot.go`: Telegram 机器人集成

#### 2.3 数据库模块 (database/)
- **模型定义 (database/model/model.go)**:
  - `User`: 用户账户
  - `Inbound`: 入站规则配置
  - `OutboundTraffics`: 出站流量统计
  - `InboundClientIps`: 客户端 IP 管理
  - `Setting`: 系统配置
  - `Client`: 客户端配置

#### 2.4 后台任务 (web/job/)
- `xray_traffic_job.go`: 流量统计
- `stats_notify_job.go`: 统计通知
- `check_xray_running_job.go`: Xray 运行状态检查
- `periodic_traffic_reset_job.go`: 定期流量重置
- `ldap_sync_job.go`: LDAP 同步
- `clear_logs_job.go`: 日志清理

#### 2.5 订阅服务 (sub/)
- 提供客户端订阅链接生成
- 支持多种客户端格式 (V2Ray, Clash, etc.)

#### 2.6 前端界面 (web/html/)
- 基于 Vue.js + Ant Design Vue
- 主要页面:
  - `index.html`: 系统监控面板
  - `inbounds.html`: 入站规则管理
  - `settings.html`: 系统设置
  - `xray.html`: Xray 配置
  - `login.html`: 登录页面

## 功能特性详解

### 1. 入站规则管理
- **协议支持**: VMess, VLESS, Trojan, Shadowsocks, HTTP, Mixed, WireGuard
- **功能**:
  - 添加/编辑/删除入站规则
  - 客户端管理 (UUID, 密码, 流量限制)
  - 流量统计和监控
  - 到期时间管理
  - IP 限制
  - 定期流量重置

### 2. 出站规则管理
- 配置上游代理
- 负载均衡
- 故障转移
- 流量统计

### 3. 系统监控
- CPU 使用率
- 内存使用率
- 磁盘使用情况
- 网络流量统计
- Xray 运行状态

### 4. 用户管理
- 多用户支持
- 权限控制
- 登录会话管理
- 双因素认证 (2FA)

### 5. 安全特性
- SSL/TLS 证书配置
- 基于域名的访问控制
- IP 白名单/黑名单
- 失败登录保护

### 6. 通知系统
- Telegram 机器人集成
- 流量告警
- 系统状态通知
- 备份推送

### 7. 国际化支持
- 多语言界面 (英语、中文、俄语、阿拉伯语等)
- 本地化时间格式

## 部署方案

### 1. 支持平台
- **Linux**: Ubuntu, Debian, CentOS, RHEL, Fedora
- **Windows**: Windows Server, Windows 10/11
- **macOS**: Intel 和 Apple Silicon
- **架构**: amd64, arm64, armv7, armv6, armv5, 386, s390x

### 2. 快速安装 (Linux)
```bash
bash <(curl -Ls https://raw.githubusercontent.com/mhsanaei/3x-ui/master/install.sh)
```

### 3. Docker 部署
```bash
# 构建镜像
docker build -t 3x-ui:latest .

# 运行容器
docker run -d \
  --name 3x-ui \
  -p 2053:2053 \
  -v /etc/x-ui:/etc/x-ui \
  3x-ui:latest
```

### 4. Docker Compose 部署
```yaml
version: '3.8'
services:
  3x-ui:
    build: .
    container_name: 3x-ui
    ports:
      - "2053:2053"
    volumes:
      - ./db/:/etc/x-ui/
    restart: unless-stopped
```

### 5. 手动编译部署

#### 5.1 环境要求
- Go 1.25.1 或更高版本
- CGO 支持 (SQLite 依赖)
- Git

#### 5.2 编译步骤

##### Linux/macOS 编译
```bash
# 克隆项目
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 设置环境变量
export CGO_ENABLED=1
export CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# 编译
go build -ldflags "-w -s" -o x-ui main.go

# 安装
sudo mv x-ui /usr/local/bin/
sudo chmod +x /usr/local/bin/x-ui
```

##### Windows 编译
```cmd
# 设置环境变量
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

# 编译
go build -ldflags "-w -s" -o x-ui.exe main.go
```

##### 交叉编译示例
```bash
# 编译 ARM64 Linux 版本
env GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
  CC=aarch64-linux-gnu-gcc \
  go build -ldflags "-w -s" -o x-ui-linux-arm64 main.go

# 编译 Windows 版本
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-w64-mingw32-gcc \
  go build -ldflags "-w -s" -o x-ui-windows-amd64.exe main.go

# 编译 macOS 版本
env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
  go build -ldflags "-w -s" -o x-ui-darwin-amd64 main.go
```

### 6. 系统服务配置

#### Linux (systemd)
```ini
# /etc/systemd/system/x-ui.service
[Unit]
Description=3X-UI Panel
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/x-ui run

[Install]
WantedBy=multi-user.target
```

启用服务:
```bash
sudo systemctl daemon-reload
sudo systemctl enable x-ui
sudo systemctl start x-ui
```

#### Windows (服务)
使用 `x-ui.service` 文件中的配置创建 Windows 服务。

### 7. 配置文件位置
- **Linux**: `/etc/x-ui/`
- **Windows**: `C:\x-ui\`
- **macOS**: `/usr/local/etc/x-ui/`

### 8. 默认配置
- **端口**: 2053
- **用户名**: admin
- **密码**: admin
- **数据库**: SQLite (`x-ui.db`)

## 安全建议

1. **更改默认凭据**: 首次安装后立即更改默认用户名和密码
2. **启用 SSL**: 配置 SSL 证书以启用 HTTPS
3. **防火墙配置**: 仅开放必要端口
4. **定期更新**: 保持软件版本最新
5. **访问控制**: 使用域名验证和 IP 白名单
6. **监控日志**: 定期检查访问日志和错误日志

## 维护操作

### 命令行工具
```bash
# 显示版本
x-ui -v

# 重置设置
x-ui setting -reset

# 显示当前设置
x-ui setting -show

# 修改端口
x-ui setting -port 2054

# 修改用户名密码
x-ui setting -username newuser -password newpass

# 重置双因素认证
x-ui setting -resetTwoFactor

# 数据库迁移
x-ui migrate
```

### 备份和恢复
```bash
# 备份数据库
cp /etc/x-ui/x-ui.db /backup/x-ui-$(date +%Y%m%d).db

# 恢复数据库
cp /backup/x-ui-20231201.db /etc/x-ui/x-ui.db
systemctl restart x-ui
```

## 开发和贡献

### 开发环境设置
```bash
# 克隆项目
git clone https://github.com/MHSanaei/3x-ui.git
cd 3x-ui

# 安装依赖
go mod download

# 运行开发服务器
go run main.go
```

### 构建和测试
```bash
# 运行测试
go test ./...

# 格式化代码
go fmt ./...

# 静态分析
go vet ./...
```

## 总结

3X-UI 是一个功能全面的 Xray 管理面板，提供了:
- 直观的 Web 界面
- 全面的协议支持
- 强大的监控和统计功能
- 灵活的部署选项
- 良好的安全特性
- 多平台支持

该项目适合个人使用，为 VPN 和代理服务的管理提供了便捷的解决方案。通过其模块化的架构和清晰的代码结构，开发者可以轻松扩展和定制功能。

---

**注意**: 此项目仅供个人使用，请勿用于非法目的，请勿在生产环境中使用。