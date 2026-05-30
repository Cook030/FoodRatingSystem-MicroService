# 美食评分系统 (Food Rating System)

基于微服务架构的美食评分系统，支持用户注册登录、餐厅管理、美食评分、地图展示等功能。采用 gRPC 进行服务间通信，etcd 实现服务注册发现与客户端负载均衡，支持 Docker 容器化一键部署。

## 技术栈

### 后端
- **语言**: Go 1.26.1
- **Web 框架**: Gin
- **RPC**: gRPC + Protocol Buffers
- **服务注册发现**: etcd + gRPC 自定义 Resolver
- **负载均衡**: gRPC Round-Robin (客户端负载均衡)
- **数据库**: PostgreSQL (pgx)
- **缓存**: Redis
- **ORM**: GORM
- **认证**: JWT (golang-jwt)
- **容器化**: Docker + Docker Compose

### 前端
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **路由**: React Router v6
- **地图**: React Leaflet
- **样式**: Tailwind CSS
- **HTTP 客户端**: Axios

## 项目结构

```
foodRatingSystem-microservice/
├── client/                    # 前端 React 应用
│   ├── Dockerfile             # 前端容器镜像 (Nginx 多阶段构建)
│   ├── nginx.conf             # Nginx 反向代理配置
│   ├── src/
│   │   ├── api/               # API 请求封装
│   │   ├── components/        # React 组件
│   │   ├── hooks/             # 自定义 Hooks
│   │   ├── pages/             # 页面组件
│   │   └── types/             # TypeScript 类型定义
│   └── package.json
│
├── server/                    # 后端微服务
│   ├── api-gateway/           # API 网关 (HTTP)
│   │   ├── Dockerfile
│   │   ├── grpc-clients/      # gRPC 客户端 (etcd 服务发现 + 负载均衡)
│   │   ├── handler/           # HTTP 处理器
│   │   ├── middleware/        # 中间件 (CORS, JWT)
│   │   └── router/            # 路由配置
│   │
│   ├── user-service/          # 用户服务 (gRPC)
│   │   ├── Dockerfile
│   │   ├── repository/        # 数据访问层
│   │   └── service/           # 业务逻辑层
│   │
│   ├── restaurant-service/    # 餐厅服务 (gRPC)
│   │   ├── Dockerfile
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── rating-service/        # 评分服务 (gRPC)
│   │   ├── Dockerfile
│   │   ├── grpc-client/       # 调用 restaurant-service (etcd 发现)
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── shared/                # 共享模块
│   │   ├── config/            # 配置管理 (环境变量驱动)
│   │   ├── database/          # 数据库连接 (PostgreSQL, Redis)
│   │   ├── model/             # 数据模型
│   │   ├── registry/          # etcd 服务注册与发现 + gRPC Resolver
│   │   └── utils/             # 工具函数 (JWT, 距离计算)
│   │
│   └── proto/                 # Protocol Buffer 定义
│       ├── user/
│       ├── restaurant/
│       └── rating/
│
├── docs/                      # 项目文档
│   ├── API.md                 # API 接口文档
│   └── DB.md                  # 数据库文档
│
├── docker-compose.yml         # Docker Compose 编排文件
└── go.mod
```

## 功能特性

- 用户注册/登录 (JWT 认证)
- 餐厅浏览与搜索 (支持按距离/评分排序)
- 附近餐厅推荐 (基于地理位置 + Haversine 距离计算)
- 美食评分与评论
- 地图可视化展示
- CORS 跨域支持
- **etcd 服务注册发现**：各微服务启动时自动注册，支持多实例水平扩展
- **客户端负载均衡**：API Gateway 通过 etcd 动态发现服务实例，gRPC Round-Robin 轮询调度
- **优雅关闭**：服务停止时自动从 etcd 注销，避免请求丢失
- **Docker 容器化一键部署**：8 个服务通过 Docker Compose 统一编排

## 架构说明

```
┌──────────────────────────────────────────────────────────────────┐
│                           宿主机                                  │
│  ┌─────────────┐     ┌──────────────────────────────────────┐   │
│  │   Client    │────>│          API Gateway (Gin)           │   │
│  │  React/Vite │     │  端口: 8080  + JWT Auth + CORS       │   │
│  └─────────────┘     └──────────────┬───────────────────────┘   │
│                                     │ gRPC (etcd 服务发现)       │
│                     ┌───────────────┼───────────────┐           │
│                     ▼               ▼               ▼           │
│               ┌─────────┐    ┌──────────┐    ┌──────────┐      │
│               │  User   │    │Restaurant│    │ Rating   │      │
│               │ Service │    │ Service  │    │ Service  │      │
│               │ :50051  │    │ :50052   │    │ :50053   │      │
│               └────┬────┘    └────┬─────┘    └────┬─────┘      │
│                    │              │               │              │
│                    └──────────────┼───────────────┘              │
│                                   ▼                             │
│                          ┌─────────────────┐                    │
│                          │   etcd (:2379)  │                    │
│                          │ 服务注册 / 发现  │                    │
│                          └─────────────────┘                    │
│                                   ▲                             │
│               ┌───────────────────┼───────────────────┐         │
│               ▼                   ▼                   ▼         │
│        ┌──────────────┐   ┌──────────────┐   ┌──────────────┐  │
│        │  PostgreSQL  │   │    Redis     │   │   Docker     │  │
│        │   (:5432)    │   │   (:6379)    │   │   Compose    │  │
│        └──────────────┘   └──────────────┘   └──────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### 服务间通信

| 调用方 | 被调用方 | 通信方式 | 发现机制 |
|--------|----------|----------|----------|
| API Gateway | user-service | gRPC | etcd 服务发现 + Round-Robin |
| API Gateway | restaurant-service | gRPC | etcd 服务发现 + Round-Robin |
| API Gateway | rating-service | gRPC | etcd 服务发现 + Round-Robin |
| rating-service | restaurant-service | gRPC | etcd 服务发现 + Round-Robin |

### etcd 服务注册与发现机制

1. **注册**：各 gRPC 服务启动时向 etcd 写入自身地址（`/services/<服务名>/<地址>`），附带租约自动续期
2. **发现**：API Gateway 通过自定义 gRPC Resolver 从 etcd 拉取服务地址列表，并建立 Watch 监听实时变化
3. **负载均衡**：gRPC 内置 Round-Robin，在多个可用实例间轮询分发请求
4. **注销**：服务收到终止信号时，主动从 etcd 删除注册信息，实现优雅下线

## API 接口

详见 [docs/API.md](docs/API.md)

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 (含各 gRPC 服务状态) |
| POST | `/api/user/register` | 用户注册 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/restaurants` | 获取餐厅列表 (支持搜索/排序) |
| GET | `/api/restaurants/:id` | 获取餐厅详情 |
| GET | `/api/restaurants/:id/ratings` | 获取餐厅评分列表 |
| GET | `/api/restaurants/nearby` | 获取附近餐厅 |
| GET | `/api/restaurants/recommend` | 获取推荐餐厅 |

