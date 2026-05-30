package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"k8s-ui-admin/internal/pkg"
	"k8s-ui-admin/internal/service"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := pkg.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ClusterPermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		role, _ := c.Get("role")
		if role == "admin" {
			c.Next()
			return
		}

		clusterID := c.Param("cluster_id")
		if clusterID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少集群ID"})
			c.Abort()
			return
		}

		var clusterIDUint uint
		_, err := fmt.Sscanf(clusterID, "%d", &clusterIDUint)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
			c.Abort()
			return
		}

		userService := service.NewUserService()
		hasPermission, err := userService.CheckClusterPermission(userID.(uint), clusterIDUint, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问该集群"})
			c.Abort()
			return
		}

		c.Next()
	}
}