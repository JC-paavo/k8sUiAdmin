package repository

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"k8s-ui-admin/internal/model"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
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

	err = createDefaultAdmin()
	if err != nil {
		log.Printf("Warning: failed to create default admin: %v", err)
	}

	return nil
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

func createDefaultAdmin() error {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@example.com",
		Role:     "admin",
		Status:   true,
	}

	return DB.Create(admin).Error
}