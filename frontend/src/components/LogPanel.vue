<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { GetLogs, GetMinerStatus, ClearLogs } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import ConfirmDialog from './ConfirmDialog.vue'
import Toast from './Toast.vue'

const logs = ref([])
const autoScroll = ref(true)
const isRunning = ref(false)
const logContainer = ref(null)
const confirmDialog = ref({
  show: false,
  title: '',
  message: '',
  type: 'warning',
  action: null
})
const toast = ref({
  show: false,
  type: 'info',
  message: ''
})

// 显示提示
const showToast = (type, message) => {
  toast.value = { show: true, type, message }
}

// 显示确认对话框
const showConfirm = (title, message, type, action) => {
  confirmDialog.value = {
    show: true,
    title,
    message,
    type,
    action
  }
}

// 处理确认
const handleConfirm = () => {
  if (confirmDialog.value.action) {
    confirmDialog.value.action()
  }
  confirmDialog.value.show = false
}

// 加载日志
const loadLogs = async () => {
  try {
    const logData = await GetLogs()
    if (logData) {
      logs.value = logData.map((line, index) => ({
        id: index,
        text: line,
        time: new Date().toLocaleTimeString()
      }))
      scrollToBottom()
    }
  } catch (err) {
    console.error('加载日志失败:', err)
  }
}

// 检查运行状态
const checkStatus = async () => {
  try {
    const status = await GetMinerStatus()
    isRunning.value = status?.running || false
  } catch (err) {
    console.error('检查状态失败:', err)
  }
}

// 清空日志
const clearLogs = async () => {
  showConfirm(
    '清空日志',
    '确定要清空所有日志吗？',
    'warning',
    async () => {
      try {
        await ClearLogs()
        logs.value = []
        showToast('success', '日志已清空')
      } catch (err) {
        showToast('error', '清空日志失败: ' + err)
      }
    }
  )
}

// 滚动到底部
const scrollToBottom = () => {
  if (autoScroll.value && logContainer.value) {
    setTimeout(() => {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }, 100)
  }
}

// 添加日志
const addLog = (logData) => {
  logs.value.push({
    id: logs.value.length,
    text: logData.line,
    time: logData.time,
    source: logData.source
  })
  
  // 限制日志数量
  if (logs.value.length > 500) {
    logs.value.shift()
  }
  
  scrollToBottom()
}

let refreshInterval = null

onMounted(() => {
  loadLogs()
  checkStatus()
  
  // 监听新日志事件
  EventsOn('miner:log', addLog)
  
  // 定期检查状态
  refreshInterval = setInterval(checkStatus, 3000)
})

onUnmounted(() => {
  EventsOff('miner:log')
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<template>
  <div class="log-panel">
    <!-- Toast 提示 -->
    <Toast
      :show="toast.show"
      :type="toast.type"
      :message="toast.message"
      @close="toast.show = false"
    />

    <!-- 确认对话框 -->
    <ConfirmDialog
      :show="confirmDialog.show"
      :title="confirmDialog.title"
      :message="confirmDialog.message"
      :type="confirmDialog.type"
      @confirm="handleConfirm"
      @cancel="confirmDialog.show = false"
      @close="confirmDialog.show = false"
    />

    <div class="log-header">
      <div class="header-left">
        <h2>📝 运行日志</h2>
        <span :class="['status-indicator', isRunning ? 'running' : 'stopped']">
          {{ isRunning ? '● 运行中' : '○ 已停止' }}
        </span>
      </div>
      <div class="header-actions">
        <label class="checkbox">
          <input v-model="autoScroll" type="checkbox" />
          <span>自动滚动</span>
        </label>
        <button class="btn btn-small" @click="clearLogs">清空日志</button>
      </div>
    </div>

    <div ref="logContainer" class="log-container">
      <div v-if="logs.length === 0" class="empty-state">
        <p>暂无日志</p>
        <small>启动挖矿后，日志将在此处显示</small>
      </div>
      <div v-else class="log-content">
        <div
          v-for="log in logs"
          :key="log.id"
          :class="['log-line', log.source === 'stderr' ? 'error' : '']"
        >
          <span class="log-time">{{ log.time }}</span>
          <span class="log-text">{{ log.text }}</span>
        </div>
      </div>
    </div>

    <div class="log-footer">
      <span class="log-count">共 {{ logs.length }} 条日志</span>
    </div>
  </div>
</template>

<style scoped>
.log-panel {
  max-width: 1200px;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 294px);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  overflow: hidden;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(0, 0, 0, 0.2);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-left h2 {
  margin: 0;
  font-size: 1.3rem;
  color: #64b5f6;
}

.status-indicator {
  padding: 0.4rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
}

.status-indicator.running {
  background: rgba(76, 175, 80, 0.2);
  color: #4caf50;
  border: 1px solid rgba(76, 175, 80, 0.5);
}

.status-indicator.stopped {
  background: rgba(158, 158, 158, 0.2);
  color: #9e9e9e;
  border: 1px solid rgba(158, 158, 158, 0.5);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.checkbox {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
}

.checkbox input[type='checkbox'] {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-1px);
}

.log-container {
  flex: 1;
  overflow-y: auto;
  background: rgba(0, 0, 0, 0.3);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.log-container::-webkit-scrollbar {
  width: 10px;
}

.log-container::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
}

.log-container::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 5px;
}

.log-container::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: rgba(255, 255, 255, 0.4);
}

.empty-state p {
  margin: 0;
  font-size: 1.2rem;
}

.empty-state small {
  margin-top: 0.5rem;
  font-size: 0.9rem;
}

.log-content {
  padding: 1rem;
}

.log-line {
  display: flex;
  gap: 1rem;
  padding: 0.4rem 0.5rem;
  margin-bottom: 0.2rem;
  border-radius: 4px;
  transition: background 0.2s ease;
  line-height: 1.5;
}

.log-line:hover {
  background: rgba(255, 255, 255, 0.05);
}

.log-line.error {
  background: rgba(244, 67, 54, 0.1);
  border-left: 3px solid #f44336;
}

.log-time {
  color: rgba(100, 181, 246, 0.8);
  font-size: 0.85rem;
  white-space: nowrap;
  min-width: 70px;
}

.log-text {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.9rem;
  word-break: break-all;
  flex: 1;
}

.log-footer {
  padding: 0.75rem 1.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-count {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.85rem;
}
</style>
