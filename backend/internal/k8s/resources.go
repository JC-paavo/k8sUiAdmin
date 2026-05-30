package k8s

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"k8s-ui-admin/internal/model"
)

func ListDeployments(cluster *model.Cluster, namespace string) ([]appsv1.Deployment, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetDeployment(cluster *model.Cluster, namespace, name string) (*appsv1.Deployment, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateDeployment(cluster *model.Cluster, deployment *appsv1.Deployment) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func UpdateDeployment(cluster *model.Cluster, deployment *appsv1.Deployment) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}

func DeleteDeployment(cluster *model.Cluster, namespace, name string) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	return client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func ScaleDeployment(cluster *model.Cluster, namespace, name string, replicas int32) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deployment.Spec.Replicas = &replicas
	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}

func ListStatefulSets(cluster *model.Cluster, namespace string) ([]appsv1.StatefulSet, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetStatefulSet(cluster *model.Cluster, namespace, name string) (*appsv1.StatefulSet, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func ListDaemonSets(cluster *model.Cluster, namespace string) ([]appsv1.DaemonSet, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetDaemonSet(cluster *model.Cluster, namespace, name string) (*appsv1.DaemonSet, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateDaemonSet(cluster *model.Cluster, daemonSet *appsv1.DaemonSet) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.AppsV1().DaemonSets(daemonSet.Namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	return err
}

func UpdateDaemonSet(cluster *model.Cluster, daemonSet *appsv1.DaemonSet) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.AppsV1().DaemonSets(daemonSet.Namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	return err
}

func DeleteDaemonSet(cluster *model.Cluster, namespace, name string) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	return client.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func ListServices(cluster *model.Cluster, namespace string) ([]corev1.Service, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetService(cluster *model.Cluster, namespace, name string) (*corev1.Service, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateService(cluster *model.Cluster, service *corev1.Service) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.CoreV1().Services(service.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	return err
}

func UpdateService(cluster *model.Cluster, service *corev1.Service) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.CoreV1().Services(service.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	return err
}

func DeleteService(cluster *model.Cluster, namespace, name string) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	return client.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func ListIngresses(cluster *model.Cluster, namespace string) ([]networkingv1.Ingress, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetIngress(cluster *model.Cluster, namespace, name string) (*networkingv1.Ingress, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateIngress(cluster *model.Cluster, ingress *networkingv1.Ingress) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.NetworkingV1().Ingresses(ingress.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	return err
}

func UpdateIngress(cluster *model.Cluster, ingress *networkingv1.Ingress) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.NetworkingV1().Ingresses(ingress.Namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	return err
}

func DeleteIngress(cluster *model.Cluster, namespace, name string) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	return client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func ListConfigMaps(cluster *model.Cluster, namespace string) ([]corev1.ConfigMap, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetConfigMap(cluster *model.Cluster, namespace, name string) (*corev1.ConfigMap, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateConfigMap(cluster *model.Cluster, configMap *corev1.ConfigMap) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.CoreV1().ConfigMaps(configMap.Namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	return err
}

func UpdateConfigMap(cluster *model.Cluster, configMap *corev1.ConfigMap) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.CoreV1().ConfigMaps(configMap.Namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	return err
}

func DeleteConfigMap(cluster *model.Cluster, namespace, name string) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	return client.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func ListSecrets(cluster *model.Cluster, namespace string) ([]corev1.Secret, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetSecret(cluster *model.Cluster, namespace, name string) (*corev1.Secret, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateSecret(cluster *model.Cluster, secret *corev1.Secret) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	return err
}

func UpdateSecret(cluster *model.Cluster, secret *corev1.Secret) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.CoreV1().Secrets(secret.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	return err
}

func DeleteSecret(cluster *model.Cluster, namespace, name string) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	return client.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func ListPods(cluster *model.Cluster, namespace string) ([]corev1.Pod, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetPod(cluster *model.Cluster, namespace, name string) (*corev1.Pod, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	return client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func DeletePod(cluster *model.Cluster, namespace, name string) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	return client.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func GetPodLogs(cluster *model.Cluster, namespace, podName, containerName string) (string, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return "", err
	}

	req := client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: containerName,
	})

	logs, err := req.Do(context.TODO()).Raw()
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

func ListNamespaces(cluster *model.Cluster) ([]corev1.Namespace, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

type NodeResourceUsage struct {
	Name               string  `json:"name"`
	CPUCores           float64 `json:"cpu_cores"`
	CPUCapacity        float64 `json:"cpu_capacity"`
	CPUPercent         float64 `json:"cpu_percent"`
	MemoryMi           float64 `json:"memory_mi"`
	MemoryCapacityMi   float64 `json:"memory_capacity_mi"`
	MemoryPercent      float64 `json:"memory_percent"`
	Status             string  `json:"status"`
	ActualCPUCores     float64 `json:"actual_cpu_cores"`
	ActualMemoryMi     float64 `json:"actual_memory_mi"`
	ActualCPUPercent   float64 `json:"actual_cpu_percent"`
	ActualMemoryPercent float64 `json:"actual_memory_percent"`
}

type ClusterResourceUsage struct {
	Nodes              int                 `json:"nodes"`
	TotalCPUCores      float64             `json:"total_cpu_cores"`
	TotalMemoryMi      float64             `json:"total_memory_mi"`
	AllocatedCPUCores  float64             `json:"allocated_cpu_cores"`
	AllocatedMemoryMi  float64             `json:"allocated_memory_mi"`
	CPUPercent         float64             `json:"cpu_percent"`
	MemoryPercent      float64             `json:"memory_percent"`
	ActualCPUCores     float64             `json:"actual_cpu_cores"`
	ActualMemoryMi     float64             `json:"actual_memory_mi"`
	ActualCPUPercent   float64             `json:"actual_cpu_percent"`
	ActualMemoryPercent float64            `json:"actual_memory_percent"`
	NodeDetails        []NodeResourceUsage `json:"node_details"`
}

func ListNodes(cluster *model.Cluster) ([]corev1.Node, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetClusterResourceUsage(cluster *model.Cluster) (*ClusterResourceUsage, error) {
	nodes, err := ListNodes(cluster)
	if err != nil {
		return nil, err
	}

	usage := &ClusterResourceUsage{
		Nodes: len(nodes),
		NodeDetails: make([]NodeResourceUsage, 0, len(nodes)),
	}

	metricsClient, metricsErr := GetMetricsClient(cluster)
	var nodeMetricsList *metricsv1beta1.NodeMetricsList
	if metricsErr == nil {
		nodeMetricsList, metricsErr = metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	}
	if metricsErr != nil {
		log.Printf("警告: metrics-server API 不可用，回退到 allocatable 计算: %v", metricsErr)
	}
	useMetrics := metricsErr == nil && nodeMetricsList != nil

	for _, node := range nodes {
		capacityCPU := node.Status.Capacity.Cpu().MilliValue()
		capacityMem := node.Status.Capacity.Memory().Value()

		allocatableCPU := node.Status.Allocatable.Cpu().MilliValue()
		allocatableMem := node.Status.Allocatable.Memory().Value()

		totalCPU := float64(capacityCPU) / 1000
		totalMem := float64(capacityMem) / (1024 * 1024)

		allocCPU := float64(allocatableCPU) / 1000
		allocMem := float64(allocatableMem) / (1024 * 1024)

		var status string
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				status = "Ready"
				break
			}
		}
		if status == "" {
			status = "NotReady"
		}

		nodeUsage := NodeResourceUsage{
			Name:             node.Name,
			CPUCores:         allocCPU,
			CPUCapacity:      totalCPU,
			MemoryMi:         allocMem,
			MemoryCapacityMi: totalMem,
			Status:           status,
		}

		if useMetrics {
			for _, nm := range nodeMetricsList.Items {
				if nm.Name == node.Name {
					actualCPU := float64(nm.Usage.Cpu().MilliValue()) / 1000
					actualMem := float64(nm.Usage.Memory().Value()) / (1024 * 1024)
					nodeUsage.ActualCPUCores = actualCPU
					nodeUsage.ActualMemoryMi = actualMem
					nodeUsage.ActualCPUPercent = safePercent(actualCPU, totalCPU)
					nodeUsage.ActualMemoryPercent = safePercent(actualMem, totalMem)
					nodeUsage.CPUPercent = nodeUsage.ActualCPUPercent
					nodeUsage.MemoryPercent = nodeUsage.ActualMemoryPercent
					break
				}
			}
		}

		if !useMetrics || (nodeUsage.ActualCPUCores == 0 && nodeUsage.ActualMemoryMi == 0) {
			nodeUsage.ActualCPUCores = allocCPU
			nodeUsage.ActualMemoryMi = allocMem
			nodeUsage.ActualCPUPercent = safePercent(allocCPU, totalCPU)
			nodeUsage.ActualMemoryPercent = safePercent(allocMem, totalMem)
			nodeUsage.CPUPercent = nodeUsage.ActualCPUPercent
			nodeUsage.MemoryPercent = nodeUsage.ActualMemoryPercent
		}

		usage.TotalCPUCores += totalCPU
		usage.TotalMemoryMi += totalMem
		usage.AllocatedCPUCores += allocCPU
		usage.AllocatedMemoryMi += allocMem
		usage.ActualCPUCores += nodeUsage.ActualCPUCores
		usage.ActualMemoryMi += nodeUsage.ActualMemoryMi
		usage.NodeDetails = append(usage.NodeDetails, nodeUsage)
	}

	usage.CPUPercent = safePercent(usage.ActualCPUCores, usage.TotalCPUCores)
	usage.MemoryPercent = safePercent(usage.ActualMemoryMi, usage.TotalMemoryMi)
	usage.ActualCPUPercent = usage.CPUPercent
	usage.ActualMemoryPercent = usage.MemoryPercent

	return usage, nil
}

func safePercent(used, total float64) float64 {
	if total == 0 {
		return 0
	}
	return used / total * 100
}

func ListEvents(cluster *model.Cluster, namespace string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetPodEvents(cluster *model.Cluster, namespace, podName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + podName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetDeploymentEvents(cluster *model.Cluster, namespace, deploymentName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + deploymentName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetServiceEvents(cluster *model.Cluster, namespace, serviceName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + serviceName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetIngressEvents(cluster *model.Cluster, namespace, ingressName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + ingressName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetConfigMapEvents(cluster *model.Cluster, namespace, configMapName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + configMapName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetSecretEvents(cluster *model.Cluster, namespace, secretName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + secretName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetStatefulSetEvents(cluster *model.Cluster, namespace, statefulSetName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + statefulSetName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetDaemonSetEvents(cluster *model.Cluster, namespace, daemonSetName string) ([]corev1.Event, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	list, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + daemonSetName,
	})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetResourceAsJSON(resource interface{}) ([]byte, error) {
	return json.MarshalIndent(resource, "", "  ")
}

type DeploymentHistoryItem struct {
	Revision         int64  `json:"revision"`
	CreationTimestamp string `json:"creationTimestamp"`
	Replicas         int32  `json:"replicas"`
	Image            string `json:"image"`
}

func GetDeploymentHistory(cluster *model.Cluster, namespace, name string) ([]DeploymentHistoryItem, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return nil, err
	}

	// 获取Deployment
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 获取该命名空间下的所有ReplicaSets
	rsList, err := client.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	history := make([]DeploymentHistoryItem, 0)
	for _, rs := range rsList.Items {
		// 检查ReplicaSet是否属于该Deployment
		isOwned := false
		for _, ownerRef := range rs.OwnerReferences {
			if ownerRef.UID == deployment.UID {
				isOwned = true
				break
			}
		}

		if !isOwned {
			continue
		}

		// 从annotation中获取revision
		revision := int64(0)
		if rs.Annotations != nil {
			if revStr, exists := rs.Annotations["deployment.kubernetes.io/revision"]; exists {
				if r, err := strconv.ParseInt(revStr, 10, 64); err == nil {
					revision = r
				}
			}
		}

		image := ""
		if len(rs.Spec.Template.Spec.Containers) > 0 {
			image = rs.Spec.Template.Spec.Containers[0].Image
		}

		history = append(history, DeploymentHistoryItem{
			Revision:          revision,
			CreationTimestamp: rs.CreationTimestamp.UTC().Format(time.RFC3339),
			Replicas:          rs.Status.Replicas,
			Image:             image,
		})
	}

	return history, nil
}