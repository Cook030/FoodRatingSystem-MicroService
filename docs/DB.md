# 数据库文档

## 数据库类型

- **主数据库**: PostgreSQL
- **缓存**: Redis
- **服务注册发现**: etcd (存储服务实例注册信息)

---

## 数据模型

### 用户表 (users)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PRIMARY KEY, AUTO_INCREMENT | 用户 ID |
| user_name | varchar(255) | UNIQUE, NOT NULL | 用户名 |
| password_hash | varchar(255) | NOT NULL | 密码哈希 |
| created_at | timestamp | AUTO_CREATE | 创建时间 |

**关联**:
- 一个用户可以有多个评分 (1:N)

---

### 餐厅表 (restaurants)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PRIMARY KEY, AUTO_INCREMENT | 餐厅 ID |
| name | varchar(255) | INDEX, NOT NULL | 餐厅名称 |
| latitude | decimal(10,7) | - | 纬度 |
| longitude | decimal(10,7) | - | 经度 |
| avg_score | decimal(3,2) | DEFAULT 0 | 平均评分 |
| category | varchar(100) | - | 餐厅分类 |
| created_at | timestamp | AUTO_CREATE | 创建时间 |
| review_count | int | DEFAULT 0 | 评论数量 |

**关联**:
- 一个餐厅可以有多个评分 (1:N)

---

### 评分表 (ratings)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PRIMARY KEY, AUTO_INCREMENT | 评分 ID |
| user_id | uint | INDEX, NOT NULL | 用户 ID (外键) |
| restaurant_id | uint | INDEX, NOT NULL | 餐厅 ID (外键) |
| stars | decimal(2,1) | - | 评分 (0-5) |
| comment | text | - | 评论内容 |
| created_at | timestamp | AUTO_CREATE | 创建时间 |

**关联**:
- 多对一: 属于一个用户
- 多对一: 属于一个餐厅

---

## 表关系图

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   users     │       │  ratings    │       │ restaurants │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ id (PK)     │──┐    │ id (PK)     │    ┌──│ id (PK)     │
│ user_name   │  │    │ user_id(FK) │◄───┘  │ name        │
│ password    │  └───►│ restaurant  │       │ latitude    │
│ created_at  │       │ _id (FK)    │──────►│ longitude   │
└─────────────┘       │ stars       │       │ avg_score   │
                      │ comment     │       │ category    │
                      │ created_at  │       │ created_at  │
                      └─────────────┘       │ review_count│
                                            └─────────────┘
```

---

## 索引

| 表 | 字段 | 类型 | 说明 |
|----|------|------|------|
| users | user_name | UNIQUE | 用户名唯一索引 |
| restaurants | name | INDEX | 餐厅名称普通索引 |
| ratings | user_id | INDEX | 用户 ID 索引 |
| ratings | restaurant_id | INDEX | 餐厅 ID 索引 |

---

## ORM

项目使用 **GORM** 作为 ORM 框架，自动创建和迁移数据表。

首次启动服务时会执行 `AutoMigrate`，自动创建以下表：
- `users`
- `restaurants`
- `ratings`

---

## Redis 缓存

Redis 用于缓存高频查询数据，减少数据库压力：

| 缓存内容 | 键格式 | 说明 |
|----------|--------|------|
| 餐厅列表 | `restaurants:*` | 搜索/排序结果缓存 |
| 餐厅详情 | `restaurant:<id>` | 单个餐厅信息缓存 |
| 附近餐厅 | `nearby:*` | 地理位置查询结果缓存 |
| 推荐餐厅 | `recommend:*` | 推荐算法结果缓存 |

### 缓存失效机制

当用户提交新的评分时，评分服务会通过 gRPC 调用餐厅服务的 `InvalidateCache` 接口，清除对应餐厅的缓存数据，确保数据一致性。

```
Rating Service ──gRPC──> Restaurant Service
                              └── 清除 Redis 缓存
```

---

## etcd 服务注册存储

etcd 作为服务注册中心，各微服务启动时向 etcd 写入注册信息：

### 键值结构

| 键 | 值 | 说明 |
|----|----|------|
| `/services/user-service/localhost:50051` | `localhost:50051` | 用户服务实例地址 |
| `/services/restaurant-service/localhost:50052` | `localhost:50052` | 餐厅服务实例地址 |
| `/services/rating-service/localhost:50053` | `localhost:50053` | 评分服务实例地址 |

### 租约机制

- 注册时创建 **10 秒 TTL 租约**
- 通过 `KeepAlive` 自动续期，服务存活期间租约持续有效
- 服务异常宕机时，租约过期后 etcd **自动删除**注册信息
- 服务优雅关闭时，**主动注销**注册信息

### 服务发现流程

1. API Gateway 启动时注册 etcd Resolver
2. Resolver 从 etcd 拉取 `/services/<服务名>/` 前缀下的所有地址
3. 建立 **Watch** 监听，实时感知服务实例上下线
4. gRPC 客户端通过 **Round-Robin** 策略在可用实例间轮询调度
