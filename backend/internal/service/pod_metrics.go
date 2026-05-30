package service

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s-ui-admin/internal/k8s"
	"k8s-ui-admin/internal/model"
	"k8s-ui-admin/internal/repository"
)

var stopCollector = make(chan struct{})

func StartMetricsCollector() {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		collectNow()

		for {
			select {
			case <-ticker.C:
				collectNow()
			case <-stopCollector:
				return
			}
		}
	}()
}

func StopMetricsCollector() {
	close(stopCollector)
}

func collectNow() {
	var clusters []model.Cluster
	err := repository.DB.Where("status = ?", "connected").Find(&clusters).Error
	if err != nil {
		log.Printf("metrics采集: 获取集群列表失败: %v", err)
		return
	}

	for _, cluster := range clusters {
		collectClusterMetrics(&cluster)
	}

	cutoff := time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	repository.DB.Where("collected_at < ?", cutoff).Delete(&model.PodMetrics{})
}

func collectClusterMetrics(cluster *model.Cluster) {
	metricsClient, err := k8s.GetMetricsClient(cluster)
	if err != nil {
		log.Printf("metrics采集: 创建Metrics客户端失败 cluster=%s: %v", cluster.Name, err)
		return
	}

	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("metrics采集: 获取Pod指标失败 cluster=%s: %v（确保集群已安装 metrics-server）", cluster.Name, err)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	var records []model.PodMetrics
	for _, pm := range podMetricsList.Items {
		for _, container := range pm.Containers {
			records = append(records, model.PodMetrics{
				CollectedAt:   now,
				ClusterID:     cluster.ID,
				PodName:       pm.Name,
				Namespace:     pm.Namespace,
				CPUMillicores: container.Usage.Cpu().MilliValue(),
				MemoryBytes:   container.Usage.Memory().Value(),
			})
		}
	}
	if len(records) > 0 {
		repository.DB.CreateInBatches(records, 100)
	}
}
