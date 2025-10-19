<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./media/3x-ui-dark.png">
    <img alt="3x-ui" src="./media/3x-ui-light.png">
  </picture>
</p>

[![Release](https://img.shields.io/github/v/release/Jiusi-pys/3x-ui.svg)](https://github.com/Jiusi-pys/3x-ui/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/Jiusi-pys/3x-ui/release.yml.svg)](https://github.com/Jiusi-pys/3x-ui/actions)
[![GO Version](https://img.shields.io/github/go-mod/go-version/Jiusi-pys/3x-ui.svg)](#)
[![Downloads](https://img.shields.io/github/downloads/Jiusi-pys/3x-ui/total.svg)](https://github.com/Jiusi-pys/3x-ui/releases/latest)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/mhsanaei/3x-ui/v2.svg)](https://pkg.go.dev/github.com/mhsanaei/3x-ui/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhsanaei/3x-ui/v2)](https://goreportcard.com/report/github.com/mhsanaei/3x-ui/v2)

**3X-UI** —— 一个基于 Web 的高阶开源控制面板，用于管理 Xray-core。此仓库在官方版本基础上增加了多项功能、修复并提供新的安装脚本。

> [!IMPORTANT]
> 本项目仅限个人合法用途，请勿用于任何违法行为或生产环境。

## 一键安装

- 运行脚本会从当前仓库下载发布包、部署面板与 Xray 内核，并生成随机的初始账号密码。
- 默认监听 `2053` 端口（面板）与 `2096` 端口（订阅）。

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Jiusi-pys/3x-ui/main/install.sh)
```

更多平台及步骤详见仓库内的《[3X-UI 编译部署指南](./COMPILATION_GUIDE.md)》，以及分平台文档：
- [Linux 指南](./docs/guides/linux.md)
- [其他 Linux 场景](./docs/guides/linux_other.md)
- [Windows 指南](./docs/guides/windows.md)
- [macOS 指南](./docs/guides/macos.md)

## 特性速览

- 全新出站规则管理模块：数据库配置自动合并至运行时模板，并在修改时自动提示重启 Xray。
- 扩展的系统服务脚本：更新 `install.sh` 及 `x-ui.sh` 支持自定义构建产物及仓库。
- 丰富的模板编辑体验：提供模板、有效配置、入出站、路由规则多种视图，便于排错。
- 兼容 Fail2ban、Telegram Bot 通知、订阅 API 等增强功能。

## 发布与下载

- [最新 Release](https://github.com/Jiusi-pys/3x-ui/releases)：提供各架构压缩包（`x-ui-linux-<arch>.tar.gz`），供脚本或手动部署使用。
- 若需自定义编译，请参考 `build/` 目录输出结构，与安装脚本保持一致。

## 文档与支持

- [COMPILATION_GUIDE.md](./COMPILATION_GUIDE.md)：覆盖 Linux/Windows/macOS 编译与手动部署步骤、常见问题。
- [docs/guides/linux.md](./docs/guides/linux.md)：Linux 平台编译、脚本与手动部署。
- [docs/guides/linux_other.md](./docs/guides/linux_other.md)：交叉编译、Docker Buildx、Alpine/OpenRC 等特殊场景。
- [docs/guides/windows.md](./docs/guides/windows.md)：Windows 平台编译、服务安装与防火墙配置。
- [docs/guides/macos.md](./docs/guides/macos.md)：macOS 编译、LaunchDaemon 配置与常见问题排查。
- 面板内置操作指南，可在“系统设置”与“Xray 设置”页面查看实时状态及日志。

## 特别感谢

- [alireza0](https://github.com/alireza0/)

## 致谢

- [Iran v2ray rules](https://github.com/chocolate4u/Iran-v2ray-rules) (许可证: **GPL-3.0**): _增强的 v2ray/xray 和 v2ray/xray-clients 路由规则，内置伊朗域名并专注于安全与广告拦截。_
- [Russia v2ray rules](https://github.com/runetfreedom/russia-v2ray-rules-dat) (许可证: **GPL-3.0**): _针对俄罗斯被封锁域名与地址生成、自动更新的 V2Ray 路由规则。_

## Support project

**If this project is helpful to you, you may wish to give it a**:star2:

<a href="https://www.buymeacoffee.com/MHSanaei" target="_blank">
<img src="./media/default-yellow.png" alt="Buy Me A Coffee" style="height: 70px !important;width: 277px !important;" >
</a>

</br>
<a href="https://nowpayments.io/donation/hsanaei" target="_blank" rel="noreferrer noopener">
   <img src="./media/donation-button-black.svg" alt="Crypto donation button by NOWPayments">
</a>

## Stargazers over Time

[![Stargazers over time](https://starchart.cc/MHSanaei/3x-ui.svg?variant=adaptive)](https://starchart.cc/MHSanaei/3x-ui)
