package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"k8s-ui-admin/internal/model"
	"k8s-ui-admin/internal/repository"
	"k8s-ui-admin/internal/pkg"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Login(username, password string) (*model.User, string, error) {
	var user model.User
	err := repository.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户名或密码错误")
		}
		return nil, "", err
	}

	if !user.Status {
		return nil, "", errors.New("用户已被禁用")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	token, err := pkg.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := repository.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(user *model.User) error {
	var existing model.User
	err := repository.DB.Unscoped().Where("username = ?", user.Username).First(&existing).Error
	if err == nil {
		return errors.New("用户名已存在")
	}

	if user.Email != "" {
		err = repository.DB.Where("email = ?", user.Email).First(&existing).Error
		if err == nil {
			return errors.New("邮箱已存在")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return repository.DB.Create(user).Error
}

func (s *UserService) UpdateUser(user *model.User) error {
	if user.ID == 1 {
		return errors.New("默认管理员账号不可编辑")
	}
	updates := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"status":   user.Status,
	}
	if user.Password != "" && len(user.Password) >= 6 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hashedPassword)
	}
	return repository.DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error
}

func (s *UserService) DeleteUser(id uint) error {
	if id == 1 {
		return errors.New("不能删除管理员账号")
	}
	repository.DB.Where("user_id = ?", id).Delete(&model.ClusterPermission{})
	return repository.DB.Unscoped().Delete(&model.User{}, id).Error
}

func (s *UserService) ListUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := repository.DB.Model(&model.User{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", kw, kw)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user model.User
	err := repository.DB.First(&user, userID).Error
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("旧密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return repository.DB.Save(&user).Error
}

func (s *UserService) GetUserClusterPermissions(userID uint) ([]model.ClusterPermission, error) {
	var permissions []model.ClusterPermission
	err := repository.DB.Preload("Cluster").Where("user_id = ?", userID).Find(&permissions).Error
	return permissions, err
}

func (s *UserService) CheckClusterPermission(userID, clusterID uint, requiredPermission string) (bool, error) {
	if requiredPermission == "read" {
		var count int64
		err := repository.DB.Model(&model.ClusterPermission{}).Where("user_id = ? AND cluster_id = ? AND permission IN (?, ?)", userID, clusterID, "read", "write").Count(&count).Error
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}

	var count int64
	err := repository.DB.Model(&model.ClusterPermission{}).Where("user_id = ? AND cluster_id = ? AND permission = ?", userID, clusterID, requiredPermission).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}