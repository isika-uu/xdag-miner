<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { StartMining, StopMining, GetMinerStatus, GetSystemInfo, LoadConfig } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import Toast from './Toast.vue'

const status = ref({
  running: false,
  hashrate: 0,
  threads: 0,
  uptime: 0,
  pool: '',
  connected: false
})

const systemInfo = ref({
  os: '',
  arch: '',
  cpuCores: 0,
  cpuModel: '',
  xmrigVersion: ''
})

const loading = ref(false)
const toast = ref({
  show: false,
  type: 'info',
  message: ''
})

const config = ref(null)

// 显示提示
const showToast = (type, message) => {
  toast.value = { show: true, type, message }
}

// 格式化哈希率
const formatHashrate = (hashrate) => {
  if (hashrate >= 1000000) {
    return (hashrate / 1000000).toFixed(2) + ' MH/s'
  } else if (hashrate >= 1000) {
    return (hashrate / 1000).toFixed(2) + ' KH/s'
  }
  return hashrate.toFixed(2) + ' H/s'
}

// 格式化运行时间
const formatUptime = (seconds) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  return `${hours}时${minutes}分${secs}秒`
}

// 刷新状态
const refreshStatus = async () => {
  try {
    const result = await GetMinerStatus()
    if (result) {
      status.value = result
    }
  } catch (err) {
    console.error('获取状态失败:', err)
  }
}

// 加载系统信息
const loadSystemInfo = async () => {
  try {
    const info = await GetSystemInfo()
    if (info) {
      systemInfo.value = info
    }
  } catch (err) {
    console.error('获取系统信息失败:', err)
  }
}

const loadConfig = async () => {
  try {
    const cfg = await LoadConfig()
    if (cfg) {
      config.value = cfg
    }
  } catch (err) {
    console.error('获取配置失败:', err)
  }
}

// 开始挖矿
const startMining = async () => {
  loading.value = true
  try {
    await StartMining()
    await refreshStatus()
    showToast('success', '挖矿已启动')
  } catch (err) {
    showToast('error', '启动失败: ' + err.toString())
  } finally {
    loading.value = false
  }
}

// 停止挖矿
const stopMining = async () => {
  loading.value = true
  try {
    await StopMining()
    await refreshStatus()
    showToast('success', '挖矿已停止')
  } catch (err) {
    showToast('error', '停止失败: ' + err.toString())
  } finally {
    loading.value = false
  }
}

let statusInterval = null

onMounted(() => {
  loadSystemInfo()
  loadConfig()
  refreshStatus()
  
  // 定期刷新状态
  statusInterval = setInterval(refreshStatus, 2000)
  
  // 监听挖矿停止事件
  EventsOn('miner:stopped', () => {
    refreshStatus()
  })
})

onUnmounted(() => {
  if (statusInterval) {
    clearInterval(statusInterval)
  }
  EventsOff('miner:stopped')
})
</script>

<template>
  <div class="dashboard">
    <!-- Toast 提示 -->
    <Toast
      :show="toast.show"
      :type="toast.type"
      :message="toast.message"
      @close="toast.show = false"
    />

    <!-- 状态卡片 -->
    <div class="cards-grid">
      <!-- 挖矿状态 -->
      <div class="card status-card">
        <div class="card-header">
          <h3>🚀 挖矿状态</h3>
          <span :class="['status-badge', status.running ? 'running' : 'stopped']">
            {{ status.running ? '运行中' : '已停止' }}
          </span>
        </div>
        <div class="card-body">
          <div class="stat-item">
            <span class="label">算力:</span>
            <span :class="['value', status.running && status.hashrate === 0 ? 'warning' : '']">
              {{ formatHashrate(status.hashrate) }}
            </span>
          </div>
          <div class="stat-item">
            <span class="label">最大线程(%):</span>
            <span class="value">{{ (config && config.cpu && config.cpu['max-threads-hint']) ? (config.cpu['max-threads-hint'] + '%') : '-' }}</span>
          </div>
          <div class="stat-item">
            <span class="label">运行时间:</span>
            <span class="value">{{ formatUptime(status.uptime) }}</span>
          </div>
          <div class="stat-item">
            <span class="label">矿池:</span>
            <span class="value small">{{ status.pool || '未配置' }}</span>
          </div>
        </div>
        
        <!-- 健康状态提示 -->
        <div v-if="status.running" class="health-status">
          <div v-if="status.connected && status.hashrate > 0" class="health-item success">
            ✓ 已连接矿池，矿工正常工作
          </div>
          <div v-else class="health-item warning">
            ⚠ 正在准备（若长时间未连接矿池或配置错误，请检查矿池地址与钱包）
          </div>
        </div>
      </div>

      <!-- 系统信息 -->
      <div class="card">
        <div class="card-header">
          <h3>💻 系统信息</h3>
        </div>
        <div class="card-body">
          <div class="stat-item">
            <span class="label">操作系统:</span>
            <span class="value">{{ systemInfo.os }}</span>
          </div>
          <div class="stat-item">
            <span class="label">架构:</span>
            <span class="value">{{ systemInfo.arch }}</span>
          </div>
          <div class="stat-item">
            <span class="label">CPU核心:</span>
            <span class="value">{{ systemInfo.cpuCores }}</span>
          </div>
          <div class="stat-item">
            <span class="label">XMRig版本:</span>
            <span class="value">{{ systemInfo.xmrigVersion }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 控制按钮 -->
    <div class="controls">
      <button
        class="btn btn-primary"
        :disabled="loading || status.running"
        @click="startMining"
      >
        ▶️ 开始挖矿
      </button>
      <button
        class="btn btn-danger"
        :disabled="loading || !status.running"
        @click="stopMining"
      >
        ⏹️ 停止挖矿
      </button>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1200px;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1.5rem;
  backdrop-filter: blur(10px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.card-header h3 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
}

.status-badge {
  padding: 0.4rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
}

.status-badge.running {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
  border: 1px solid rgba(76, 175, 80, 0.5);
}

.status-badge.stopped {
  background: rgba(158, 158, 158, 0.2);
  color: #9e9e9e;
  border: 1px solid rgba(158, 158, 158, 0.5);
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-item .label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.95rem;
}

.stat-item .value {
  color: #fff;
  font-weight: 600;
  font-size: 1.1rem;
}

.stat-item .value.warning {
  color: #ff9800;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
  }
}

.stat-item .value.small {
  font-size: 0.9rem;
}

.health-status {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.health-item {
  padding: 0.75rem 1rem;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.health-item.success {
  background: rgba(76, 175, 80, 0.15);
  color: #4caf50;
  border: 1px solid rgba(76, 175, 80, 0.3);
}

.health-item.warning {
  background: rgba(255, 152, 0, 0.15);
  color: #ff9800;
  border: 1px solid rgba(255, 152, 0, 0.3);
  animation: warningPulse 2s ease-in-out infinite;
}

@keyframes warningPulse {
  0%, 100% {
    border-color: rgba(255, 152, 0, 0.3);
  }
  50% {
    border-color: rgba(255, 152, 0, 0.6);
  }
}

.controls {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.btn {
  padding: 1rem 2.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  min-width: 180px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #4caf50 0%, #45a049 100%);
  color: white;
}

.btn-primary:not(:disabled):hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(76, 175, 80, 0.4);
}

.btn-danger {
  background: linear-gradient(135deg, #f44336 0%, #d32f2f 100%);
  color: white;
}

.btn-danger:not(:disabled):hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(244, 67, 54, 0.4);
}
</style>
