<template>
  <div class="resource-detail-page">
    <div v-if="loading" class="loading-container">
      <el-icon class="loading-icon" :size="40"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    
    <div v-else-if="error" class="error-container">
      <el-icon :size="60" color="#f56c6c"><WarnTriangleFilled /></el-icon>
      <h2>加载失败</h2>
      <p>{{ error }}</p>
      <el-button type="primary" @click="fetchData">重试</el-button>
    </div>
    
    <div v-else class="detail-content">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-left">
          <el-button @click="goBack" class="back-btn">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <div class="title-section">
            <div class="title-header">
              <div class="title-icon">
                <component :is="resourceIcon" :size="28" />
              </div>
              <div>
                <h1>{{ resourceName }}</h1>
                <div class="resource-meta">
                  <el-tag>{{ resourceNamespace }}</el-tag>
                  <el-tag :type="statusType">{{ statusText }}</el-tag>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="header-right">
          <el-button @click="handleRefresh" :loading="refreshing">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button type="primary" @click="showEditDialog = true">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-button type="danger" @click="handleDelete">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>
      
      <!-- 资源详情卡片 -->
      <div class="detail-cards">
        <!-- 基本信息 -->
        <el-card class="detail-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>基本信息</span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="名称">{{ resource?.metadata?.name }}</el-descriptions-item>
            <el-descriptions-item label="命名空间">{{ resource?.metadata?.namespace }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatTime(resource?.metadata?.creationTimestamp) }}</el-descriptions-item>
            <el-descriptions-item label="UID">{{ resource?.metadata?.uid }}</el-descriptions-item>
            <el-descriptions-item label="资源版本">{{ resource?.metadata?.resourceVersion }}</el-descriptions-item>
            <el-descriptions-item label="生成版本">{{ resource?.metadata?.generation }}</el-descriptions-item>
            
            <!-- Deployment专属信息 -->
            <template v-if="resourceType === 'deployments'">
              <el-descriptions-item label="副本数">
                <div class="replica-control">
                  <el-button size="small" @click="adjustReplicas(-1)" :disabled="resource?.spec?.replicas <= 0">-</el-button>
                  <span class="replica-number">{{ resource?.spec?.replicas || 0 }}</span>
                  <el-button size="small" @click="adjustReplicas(1)">+</el-button>
                </div>
              </el-descriptions-item>
              <el-descriptions-item label="当前副本">
                <span class="replica-number">{{ resource?.status?.replicas || 0 }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="就绪副本">
                <span class="replica-number ready">{{ resource?.status?.readyReplicas || 0 }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="可用副本">
                <span class="replica-number">{{ resource?.status?.availableReplicas || 0 }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="更新副本">
                <span class="replica-number">{{ resource?.status?.updatedReplicas || 0 }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="生成版本">
                <span class="generation-number">{{ resource?.metadata?.generation }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="就绪时间">
                {{ resource?.status?.conditions?.find(c => c.type === 'Available')?.lastUpdateTime ? formatTime(resource?.status?.conditions?.find(c => c.type === 'Available')?.lastUpdateTime) : '-' }}
              </el-descriptions-item>
            </template>
            
            <!-- StatefulSet专属信息 -->
            <template v-if="resourceType === 'statefulsets'">
              <el-descriptions-item label="副本数">
                {{ resource?.spec?.replicas || 0 }}</el-descriptions-item>
              <el-descriptions-item label="当前副本">
                {{ resource?.status?.replicas || 0 }}</el-descriptions-item>
              <el-descriptions-item label="就绪副本">
                {{ resource?.status?.readyReplicas || 0 }}</el-descriptions-item>
            </template>
            
            <!-- DaemonSet专属信息 -->
            <template v-if="resourceType === 'daemonsets'">
              <el-descriptions-item label="期望副本">
                {{ resource?.status?.desiredNumberScheduled || 0 }}</el-descriptions-item>
              <el-descriptions-item label="当前副本">
                {{ resource?.status?.currentNumberScheduled || 0 }}</el-descriptions-item>
              <el-descriptions-item label="就绪副本">
                {{ resource?.status?.numberReady || 0 }}</el-descriptions-item>
            </template>
            
            <!-- Pod专属信息 -->
            <template v-if="resourceType === 'pods'">
              <el-descriptions-item label="Pod IP">
                <code class="pod-ip-code">{{ resource?.status?.podIP || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item label="节点">
                <el-tag size="small">{{ resource?.spec?.nodeName || '-' }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="QoS 等级">
                <el-tag size="small" type="info">{{ resource?.status?.qosClass || '-' }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="服务账号">
                {{ resource?.spec?.serviceAccountName || 'default' }}
              </el-descriptions-item>
              <el-descriptions-item label="重启策略">
                {{ resource?.spec?.restartPolicy || '-' }}
              </el-descriptions-item>
            </template>
          </el-descriptions>
        </el-card>
        
        <!-- 标签和注解 -->
        <el-card class="detail-card">
          <template #header>
            <div class="card-header">
              <el-icon><PriceTag /></el-icon>
              <span>标签 & 注解</span>
            </div>
          </template>
          <div class="labels-section">
            <h4>标签</h4>
            <div class="label-list" v-if="Object.keys(resource?.metadata?.labels || {}).length > 0">
              <el-tag v-for="(value, key) in resource?.metadata?.labels" :key="key" class="label-tag">
                {{ key }}: {{ value }}
              </el-tag>
            </div>
            <p v-else class="empty-text">无标签</p>
            
            <h4>注解</h4>
            <div class="label-list" v-if="Object.keys(resource?.metadata?.annotations || {}).length > 0">
              <el-tag v-for="(value, key) in resource?.metadata?.annotations" :key="key" type="info" class="label-tag annotation">
                {{ key }}: {{ value }}
              </el-tag>
            </div>
            <p v-else class="empty-text">无注解</p>
          </div>
        </el-card>
        
        <!-- Pod 资源监控图表 -->
        <template v-if="resourceType === 'pods' && resource">
          <PodMetricsChart
            :clusterId="clusterId"
            :namespace="resource.metadata.namespace"
            :podName="resource.metadata.name"
          />
        </template>
        
        <!-- ConfigMap 数据内容 -->
        <el-card class="detail-card" v-if="resourceType === 'configmaps'">
          <template #header>
            <div class="card-header">
              <IconConfigMap :size="18" />
              <span>配置数据</span>
              <el-button size="small" type="primary" class="card-header-btn" @click="addDataItem">
                <el-icon><Plus /></el-icon>
                添加条目
              </el-button>
            </div>
          </template>
          <div v-if="configMapEntries.length === 0" class="empty-text">暂无配置数据</div>
          <el-table v-else :data="configMapEntries" border stripe size="small">
            <el-table-column prop="key" label="键" min-width="200">
              <template #default="scope">
                <span class="data-key">{{ scope.row.key }}</span>
              </template>
            </el-table-column>
            <el-table-column label="值" min-width="400">
              <template #default="scope">
                <div class="data-value-cell">
                  <pre class="data-value">{{ scope.row.value }}</pre>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="scope">
                <el-button size="small" type="primary" @click="editDataItem(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteDataItem(scope.row.key)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- Secret 数据内容 -->
        <el-card class="detail-card" v-if="resourceType === 'secrets'">
          <template #header>
            <div class="card-header">
              <IconSecret :size="18" />
              <span>密钥数据</span>
              <el-button size="small" type="primary" class="card-header-btn" @click="addDataItem">
                <el-icon><Plus /></el-icon>
                添加条目
              </el-button>
            </div>
          </template>
          <div v-if="secretEntries.length === 0" class="empty-text">暂无密钥数据</div>
          <el-table v-else :data="secretEntries" border stripe size="small">
            <el-table-column prop="key" label="键" min-width="180">
              <template #default="scope">
                <span class="data-key">{{ scope.row.key }}</span>
              </template>
            </el-table-column>
            <el-table-column label="值（已解密）" min-width="350">
              <template #default="scope">
                <div class="data-value-cell">
                  <pre class="data-value secret-value">{{ scope.row.decoded }}</pre>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="scope">
                <el-button size="small" type="primary" @click="editDataItem(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteDataItem(scope.row.key)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- Ingress 路由配置 -->
        <template v-if="resourceType === 'ingresses'">
          <el-card class="detail-card">
            <template #header>
              <div class="card-header">
                <IconIngress :size="18" />
                <span>IngressClass</span>
              </div>
            </template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="Ingress Class">
                <el-tag>{{ resource?.spec?.ingressClassName || '默认' }}</el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card class="detail-card">
            <template #header>
              <div class="card-header">
                <IconIngress :size="18" />
                <span>路由配置</span>
              </div>
            </template>
            <div v-if="!resource?.spec?.rules || resource.spec.rules.length === 0" class="empty-text">暂无路由规则</div>
            <template v-else>
              <div v-for="(rule, ruleIdx) in resource.spec.rules" :key="ruleIdx" class="routing-host-group">
                <div class="routing-host-header">
                  <el-icon><Link /></el-icon>
                  <span class="host-name">域名: {{ rule.host || '*' }}</span>
                </div>
                <el-table 
                  :data="(rule.http?.paths || []).map(p => ({ ...p, _ruleIdx: ruleIdx }))" 
                  border 
                  stripe 
                  size="small"
                  class="routing-table"
                  :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-primary)', textAlign: 'center' }"
                >
                  <el-table-column label="URL 匹配方式" width="150">
                    <template #default="scope">
                      <el-tag size="small" :type="scope.row.pathType === 'Prefix' ? 'primary' : scope.row.pathType === 'Exact' ? 'success' : 'info'">
                        {{ scope.row.pathType || '-' }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column label="URL" min-width="200">
                    <template #default="scope">
                      <code class="path-code">{{ scope.row.path || '/' }}</code>
                    </template>
                  </el-table-column>
                  <el-table-column label="服务名称" min-width="180">
                    <template #default="scope">
                      <span>{{ scope.row.backend?.service?.name || '-' }}</span>
                    </template>
                  </el-table-column>
                  <el-table-column label="服务端口" width="120">
                    <template #default="scope">
                      <el-tag size="small" type="warning">
                        {{ scope.row.backend?.service?.port?.number || scope.row.backend?.service?.port?.name || '-' }}
                      </el-tag>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </template>
          </el-card>
        </template>

        <!-- Service 端口信息 -->
        <template v-if="resourceType === 'services'">
          <el-card class="detail-card">
            <template #header>
              <div class="card-header">
                <IconService :size="18" />
              <span>端口信息</span>
              </div>
            </template>
            <el-table :data="servicePorts" border stripe size="small">
              <el-table-column label="端口名称" width="150">
                <template #default="scope">{{ scope.row.name || '-' }}</template>
              </el-table-column>
              <el-table-column label="协议" width="100">
                <template #default="scope">
                  <el-tag size="small">{{ scope.row.protocol || 'TCP' }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="端口" width="100">
                <template #default="scope">
                  <el-tag size="small" type="primary">{{ scope.row.port }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="目标端口" width="120">
                <template #default="scope">{{ scope.row.targetPort || '-' }}</template>
              </el-table-column>
              <el-table-column label="NodePort" width="120">
                <template #default="scope">
                  <el-tag v-if="scope.row.nodePort" size="small" type="warning">{{ scope.row.nodePort }}</el-tag>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </template>

        <!-- Service 关联 Pod -->
        <el-card class="detail-card" v-if="resourceType === 'services'">
          <template #header>
            <div class="card-header">
              <IconPod :size="18" />
              <span>关联 Pod</span>
            </div>
          </template>
          <div v-if="loadingRelated" class="loading-text">
            <el-icon class="is-loading"><Loading /></el-icon> 加载中...
          </div>
          <div v-else-if="relatedPods.length === 0" class="empty-text">暂无关联 Pod</div>
          <div v-else class="pod-table">
            <el-table :data="relatedPods" border stripe size="small">
              <el-table-column prop="metadata.name" label="名称" min-width="200">
                <template #default="scope">
                  <span class="clickable-name" @click="goToPodDetail(scope.row)">
                    {{ scope.row.metadata.name }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column label="Pod IP" width="150">
                <template #default="scope">
                  <code class="pod-ip">{{ scope.row.status?.podIP || '-' }}</code>
                </template>
              </el-table-column>
              <el-table-column prop="status.phase" label="状态" width="100">
                <template #default="scope">
                  <el-tag size="small" :type="getPodStatusType(scope.row)">
                    {{ scope.row.status?.phase || 'Unknown' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="启动时间" width="180">
                <template #default="scope">
                  {{ scope.row.status?.startTime ? formatTime(scope.row.status.startTime) : '-' }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100" fixed="right">
                <template #default="scope">
                  <el-button size="small" @click="goToPodDetail(scope.row)">详情</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-card>

        <!-- Deployment/StatefulSet/DaemonSet关联资源 -->
        <el-card class="detail-card" v-if="resourceType === 'deployments' || resourceType === 'statefulsets' || resourceType === 'daemonsets'">
          <template #header>
            <div class="card-header">
              <IconDeployment :size="18" />
              <span>关联资源</span>
            </div>
          </template>
          
          <el-tabs v-model="activeRelatedTab" class="related-tabs">
            <!-- Pod列表 -->
            <el-tab-pane label="Pod列表" name="pods">
              <div class="related-resources">
                <div class="resource-section">
                  <div v-if="loadingRelated" class="loading-text">
                    <el-icon class="is-loading"><Loading /></el-icon> 加载中...
                  </div>
                  <div v-else-if="relatedPods.length === 0" class="empty-text">
                    暂无关联 Pod
                  </div>
                  <div v-else class="pod-table">
                    <el-table :data="relatedPods" border stripe size="small">
                      <el-table-column prop="metadata.name" label="名称" min-width="200">
                        <template #default="scope">
                          <span class="clickable-name" @click="goToPodDetail(scope.row)">
                            {{ scope.row.metadata.name }}
                          </span>
                        </template>
                      </el-table-column>
                      <el-table-column prop="status.phase" label="状态" width="100">
                        <template #default="scope">
                          <el-tag size="small" :type="getPodStatusType(scope.row)">
                            {{ scope.row.status?.phase || 'Unknown' }}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="镜像" min-width="280" show-overflow-tooltip>
                        <template #default="scope">
                          <div v-if="scope.row.spec?.containers && scope.row.spec.containers.length > 0">
                            <div v-for="(container, idx) in scope.row.spec.containers" :key="idx" class="pod-container">
                              <span class="container-name">{{ container.name }}:</span>
                              <span class="container-image">{{ container.image }}</span>
                            </div>
                          </div>
                          <span v-else>-</span>
                        </template>
                      </el-table-column>
                      <el-table-column prop="status.ready" label="就绪" width="80">
                        <template #default="scope">
                          {{ scope.row.status?.conditions?.find(c => c.type === 'Ready')?.status || 'Unknown' }}
                        </template>
                      </el-table-column>
                      <el-table-column prop="status.startTime" label="启动时间" width="180">
                        <template #default="scope">
                          {{ scope.row.status?.startTime ? formatTime(scope.row.status.startTime) : '-' }}
                        </template>
                      </el-table-column>
                      <el-table-column label="操作" width="120" fixed="right">
                        <template #default="scope">
                          <el-button size="small" @click="goToPodDetail(scope.row)">详情</el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>
                </div>
              </div>
            </el-tab-pane>
            
            <!-- Init容器 -->
            <el-tab-pane label="Init容器" name="initContainers">
              <div class="container-info" v-if="resource?.spec?.template?.spec?.initContainers?.length > 0">
                <div v-for="(container, index) in resource?.spec?.template?.spec?.initContainers" :key="index" class="container-item">
                  <div class="container-header">
                    <h4>{{ container.name }}</h4>
                    <el-tag size="small" type="info">Init容器</el-tag>
                  </div>
                  <div class="container-details">
                    <div class="detail-row">
                      <span class="label">镜像:</span>
                      <span class="value">{{ container.image || '-' }}</span>
                    </div>
                    <div class="detail-row" v-if="container.command && Array.isArray(container.command)">
                      <span class="label">命令:</span>
                      <span class="value">{{ container.command.join(' ') }}</span>
                    </div>
                    <div class="detail-row" v-if="container.args && Array.isArray(container.args)">
                      <span class="label">参数:</span>
                      <span class="value">{{ container.args.join(' ') }}</span>
                    </div>
                  </div>
                  <div class="container-actions">
                    <el-button size="small" @click="viewInitContainerLogs(container.name)">查看日志</el-button>
                  </div>
                </div>
              </div>
              <div v-else class="empty-text">
                暂无 Init 容器
              </div>
            </el-tab-pane>
            
            <!-- 历史版本 -->
            <el-tab-pane v-if="resourceType === 'deployments'" label="历史版本" name="history">
              <div class="resource-section">
                <div v-if="loadingHistory" class="loading-text">
                  <el-icon class="is-loading"><Loading /></el-icon> 加载中...
                </div>
                <div v-else-if="historyVersions.length === 0" class="empty-text">
                  暂无历史版本
                </div>
                <div v-else class="history-table">
                  <el-table 
                    :data="historyVersions" 
                    border 
                    stripe 
                    size="small"
                    style="width: 100%"
                  >
                    <el-table-column prop="revision" label="版本号" width="100" sortable></el-table-column>
                    <el-table-column prop="creationTimestamp" label="创建时间" width="180">
                      <template #default="scope">
                        {{ formatTime(scope.row.creationTimestamp) }}
                      </template>
                    </el-table-column>
                    <el-table-column prop="replicas" label="副本数" width="100"></el-table-column>
                    <el-table-column prop="image" label="镜像" show-overflow-tooltip></el-table-column>
                  </el-table>
                </div>
              </div>
            </el-tab-pane>
            
            <!-- 关联Service -->
            <el-tab-pane label="关联Service" name="services">
              <div class="resource-section">
                <div v-if="loadingServices" class="loading-text">
                  <el-icon class="is-loading"><Loading /></el-icon> 加载中...
                </div>
                <div v-else-if="relatedServices.length === 0" class="empty-text">
                  暂无关联 Service
                </div>
                <div v-else class="service-table">
                  <el-table :data="relatedServices" border stripe size="small">
                    <el-table-column prop="metadata.name" label="名称" min-width="150">
                      <template #default="scope">
                        <span class="clickable-name" @click="goToServiceDetail(scope.row)">
                          {{ scope.row.metadata.name }}
                        </span>
                      </template>
                    </el-table-column>
                    <el-table-column prop="spec.type" label="类型" width="100" />
                    <el-table-column prop="spec.clusterIP" label="ClusterIP" width="120" />
                    <el-table-column label="端口" min-width="200">
                      <template #default="scope">
                        <div v-if="scope.row.spec?.ports && Array.isArray(scope.row.spec.ports)">
                          <div v-for="(port, idx) in scope.row.spec.ports" :key="idx" class="port-info">
                            <span>{{ port.port }}:{{ port.targetPort }}/{{ port.protocol }}</span>
                          </div>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="选择器" min-width="150">
                      <template #default="scope">
                        <div v-if="scope.row.spec?.selector" class="selector-tags">
                          <el-tag v-for="(value, key) in scope.row.spec.selector" :key="key" size="small">
                            {{ key }}: {{ value }}
                          </el-tag>
                        </div>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </div>
    </div>
    
    <!-- 编辑对话框 -->
    <el-dialog 
      title="编辑资源" 
      v-model="showEditDialog" 
      width="80%" 
      :close-on-click-modal="false"
      class="edit-dialog"
    >
      <div class="editor-container">
        <el-input
          v-model="editYaml"
          type="textarea"
          :rows="30"
          class="yaml-editor"
          placeholder="请输入 YAML 内容"
        />
      </div>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="updating">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- Init容器日志对话框 -->
    <el-dialog 
      :title="`Init容器日志 - ${selectedInitContainer}`" 
      v-model="showInitContainerLogs" 
      width="80%"
      class="logs-dialog"
    >
      <div class="logs-container">
        <div v-if="loadingInitContainerLogs" class="loading-text">
          <el-icon class="is-loading"><Loading /></el-icon> 加载中...
        </div>
        <div v-else-if="initContainerLogs" class="logs-content">
          <pre>{{ initContainerLogs }}</pre>
        </div>
        <div v-else class="empty-text">
          暂无日志
        </div>
      </div>
      <template #footer>
        <el-button @click="showInitContainerLogs = false">关闭</el-button>
        <el-button type="primary" @click="fetchInitContainerLogs">刷新</el-button>
      </template>
    </el-dialog>
    
    <!-- 数据条目编辑对话框 -->
    <el-dialog 
      :title="editingExistingKey ? '编辑数据条目' : '添加数据条目'"
      v-model="showDataEditDialog" 
      width="650px"
      :close-on-click-modal="false"
    >
      <el-form label-width="60px" @submit.prevent="saveDataItem">
        <el-form-item label="键">
          <el-input v-model="editDataKey" placeholder="请输入键名" :disabled="editingExistingKey" />
        </el-form-item>
        <el-form-item label="值">
          <el-input 
            v-model="editDataValue" 
            type="textarea" 
            :rows="8"
            placeholder="请输入值内容"
          />
        </el-form-item>
        <el-form-item v-if="resourceType === 'secrets'" label="提示">
          <span class="form-tip">Secret 值将自动进行 Base64 编码存储</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDataEditDialog = false">取消</el-button>
        <el-button type="primary" @click="saveDataItem" :loading="savingData">保存</el-button>
      </template>
    </el-dialog>
    

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { k8sAPI, clusterAPI } from '@/utils/api'
import { getK8sIcon } from '@/assets/icons/index.js'
import IconService from '@/assets/icons/IconService.vue'
import IconPod from '@/assets/icons/IconPod.vue'
import IconDeployment from '@/assets/icons/IconDeployment.vue'
import IconIngress from '@/assets/icons/IconIngress.vue'
import IconSecret from '@/assets/icons/IconSecret.vue'
import IconConfigMap from '@/assets/icons/IconConfigMap.vue'
import PodMetricsChart from '@/components/PodMetricsChart.vue'
import jsYaml from 'js-yaml'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const refreshing = ref(false)
const updating = ref(false)
const error = ref('')
const resource = ref(null)
const showEditDialog = ref(false)
const editYaml = ref('')
const relatedPods = ref([])
const relatedReplicaSets = ref([])
const loadingRelated = ref(false)
const relatedServices = ref([])
const loadingServices = ref(false)
const activeRelatedTab = ref('pods')
const showInitContainerLogs = ref(false)
const selectedInitContainer = ref('')
const initContainerLogs = ref('')
const loadingInitContainerLogs = ref(false)
const historyVersions = ref([])
const loadingHistory = ref(false)
const rollbackRevision = ref('')
const showDataEditDialog = ref(false)
const editDataKey = ref('')
const editDataValue = ref('')
const editingExistingKey = ref(false)
const savingData = ref(false)
let refreshInterval = null

const clusterId = computed(() => route.params.id)

// 根据路由名称推断资源类型
const resourceType = computed(() => {
  if (route.params.type) return route.params.type
  
  const routeName = route.name
  switch (routeName) {
    case 'DeploymentDetail': return 'deployments'
    case 'PodDetail': return 'pods'
    case 'ServiceDetail': return 'services'
    case 'IngressDetail': return 'ingresses'
    case 'ConfigMapDetail': return 'configmaps'
    case 'SecretDetail': return 'secrets'
    case 'StatefulSetDetail': return 'statefulsets'
    case 'DaemonSetDetail': return 'daemonsets'
    default: return null
  }
})

const resourceNamespace = computed(() => route.params.namespace)
const resourceName = computed(() => route.params.name)

// 获取资源类型的中文名称
const resourceTypeLabel = computed(() => {
  const labels = {
    'deployments': 'Deployment',
    'pods': 'Pod',
    'services': 'Service',
    'ingresses': 'Ingress',
    'configmaps': 'ConfigMap',
    'secrets': 'Secret',
    'statefulsets': 'StatefulSet',
    'daemonsets': 'DaemonSet'
  }
  return labels[resourceType.value] || '资源'
})

const statusType = computed(() => {
  const status = resource.value?.status?.phase || resource.value?.status?.conditions?.[0]?.type
  if (status === 'Ready' || status === 'Running') return 'success'
  if (status === 'Pending') return 'warning'
  if (status === 'Failed') return 'danger'
  return 'info'
})

const statusText = computed(() => {
  const status = resource.value?.status?.phase || resource.value?.status?.conditions?.[0]?.type
  return status || 'Unknown'
})

const resourceIcon = computed(() => {
  return getK8sIcon(resourceType.value)
})

const servicePorts = computed(() => {
  if (!resource.value?.spec?.ports) return []
  return resource.value.spec.ports
})

const yamlContent = computed(() => {
  if (!resource.value) return ''
  try {
    return jsYaml.dump(resource.value)
  } catch (e) {
    return JSON.stringify(resource.value, null, 2)
  }
})

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString('zh-CN')
}

const configMapEntries = computed(() => {
  if (!resource.value?.data) return []
  return Object.entries(resource.value.data).map(([key, value]) => ({
    key,
    value
  }))
})

const secretEntries = computed(() => {
  if (!resource.value?.data) return []
  return Object.entries(resource.value.data).map(([key, value]) => {
    let decoded = ''
    try {
      decoded = atob(value)
    } catch {
      decoded = value
    }
    return { key, value, decoded }
  })
})

const base64Encode = (str) => {
  try {
    return btoa(unescape(encodeURIComponent(str)))
  } catch {
    return btoa(str)
  }
}

const addDataItem = () => {
  editingExistingKey.value = false
  editDataKey.value = ''
  editDataValue.value = ''
  showDataEditDialog.value = true
}

const editDataItem = (item) => {
  editingExistingKey.value = true
  editDataKey.value = item.key
  if (resourceType.value === 'secrets') {
    editDataValue.value = item.decoded
  } else {
    editDataValue.value = item.value
  }
  showDataEditDialog.value = true
}

const saveDataItem = async () => {
  if (!editDataKey.value.trim()) {
    ElMessage.error('键名不能为空')
    return
  }

  savingData.value = true
  try {
    const id = clusterId.value
    const ns = resourceNamespace.value
    const name = resourceName.value
    const dataValue = resourceType.value === 'secrets' 
      ? base64Encode(editDataValue.value) 
      : editDataValue.value

    const updatedResource = JSON.parse(JSON.stringify(resource.value))
    if (!updatedResource.data) {
      updatedResource.data = {}
    }
    updatedResource.data[editDataKey.value] = dataValue

    if (resourceType.value === 'configmaps') {
      await k8sAPI.updateConfigMap(id, ns, name, updatedResource)
    } else if (resourceType.value === 'secrets') {
      await k8sAPI.updateSecret(id, ns, name, updatedResource)
    }

    ElMessage.success(editingExistingKey.value ? '更新成功' : '添加成功')
    showDataEditDialog.value = false
    await fetchData()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || err.message || '保存失败')
  } finally {
    savingData.value = false
  }
}

const deleteDataItem = async (key) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除数据条目 "${key}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const id = clusterId.value
    const ns = resourceNamespace.value
    const name = resourceName.value

    const updatedResource = JSON.parse(JSON.stringify(resource.value))
    delete updatedResource.data[key]

    if (resourceType.value === 'configmaps') {
      await k8sAPI.updateConfigMap(id, ns, name, updatedResource)
    } else if (resourceType.value === 'secrets') {
      await k8sAPI.updateSecret(id, ns, name, updatedResource)
    }

    ElMessage.success('删除成功')
    await fetchData()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(err.response?.data?.error || err.message || '删除失败')
    }
  }
}

const flattenStatus = (status) => {
  if (!status) return {}
  const result = {}
  
  const flatten = (obj, prefix = '') => {
    for (const [key, value] of Object.entries(obj)) {
      if (value === null || value === undefined) continue
      if (typeof value === 'object' && !Array.isArray(value)) {
        flatten(value, prefix + key + '.')
      } else if (Array.isArray(value)) {
        result[prefix + key] = `[${value.length} items]`
      } else {
        result[prefix + key] = String(value)
      }
    }
  }
  
  flatten(status)
  return result
}

const fetchData = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const id = clusterId.value
    const type = resourceType.value
    const ns = resourceNamespace.value
    const name = resourceName.value
    
    if (!id || !type || !ns || !name) {
      throw new Error('缺少必要的参数')
    }
    
    let response
    switch (type) {
      case 'deployments':
        response = await k8sAPI.getDeployment(id, ns, name)
        break
      case 'pods':
        response = await k8sAPI.getPod(id, ns, name)
        break
      case 'services':
        response = await k8sAPI.getService(id, ns, name)
        break
      case 'configmaps':
        response = await k8sAPI.getConfigMap(id, ns, name)
        break
      case 'secrets':
        response = await k8sAPI.getSecret(id, ns, name)
        break
      case 'ingresses':
        response = await k8sAPI.getIngress(id, ns, name)
        break
      case 'statefulsets':
        response = await k8sAPI.getStatefulSet(id, ns, name)
        break
      case 'daemonsets':
        response = await k8sAPI.getDaemonSet(id, ns, name)
        break
      default:
        throw new Error(`未知的资源类型: ${type}`)
    }
    
    resource.value = response.data
    editYaml.value = yamlContent.value
    
    // 如果是Deployment，加载关联资源
    if (type === 'deployments' || type === 'statefulsets' || type === 'daemonsets' || type === 'services') {
      fetchRelatedResources()
    }
  } catch (err) {
    error.value = err.response?.data?.error || err.message || '加载失败'
  } finally {
    loading.value = false
  }
}

// 获取Deployment关联的资源
const fetchRelatedResources = async () => {
  loadingRelated.value = true
  loadingServices.value = true
  
  try {
    const id = clusterId.value
    const ns = resourceNamespace.value
    
    // 获取该命名空间下的所有Pod
    const podsResponse = await k8sAPI.listPods(id, ns)
    const allPods = podsResponse.data || []
    
    // 根据selector匹配Pod
    if (resource.value?.spec?.selector?.matchLabels) {
      const matchLabels = resource.value.spec.selector.matchLabels
      relatedPods.value = allPods.filter(pod => {
        const podLabels = pod.metadata?.labels || {}
        return Object.entries(matchLabels).every(
          ([key, value]) => podLabels[key] === value
        )
      })
    } else {
      relatedPods.value = allPods
    }
    
    // 获取该命名空间下的所有Service
    const servicesResponse = await k8sAPI.listServices(id, ns)
    const allServices = servicesResponse.data || []
    
    // 根据selector匹配Service
    if (resource.value?.spec?.selector?.matchLabels) {
      const matchLabels = resource.value.spec.selector.matchLabels
      relatedServices.value = allServices.filter(service => {
        const serviceSelector = service.spec?.selector || {}
        return Object.entries(matchLabels).every(
          ([key, value]) => serviceSelector[key] === value
        )
      })
    } else {
      relatedServices.value = []
    }
    
    // 这里暂时只显示Pod和Service
    relatedReplicaSets.value = []
    
  } catch (err) {
    console.error('加载关联资源失败:', err)
    relatedPods.value = []
    relatedServices.value = []
    relatedReplicaSets.value = []
  } finally {
    loadingRelated.value = false
    loadingServices.value = false
  }
}

// 查看Init容器日志
const viewInitContainerLogs = async (containerName) => {
  selectedInitContainer.value = containerName
  showInitContainerLogs.value = true
  await fetchInitContainerLogs()
}

// 获取Init容器日志
const fetchInitContainerLogs = async () => {
  if (!selectedInitContainer.value) return
  
  loadingInitContainerLogs.value = true
  initContainerLogs.value = ''
  
  try {
    const id = clusterId.value
    const ns = resourceNamespace.value
    
    // 由于Init容器日志获取可能需要特定API，这里暂时显示提示信息
    // 实际实现可能需要Pod日志API
    initContainerLogs.value = '提示: Init容器日志获取可能需要在Pod详情页查看\n\nInit容器: ' + selectedInitContainer.value
    
  } catch (err) {
    console.error('获取Init容器日志失败:', err)
    initContainerLogs.value = '获取日志失败: ' + (err.message || '未知错误')
  } finally {
    loadingInitContainerLogs.value = false
  }
}

// 跳转到Pod详情页
const goToPodDetail = (pod) => {
  router.push(`/clusters/${clusterId.value}/pods/${pod.metadata.namespace}/${pod.metadata.name}`)
}

// 跳转到Service详情页
const goToServiceDetail = (service) => {
  router.push(`/clusters/${clusterId.value}/services/${service.metadata.namespace}/${service.metadata.name}`)
}

// 调整副本数
const adjustReplicas = async (delta) => {
  if (!resource.value) return
  
  const newReplicas = Math.max(0, (resource.value.spec.replicas || 0) + delta)
  try {
    const id = clusterId.value
    await k8sAPI.scaleDeployment(
      id, 
      resource.value.metadata.namespace, 
      resource.value.metadata.name, 
      newReplicas
    )
    await fetchData()
    ElMessage.success('副本数更新成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || error.message || '更新失败')
  }
}



// 启动实时更新 - 只在需要时刷新数据，不强制每5秒刷新
const startLiveUpdate = () => {
  // 移除自动刷新，使用用户主动刷新按钮
  // 如果需要更高效的响应式更新，未来可以考虑WebSocket等方案
}

// 停止实时更新
const stopLiveUpdate = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
}

const handleRefresh = async () => {
  refreshing.value = true
  await fetchData()
  refreshing.value = false
}

const handleUpdate = async () => {
  updating.value = true
  
  try {
    const updatedResource = jsYaml.load(editYaml.value)
    
    const id = clusterId.value
    const type = resourceType.value
    const ns = resourceNamespace.value
    const name = resourceName.value
    
    switch (type) {
      case 'deployments':
        await k8sAPI.updateDeployment(id, ns, name, updatedResource)
        break
      case 'statefulsets':
        await k8sAPI.updateStatefulSet(id, ns, name, updatedResource)
        break
      case 'daemonsets':
        await k8sAPI.updateDaemonSet(id, ns, name, updatedResource)
        break
      case 'services':
        await k8sAPI.updateService(id, ns, name, updatedResource)
        break
      case 'configmaps':
        await k8sAPI.updateConfigMap(id, ns, name, updatedResource)
        break
      case 'ingresses':
        await k8sAPI.updateIngress(id, ns, name, updatedResource)
        break
      case 'secrets':
        await k8sAPI.updateSecret(id, ns, name, updatedResource)
        break
      default:
        throw new Error('该资源类型不支持编辑')
    }
    
    ElMessage.success('更新成功')
    showEditDialog.value = false
    await fetchData()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || err.message || '更新失败')
  } finally {
    updating.value = false
  }
}

const handleDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${resourceName.value} 吗？此操作不可恢复。`,
      '警告',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const id = clusterId.value
    const type = resourceType.value
    const ns = resourceNamespace.value
    const name = resourceName.value
    
    switch (type) {
      case 'deployments':
        await k8sAPI.deleteDeployment(id, ns, name)
        break
      case 'pods':
        await k8sAPI.deletePod(id, ns, name)
        break
      case 'services':
        await k8sAPI.deleteService(id, ns, name)
        break
      case 'configmaps':
        await k8sAPI.deleteConfigMap(id, ns, name)
        break
      case 'secrets':
        await k8sAPI.deleteSecret(id, ns, name)
        break
      case 'ingresses':
        await k8sAPI.deleteIngress(id, ns, name)
        break
      case 'statefulsets':
        await k8sAPI.deleteStatefulSet(id, ns, name)
        break
      case 'daemonsets':
        await k8sAPI.deleteDaemonSet(id, ns, name)
        break
      default:
        throw new Error('未知的资源类型')
    }
    
    ElMessage.success('删除成功')
    router.back()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(err.response?.data?.error || err.message || '删除失败')
    }
  }
}

// 获取Pod状态标签类型
const getPodStatusType = (pod) => {
  const phase = pod.status?.phase
  if (phase === 'Running') return 'success'
  if (phase === 'Pending') return 'warning'
  if (phase === 'Failed') return 'danger'
  return 'info'
}

const copyYaml = async () => {
  try {
    await navigator.clipboard.writeText(yamlContent.value)
    ElMessage.success('已复制到剪贴板')
  } catch (err) {
    ElMessage.error('复制失败')
  }
}

// 处理历史版本选择
const handleHistorySelection = (selection) => {
  if (selection.length > 0) {
    rollbackRevision.value = selection[0].revision
  } else {
    rollbackRevision.value = ''
  }
}

const goBack = () => {
  router.back()
}



onMounted(() => {
  fetchData()
  startLiveUpdate()
})

onUnmounted(() => {
  stopLiveUpdate()
})
</script>

<style scoped>
.resource-detail-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px;
}

.loading-container,
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 20px;
  color: var(--text-tertiary);
}

.loading-icon {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.error-container h2 {
  margin: 20px 0 10px;
  color: var(--text-primary);
}

.error-container p {
  margin-bottom: 20px;
  color: var(--text-secondary);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px;
  background: var(--bg-card);
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(0, 0, 0, 0.2);
  color: black;
  display: flex;
  align-items: center;
  gap: 4px;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 1);
}

/* 关联资源样式 */
.related-resources {
  padding: 8px 0;
}

.resource-section h4 {
  margin: 0 0 12px 0;
  color: var(--text-primary);
  font-size: 15px;
}

.loading-text {
  color: var(--text-tertiary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.pod-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pod-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.pod-item:hover {
  background: var(--bg-card);
  border-color: var(--text-secondary);
  transform: translateX(4px);
}

.pod-info {
  flex: 1;
}

.pod-name {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.pod-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pod-namespace {
  color: var(--text-tertiary);
  font-size: 13px;
}

.arrow-icon {
  color: var(--text-tertiary);
  font-size: 16px;
}

.title-section h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.resource-meta {
  display: flex;
  gap: 8px;
}

.header-right {
  display: flex;
  gap: 12px;
}

.title-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #e8f4fd;
  color: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.detail-cards {
  display: grid;
  gap: 20px;
  margin-bottom: 20px;
}

.detail-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.copy-btn {
  margin-left: auto;
}

.labels-section h4 {
  margin: 16px 0 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.labels-section h4:first-child {
  margin-top: 0;
}

.label-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.label-tag {
  margin: 4px;
}

.label-tag.annotation {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-text {
  color: var(--text-tertiary);
  font-size: 14px;
}

.replica-control {
  display: flex;
  align-items: center;
  gap: 12px;
}

.replica-number {
  font-size: 18px;
  font-weight: 600;
  min-width: 30px;
  text-align: center;
}

.replica-number.ready {
  color: #67c23a;
}

.generation-control {
  display: flex;
  align-items: center;
  gap: 12px;
}

.generation-number {
  font-size: 16px;
  font-weight: 600;
  font-family: 'Courier New', monospace;
}

.history-list {
  margin-top: 12px;
}

.rollback-container {
  padding: 8px 0;
}

.loading-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--text-tertiary);
  padding: 40px 0;
}

.yaml-card {
  border-radius: 12px;
}

.yaml-content {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 20px;
  border-radius: 8px;
  overflow-x: auto;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  max-height: 500px;
  overflow-y: auto;
  white-space: pre-wrap;
}

.editor-container {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
}

.yaml-editor :deep(textarea) {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  background: #1e1e1e;
  color: #d4d4d4;
  border: none;
}

.related-tabs {
  margin-top: 8px;
}

.pod-table {
  margin-top: 12px;
}

.service-table {
  margin-top: 12px;
}

.port-info {
  margin-bottom: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.selector-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.pod-container {
  display: flex;
  gap: 6px;
  margin-bottom: 4px;
  font-size: 12px;
}

.pod-container:last-child {
  margin-bottom: 0;
}

.container-name {
  font-weight: 500;
  color: var(--text-secondary);
  white-space: nowrap;
}

.container-image {
  color: var(--text-primary);
  font-family: 'Courier New', monospace;
  word-break: break-all;
}

.clickable-name {
  cursor: pointer;
  color: #409eff;
  transition: color 0.2s;
}

.clickable-name:hover {
  color: #66b1ff;
  text-decoration: underline;
}

.container-info {
  margin-top: 12px;
}

.container-item {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
  background: var(--bg-secondary);
}

.container-item:last-child {
  margin-bottom: 0;
}

.container-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.container-header h4 {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
}

.container-details {
  margin-bottom: 12px;
}

.detail-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-row .label {
  font-weight: 500;
  color: var(--text-secondary);
  min-width: 60px;
}

.detail-row .value {
  color: var(--text-primary);
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.container-actions {
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.logs-container {
  background: #1e1e1e;
  border-radius: 8px;
  padding: 16px;
  max-height: 500px;
  overflow-y: auto;
}

.logs-content pre {
  margin: 0;
  color: #d4d4d4;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
}

.card-header-btn {
  margin-left: auto;
}

.data-key {
  font-weight: 600;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: var(--text-primary);
}

.data-value-cell {
  max-height: 200px;
  overflow-y: auto;
}

.data-value {
  margin: 0;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-radius: 6px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: var(--text-primary);
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.5;
}

.secret-value {
  background: #1e1e1e;
  color: #d4d4d4;
  border: 1px solid var(--border-color);
}

.form-tip {
  color: var(--text-tertiary);
  font-size: 12px;
}

.routing-host-group {
  margin-bottom: 20px;
}

.routing-host-group:last-child {
  margin-bottom: 0;
}

.routing-host-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  padding: 10px 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.host-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--text-primary);
  font-family: 'Courier New', monospace;
}

.routing-table {
  margin-top: 0;
}

.path-code {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: var(--color-primary);
  background: var(--bg-secondary);
  padding: 2px 8px;
  border-radius: 4px;
}

.pod-ip {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  font-weight: 600;
  color: #67c23a;
  background: #f0f9eb;
  padding: 2px 8px;
  border-radius: 4px;
}

.pod-ip-code {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  font-weight: 600;
  color: #67c23a;
}
</style>
