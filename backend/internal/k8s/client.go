package k8s

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8s-ui-admin/internal/model"
)

var clients = make(map[uint]*kubernetes.Clientset)

func GetClient(cluster *model.Cluster) (*kubernetes.Clientset, error) {
	if client, ok := clients[cluster.ID]; ok {
		return client, nil
	}

	config, err := buildConfig(cluster)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	clients[cluster.ID] = clientset
	return clientset, nil
}

func RemoveClient(clusterID uint) {
	delete(clients, clusterID)
}

var metricsClients = make(map[uint]*metricsclientset.Clientset)

func GetMetricsClient(cluster *model.Cluster) (*metricsclientset.Clientset, error) {
	if client, ok := metricsClients[cluster.ID]; ok {
		return client, nil
	}

	config, err := buildConfig(cluster)
	if err != nil {
		return nil, err
	}

	clientset, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	metricsClients[cluster.ID] = clientset
	return clientset, nil
}

func RemoveMetricsClient(clusterID uint) {
	delete(metricsClients, clusterID)
}

// Kubeconfig represents the structure of a kubeconfig file
type Kubeconfig struct {
	APIVersion     string            `yaml:"apiVersion"`
	Kind           string            `yaml:"kind"`
	Clusters       []ClusterEntry    `yaml:"clusters"`
	Contexts       []ContextEntry    `yaml:"contexts"`
	CurrentContext string            `yaml:"current-context"`
	Users         []UserEntry       `yaml:"users"`
}

type ClusterEntry struct {
	Name    string      `yaml:"name"`
	Cluster ClusterData `yaml:"cluster"`
}

type ClusterData struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

type ContextEntry struct {
	Name    string     `yaml:"name"`
	Context ContextData `yaml:"context"`
}

type ContextData struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type UserEntry struct {
	Name string    `yaml:"name"`
	User UserData `yaml:"user"`
}

type UserData struct {
	Token                 string `yaml:"token,omitempty"`
	Username             string `yaml:"username,omitempty"`
	Password             string `yaml:"password,omitempty"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData        string `yaml:"client-key-data,omitempty"`
}

func buildConfig(cluster *model.Cluster) (*rest.Config, error) {
	// 如果有 kubeconfig，优先使用 kubeconfig
	if cluster.Kubeconfig != "" {
		return buildConfigFromKubeconfig(cluster.Kubeconfig, cluster.ID)
	}

	// 否则使用 Token 或其他认证方式
	return buildConfigFromCredentials(cluster)
}

// buildConfigFromKubeconfig 从 kubeconfig 构建配置
func buildConfigFromKubeconfig(kubeconfigContent string, clusterID uint) (*rest.Config, error) {
	var kubeconfig Kubeconfig
	err := yaml.Unmarshal([]byte(kubeconfigContent), &kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("解析 kubeconfig 失败: %v", err)
	}

	// 查找当前上下文
	var currentContext ContextData
	found := false
	for _, ctx := range kubeconfig.Contexts {
		if ctx.Name == kubeconfig.CurrentContext {
			currentContext = ctx.Context
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("未找到当前上下文: %s", kubeconfig.CurrentContext)
	}

	// 查找集群信息
	var clusterData ClusterData
	for _, cluster := range kubeconfig.Clusters {
		if cluster.Name == currentContext.Cluster {
			clusterData = cluster.Cluster
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("未找到集群: %s", currentContext.Cluster)
	}

	// 查找用户信息
	var userData UserData
	for _, user := range kubeconfig.Users {
		if user.Name == currentContext.User {
			userData = user.User
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("未找到用户: %s", currentContext.User)
	}

	// 构建配置
	config := &rest.Config{
		Host: clusterData.Server,
	}

	// 处理 TLS 配置
	if clusterData.CertificateAuthorityData != "" {
		caData, err := base64.StdEncoding.DecodeString(clusterData.CertificateAuthorityData)
		if err != nil {
			return nil, fmt.Errorf("解码 CA 证书失败: %v", err)
		}
		config.TLSClientConfig.CAData = caData
	} else if clusterData.InsecureSkipTLSVerify {
		config.TLSClientConfig.Insecure = true
	}

	// 处理认证
	if userData.Token != "" {
		config.BearerToken = userData.Token
	} else if userData.ClientCertificateData != "" && userData.ClientKeyData != "" {
		certData, err := base64.StdEncoding.DecodeString(userData.ClientCertificateData)
		if err != nil {
			return nil, fmt.Errorf("解码客户端证书失败: %v", err)
		}
		keyData, err := base64.StdEncoding.DecodeString(userData.ClientKeyData)
		if err != nil {
			return nil, fmt.Errorf("解码客户端密钥失败: %v", err)
		}
		config.TLSClientConfig.CertData = certData
		config.TLSClientConfig.KeyData = keyData
	} else if userData.Username != "" && userData.Password != "" {
		config.Username = userData.Username
		config.Password = userData.Password
	} else {
		return nil, fmt.Errorf("kubeconfig 中未找到有效的认证信息")
	}

	return config, nil
}

// buildConfigFromCredentials 从凭据构建配置
func buildConfigFromCredentials(cluster *model.Cluster) (*rest.Config, error) {
	tlsClientConfig := rest.TLSClientConfig{}

	if cluster.CACert != "" {
		caCert, err := base64.StdEncoding.DecodeString(cluster.CACert)
		if err != nil {
			return nil, fmt.Errorf("解码 CA 证书失败: %v", err)
		}
		tlsClientConfig.CAData = caCert
	} else {
		tlsClientConfig.Insecure = true
	}

	config := &rest.Config{
		Host:            cluster.Server,
		TLSClientConfig: tlsClientConfig,
	}

	if cluster.Token != "" {
		config.BearerToken = cluster.Token
	}

	return config, nil
}

func TestConnection(cluster *model.Cluster) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	_, err = client.ServerVersion()
	if err != nil {
		RemoveClient(cluster.ID)
		return err
	}

	return nil
}

func GetVersion(cluster *model.Cluster) (string, error) {
	client, err := GetClient(cluster)
	if err != nil {
		return "", err
	}

	version, err := client.ServerVersion()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", version.GitVersion), nil
}