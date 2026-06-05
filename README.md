# Fortress-pingan

**Fortress：轻量级运维堡垒机** — 基于 WebSocket + SSH 的 Web 运维安全审计网关

## 架构

```
Browser (xterm.js) ←→ WebSocket ←→ Go Backend ←→ SSH → Target Server
                                      ↓
                                 Audit Recorder → MySQL (audit_logs)
                                      ↓
                                 REST API (Gin) → Asset / Session / User CRUD
```

```
                    ┌────────────────────────────────────┐
                    │          WS Handler (4 goroutines)  │
                    │                                    │
  Browser           │  sshReadLoop: TeeReader → WS + Audit│   Target
  xterm.js  ◄──WS──┼── wsReadLoop: WS → SSH stdin       │──SSH──► Server
                    │  auditLoop:  batch insert (50/500ms)│
                    │  heartbeat:  30s ping              │
                    └────────────────────────────────────┘
```

## 技术栈

| 组件 | 选型 |
|------|------|
| 后端框架 | Go 1.25 + Gin (REST) |
| WebSocket | gorilla/websocket v1.5 |
| SSH 协议 | golang.org/x/crypto/ssh |
| Web 终端 | xterm.js + FitAddon |
| 前端 | React 19 + TypeScript + Vite 6 + Tailwind CSS 3.4 |
| 数据库 | MySQL 8.0 (sqlx + 手写 SQL) |
| 加密 | AES-256-GCM (凭证加密) + bcrypt (密码哈希) |
| 认证 | JWT (golang-jwt/v5) |

## 核心功能

- **Web SSH 终端**：基于 xterm.js 的全功能终端，支持 resize、ANSI 转义、256 色
- **WebSocket-SSH 桥接**：4 goroutine 架构，JSON 控制 + 二进制 I/O 混合协议
- **全量命令审计**：`io.TeeReader` 分流 → `bufio.Scanner` 分行 → 批量入库（50条/500ms 刷新）
- **资产凭证加密**：AES-256-GCM 加密存储，解密仅在使用时短暂持有
- **会话管理**：实时会话列表、强制终止、自动清理
- **审计日志检索**：支持按用户/资产/时间/关键字四维搜索和分页
- **JWT 认证**：bcrypt 密码哈希 + JWT Bearer Token + 角色控制（admin/operator）

## 项目结构

```
├── schema/             # MySQL DDL（users, assets, sessions, audit_logs）
├── backend/            # Go 后端
│   ├── cmd/server/     # 入口：DB 连接, 依赖注入, 信号处理
│   └── internal/
│       ├── model/      # User, Asset, Session, AuditLog 实体
│       ├── repository/ # sqlx 数据访问层
│       ├── service/    # 业务逻辑（认证/资产/会话/审计）
│       ├── api/rest/   # Gin REST handlers + JWT 中间件
│       ├── api/ws/     # WebSocket handler（核心桥接）+ Hub 会话管理
│       ├── ssh/        # SSH 客户端封装 + PTY 会话
│       └── audit/      # 审计记录器（TeeReader → Scanner → Batch Insert）
├── frontend/           # React 前端
│   └── src/
│       ├── api/        # Fetch 封装 + JWT 注入
│       ├── components/ # Layout（侧边栏 + 暗黑模式）
│       ├── hooks/      # useWebSocket（自动重连）
│       └── pages/      # Login, Dashboard, Hosts, Terminal, Sessions, AuditLog
└── deploy/             # docker-compose.yml
```

## 快速开始

```bash
# 启动全部服务（MySQL + Backend :8080 + Frontend :80）
cd deploy && docker compose up -d

# 浏览器访问 http://localhost
# 默认账号：admin / admin123
```

本地开发：

```bash
# 后端
cd backend && go run ./cmd/server/

# 前端
cd frontend && npm install && npm run dev
# 访问 http://localhost:3000，API 自动代理到 :8080
```

## 数据库

| 表 | 说明 |
|----|------|
| `users` | 用户账号，bcrypt 哈希，角色 admin/operator |
| `assets` | 管理资产，AES-GCM 加密凭证 |
| `sessions` | SSH 会话记录，active/closed 状态 |
| `audit_logs` | 命令审计日志，按 session/user/asset/time 索引 |

## 安全设计

- 凭证：AES-256-GCM 加密 + Base64 编码存储，密钥通过环境变量注入
- 密码：bcrypt(cost=10) 哈希，JWT 24h 过期
- 审计：日志表只读权限，前端不可编辑
- 生产增强：mTLS 到目标机、KMS 密钥管理、审计日志 HMAC 防篡改

## 性能特征

- Go 后端二进制：<15MB（CGO_ENABLED=0，静态编译）
- 前端：Nginx 反向代理 + WebSocket 升级支持
- 审计写入：批量 50 条/500ms 窗口，DB 写放大降低 50 倍

## License

MIT
