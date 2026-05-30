package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"k8s-ui-admin/internal/api"
	"k8s-ui-admin/internal/middleware"
	"k8s-ui-admin/internal/repository"
	"k8s-ui-admin/internal/service"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
}

func main() {
	// 优先使用当前目录下的数据库文件，避免环境变量干扰
	dbPath := "k8s_ui_admin.db"
	
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Failed to get working directory: %v", err)
		wd = "."
	}
	
	fullPath := filepath.Join(wd, dbPath)
	log.Printf("Using database file: %s", fullPath)

	err = repository.InitDB(fullPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	userAPI := api.NewUserAPI()
	clusterAPI := api.NewClusterAPI()
	k8sAPI := api.NewK8sAPI()

	r.POST("/api/login", userAPI.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user", userAPI.GetUser)
		auth.POST("/user/change-password", userAPI.ChangePassword)
		auth.GET("/user/permissions", userAPI.GetUserClusterPermissions)
		auth.GET("/user/clusters/:cluster_id/permission", userAPI.CheckClusterPermission)

		auth.GET("/clusters", clusterAPI.ListClusters)
		auth.GET("/clusters/:id", clusterAPI.GetCluster)
		auth.POST("/clusters/test/:id", clusterAPI.TestConnection)
		auth.POST("/clusters/refresh/:id", clusterAPI.RefreshStatus)

		admin := auth.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/users", userAPI.CreateUser)
			admin.GET("/users", userAPI.ListUsers)
			admin.GET("/users/:id", userAPI.GetUser)
			admin.PUT("/users/:id", userAPI.UpdateUser)
			admin.DELETE("/users/:id", userAPI.DeleteUser)

			admin.POST("/clusters", clusterAPI.CreateCluster)
			admin.PUT("/clusters/:id", clusterAPI.UpdateCluster)
			admin.DELETE("/clusters/:id", clusterAPI.DeleteCluster)
			admin.GET("/clusters/:id/permissions", clusterAPI.GetClusterPermissions)
			admin.POST("/clusters/:id/permissions", clusterAPI.AddUserPermission)
			admin.DELETE("/clusters/:id/permissions/:user_id", clusterAPI.RemoveUserPermission)
		}

		read := auth.Group("/k8s/:cluster_id")
	read.Use(middleware.ClusterPermissionMiddleware("read"))
	{
		read.GET("/namespaces", k8sAPI.ListNamespaces)
		read.GET("/cluster/resource-usage", k8sAPI.GetClusterResourceUsage)
		read.GET("/deployments", k8sAPI.ListDeployments)
		read.GET("/deployments/:namespace/:name", k8sAPI.GetDeployment)
		read.GET("/deployments/:namespace/:name/events", k8sAPI.GetDeploymentEvents)
		read.GET("/deployments/:namespace/:name/history", k8sAPI.GetDeploymentHistory)
		read.GET("/statefulsets", k8sAPI.ListStatefulSets)
		read.GET("/statefulsets/:namespace/:name", k8sAPI.GetStatefulSet)
		read.GET("/statefulsets/:namespace/:name/events", k8sAPI.GetStatefulSetEvents)
		read.GET("/daemonsets", k8sAPI.ListDaemonSets)
		read.GET("/daemonsets/:namespace/:name", k8sAPI.GetDaemonSet)
		read.GET("/daemonsets/:namespace/:name/events", k8sAPI.GetDaemonSetEvents)
		read.GET("/services", k8sAPI.ListServices)
		read.GET("/services/:namespace/:name", k8sAPI.GetService)
		read.GET("/services/:namespace/:name/events", k8sAPI.GetServiceEvents)
		read.GET("/ingresses", k8sAPI.ListIngresses)
		read.GET("/ingresses/:namespace/:name", k8sAPI.GetIngress)
		read.GET("/ingresses/:namespace/:name/events", k8sAPI.GetIngressEvents)
		read.GET("/configmaps", k8sAPI.ListConfigMaps)
		read.GET("/configmaps/:namespace/:name", k8sAPI.GetConfigMap)
		read.GET("/configmaps/:namespace/:name/events", k8sAPI.GetConfigMapEvents)
		read.GET("/secrets", k8sAPI.ListSecrets)
		read.GET("/secrets/:namespace/:name", k8sAPI.GetSecret)
		read.GET("/secrets/:namespace/:name/events", k8sAPI.GetSecretEvents)
		read.GET("/pods", k8sAPI.ListPods)
		read.GET("/pods/:namespace/:name", k8sAPI.GetPod)
		read.GET("/pods/:namespace/:name/logs", k8sAPI.GetPodLogs)
		read.GET("/pods/:namespace/:name/metrics", k8sAPI.GetPodMetrics)
		read.GET("/pods/:namespace/:name/events", k8sAPI.GetPodEvents)
		read.GET("/events", k8sAPI.ListEvents)
	}

		write := auth.Group("/k8s/:cluster_id")
		write.Use(middleware.ClusterPermissionMiddleware("write"))
		{
			write.POST("/deployments", k8sAPI.CreateDeployment)
			write.PUT("/deployments/:namespace/:name", k8sAPI.UpdateDeployment)
			write.DELETE("/deployments/:namespace/:name", k8sAPI.DeleteDeployment)
			write.POST("/deployments/:namespace/:name/scale", k8sAPI.ScaleDeployment)

			write.POST("/daemonsets", k8sAPI.CreateDaemonSet)
			write.PUT("/daemonsets/:namespace/:name", k8sAPI.UpdateDaemonSet)
			write.DELETE("/daemonsets/:namespace/:name", k8sAPI.DeleteDaemonSet)

			write.POST("/services", k8sAPI.CreateService)
			write.PUT("/services/:namespace/:name", k8sAPI.UpdateService)
			write.DELETE("/services/:namespace/:name", k8sAPI.DeleteService)

			write.POST("/ingresses", k8sAPI.CreateIngress)
			write.PUT("/ingresses/:namespace/:name", k8sAPI.UpdateIngress)
			write.DELETE("/ingresses/:namespace/:name", k8sAPI.DeleteIngress)

			write.POST("/configmaps", k8sAPI.CreateConfigMap)
			write.PUT("/configmaps/:namespace/:name", k8sAPI.UpdateConfigMap)
			write.DELETE("/configmaps/:namespace/:name", k8sAPI.DeleteConfigMap)

			write.POST("/secrets", k8sAPI.CreateSecret)
			write.PUT("/secrets/:namespace/:name", k8sAPI.UpdateSecret)
			write.DELETE("/secrets/:namespace/:name", k8sAPI.DeleteSecret)

			write.DELETE("/pods/:namespace/:name", k8sAPI.DeletePod)
		}
	}

	serveFrontend(r)

	service.StartMetricsCollector()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}