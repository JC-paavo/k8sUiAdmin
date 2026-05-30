package service

import (
	"errors"

	"k8s-ui-admin/internal/k8s"
	"k8s-ui-admin/internal/model"
	"k8s-ui-admin/internal/repository"
)

type ClusterService struct{}

func NewClusterService() *ClusterService {
	return &ClusterService{}
}

func (s *ClusterService) CreateCluster(cluster *model.Cluster) error {
	var existing model.Cluster
	err := repository.DB.Where("name = ?", cluster.Name).First(&existing).Error
	if err == nil {
		return errors.New("集群名称已存在")
	}

	cluster.Status = "pending"
	err = repository.DB.Create(cluster).Error
	if err != nil {
		return err
	}

	version, err := k8s.GetVersion(cluster)
	if err != nil {
		cluster.Status = "error"
		cluster.Version = ""
	} else {
		cluster.Status = "connected"
		cluster.Version = version
	}

	return repository.DB.Save(cluster).Error
}

func (s *ClusterService) GetClusterByID(id uint) (*model.Cluster, error) {
	var cluster model.Cluster
	err := repository.DB.First(&cluster, id).Error
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (s *ClusterService) UpdateCluster(cluster *model.Cluster) error {
	oldCluster, err := s.GetClusterByID(cluster.ID)
	if err != nil {
		return err
	}

	if oldCluster.Name != cluster.Name {
		var existing model.Cluster
		err := repository.DB.Where("name = ? AND id != ?", cluster.Name, cluster.ID).First(&existing).Error
		if err == nil {
			return errors.New("集群名称已存在")
		}
	}

	k8s.RemoveClient(cluster.ID)

	err = repository.DB.Save(cluster).Error
	if err != nil {
		return err
	}

	err = k8s.TestConnection(cluster)
	if err != nil {
		cluster.Status = "error"
	} else {
		cluster.Status = "connected"
		version, _ := k8s.GetVersion(cluster)
		cluster.Version = version
	}

	return repository.DB.Save(cluster).Error
}

func (s *ClusterService) DeleteCluster(id uint) error {
	k8s.RemoveClient(id)
	return repository.DB.Delete(&model.Cluster{}, id).Error
}

func (s *ClusterService) ListClusters() ([]model.Cluster, error) {
	var clusters []model.Cluster
	err := repository.DB.Find(&clusters).Error
	return clusters, err
}

func (s *ClusterService) ListClustersByUser(userID uint) ([]model.Cluster, error) {
	var clusters []model.Cluster
	err := repository.DB.Joins("JOIN cluster_permissions ON clusters.id = cluster_permissions.cluster_id").
		Where("cluster_permissions.user_id = ?", userID).
		Find(&clusters).Error
	return clusters, err
}

func (s *ClusterService) TestClusterConnection(id uint) error {
	cluster, err := s.GetClusterByID(id)
	if err != nil {
		return err
	}

	return k8s.TestConnection(cluster)
}

func (s *ClusterService) RefreshClusterStatus(id uint) error {
	cluster, err := s.GetClusterByID(id)
	if err != nil {
		return err
	}

	k8s.RemoveClient(id)

	err = k8s.TestConnection(cluster)
	if err != nil {
		cluster.Status = "error"
	} else {
		cluster.Status = "connected"
		version, _ := k8s.GetVersion(cluster)
		cluster.Version = version
	}

	return repository.DB.Save(cluster).Error
}

func (s *ClusterService) AddUserPermission(userID, clusterID uint, permission string) error {
	var existing model.ClusterPermission
	err := repository.DB.Where("user_id = ? AND cluster_id = ?", userID, clusterID).First(&existing).Error
	if err == nil {
		existing.Permission = permission
		return repository.DB.Save(&existing).Error
	}

	perm := &model.ClusterPermission{
		UserID:     userID,
		ClusterID:  clusterID,
		Permission: permission,
	}

	return repository.DB.Create(perm).Error
}

func (s *ClusterService) RemoveUserPermission(userID, clusterID uint) error {
	return repository.DB.Where("user_id = ? AND cluster_id = ?", userID, clusterID).Delete(&model.ClusterPermission{}).Error
}

func (s *ClusterService) GetClusterPermissions(clusterID uint) ([]model.ClusterPermission, error) {
	var permissions []model.ClusterPermission
	err := repository.DB.Preload("User").Where("cluster_id = ?", clusterID).Find(&permissions).Error
	return permissions, err
}