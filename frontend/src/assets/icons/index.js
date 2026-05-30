import IconCluster from './IconCluster.vue'
import IconConfigMap from './IconConfigMap.vue'
import IconContainer from './IconContainer.vue'
import IconDaemonSet from './IconDaemonSet.vue'
import IconDeployment from './IconDeployment.vue'
import IconIngress from './IconIngress.vue'
import IconNamespace from './IconNamespace.vue'
import IconNode from './IconNode.vue'
import IconPod from './IconPod.vue'
import IconSecret from './IconSecret.vue'
import IconService from './IconService.vue'
import IconStatefulSet from './IconStatefulSet.vue'

export const k8sIcons = {
  cluster: IconCluster,
  configmaps: IconConfigMap,
  container: IconContainer,
  daemonsets: IconDaemonSet,
  deployments: IconDeployment,
  ingresses: IconIngress,
  namespace: IconNamespace,
  node: IconNode,
  pods: IconPod,
  secrets: IconSecret,
  services: IconService,
  statefulsets: IconStatefulSet,
}

export function getK8sIcon(type) {
  return k8sIcons[type] || IconPod
}
