# 数据库文档

## 数据库类型

- **主数据库**: PostgreSQL
- **缓存**: Redis

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

## 索引

| 表 | 字段 | 类型 | 说明 |
|----|------|------|------|
| users | user_name | UNIQUE | 用户名唯一索引 |
| restaurants | name | INDEX | 餐厅名称普通索引 |
| ratings | user_id | INDEX | 用户 ID 索引 |
| ratings | restaurant_id | INDEX | 餐厅 ID 索引 |

## Redis 缓存

Redis 用于缓存以下数据：

- 餐厅列表
- 餐厅详情
- 附近餐厅查询结果
- 用户会话信息

## ORM

项目使用 **GORM** 作为 ORM 框架，自动创建和迁移数据表。
