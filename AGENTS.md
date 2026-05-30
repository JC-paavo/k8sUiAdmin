# K8s UI Admin - AI 编码助手配置

## 项目概述

**项目名称：** K8s UI Admin\
**项目类型：** 全栈 Kubernetes 管理平台\
**核心功能：** 基于 Web 的界面，用于管理 Kubernetes 集群、用户和资源\
**目标用户：** DevOps 工程师、系统管理员和平台团队

***

## 技术栈

### 前端

- **框架：** Vue 3 (组合式 API)
- **构建工具：** Vite
- **UI 框架：** Element Plus（含 Element Plus Icons Vue）
- **图表库：** ECharts（Pod 监控曲线图）
- **状态管理：** Pinia
- **路由：** Vue Router 4
- **HTTP 客户端：** Axios
- **语言：** JavaScript/Vue
- **图标：** 阿里巴巴矢量图标库（自定义 SVG 组件）+ Element Plus 图标
- **YAML 处理：** js-yaml

### 后端

- **框架：** Gin (Go Web 框架)
- **语言：** Go 1.21+
- **ORM：** GORM
- **数据库：** SQLite（WAL 模式，`_busy_timeout=5000`）
- **认证：** JWT (golang-jwt)
- **Kubernetes 客户端：** client-go
- **Metrics 客户端：** client-go metrics API (`MetricsV1beta1`)
- **密码哈希：** bcrypt (golang.org/x/crypto)

***

## 项目结构

```
k8sUiAdmin/
├── backend/                    # Go 后端服务
│   ├── cmd/
│   │   ├── main.go            # 应用入口
│   │   └── embed.go           # 前端静态资源嵌入
│   ├── internal/
│   │   ├── api/               # HTTP API 处理器
│   │   │   ├── k8s_api.go    # Kubernetes 资源 API + Pod Metrics API
│   │   │   ├── user_api.go   # 用户管理 API
│   │   │   └── cluster_api.go # 集群管理 API（含 checkClusterAccess 权限辅助方法）
│   │   ├── middleware/        # HTTP 中间件
│   │   │   └── auth.go       # JWT 认证中间件
│   │   ├── model/            # 数据模型
│   │   │   └── models.go     # 数据库模型（含 PodMetrics、gorm.DeletedAt 软删除）
│   │   ├── repository/        # 数据访问层
│   │   │   └── db.go         # 数据库初始化、AutoMigrate、ensureIndexes、migrateDeletedAt
│   │   ├── service/          # 业务逻辑层
│   │   │   ├── user_service.go
│   │   │   ├── cluster_service.go
│   │   │   └── pod_metrics.go # Pod Metrics 定时采集器（含配置项和过期清理）
│   │   ├── k8s/              # Kubernetes 客户端
│   │   │   ├── client.go     # K8s 客户端初始化（含 GetMetricsClient）
│   │   │   └── resources.go  # K8s 资源操作
│   │   └── pkg/              # 共享包
│   │       └── jwt.go        # JWT 工具
│   ├── data/                 # 数据库文件
│   │   └── k8s_ui_admin.db
│   ├── go.mod
│   ├── go.sum
│   └── .env                  # 环境配置（含 metrics 采集配置）
│
├── frontend/                  # Vue 3 前端
│   ├── src/
│   │   ├── api/             # (遗留)
│   │   ├── assets/
│   │   │   ├── icons/       # 自定义 K8s 资源 SVG 图标
│   │   │   │   ├── IconPod.vue / IconDeployment.vue / ...
│   │   │   │   └── index.js # 图标导出映射
│   │   │   └── styles/
│   │   │       └── variables.css # 全局 CSS 变量
│   │   ├── components/      # 可复用组件
│   │   │   └── PodMetricsChart.vue # Pod CPU/Memory ECharts 监控图表
│   │   ├── layout/
│   │   │   ├── Layout.vue   # 带侧边栏的主布局
│   │   │   └── ChangePassword.vue # 修改密码对话框
│   │   ├── router/
│   │   │   └── index.js     # Vue Router 配置
│   │   ├── stores/          # Pinia 状态库
│   │   │   ├── auth.js      # 认证状态
│   │   │   └── cluster.js   # 集群状态
│   │   ├── utils/
│   │   │   └── api.js       # 带拦截器的 API 客户端
│   │   ├── views/           # 页面组件
│   │   │   ├── Login.vue
│   │   │   ├── Dashboard.vue
│   │   │   ├── cluster/
│   │   │   │   ├── ClusterList.vue
│   │   │   │   ├── ClusterDetail.vue
│   │   │   │   ├── ResourceList.vue
│   │   │   │   ├── ResourceDetail.vue
│   │   │   │   ├── DeploymentList.vue / PodList.vue / ...
│   │   │   │   └── (各资源独立列表页面)
│   │   │   └── user/
│   │   │       ├── UserList.vue
│   │   │       ├── UserCreate.vue
│   │   │       └── UserEdit.vue
│   │   ├── App.vue
│   │   └── main.js
│   ├── index.html
│   ├── package.json
│   └── vite.config.js       # Vite 配置（base: './' 支持离线部署）
│
├── sql/                      # 数据库初始化脚本（参考，实际由 GORM 自动管理）
│   ├── 01_create_tables.sql
│   └── 02_create_indexes.sql
│
├── Dockerfile                # Docker 构建文件
├── build.sh                  # 本地构建脚本
├── AGENTS.md                 # 本文件
└── README.md
```

