package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"k8s-ui-admin/internal/model"
	"k8s-ui-admin/internal/service"
)

type ClusterAPI struct {
	clusterService *service.ClusterService
	userService    *service.UserService
}

func NewClusterAPI() *ClusterAPI {
	return &ClusterAPI{
		clusterService: service.NewClusterService(),
		userService:    service.NewUserService(),
	}
}

func (api *ClusterAPI) CreateCluster(c *gin.Context) {
	var cluster model.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.clusterService.CreateCluster(&cluster)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "集群创建成功", "cluster": cluster})
}

func (api *ClusterAPI) GetCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	if !api.checkClusterAccess(c, uint(id)) {
		return
	}

	cluster, err := api.clusterService.GetClusterByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "集群不存在"})
		return
	}

	c.JSON(http.StatusOK, cluster)
}

func (api *ClusterAPI) UpdateCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	var cluster model.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cluster.ID = uint(id)
	err = api.clusterService.UpdateCluster(&cluster)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "集群更新成功"})
}

func (api *ClusterAPI) DeleteCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	err = api.clusterService.DeleteCluster(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "集群删除成功"})
}

func (api *ClusterAPI) ListClusters(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var clusters []model.Cluster
	var err error

	if role == "admin" {
		clusters, err = api.clusterService.ListClusters()
	} else {
		userID, _ := c.Get("user_id")
		clusters, err = api.clusterService.ListClustersByUser(userID.(uint))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clusters)
}

func (api *ClusterAPI) TestConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	if !api.checkClusterAccess(c, uint(id)) {
		return
	}

	err = api.clusterService.TestClusterConnection(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "连接测试成功"})
}

func (api *ClusterAPI) RefreshStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	if !api.checkClusterAccess(c, uint(id)) {
		return
	}

	err = api.clusterService.RefreshClusterStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态刷新成功"})
}

func (api *ClusterAPI) AddUserPermission(c *gin.Context) {
	clusterIDStr := c.Param("id")
	clusterID, err := strconv.ParseUint(clusterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	var req struct {
		UserID     uint   `json:"user_id" binding:"required"`
		Permission string `json:"permission" binding:"required,oneof=read write"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = api.clusterService.AddUserPermission(req.UserID, uint(clusterID), req.Permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限添加成功"})
}

func (api *ClusterAPI) RemoveUserPermission(c *gin.Context) {
	clusterIDStr := c.Param("id")
	clusterID, err := strconv.ParseUint(clusterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID格式错误"})
		return
	}

	err = api.clusterService.RemoveUserPermission(uint(userID), uint(clusterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限移除成功"})
}

func (api *ClusterAPI) GetClusterPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	permissions, err := api.clusterService.GetClusterPermissions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

func (api *ClusterAPI) checkClusterAccess(c *gin.Context, clusterID uint) bool {
	role, _ := c.Get("role")
	if role == "admin" {
		return true
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return false
	}
	hasPermission, err := api.userService.CheckClusterPermission(userID.(uint), clusterID, "read")
	if err != nil || !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问该集群"})
		return false
	}
	return true
}