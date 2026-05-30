# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: JWT Bearer Token

## 通用响应格式

成功响应通常返回 `200 OK` 或 `201 Created`，失败响应返回对应的 HTTP 状态码和错误信息：

```json
{
  "error": "错误描述"
}
```

---

## 健康检查

### 获取系统健康状态

```
GET /api/health
```

返回 API Gateway 及各 gRPC 微服务的健康状态。

**响应** (`200 OK` / `503 Service Unavailable`):

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

- `status`: `healthy` (全部正常) 或 `degraded` (有服务异常)
- `services`: 各 gRPC 服务的健康检查结果，可能的状态有 `SERVING` / `NOT_SERVING` / `UNKNOWN`

---

## 认证

### 用户注册

```
POST /api/user/register
```

**请求体**:

```json
{
  "username": "string",
  "password": "string"
}
```

**响应** (`200 OK`):

```json
{
  "message": "注册成功",
  "user": {
    "id": 1,
    "username": "string"
  },
  "token": "string"
}
```

### 用户登录

```
POST /api/user/login
```

**请求体**:

```json
{
  "username": "string",
  "password": "string"
}
```

**响应** (`200 OK`):

```json
{
  "message": "登录成功",
  "user": {
    "id": 1,
    "username": "string"
  },
  "token": "string"
}
```

---

## 餐厅接口

### 获取餐厅列表

```
GET /api/restaurants
```

支持按地理位置、搜索关键词排序查询。

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| lat | float | 否 | 纬度，用于计算距离 |
| lon | float | 否 | 经度，用于计算距离 |
| search | string | 否 | 搜索关键词（餐厅名称） |
| sort | string | 否 | 排序方式：`distance` (默认) / `score` |

**响应** (`200 OK`):

```json
[
  {
    "id": 1,
    "name": "餐厅名称",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "avg_score": 4.5,
    "category": "中餐",
    "review_count": 100,
    "distance": 1.25,
    "final_score": 4.8
  }
]
```

### 获取餐厅详情

```
GET /api/restaurants/:id
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 餐厅 ID |

**响应** (`200 OK`):

```json
{
  "id": 1,
  "name": "餐厅名称",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "avg_score": 4.5,
  "category": "中餐",
  "review_count": 100
}
```

### 获取餐厅评分列表

```
GET /api/restaurants/:id/ratings
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 餐厅 ID |

**响应** (`200 OK`):

```json
[
  {
    "id": 1,
    "user_id": 1,
    "restaurant_id": 1,
    "stars": 4.5,
    "comment": "很好吃",
    "created_at": "2024-01-01T00:00:00Z",
    "username": "用户名"
  }
]
```

### 获取附近餐厅

```
GET /api/restaurants/nearby
```

基于用户当前地理位置获取附近餐厅，按距离排序。

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| lat | float | 是 | 纬度 |
| lon | float | 是 | 经度 |

**响应** (`200 OK`):

```json
[
  {
    "id": 1,
    "name": "餐厅名称",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "avg_score": 4.5,
    "category": "中餐",
    "review_count": 100,
    "distance": 0.85
  }
]
```

### 获取推荐餐厅

```
GET /api/restaurants/recommend
```

基于地理位置和综合评分算法返回推荐餐厅列表。

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| lat | float | 是 | 纬度 |
| lon | float | 是 | 经度 |

**响应** (`200 OK`):

```json
[
  {
    "id": 1,
    "name": "餐厅名称",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "avg_score": 4.5,
    "category": "中餐",
    "review_count": 100,
    "distance": 1.2,
    "final_score": 4.85
  }
]
```

### 创建餐厅（需要认证）

```
POST /api/restaurants
Authorization: Bearer <token>
```

**请求体**:

```json
{
  "name": "餐厅名称",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "category": "中餐"
}
```

**响应** (`201 Created`):

```json
{
  "id": 1,
  "name": "餐厅名称",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "avg_score": 0,
  "category": "中餐",
  "review_count": 0
}
```

---

## 评分接口

### 提交评分（需要认证）

```
POST /api/rating
Authorization: Bearer <token>
```

**请求体**:

```json
{
  "restaurant_id": 1,
  "restaurant_name": "餐厅名称",
  "stars": 4.5,
  "comment": "很好吃"
}
```

- `restaurant_id` 和 `restaurant_name` 至少提供一个
- `stars`: 评分星级，范围 0-5
- `comment`: 评论内容

**响应** (`200 OK`):

```json
{
  "message": "评价成功！"
}
```

---

## 认证相关说明

### 获取 Token

通过 `/api/user/register` 或 `/api/user/login` 接口获取 `token`。

### 使用 Token

在请求头中携带：

```
Authorization: Bearer <token>
```

### Token 过期处理

当 Token 过期或无效时，受保护接口会返回 `401 Unauthorized`：

```json
{
  "error": "Invalid or expired token"
}
```

前端应清除本地存储的 Token 和用户信息，并重定向到登录页面。

### 受保护接口列表

以下接口需要在请求头中携带有效的 JWT Token：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/rating` | 提交评分 |
| POST | `/api/restaurants` | 创建餐厅 |