### 需要认证的接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/rating` | 提交评分 |
| POST | `/api/restaurants` | 创建餐厅 |

## 快速开始

### 方式一：Docker Compose 一键部署（推荐）

确保已安装 Docker 和 Docker Compose。

```bash
# 克隆项目后，在项目根目录执行
docker compose up --build -d

# 查看服务日志
docker compose logs -f

# 停止所有服务
docker compose down

# 停止并清除数据卷
docker compose down -v
```

启动后访问：
- 前端：`http://localhost`
- API 网关：`http://localhost:8080/api`
- PostgreSQL：`localhost:5432`
- Redis：`localhost:6379`
- etcd：`localhost:2379`

### 方式二：本地开发模式

#### 环境要求
- Go 1.26+
- Node.js 18+
- PostgreSQL
- Redis
- etcd

#### 启动基础设施

```bash
# 启动 PostgreSQL
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=postgres postgres:16-alpine

# 启动 Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 启动 etcd
docker run -d --name etcd -p 2379:2379 \
  quay.io/coreos/etcd:v3.6.11 \
  etcd --advertise-client-urls http://0.0.0.0:2379 \
       --listen-client-urls http://0.0.0.0:2379
```

#### 启动后端

```bash
cd server

# 安装依赖
go mod download

# 启动用户服务
go run user-service/main.go

# 启动餐厅服务
go run restaurant-service/main.go

# 启动评分服务
go run rating-service/main.go

# 启动 API 网关
go run api-gateway/main.go
```

#### 启动前端

```bash
cd client

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端默认运行在 `http://localhost:5173`，API 网关运行在 `http://localhost:8080`。

### 水平扩展示例

利用 etcd 服务发现实现同一服务的多实例部署：

```bash
# 终端 1：启动 user-service 实例 1
SERVICE_ADDR=localhost:50051 SERVICE_PORT=50051 go run server/user-service/main.go

# 终端 2：启动 user-service 实例 2
SERVICE_ADDR=localhost:50054 SERVICE_PORT=50054 go run server/user-service/main.go
```

API Gateway 的 etcd + Round-Robin 负载均衡会自动发现两个实例并均匀分发请求。

## 环境变量配置

所有配置均通过环境变量注入，无需修改代码即可适配不同部署环境。

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `API_GATEWAY_PORT` | API 网关端口 | `8080` |
| `DB_HOST` | PostgreSQL 地址 | `localhost` |
| `DB_PORT` | PostgreSQL 端口 | `5432` |
| `DB_USER` | PostgreSQL 用户名 | `postgres` |
| `DB_PASSWORD` | PostgreSQL 密码 | `root` |
| `DB_NAME` | PostgreSQL 数据库名 | `postgres` |
| `DB_SSLMODE` | PostgreSQL SSL 模式 | `disable` |
| `REDIS_HOST` | Redis 地址 | `localhost` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `ETCD_ENDPOINTS` | etcd 地址列表（逗号分隔） | `localhost:2379` |
| `SERVICE_NAME` | 服务名称 | `""` |
| `SERVICE_ADDR` | 服务注册地址（供 etcd 使用） | `""` |
| `SERVICE_PORT` | 服务监听端口 | `""` |
| `JWT_SECRET` | JWT 密钥 | `your-secret-key-change-this-in-production` |
| `VITE_API_BASE_URL` | 前端 API 基础地址 | `http://localhost:8080/api` |

> ⚠️ **生产环境注意**：请务必修改 `JWT_SECRET`，建议使用随机生成的强密钥并通过环境变量注入。

## 认证说明

系统使用 JWT (JSON Web Token) 进行身份认证：

1. 用户登录/注册后获得 Token
2. 受保护的接口需要在请求头中携带 Token：
   ```
   Authorization: Bearer <token>
   ```
3. Token 有效期为 24 小时

## 数据库设计

详见 [docs/DB.md](docs/DB.md)

主要数据表：
- `users` - 用户表
- `restaurants` - 餐厅表
- `ratings` - 评分表

项目使用 GORM 自动迁移数据表结构，首次启动时会自动创建表。

## 健康检查

API Gateway 提供 `/api/health` 端点，返回各 gRPC 微服务的健康状态：

```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "services": [
    {"name": "user-service", "status": "SERVING"},
    {"name": "restaurant-service", "status": "SERVING"},
    {"name": "rating-service", "status": "SERVING"}
  ]
}
```

当任一后端服务不可用时，`status` 变为 `degraded`，HTTP 状态码变为 `503`。
