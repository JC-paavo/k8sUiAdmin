-- K8s UI Admin - 建表脚本
-- 数据库: SQLite3

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    status BOOLEAN DEFAULT 1
);

-- 集群表
CREATE TABLE IF NOT EXISTS clusters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    name VARCHAR(255) NOT NULL UNIQUE,
    alias VARCHAR(255),
    server VARCHAR(255) NOT NULL,
    kubeconfig TEXT,
    ca_cert TEXT,
    token TEXT,
    version VARCHAR(100),
    status VARCHAR(50) DEFAULT 'pending',
    description TEXT
);

-- 集群权限表
CREATE TABLE IF NOT EXISTS cluster_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    user_id INTEGER NOT NULL,
    cluster_id INTEGER NOT NULL,
    permission VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (cluster_id) REFERENCES clusters(id)
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    user_id INTEGER NOT NULL,
    cluster_id INTEGER,
    action VARCHAR(255) NOT NULL,
    resource VARCHAR(255),
    resource_name VARCHAR(255),
    namespace VARCHAR(255),
    details TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 插入默认管理员（密码: admin, 使用 bcrypt 哈希）
INSERT OR IGNORE INTO users (id, username, password, email, role, status)
VALUES (1, 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin@example.com', 'admin', 1);
