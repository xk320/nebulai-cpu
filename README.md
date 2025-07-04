# Nebulai-CPU 挖矿工具

## 项目简介

Nebulai-CPU 是一个用于 Nebulai Network 平台的自动化算力提交工具，支持多账号、代理、CPU/GPU 运算（当前仅实现CPU），可自动完成任务领取、计算与提交。

## 目录结构

- `bin/` 目录下为已编译好的多平台二进制文件（Windows/Linux/macOS，支持 x86_64/arm64）。
- `main.go` 主程序入口。
- `config.json.example` 配置文件示例。
- `apis/`、`matrix/`、`logger/` 为核心功能模块。

## 快速开始

### 方式一：直接运行已编译二进制

1. 复制 `config.json.example` 为 `config.json`，并根据实际账号信息填写：

```json
{
  "accounts": [
    {
      "token": "authorization token",
      "jwt_token": "token",
      "proxy": "http://xxxxx:xxx"
    }
  ]
}
```
- `token`、`jwt_token` 可在 Nebulai 官网获取。
- `proxy` 支持 http/https/socks5，若不需要可留空。

2. 选择对应平台的二进制文件（如 macOS x86_64 用 `bin/nebulai-cpu-darwin-amd64`，Windows 用 `.exe` 文件等），命令行运行：

```bash
./bin/nebulai-cpu-darwin-amd64
```

### 方式二：源码编译运行

1. 安装 Go 1.24 及以上版本。
2. 拉取源码并安装依赖：

```bash
go mod tidy
```

3. 编译当前系统版本：

```bash
make build
```

4. 运行：

```bash
./nebulai-cpu
```

或使用 Makefile 一键打包所有平台：

```bash
make release
```

## 配置说明

- `config.json` 支持多账号，格式如下：

```json
{
  "accounts": [
    { "token": "...", "jwt_token": "...", "proxy": "..." }
  ],
  "gpu_enabled": false
}
```
- `gpu_enabled` 目前仅为预留，当前版本仅支持CPU。

## 依赖说明

- 仅依赖 `golang.org/x/net`，其余为标准库。
- 推荐 Go 1.24 及以上。

## 常见问题

1. **如何获取 token 和 jwt_token？**
   - 登录 Nebulai 官网，浏览器开发者工具中可获取。
2. **如何配置代理？**
   - 支持 http/https/socks5，格式如 `http://ip:port` 或 `socks5://ip:port`。
3. **日志输出乱码？**
   - 日志带有彩色和图标，部分终端不支持可忽略。
4. **多账号并发？**
   - 支持，所有账号任务并发执行。

## 免责声明

本工具仅供学习与研究，使用风险自负。 