***

## 核心功能

### 1. 用户管理

- 用户注册和登录
- 密码管理（修改密码需要校验旧密码，新密码强度要求 >=6 位）
- 基于角色的访问控制（管理员/普通用户）
- JWT 会话管理
- 默认管理员（admin, ID=1）不可删除、不可编辑用户名/角色/状态
- 管理员可为每个集群分配权限（只读/读写）
- 用户列表支持关键字搜索（用户名/邮箱模糊匹配）
- 4 个统计卡片：用户总数 / 管理员数 / 已启用 / 已禁用

### 2. 集群管理

- 添加/编辑/删除 Kubernetes 集群
- 测试集群连接
- 刷新集群状态
- 版本检测
- 集群 CPU/内存资源利用率展示（基于 metrics-server API）
- 搜索过滤和分页
- 权限访问控制：管理员可见全部集群，普通用户仅可见已授权的集群

### 3. 资源管理

- **命名空间：** 列表和切换
- **部署 (Deployment)：** 创建、读取、更新、删除、扩缩容、历史版本
- **有状态副本集 (StatefulSet)：** 查询、事件
- **守护进程集 (DaemonSet)：** CRUD、事件
- **服务 (Service)：** 创建、读取、更新、删除、NodePort 展示
- **入口 (Ingress)：** 创建、读取、更新、删除、路由规则
- **配置映射 (ConfigMap)：** 创建、读取、编辑、删除
- **密钥 (Secret)：** 创建、读取、编辑、删除
- **容器组 (Pod)：** 列表、删除、多容器日志、日志下载、事件
- **事件 (Events)：** 资源事件查看
- 所有资源使用自定义 Kubernetes 图标（阿里巴巴矢量图标库 SVG）

### 4. Pod 监控（Metrics）

- **数据源：** 通过 K8s metrics-server API (`MetricsV1beta1.PodMetricses()`) 获取 Pod 实时 CPU/Memory
- **定时采集器：** 后端 `StartMetricsCollector()` 定时采集所有已连接集群的 Pod 指标并写入数据库
- **监控图表：** Pod 详情页嵌入 `PodMetricsChart.vue` 组件，使用 ECharts 绘制 CPU（cores）和 Memory（MiB）时间序列曲线
- **自动刷新：** 前端每 60 秒重新拉取数据，图表实时更新
- **实时回退：** 数据库无历史数据时，API 自动实时拉取当前指标返回（`fetchLiveMetrics`）
- **配置项（.env）：**

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `METRICS_COLLECT_INTERVAL` | 60 | 采集间隔（秒） |
| `METRICS_RETENTION_MINUTES` | 10 | 数据保留时长（分钟） |

### 5. 权限系统

- 集群级权限（只读/读写）
- 管理员 vs 普通用户角色
- 管理员自动拥有所有集群的所有权限
- JWT 中间件 + 权限检查中间件
- 集群级别 API 端点（GetCluster / TestConnection / RefreshStatus）对非管理员进行权限校验（`checkClusterAccess`）
- 用户管理页面：admin (ID=1) 不可删除、不可编辑用户名/角色/状态、不可自删
- 集群权限列表正常显示授权账号名称（`Preload("User")`）

### 6. 资源利用率

- 基于 metrics-server API 获取 Pod 实际资源使用量（CPU/Memory）
- 裸机 fallback：metrics-server 不可用时回退到 Allocatable 计算
- 节点级别和集群级别资源使用率展示

