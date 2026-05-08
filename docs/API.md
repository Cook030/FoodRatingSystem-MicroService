# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: JWT Bearer Token

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

**响应**:
```json
{
  "message": "注册成功",
  "user": {
    "id": "string",
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

**响应**:
```json
{
  "message": "登录成功",
  "user": {
    "id": "string",
    "username": "string"
  },
  "token": "string"
}
```

## 餐厅接口

### 获取餐厅列表

```
GET /api/restaurants
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| size | int | 否 | 每页数量，默认 20 |
| category | string | 否 | 餐厅分类 |

**响应**:
```json
{
  "restaurants": [...],
  "total": 100
}
```

### 获取餐厅详情

```
GET /api/restaurants/:id
```

**响应**:
```json
{
  "id": 1,
  "name": "餐厅名称",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "avg_score": 4.5,
  "category": "中餐",
  "created_at": "2024-01-01T00:00:00Z",
  "review_count": 100
}
```

### 获取餐厅评分

```
GET /api/restaurants/:id/ratings
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**响应**:
```json
{
  "ratings": [
    {
      "id": 1,
      "user_id": 1,
      "stars": 4.5,
      "comment": "很好吃",
      "created_at": "2024-01-01T00:00:00Z",
      "user": {
        "user_name": "用户名"
      }
    }
  ]
}
```

### 获取附近餐厅

```
GET /api/restaurants/nearby
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| lat | float | 是 | 纬度 |
| lng | float | 是 | 经度 |
| radius | float | 否 | 半径（公里），默认 5 |

### 获取推荐餐厅

```
GET /api/restaurants/recommend
```

### 创建餐厅 (需要认证)

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

## 评分接口

### 提交评分 (需要认证)

```
POST /api/rating
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "restaurant_id": 1,
  "stars": 4.5,
  "comment": "很好吃"
}
```

## 其他接口

### 健康检查

```
GET /api/health
```

**响应**:
```json
{
  "status": "ok"
}
```

## 错误响应

所有错误响应格式统一为：

```json
{
  "error": "错误信息"
}
```

**常见状态码**:
| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 无效 |
| 500 | 服务器内部错误 |
