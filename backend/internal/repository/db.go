package repository

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"k8s-ui-admin/internal/model"
)

var DB *gorm.DB

func InitDB(dbPath string, adminPassword string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath+"?_journal_mode=WAL&_busy_timeout=5000"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	// 清理旧的空字符串 deleted_at 值（兼容 string → gorm.DeletedAt 类型变更）
	migrateDeletedAt()

	err = DB.AutoMigrate(
		&model.User{},
		&model.Cluster{},
		&model.ClusterPermission{},
		&model.AuditLog{},
		&model.PodMetrics{},
	)
	if err != nil {
		return err
	}

	// 确保集群表有必要的列
	err = ensureClusterColumns()
	if err != nil {
		log.Printf("Warning: failed to ensure cluster columns: %v", err)
	}

	err = ensureIndexes()
	if err != nil {
		log.Printf("Warning: failed to ensure indexes: %v", err)
	}

	err = ensureDefaultAdmin(adminPassword)
	if err != nil {
		log.Printf("Warning: failed to ensure default admin: %v", err)
	}

	return nil
}

func migrateDeletedAt() {
	tables := []string{"users", "clusters", "cluster_permissions", "audit_logs"}
	for _, table := range tables {
		DB.Exec("UPDATE " + table + " SET deleted_at = NULL WHERE deleted_at = ''")
	}
}

// ensureClusterColumns 确保集群表有必要的列（GORM AutoMigrate 可能不会为已有表添加新列）
func ensureClusterColumns() error {
	// 检查 kubeconfig 列是否存在，不存在则添加
	var count int64
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('clusters') WHERE name = 'kubeconfig'").Scan(&count)
	if count == 0 {
		log.Println("Adding kubeconfig column to clusters table")
		DB.Exec("ALTER TABLE clusters ADD COLUMN kubeconfig TEXT")
	}

	// 检查 ca_cert 列是否存在，不存在则添加
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('clusters') WHERE name = 'ca_cert'").Scan(&count)
	if count == 0 {
		log.Println("Adding ca_cert column to clusters table")
		DB.Exec("ALTER TABLE clusters ADD COLUMN ca_cert TEXT")
	}

	return nil
}

func ensureDefaultAdmin(password string) error {
	var admin model.User
	result := DB.Where("username = ?", "admin").First(&admin)
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if result.Error == gorm.ErrRecordNotFound {
		admin = model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Role:     "admin",
			Status:   true,
		}
		return DB.Create(&admin).Error
	}

	if result.Error != nil {
		return result.Error
	}

	admin.Password = string(hashedPassword)
	return DB.Save(&admin).Error
}

func ensureIndexes() error {
	indexes := []string{
		// 集群权限查询: 按用户ID查集群/按集群ID查用户
		"CREATE INDEX IF NOT EXISTS idx_cluster_permissions_user_id ON cluster_permissions(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_cluster_permissions_cluster_id ON cluster_permissions(cluster_id)",
		"CREATE INDEX IF NOT EXISTS idx_cluster_permissions_user_cluster ON cluster_permissions(user_id, cluster_id)",

		// Pod指标查询: GetPodMetrics 按 cluster+namespace+pod+时间 组合查询
		"CREATE INDEX IF NOT EXISTS idx_pod_metrics_lookup ON pod_metrics(cluster_id, namespace, pod_name, collected_at)",

		// Pod指标清理: 定期按时间删除过期数据
		"CREATE INDEX IF NOT EXISTS idx_pod_metrics_collected_at ON pod_metrics(collected_at)",

		// 集群状态查询: 仪表盘统计/采集器过滤已连接集群
		"CREATE INDEX IF NOT EXISTS idx_clusters_status ON clusters(status)",

		// 审计日志: 按用户+时间查询操作记录
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_user_time ON audit_logs(user_id, created_at)",
	}

	for _, idx := range indexes {
		if err := DB.Exec(idx).Error; err != nil {
			log.Printf("Warning: create index failed: %v", err)
		}
	}

	return nil
}