***

## API 端点

### 认证

- `POST /api/login` - 用户登录
- `GET /api/user` - 获取当前用户
- `POST /api/user/change-password` - 修改密码（校验旧密码，新密码>=6位且不同于旧密码）
- `GET /api/user/clusters/:id/permission` - 检查集群权限
- `GET /api/user/permissions` - 获取用户所有集群权限列表

### 集群（需要认证，非管理员过滤到已授权集群）

- `GET /api/clusters` - 列出集群（管理员全部，普通用户仅已授权）
- `GET /api/clusters/:id` - 获取集群详情（非管理员需 read 权限）
- `POST /api/clusters/test/:id` - 测试连接（非管理员需 read 权限）
- `POST /api/clusters/refresh/:id` - 刷新状态（非管理员需 read 权限）

### 管理员 API（需要管理员角色）

- `POST /api/admin/clusters` - 创建集群
- `PUT /api/admin/clusters/:id` - 更新集群
- `DELETE /api/admin/clusters/:id` - 删除集群
- `GET /api/admin/clusters/:id/permissions` - 获取权限列表（含 Preload User）
- `POST /api/admin/clusters/:id/permissions` - 添加权限
- `DELETE /api/admin/clusters/:id/permissions/:user_id` - 移除权限
- `GET /api/admin/users` - 列出用户（支持 ?keyword= 搜索用户名/邮箱，分页 ?page=&pageSize=）
- `POST /api/admin/users` - 创建用户（密码独立接收，最少6位）
- `PUT /api/admin/users/:id` - 更新用户（ID=1 admin 受保护不可编辑）
- `DELETE /api/admin/users/:id` - 删除用户（ID=1 admin 受保护不可删除）

### Kubernetes 资源（需要集群权限）

- `GET /api/k8s/:cluster_id/cluster/resource-usage` - 集群资源使用率（metrics-server）
- `GET /api/k8s/:cluster_id/namespaces` - 列出命名空间
- `GET /api/k8s/:cluster_id/deployments` - 列出部署
- `POST /api/k8s/:cluster_id/deployments` - 创建部署
- `GET /api/k8s/:cluster_id/deployments/:namespace/:name` - 获取部署
- `PUT /api/k8s/:cluster_id/deployments/:namespace/:name` - 更新部署
- `DELETE /api/k8s/:cluster_id/deployments/:namespace/:name` - 删除部署
- `POST /api/k8s/:cluster_id/deployments/:namespace/:name/scale` - 扩缩容部署
- `GET /api/k8s/:cluster_id/deployments/:namespace/:name/history` - 部署历史
- `GET /api/k8s/:cluster_id/deployments/:namespace/:name/events` - 部署事件
- 各类资源的事件查询（StatefulSet、DaemonSet、Service、Ingress、ConfigMap 等 .../events）
- `GET /api/k8s/:cluster_id/pods/:namespace/:name/metrics` - Pod Metrics 时间序列（CPU/Memory）
- `GET /api/k8s/:cluster_id/pods/:namespace/:name/events` - Pod 事件
- Pod 日志（支持指定容器）：`GET .../pods/:namespace/:name/logs`
- 全局事件列表：`GET /api/k8s/:cluster_id/events`

***

## 数据库 Schema

### 用户表 (Users)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 主键，admin 固定 ID=1 |
| username | VARCHAR(255) UNIQUE | 用户名 |
| password | VARCHAR(255) | bcrypt 哈希密码 |
| email | VARCHAR(255) UNIQUE | 邮箱 |
| role | VARCHAR(50) DEFAULT 'user' | admin / user |
| status | BOOLEAN DEFAULT 1 | 启用/禁用 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| deleted_at | DATETIME | 软删除时间（gorm.DeletedAt，NULL 表示未删除） |

### 集群表 (Clusters)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 主键 |
| name | VARCHAR(255) UNIQUE | 集群名称 |
| alias | VARCHAR(255) | 集群别名 |
| server | VARCHAR(255) | K8s API 服务器地址 |
| kubeconfig | TEXT | base64 编码的 kubeconfig |
| ca_cert | TEXT | CA 证书 |
| token | TEXT | 访问 Token |
| version | VARCHAR(100) | K8s 版本 |
| status | VARCHAR(50) DEFAULT 'pending' | connected/pending/error |
| description | TEXT | 描述 |

