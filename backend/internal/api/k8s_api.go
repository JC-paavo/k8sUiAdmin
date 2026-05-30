package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s-ui-admin/internal/k8s"
	"k8s-ui-admin/internal/model"
	"k8s-ui-admin/internal/repository"
	"k8s-ui-admin/internal/service"
)

type K8sAPI struct {
	clusterService *service.ClusterService
}

func NewK8sAPI() *K8sAPI {
	return &K8sAPI{
		clusterService: service.NewClusterService(),
	}
}

func (api *K8sAPI) getCluster(c *gin.Context) (*model.Cluster, error) {
	clusterIDStr := c.Param("cluster_id")
	clusterID, err := strconv.ParseUint(clusterIDStr, 10, 32)
	if err != nil {
		return nil, err
	}
	return api.clusterService.GetClusterByID(uint(clusterID))
}

func (api *K8sAPI) ListNamespaces(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespaces, err := k8s.ListNamespaces(cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespaces)
}

func (api *K8sAPI) ListDeployments(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	deployments, err := k8s.ListDeployments(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployments)
}

func (api *K8sAPI) GetDeployment(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	deployment, err := k8s.GetDeployment(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

func (api *K8sAPI) CreateDeployment(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k8s.CreateDeployment(cluster, &deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Deployment创建成功", "deployment": deployment})
}

func (api *K8sAPI) UpdateDeployment(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deployment.Namespace = namespace
	deployment.Name = name

	err = k8s.UpdateDeployment(cluster, &deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment更新成功", "deployment": deployment})
}

func (api *K8sAPI) DeleteDeployment(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	err = k8s.DeleteDeployment(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment删除成功"})
}

func (api *K8sAPI) ScaleDeployment(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	replicasStr := c.Query("replicas")
	replicas, err := strconv.ParseInt(replicasStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "副本数格式错误"})
		return
	}

	err = k8s.ScaleDeployment(cluster, namespace, name, int32(replicas))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment缩放成功"})
}

func (api *K8sAPI) ListStatefulSets(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	sts, err := k8s.ListStatefulSets(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sts)
}

func (api *K8sAPI) GetStatefulSet(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	statefulSet, err := k8s.GetStatefulSet(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statefulSet)
}

func (api *K8sAPI) ListDaemonSets(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	daemonSets, err := k8s.ListDaemonSets(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, daemonSets)
}

func (api *K8sAPI) GetClusterResourceUsage(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	usage, err := k8s.GetClusterResourceUsage(cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usage)
}

func (api *K8sAPI) CreateDaemonSet(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	var daemonSet appsv1.DaemonSet
	if err := c.ShouldBindJSON(&daemonSet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k8s.CreateDaemonSet(cluster, &daemonSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "DaemonSet创建成功", "daemonSet": daemonSet})
}

func (api *K8sAPI) UpdateDaemonSet(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	var daemonSet appsv1.DaemonSet
	if err := c.ShouldBindJSON(&daemonSet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	daemonSet.Namespace = namespace
	daemonSet.Name = name

	err = k8s.UpdateDaemonSet(cluster, &daemonSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DaemonSet更新成功", "daemonSet": daemonSet})
}

func (api *K8sAPI) GetDaemonSet(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	daemonSet, err := k8s.GetDaemonSet(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, daemonSet)
}

func (api *K8sAPI) DeleteDaemonSet(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	err = k8s.DeleteDaemonSet(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DaemonSet删除成功"})
}

func (api *K8sAPI) ListServices(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	services, err := k8s.ListServices(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

func (api *K8sAPI) GetService(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	service, err := k8s.GetService(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

func (api *K8sAPI) CreateService(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	var service corev1.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k8s.CreateService(cluster, &service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service创建成功", "service": service})
}

func (api *K8sAPI) UpdateService(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	var service corev1.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Namespace = namespace
	service.Name = name

	err = k8s.UpdateService(cluster, &service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service更新成功", "service": service})
}

func (api *K8sAPI) DeleteService(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	err = k8s.DeleteService(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service删除成功"})
}

func (api *K8sAPI) ListIngresses(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	ingresses, err := k8s.ListIngresses(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingresses)
}

func (api *K8sAPI) GetIngress(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	ingress, err := k8s.GetIngress(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingress)
}

func (api *K8sAPI) CreateIngress(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	var ingress networkingv1.Ingress
	if err := c.ShouldBindJSON(&ingress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k8s.CreateIngress(cluster, &ingress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ingress创建成功", "ingress": ingress})
}

func (api *K8sAPI) UpdateIngress(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	var ingress networkingv1.Ingress
	if err := c.ShouldBindJSON(&ingress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingress.Namespace = namespace
	ingress.Name = name

	err = k8s.UpdateIngress(cluster, &ingress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ingress更新成功", "ingress": ingress})
}

func (api *K8sAPI) DeleteIngress(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	err = k8s.DeleteIngress(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ingress删除成功"})
}

func (api *K8sAPI) ListConfigMaps(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	configMaps, err := k8s.ListConfigMaps(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configMaps)
}

func (api *K8sAPI) GetConfigMap(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	configMap, err := k8s.GetConfigMap(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configMap)
}

func (api *K8sAPI) CreateConfigMap(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	var configMap corev1.ConfigMap
	if err := c.ShouldBindJSON(&configMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k8s.CreateConfigMap(cluster, &configMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ConfigMap创建成功", "configMap": configMap})
}

func (api *K8sAPI) UpdateConfigMap(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	var configMap corev1.ConfigMap
	if err := c.ShouldBindJSON(&configMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configMap.Namespace = namespace
	configMap.Name = name

	err = k8s.UpdateConfigMap(cluster, &configMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ConfigMap更新成功", "configMap": configMap})
}

func (api *K8sAPI) DeleteConfigMap(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	err = k8s.DeleteConfigMap(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ConfigMap删除成功"})
}

func (api *K8sAPI) ListSecrets(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	secrets, err := k8s.ListSecrets(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, secrets)
}

func (api *K8sAPI) GetSecret(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	secret, err := k8s.GetSecret(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, secret)
}

func (api *K8sAPI) CreateSecret(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	var secret corev1.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k8s.CreateSecret(cluster, &secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Secret创建成功", "secret": secret})
}

func (api *K8sAPI) UpdateSecret(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	var secret corev1.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secret.Namespace = namespace
	secret.Name = name

	err = k8s.UpdateSecret(cluster, &secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Secret更新成功", "secret": secret})
}

func (api *K8sAPI) DeleteSecret(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	err = k8s.DeleteSecret(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Secret删除成功"})
}

func (api *K8sAPI) ListPods(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	pods, err := k8s.ListPods(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pods)
}

func (api *K8sAPI) GetPod(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	pod, err := k8s.GetPod(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}

func (api *K8sAPI) DeletePod(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	err = k8s.DeletePod(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pod删除成功"})
}

func (api *K8sAPI) GetPodEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetPodEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetDeploymentEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetDeploymentEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) ListEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.DefaultQuery("namespace", "")

	events, err := k8s.ListEvents(cluster, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetPodLogs(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	podName := c.Param("name")
	containerName := c.DefaultQuery("container", "")

	logs, err := k8s.GetPodLogs(cluster, namespace, podName, containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (api *K8sAPI) GetServiceEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetServiceEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetIngressEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetIngressEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetConfigMapEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetConfigMapEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetSecretEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetSecretEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetStatefulSetEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetStatefulSetEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetDaemonSetEvents(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	events, err := k8s.GetDaemonSetEvents(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (api *K8sAPI) GetDeploymentHistory(c *gin.Context) {
	cluster, err := api.getCluster(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	history, err := k8s.GetDeploymentHistory(cluster, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (api *K8sAPI) GetPodMetrics(c *gin.Context) {
	clusterIDStr := c.Param("cluster_id")
	clusterID, err := strconv.ParseUint(clusterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID格式错误"})
		return
	}

	namespace := c.Param("namespace")
	podName := c.Param("name")

	cutoff := time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")

	var metrics []model.PodMetrics
	err = repository.DB.Where(
		"cluster_id = ? AND namespace = ? AND pod_name = ? AND collected_at >= ?",
		uint(clusterID), namespace, podName, cutoff,
	).Order("collected_at asc").Find(&metrics).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(metrics) == 0 {
		var cluster model.Cluster
		if err := repository.DB.First(&cluster, clusterID).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
			return
		}
		livePoint := fetchLiveMetrics(&cluster, namespace, podName)
		if livePoint == nil {
			c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": []TimelinePoint{*livePoint}})
		return
	}

	type TimelinePoint struct {
		Time   string  `json:"time"`
		CPU    float64 `json:"cpu"`
		Memory float64 `json:"memory"`
	}

	grouped := make(map[string][]model.PodMetrics)
	for _, m := range metrics {
		grouped[m.CollectedAt] = append(grouped[m.CollectedAt], m)
	}

	timeline := make([]TimelinePoint, 0, len(grouped))
	for ts, items := range grouped {
		var totalCPU, totalMem int64
		for _, item := range items {
			totalCPU += item.CPUMillicores
			totalMem += item.MemoryBytes
		}
		timeline = append(timeline, TimelinePoint{
			Time:   ts,
			CPU:    float64(totalCPU) / 1000.0,
			Memory: float64(totalMem) / (1024 * 1024),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": timeline})
}

type TimelinePoint struct {
	Time   string  `json:"time"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

func fetchLiveMetrics(cluster *model.Cluster, namespace, podName string) *TimelinePoint {
	metricsClient, err := k8s.GetMetricsClient(cluster)
	if err != nil {
		return nil
	}

	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + podName,
	})
	if err != nil {
		return nil
	}

	for _, pm := range podMetricsList.Items {
		if pm.Name == podName {
			var totalCPU, totalMem int64
			for _, container := range pm.Containers {
				totalCPU += container.Usage.Cpu().MilliValue()
				totalMem += container.Usage.Memory().Value()
			}
			now := time.Now().Format("2006-01-02 15:04:05")
			return &TimelinePoint{
				Time:   now,
				CPU:    float64(totalCPU) / 1000.0,
				Memory: float64(totalMem) / (1024 * 1024),
			}
		}
	}
	return nil
}

