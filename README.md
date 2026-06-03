# SimpleBank

基于 Go 的微服务银行核心系统，支持用户管理、双分录转账、JWT/PASETO 认证、异步任务队列，提供 gRPC + REST 双协议 API。

## 技术栈

| 层面 | 技术 |
|------|------|
| 语言 | Go 1.26 |
| RPC | gRPC + Protobuf |
| HTTP | grpc-gateway（自动生成 REST API） |
| 数据库 | PostgreSQL 16 + pgx/v5 |
| 代码生成 | sqlc（SQL → 类型安全 Go 代码） |
| 迁移 | golang-migrate |
| 认证 | PASETO V2 + JWT（双方案可切换） |
| 异步任务 | Asynq（基于 Redis） |
| 邮件 | SMTP（Gmail） |
| 可观测性 | Prometheus /metrics + 结构化日志 + Request ID |
| 容器化 | Docker 多阶段构建 |
| 编排 | Kubernetes Deployment + Service |
| CI/CD | GitHub Actions（test → build） |

## 快速开始

### 前置要求

- Go 1.26+
- Docker & Docker Compose

### 启动

```bash
# 1. 启动基础设施
docker compose up -d postgres redis

# 2. 生成代码 & 迁移
make sqlc
make migrateup

# 3. 启动服务
make server
```

### 验证

```bash
# 健康检查
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123","full_name":"Alice","email":"alice@test.com"}'

# 登录
curl -X POST http://localhost:8080/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}'

# Prometheus 指标
curl http://localhost:8080/metrics

# Swagger UI
open http://localhost:8080/swagger/
```

### 运行测试

```bash
make test
```

## 项目结构

```
├── cmd/server/         # 主入口，组装所有组件
├── proto/              # Protobuf 服务定义
├── pb/                 # protoc 自动生成的 Go 代码
├── gapi/               # gRPC 服务实现
├── db/
│   ├── migration/      # 数据库迁移文件
│   ├── query/          # sqlc SQL 查询
│   ├── sqlc/           # sqlc 生成的 CRUD 代码
│   └── mock/           # Mock Store
├── token/              # JWT + PASETO Token
├── worker/             # Asynq 异步任务
├── mail/               # 邮件发送
├── middleware/          # HTTP 中间件
├── health/             # 健康检查
├── doc/swagger/        # Swagger UI
├── k8s/                # Kubernetes 部署配置
├── .github/workflows/  # CI 流水线
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```

## 核心设计

### 双分录记账

转账操作在一个数据库事务中完成 5 个步骤：

```
创建 Transfer 记录
→ From Entry（出账 -100）＋ To Entry（入账 +100）
→ From Account balance -= 100 ＋ To Account balance += 100
```

每一笔资金变动都有 `entries` 记录可追溯，满足会计恒等式。

### 并发死锁预防

两个账户的余额更新按 `account_id` 大小排序后加 `FOR NO KEY UPDATE` 行级锁，保证 A→B 和 B→A 并发时不会死锁。

### 双 Token 认证

- **Access Token**（15 分钟）：服务间鉴权
- **Refresh Token**（24 小时）：换取新 Access Token
- PASETO V2 默认方案，JWT 可切换（Maker 接口）

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /v1/users | 注册用户 |
| POST | /v1/users/login | 登录 |
| PATCH | /v1/users/{username} | 更新用户信息 |
| GET | /v1/verify_email | 验证邮箱 |
| GET | /health | 存活检查 |
| GET | /health/ready | 就绪检查 |
| GET | /metrics | Prometheus 指标 |
| GET | /swagger/ | Swagger UI |

## 许可证

MIT

## 致谢

本项目设计灵感来源于 [techschool/simplebank](https://github.com/techschool/simplebank)，
在原项目基础上进行了架构优化、依赖升级和功能扩展。