### 集群权限表 (Cluster Permissions)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 主键 |
| user_id | INTEGER FK | 关联用户 |
| cluster_id | INTEGER FK | 关联集群 |
| permission | VARCHAR(50) | read / write |
| created_at/updated_at/deleted_at | DATETIME | 时间戳 |

### 审计日志表 (Audit Logs)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 主键 |
| user_id | INTEGER FK | 操作用户 |
| cluster_id | INTEGER | 关联集群 |
| action | VARCHAR(255) | 操作类型 |
| resource | VARCHAR(255) | 资源类型 |
| resource_name | VARCHAR(255) | 资源名称 |
| namespace | VARCHAR(255) | 命名空间 |
| details | TEXT (JSON) | 详细信息 |

### Pod Metrics 表 (Pod Metrics)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 主键 |
| cluster_id | INTEGER | 关联集群 ID |
| namespace | VARCHAR(255) | 命名空间 |
| pod_name | VARCHAR(255) | Pod 名称 |
| cpu_millicores | BIGINT | CPU 使用量（毫核） |
| memory_bytes | BIGINT | 内存使用量（字节） |
| collected_at | VARCHAR(50) | 采集时间戳（格式：2006-01-02 15:04:05） |

***

## 数据库索引

以下索引在服务启动时由 `ensureIndexes()` 自动创建（`CREATE INDEX IF NOT EXISTS`）：

| 索引名 | 目标表 | 列 | 用途 |
|--------|--------|---|------|
| `idx_cluster_permissions_user_id` | cluster_permissions | `user_id` | 按用户查权限 |
| `idx_cluster_permissions_cluster_id` | cluster_permissions | `cluster_id` | 按集群查权限列表 |
| `idx_cluster_permissions_user_cluster` | cluster_permissions | `(user_id, cluster_id)` | 权限检查复合查询 |
| `idx_pod_metrics_lookup` | pod_metrics | `(cluster_id, namespace, pod_name, collected_at)` | GetPodMetrics 核心查询 |
| `idx_pod_metrics_collected_at` | pod_metrics | `collected_at` | 定期清理过期数据 |
| `idx_clusters_status` | clusters | `status` | 采集器筛选已连接集群 |
| `idx_audit_logs_user_time` | audit_logs | `(user_id, created_at)` | 用户操作日志查询 |

***

## 环境变量配置

配置文件：`backend/.env`

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8080 | 服务端口 |
| `DB_PATH` | `./data/k8s_ui_admin.db` | SQLite 数据库文件路径 |
| `JWT_SECRET` | (预置) | JWT 签名密钥 |
| `JWT_EXPIRE_HOURS` | 24 | JWT 过期时长（小时） |
| `METRICS_COLLECT_INTERVAL` | 60 | Metrics 采集间隔（秒，最小 10） |
| `METRICS_RETENTION_MINUTES` | 10 | Metrics 数据保留时长（分钟，最小 1） |

***

## 认证流程

1. 用户提交用户名/密码到 `/api/login`
2. 后端验证凭证并生成 JWT token（24小时有效期）
3. 前端将 token 存储在 localStorage
4. 所有后续请求包含 `Authorization: Bearer <token>` 请求头
5. JWT 中间件验证 token 并提取用户信息（`user_id`, `role`, `username`）
6. 管理员绕过集群权限检查；普通用户需有集群权限
7. 拦截器自动处理 401 响应并跳转到登录页

***

## 默认凭证

- **用户名：** admin
- **密码：** admin
- **角色：** admin
- **注意：** admin 账号（ID=1）不可删除、不可编辑用户名/角色/状态

***

## 运行应用

### 开发模式

```bash
# 后端
cd backend
go run cmd/main.go
# 服务运行在 http://localhost:8080

# 前端（需要后端同时运行）
cd frontend
npm install
npm run dev
# 服务运行在 http://localhost:5173
```

### 生产/单文件部署（推荐）

```bash
# 使用构建脚本
chmod +x build.sh
./build.sh

# 运行
./k8s-ui-admin
# 访问 http://localhost:8080
```

### Docker 部署

```bash
docker build -t k8s-ui-admin .
docker run -d -p 8080:8080 -v k8s-data:/app/data k8s-ui-admin
# 访问 http://localhost:8080
```

***

## 离线部署说明

- 前端构建时使用 `base: './'` 相对路径
- 构建产物通过 Go embed 嵌入后端二进制文件中
- 单文件运行，无需额外部署前端静态文件
- 所有前端资源（Element Plus、Vue、ECharts 等）在构建时打包到 dist 目录

