package model

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	Username          string               `gorm:"unique;not null" json:"username"`
	Password          string               `gorm:"not null" json:"-"`
	Email             string               `gorm:"unique" json:"email"`
	Role              string               `gorm:"not null;default:'user'" json:"role"`
	Status            bool                 `gorm:"default:true" json:"status"`
	ClusterPermissions []ClusterPermission `gorm:"foreignKey:UserID"`
}

type Cluster struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	Name        string               `gorm:"unique;not null" json:"name"`
	Alias       string               `json:"alias"`
	Server      string               `gorm:"not null" json:"server"`
	Kubeconfig  string               `json:"kubeconfig"`
	CACert      string               `json:"caCert"`
	Token       string               `json:"-"`
	Version     string               `json:"version"`
	Status      string               `gorm:"default:'pending'" json:"status"`
	Description string               `json:"description"`
	Permissions []ClusterPermission  `gorm:"foreignKey:ClusterID"`
}

type ClusterPermission struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	UserID     uint   `gorm:"not null" json:"user_id"`
	ClusterID  uint   `gorm:"not null" json:"cluster_id"`
	Permission string `gorm:"not null" json:"permission"`
	User       User   `gorm:"foreignKey:UserID"`
	Cluster    Cluster `gorm:"foreignKey:ClusterID"`
}

type AuditLog struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	ClusterID   uint      `json:"cluster_id"`
	Action      string    `gorm:"not null" json:"action"`
	Resource    string    `json:"resource"`
	ResourceName string   `json:"resource_name"`
	Namespace   string    `json:"namespace"`
	Details     string    `json:"details"`
	User        User      `gorm:"foreignKey:UserID"`
}