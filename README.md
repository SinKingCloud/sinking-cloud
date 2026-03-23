<!-- markdownlint-disable MD033 -->
<div align="center">

# Sinking-Cloud

基于 Go 语言构建的微服务云应用框架，专为分布式、可扩展系统设计。

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/github/license/SinKingCloud/sinking-cloud)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)
![Architecture](https://img.shields.io/badge/Architecture-Microservices-green)

</div>

## 简介

Sinking-Cloud 是一个完整的微服务框架，提供认证、用户管理、支付、WebSocket 通信、定时任务等功能。基于自研的 [sinking-go](https://github.com/SinKingCloud/sinking-go) 框架构建，通过 Consul 实现服务发现和 RPC 通信。

## 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                     service-gateway                         │
│                      (API 网关 :1000)                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
        ▼             ▼             ▼
   ┌─────────┐  ┌─────────┐   ┌─────────┐
   │  auth   │  │  user   │   │  admin  │
   │  :1004  │  │  :1007  │   │  :1008  │
   └─────────┘  └─────────┘   └─────────┘
                                             
        └─────────────┼─────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
        ▼             ▼             ▼
   ┌─────────┐  ┌─────────┐   ┌─────────┐
   │  socket │  │   pay   │   │  cron   │
   │  :1001  │  │  :1005  │   │  :1002  │
   └─────────┘  └─────────┘   └─────────┘
```

## 功能特性

| 类别 | 功能 |
|------|------|
| 🔐 **认证授权** | 基于 JWT 的登录/注销，支持密码和短信验证码 |
| 👤 **用户管理** | 完整的用户 CRUD、会话管理、余额管理 |
| ⚙️ **后台管理** | RBAC 角色权限管理、菜单访问控制 |
| 🧩 **验证码** | 滑块验证码，保护表单安全 |
| 📱 **短信/邮件** | 阿里云短信和邮件服务集成 |
| 💳 **支付功能** | 支付宝集成，支持支付订单 |
| 🔌 **WebSocket** | 实时推送通知 |
| ⏰ **定时任务** | Cron 定时任务支持 |
| 📦 **对象存储** | 阿里云 OSS 文件存储 |
| 🔍 **服务发现** | 基于 Consul 的微服务通信 |
| ⚡ **配置管理** | 系统配置和用户配置，Redis 缓存 |

## 技术栈

<div align="center">

| 技术 | 用途 |
|------|------|
| ![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go) | 编程语言 |
| ![GORM](https://img.shields.io/badge/GORM-ORM-green) | ORM 框架 |
| ![Redis](https://img.shields.io/badge/Redis-Cache-red) | 缓存层 |
| ![Consul](https://img.shields.io/badge/Consul-Service%20Discovery-yellow) | 服务发现 |
| ![NATS](https://img.shields.io/badge/NATS-Message%20Queue-green) | 消息队列 |
| ![Alibaba Cloud](https://img.shields.io/badge/Alibaba%20Cloud-Cloud%20Services-orange) | 云服务 |

</div>

## 微服务

| 服务 | 端口 | 说明 |
|------|------|------|
| 🌀 gateway | 1000 | API 网关，反向代理和路由 |
| 🔐 auth | 1004 | 认证服务，验证码，短信/邮件验证 |
| 👤 user | 1007 | 用户服务，会话、配置管理 |
| ⚙️ admin | 1008 | 后台管理，RBAC 权限控制 |
| 🔌 socket | 1001 | WebSocket 实时推送 |
| 💳 pay | 1005 | 支付服务 |
| ⏰ cron | 1002 | 定时任务 |
| 📦 driver | 1003 | 存储驱动 |

## 快速开始

### 环境要求

- Go 1.25+
- Consul
- MySQL/PostgreSQL
- Redis
- （可选）阿里云账号用于短信/邮件/OSS

### 克隆和编译

```bash
git clone https://github.com/SinKingCloud/sinking-cloud.git
cd sinking-cloud/cloud-server/cloud-service

# 编译各服务
cd service-gateway && go build -o gateway .
cd ../service-auth && go build -o auth .
# ... 编译其他服务
```

### 配置

每个服务都有独立的 `config/application.yml`：

```yaml
server:
  mode: dev
  host: 0.0.0.0
  port: 1004

consul:
  group: sinking-cloud
  address: 127.0.0.1:8500
  token: your-consul-token
  local: 127.0.0.1:1004
```

### 运行

```bash
# 启动 Consul
consul agent -dev

# 启动各��务（每个终端运行一个）
./service-gateway/gateway
./service-auth/auth
./service-user/user
./service-admin/admin
```

## 项目结构

```
sinking-cloud/
├── cloud-server/
│   ├── cloud-common/          # 公共库
│   │   ├── bootstrap/        # 框架初始化
│   │   ├── constant/         # 常量定义（服务名、JWT等）
│   │   ├── global/           # 全局应用状态
│   │   ├── middleware/       # HTTP 中间件（CORS、安全）
│   │   ├── service/         # 公共业务服务
│   │   └── util/            # 工具（缓存、数据库、JWT等）
│   │
│   └── cloud-service/        # 微服务集合
│       ├── service-gateway/  # API 网关
│       ├── service-auth/     # 认证服务
│       ├── service-user/     # 用户服务
│       ├── service-admin/    # 后台管理
│       ├── service-socket/   # WebSocket
│       ├── service-cron/     # 定时任务
│       ├── service-driver/   # 存储驱动
│       └── service-pay/      # 支付服务
│
└── README.md
```

## RPC 通信

服务间通过 RPC 进行通信：

```go
// 调用 auth 服务的验证码验证
result, err := rpc.Call("Auth.CheckCaptcha", map[string]interface{}{
    "key": "xxx",
    "x":   100,
    "y":   50,
})
```

## 贡献指南

欢迎提交 Pull Request！

## 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。
