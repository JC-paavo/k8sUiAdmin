# K8s UI Admin

一个基于 Web 的 Kubernetes 集群管理平台，提供直观的图形化界面来管理集群资源、用户权限和 Pod 监控。

## 功能概览

### 用户管理
- 用户注册登录与 JWT 会话管理
- 基于角色的访问控制（管理员 / 普通用户）
- 用户列表搜索、创建、编辑、物理删除
- 管理员账户（admin）受保护不可删除

### 集群管理
- 添加 / 编辑 / 删除 Kubernetes 集群
- 测试连通性、刷新集群状态
- 集群版本检测与资源利用率展示
- 按名称搜索过滤和分页

### 资源管理
- **工作负载：** Deployment（扩缩容、历史版本）、StatefulSet、DaemonSet、Pod
- **网络：** Service、Ingress
- **配置：** ConfigMap、Secret
- **运维：** Pod 日志（多容器支持 + 下载）、事件查看、Pod 终端（WebShell）

### Pod 监控
- 基于 metrics-server 的 Pod CPU / Memory 实时采集
- ECharts 时间序列曲线图
- 30 秒自动刷新
- 采集间隔和数据保留时长可配置

### 权限系统
- 集群级访问控制（只读 / 读写）
- 管理员自动拥有所有集群权限
- 普通用户仅可见已授权集群

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3（组合式 API）、Vite、Element Plus、ECharts、xterm.js |
| 后端 | Go 1.21+、Gin、GORM、client-go、gorilla/websocket |
| 数据库 | SQLite（WAL 模式） |
| 认证 | JWT（golang-jwt） |
| 部署 | 单文件二进制（Go embed 嵌入前端）、Docker |

## 快速开始

### 方式一：单文件运行（推荐）

```bash
# 构建
./build.sh

# 运行
./k8s-ui-admin
# 访问 http://localhost:8080
```

### 方式二：开发模式

```bash
# 终端 1 - 后端
cd backend
go run ./cmd/
# http://localhost:8080

# 终端 2 - 前端
cd frontend
npm install
npm run dev
# http://localhost:5173
```

### 方式三：Docker

```bash
docker build -t k8s-ui-admin .
docker run -d -p 8080:8080 -v k8s-data:/app/data k8s-ui-admin
# 访问 http://localhost:8080
```

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin | 管理员 |

## 环境变量

配置位于 `backend/.env`：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8080 | 服务端口 |
| `DB_PATH` | `./data/k8s_ui_admin.db` | 数据库路径 |
| `JWT_SECRET` |（预置）| JWT 签名密钥 |
| `JWT_EXPIRE_HOURS` | 24 | Token 过期时长 |
| `METRICS_COLLECT_INTERVAL` | 60 | 采集间隔（秒） |
| `METRICS_RETENTION_MINUTES` | 10 | 数据保留时长（分钟） |

## 项目结构

```
k8sUiAdmin/
├── backend/
│   ├── cmd/                    # 应用入口（main + embed）
│   ├── internal/
│   │   ├── api/                # HTTP API 处理器
│   │   ├── k8s/                # K8s 客户端（client-go + metrics + exec）
│   │   ├── middleware/         # JWT 认证中间件
│   │   ├── model/              # 数据模型
│   │   ├── repository/         # 数据库初始化与自动迁移
│   │   ├── service/            # 业务逻辑层（用户/集群/Pod 采集器）
│   │   └── pkg/                # JWT 工具包
│   ├── data/                   # SQLite 数据库文件
│   └── .env                    # 环境配置
├── frontend/
│   ├── src/
│   │   ├── assets/             # 自定义 K8s SVG 图标 + 全局样式
│   │   ├── components/         # PodMetricsChart、WebTerminal 等组件
│   │   ├── layout/             # 主布局
│   │   ├── router/             # 路由配置
│   │   ├── stores/             # Pinia 状态管理
│   │   ├── utils/              # Axios API 客户端
│   │   └── views/              # 登录/仪表盘/集群/用户页面
│   └── vite.config.js
├── Dockerfile
├── build.sh
└── AGENTS.md
```

## 功能截图

<!-- 在部署后可使用实际截图替换 -->

## API 概览

### 认证
- `POST /api/login` — 登录
- `GET /api/user` — 获取当前用户
- `POST /api/user/change-password` — 修改密码

### 集群
- `GET /api/clusters` — 列出集群
- `GET /api/clusters/:id` — 获取详情

### 管理员
- `GET /api/admin/users` — 用户列表（关键字搜索、分页）
- `POST /api/admin/users` — 创建用户
- `PUT /api/admin/users/:id` — 更新用户（含密码重置）
- `DELETE /api/admin/users/:id` — 删除用户
- `POST /api/admin/clusters` — 添加集群
- `GET /api/admin/clusters/:id/permissions` — 权限列表
- `POST /api/admin/clusters/:id/permissions` — 添加权限
- `DELETE /api/admin/clusters/:id/permissions/:user_id` — 移除权限

### Kubernetes 资源
- `GET /api/k8s/:cluster_id/deployments` — 部署列表
- `GET /api/k8s/:cluster_id/pods` — Pod 列表
- `GET /api/k8s/:cluster_id/pods/:ns/:name/logs` — Pod 日志
- `GET /api/k8s/:cluster_id/pods/:ns/:name/metrics` — Pod 监控数据
- `GET /api/k8s/:cluster_id/pods/:ns/:name/exec` — Pod WebSocket 终端
- `GET /api/k8s/:cluster_id/events` — 全部事件

## 故障排查

### Pod 监控无数据
确保集群已安装 [metrics-server](https://github.com/kubernetes-sigs/metrics-server)。

### 数据库锁
SQLite WAL 模式已解决多写冲突，如遇异常可删除 WAL 文件后重启：

```bash
rm backend/k8s_ui_admin.db-shm backend/k8s_ui_admin.db-wal
```

### 重置数据
删除数据库文件后重启，系统会自动创建表结构和默认管理员：

```bash
rm backend/k8s_ui_admin.db*
```

---

*项目仍在持续开发中*
