<template>
  <div class="terminal-wrapper" ref="terminalWrapper">
    <div class="terminal-container" ref="terminalContainer"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'

const props = defineProps({
  clusterId: { type: [Number, String], required: true },
  namespace: { type: String, required: true },
  podName: { type: String, required: true },
  container: { type: String, required: true }
})

const terminalWrapper = ref(null)
const terminalContainer = ref(null)

let term = null
let fitAddon = null
let ws = null

const connect = () => {
  const token = localStorage.getItem('token')
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const url = `${protocol}//${host}/api/k8s/${props.clusterId}/pods/${props.namespace}/${props.podName}/exec?container=${encodeURIComponent(props.container)}&token=${encodeURIComponent(token)}`

  ws = new WebSocket(url)
  ws.binaryType = 'arraybuffer'

  ws.onopen = () => {
    if (term && fitAddon) {
      fitAddon.fit()
      const dims = { cols: term.cols, rows: term.rows }
      ws.send(JSON.stringify(dims))
    }
  }

  ws.onmessage = (event) => {
    if (event.data instanceof ArrayBuffer) {
      const data = new Uint8Array(event.data)
      term.write(data)
    } else {
      term.write(event.data)
    }
  }

  ws.onclose = () => {
    term.write('\r\n\x1b[33m连接已断开\x1b[0m\r\n')
  }

  ws.onerror = () => {
    term.write('\r\n\x1b[31m连接错误\x1b[0m\r\n')
  }
}

const handleResize = () => {
  if (fitAddon) {
    fitAddon.fit()
    if (ws && ws.readyState === WebSocket.OPEN) {
      const dims = { cols: term.cols, rows: term.rows }
      ws.send(JSON.stringify(dims))
    }
  }
}

onMounted(() => {
  term = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: "'Cascadia Code', 'Fira Code', 'Courier New', monospace",
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
      cursor: '#d4d4d4'
    },
    allowProposedApi: true,
    allowTransparency: false
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.open(terminalContainer.value)

  term.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data)
    }
  })

  fitAddon.fit()
  connect()

  setTimeout(() => {
    if (fitAddon) fitAddon.fit()
  }, 200)

  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (ws) {
    ws.close()
    ws = null
  }
  if (term) {
    term.dispose()
    term = null
  }
})
</script>

<style scoped>
.terminal-wrapper {
  width: 100%;
  height: 60vh;
  min-height: 400px;
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
}

.terminal-container {
  width: 100%;
  height: 100%;
  padding: 8px;
}
</style>
