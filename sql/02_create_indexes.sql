-- K8s UI Admin - 索引优化脚本
-- 数据库: SQLite3

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 集群表索引
CREATE INDEX IF NOT EXISTS idx_clusters_name ON clusters(name);
CREATE INDEX IF NOT EXISTS idx_clusters_status ON clusters(status);
CREATE INDEX IF NOT EXISTS idx_clusters_deleted_at ON clusters(deleted_at);

-- 集群权限表索引
CREATE INDEX IF NOT EXISTS idx_cluster_permissions_user_id ON cluster_permissions(user_id);
CREATE INDEX IF NOT EXISTS idx_cluster_permissions_cluster_id ON cluster_permissions(cluster_id);
CREATE INDEX IF NOT EXISTS idx_cluster_permissions_permission ON cluster_permissions(permission);
CREATE INDEX IF NOT EXISTS idx_cluster_permissions_user_cluster ON cluster_permissions(user_id, cluster_id);
CREATE INDEX IF NOT EXISTS idx_cluster_permissions_deleted_at ON cluster_permissions(deleted_at);

-- 审计日志表索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_cluster_id ON audit_logs(cluster_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_action ON audit_logs(user_id, action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_deleted_at ON audit_logs(deleted_at);