***

## 重要约定

### 代码风格

- 变量/函数名使用英文
- 注释使用中文
- Vue 组件使用组合式 API 和 `<script setup>`
- Go 代码遵循标准 Go 约定
- 页面使用中文文案

### API 约定

- RESTful API 设计
- JSON 请求/响应格式
- 正确的 HTTP 状态码
- 错误消息使用中文
- 创建用户时密码单独接收（不与 model.User 绑定，因 Password 字段标记 `json:"-"`）

### 数据库

- 使用 GORM 进行数据库操作
- `model.User.DeletedAt` 类型为 `gorm.DeletedAt`（JSON 隐藏），空值存 NULL（非 `""`）
- 启用 WAL 模式 + `_busy_timeout=5000`（`db.go InitDB`）
- `SetMaxOpenConns(1)` 避免 SQLite 多 writer 冲突
- `SkipDefaultTransaction: true` 避免不必要的长事务
- 首次启动自动创建默认管理员 + 自动建索引
- 数据库文件从 `./data/` 目录读取（优先），回退到旧路径
- 自动迁移 `migrateDeletedAt()` 清理旧空字符串 deleted_at 值

### 前端状态管理

- 使用 Pinia stores 进行全局状态管理
- 使用 `ref` 和 `reactive` 进行组件状态管理
- 使用 `computed` 进行派生状态管理
- 权限 UI：管理员可见管理功能按钮，普通用户隐藏

### 权限设计

- 后端 `cluster_api.go` 统一使用 `checkClusterAccess` 方法验证集群访问权限
- 管理员自动通过，普通用户查询 `cluster_permissions` 表验证 read 权限
- `ListClusters` 对普通用户进行 INNER JOIN 过滤
- 前端 Dashboard 对非管理员隐藏"添加集群"按钮，空状态提示"联系管理员"

### Git 约定

- 提交信息使用中文
- 使用语义化提交格式（feat, fix, docs 等）

***

## 开发指南

### 添加新的 Kubernetes 资源

1. 在 `backend/internal/k8s/resources.go` 添加 Go 函数
2. 在 `backend/internal/api/k8s_api.go` 添加 API 处理器
3. 在 `backend/cmd/main.go` 注册路由
4. 在 `frontend/src/utils/api.js` 添加 API 方法
5. 根据需要更新前端组件

### 添加新功能

1. 遵循代码库中的现有模式
2. 保持与 API 约定的一致性
3. 添加新的端点或功能时更新此文件
4. 提交前进行全面测试

### 构建部署

1. 运行 `./build.sh` 或 `docker build -t k8s-ui-admin .`
2. 构建产物为单文件 `k8s-ui-admin`（Go 二进制 + 嵌入前端）
3. 运行 `./k8s-ui-admin` 自动创建数据库和默认管理员

***

## 故障排除

### 401 未授权错误

- 检查 JWT token 是否正确存储和发送
- 验证 token 是否过期（24小时有效期）
- 检查认证中间件配置

### Kubernetes 连接问题

- 验证 kubeconfig 是否有效
- 检查集群连接状态
- 确保正确的 RBAC 权限

### 资源利用率 / Pod 监控

- 确保集群已安装 metrics-server（`kubectl top pod` 需要 metrics-server）
- 无 metrics-server 时系统回退到 Allocatable 计算
- 采集器静默跳过 metrics-server 不可用的集群（打印日志）
- 监控数据保留时长由 `METRICS_RETENTION_MINUTES` 控制（默认10分钟）

### 新增用户无法登录

- 确认密码长度 >= 6 位
- 确认后端 `CreateUser` 使用独立请求结构体（非 `model.User`）正确接收密码
- `model.User.Password` 标记 `json:"-"`，创建时需单独处理

### 数据库问题

- 删除 `k8s_ui_admin.db` 及 WAL 文件（`.db-shm`、`.db-wal`）重置数据库
- 重启后端以使用默认管理员用户重新创建
- 表结构和索引由 GORM AutoMigrate + ensureIndexes 自动管理
- SQLite "database is locked" 已通过 WAL 模式 + 单连接 + batch insert 解决

### 前端图表不显示

- 确认 `metricsData.value = res.data?.data || []` 正确解包 axios 响应
- ECharts 实例在 `onBeforeUnmount` 中 dispose，避免内存泄漏
- 响应式图表在 `window.resize` 时自动重绘

***

*最后更新：2026-05-31*
