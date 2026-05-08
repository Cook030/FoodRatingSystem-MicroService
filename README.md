# 美食评分系统 (Food Rating System)

基于微服务架构的美食评分系统，支持用户注册登录、餐厅管理、美食评分、地图展示等功能。

## 技术栈

### 后端
- **语言**: Go 1.26.1
- **Web 框架**: Gin
- **RPC**: gRPC + Protocol Buffers
- **数据库**: PostgreSQL (pgx)
- **缓存**: Redis
- **ORM**: GORM
- **认证**: JWT (golang-jwt)

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
│   │   ├── grpc-clients/      # gRPC 客户端
│   │   ├── handler/           # HTTP 处理器
│   │   ├── middleware/        # 中间件 (CORS, JWT)
│   │   └── router/            # 路由配置
│   │
│   ├── user-service/          # 用户服务 (gRPC)
│   │   ├── repository/        # 数据访问层
│   │   └── service/           # 业务逻辑层
│   │
│   ├── restaurant-service/    # 餐厅服务 (gRPC)
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── rating-service/        # 评分服务 (gRPC)
│   │   ├── grpc-client/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── shared/                # 共享模块
│   │   ├── config/            # 配置管理
│   │   ├── database/          # 数据库连接 (PostgreSQL, Redis)
│   │   ├── model/             # 数据模型
│   │   └── utils/             # 工具函数 (JWT, 距离计算)
│   │
│   └── proto/                 # Protocol Buffer 定义
│       ├── user/
│       ├── restaurant/
│       └── rating/
│
└── go.mod
```

## 功能特性

- 用户注册/登录 (JWT 认证)
- 餐厅浏览与搜索
- 附近餐厅推荐 (基于地理位置)
- 美食评分与评论
- 地图可视化展示
- CORS 跨域支持

## API 接口

### 公开接口
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/user/register` | 用户注册 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/restaurants` | 获取餐厅列表 |
| GET | `/api/restaurants/:id` | 获取餐厅详情 |
| GET | `/api/restaurants/:id/ratings` | 获取餐厅评分 |
| GET | `/api/restaurants/nearby` | 获取附近餐厅 |
| GET | `/api/restaurants/recommend` | 获取推荐餐厅 |

### 需要认证的接口
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/rating` | 提交评分 |
| POST | `/api/restaurants` | 创建餐厅 |

## 快速开始

### 环境要求
- Go 1.26+
- Node.js 18+
- PostgreSQL
- Redis

### 后端启动

```bash
cd server

# 安装依赖
go mod download

# 启动 API 网关
go run api-gateway/main.go

# 启动用户服务
go run user-service/main.go

# 启动餐厅服务
go run restaurant-service/main.go

# 启动评分服务
go run rating-service/main.go
```

### 前端启动

```bash
cd client

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 认证说明

系统使用 JWT (JSON Web Token) 进行身份认证：

1. 用户登录/注册后获得 Token
2. 受保护的接口需要在请求头中携带 Token：
   ```
   Authorization: Bearer <token>
   ```
3. Token 有效期为 24 小时

> ⚠️ 注意：生产环境请修改 JWT Secret，建议从环境变量读取。

## 架构说明

```
┌─────────────┐
│   Client    │  React 前端
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────┐
│ API Gateway │  Gin (端口: 8080)
│  + JWT Auth │
└──────┬──────┘
       │ gRPC
       ▼
┌──────┬──────┬──────────┐
│ User │ Rest │  Rating  │  微服务
│ Svc  │ Svc  │  Svc     │
└──────┴──────┴──────────┘
       │
       ▼
┌─────────────┐
│ PostgreSQL  │  数据库
│   + Redis   │  缓存
└─────────────┘
```
