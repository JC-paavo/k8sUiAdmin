# K8s UI Admin - Kubernetes 部署指南

## 前置条件

1. 已配置好的 Kubernetes 集群
2. `kubectl` 命令行工具已配置连接到集群
3. `k8s-ui-admin` 镜像已推送到集群可访问的镜像仓库（或者使用本地镜像）

## 使用本地镜像

如果你的镜像在本地（使用刚才构建的 `k8s-ui-admin:latest`），需要先将镜像加载到 Kubernetes 节点上，或者推送到一个 registry。

### 方式一：推送到 Docker Hub 或私有镜像仓库

```bash
# 重新打标签
docker tag k8s-ui-admin:latest your-registry/k8s-ui-admin:latest

# 推送
docker push your-registry/k8s-ui-admin:latest
```

然后修改 `k8s-ui-admin.yaml` 中的 `image` 字段为你的镜像地址：
```yaml
image: your-registry/k8s-ui-admin:latest
```

### 方式二：使用 kind/k3s/minikube 等本地集群

对于本地开发集群，可以直接加载镜像：

**kind:**
```bash
kind load docker-image k8s-ui-admin:latest
```

**minikube:**
```bash
minikube image load k8s-ui-admin:latest
```

## 部署步骤

### 1. 配置管理员密码（可选）

编辑 `k8s-ui-admin.yaml` 中的 ConfigMap 部分，修改密码：
```yaml
data:
  ADMIN_PASSWORD: "your-secure-password-here"  # 修改为你想要的密码
```

### 2. 部署到 Kubernetes

```bash
kubectl apply -f k8s-ui-admin.yaml
```

### 3. 验证部署

检查 Pod 状态：
```bash
kubectl get pods -n k8s-ui-admin
```

检查 Service：
```bash
kubectl get svc -n k8s-ui-admin
```

查看日志：
```bash
kubectl logs -f deployment/k8s-ui-admin -n k8s-ui-admin
```

## 访问应用

### 方式一：端口转发（快速测试）

```bash
kubectl port-forward svc/k8s-ui-admin 8080:8080 -n k8s-ui-admin
```

然后访问：http://localhost:8080

### 方式二：使用 NodePort（可选）

如果需要从外部访问，取消 `k8s-ui-admin.yaml` 中 NodePort Service 部分的注释，然后重新应用：

```bash
kubectl apply -f k8s-ui-admin.yaml
```

访问地址：http://<node-ip>:30880

### 方式三：使用 Ingress（推荐用于生产）

取消 Ingress 部分的注释，修改主机名，然后应用：

```bash
kubectl apply -f k8s-ui-admin.yaml
```

## 管理

### 更新镜像

```bash
# 更新镜像
kubectl set image deployment/k8s-ui-admin k8s-ui-admin=your-registry/k8s-ui-admin:new-version -n k8s-ui-admin

# 或者使用滚动重启
kubectl rollout restart deployment/k8s-ui-admin -n k8s-ui-admin
```

### 查看部署历史

```bash
kubectl rollout history deployment/k8s-ui-admin -n k8s-ui-admin
```

### 回滚

```bash
kubectl rollout undo deployment/k8s-ui-admin -n k8s-ui-admin
```

### 清理

```bash
kubectl delete -f k8s-ui-admin.yaml
```

## 默认登录信息

- **用户名**：`admin`
- **密码**：`admin`（或你在 ConfigMap 中配置的密码）
