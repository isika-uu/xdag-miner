<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { LoadConfig, SaveConfig, GetDefaultConfig, GetMinerStatus, GetSystemInfo } from '../../wailsjs/go/main/App'
import ConfigHelp from './ConfigHelp.vue'
import Toast from './Toast.vue'
import ConfirmDialog from './ConfirmDialog.vue'

const config = ref(null)
const loading = ref(false)
const saving = ref(false)
const toast = ref({
  show: false,
  type: 'info',
  message: ''
})
const confirmDialog = ref({
  show: false,
  title: '',
  message: '',
  type: 'warning',
  action: null
})
const status = ref({
  running: false
})
const formDisabled = computed(() => !!(status.value && status.value.running))
let statusTimer = null
const systemInfo = ref({ arch: '' })

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
const handleConfirm = async () => {
  if (confirmDialog.value.action) {
    await confirmDialog.value.action()
  }
  confirmDialog.value.show = false
}

// 加载配置
const loadConfig = async () => {
  loading.value = true
  try {
    const cfg = await LoadConfig()
    if (cfg) {
      config.value = cfg
    }
  } catch (err) {
    showToast('error', '加载配置失败: ' + err)
  } finally {
    loading.value = false
  }
}

// 保存配置
const saveConfig = async () => {
  saving.value = true
  try {
    status.value = await GetMinerStatus()
    if (status.value && status.value.running) {
      showToast('warning', '当前挖矿运行中，禁止修改配置')
      return
    }
    await SaveConfig(config.value)
    showToast('success', '配置保存成功！')
  } catch (err) {
    showToast('error', '保存失败: ' + err)
  } finally {
    saving.value = false
  }
}

// 重置为默认配置
const resetToDefault = () => {
  showConfirm(
    '重置配置',
    '确定要重置为默认配置吗？这将清除所有自定义设置。',
    'warning',
    async () => {
      try {
        const defaultConfig = await GetDefaultConfig()
        config.value = defaultConfig
        showToast('info', '已重置为默认配置，请点击保存按钮保存更改')
      } catch (err) {
        showToast('error', '重置失败: ' + err)
      }
    }
  )
}

// 添加矿池
const addPool = () => {
  // 检查是否所有矿池都为空
  const hasEmptyPool = config.value.pools && config.value.pools.some(pool => 
    !pool.url.trim() && !pool.user.trim()
  );
  
  if (hasEmptyPool) {
    showToast('warning', '请先填写现有矿池配置，不能添加空矿池');
    return;
  }
  
  if (!config.value.pools) {
    config.value.pools = []
  }
  config.value.pools.push({
    url: '',
    user: '',
    pass: 'x',
    enabled: true,
    nicehash: false,
    tls: false
  })
  showToast('success', '已添加新矿池配置')
}

// 删除矿池
const removePool = (index) => {
  // 检查是否是最后一个矿池
  if (config.value.pools.length <= 1) {
    showConfirm(
      '删除矿池',
      '这是最后一个矿池配置，删除后将无法挖矿。确定要删除吗？',
      'danger',
      () => {
        config.value.pools.splice(index, 1)
        showToast('success', '矿池配置已删除')
        
        // 添加一个空的矿池配置，确保至少有一个
        config.value.pools.push({
          url: '',
          user: '',
          pass: 'x',
          enabled: true,
          nicehash: false,
          tls: false
        })
        showToast('info', '已添加一个新的空矿池配置')
      }
    )
  } else {
    showConfirm(
      '删除矿池',
      '确定要删除这个矿池配置吗？',
      'danger',
      () => {
        config.value.pools.splice(index, 1)
        showToast('success', '矿池配置已删除')
      }
    )
  }
}

const refreshStatus = async () => {
  try {
    const s = await GetMinerStatus()
    status.value = s || { running: false }
  } catch (_) {
  }
}

const loadSystemInfo = async () => {
  try {
    const info = await GetSystemInfo()
    if (info) {
      systemInfo.value = info
    }
  } catch (_) {
  }
}

onMounted(() => {
  loadConfig()
  refreshStatus()
  loadSystemInfo()
  statusTimer = setInterval(refreshStatus, 2000)
})

onUnmounted(() => {
  if (statusTimer) {
    clearInterval(statusTimer)
    statusTimer = null
  }
})
</script>

<template>
  <div class="config-panel">
    <div v-if="loading" class="loading">
      加载配置中...
    </div>

    <div v-else-if="config" class="config-content">
      <!-- 帮助按钮 -->
      <ConfigHelp />
      
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

      <div v-if="formDisabled" class="lock-banner">⚠️ 挖矿运行中，配置已锁定</div>

      <!-- 矿池配置 -->
      <section class="config-section">
        <div class="section-header">
          <h2>⛏️ 矿池配置</h2>
          <button class="btn btn-small" :disabled="formDisabled" @click="addPool">+ 添加矿池</button>
        </div>

        <div v-for="(pool, index) in config.pools" :key="index" class="pool-item">
          <div class="pool-header">
            <h4>矿池 #{{ index + 1 }}</h4>
            <button class="btn-remove" :disabled="formDisabled" @click="removePool(index)">删除</button>
          </div>

          <div class="form-grid">
            <div class="form-group">
              <label>矿池地址 *</label>
              <input
                v-model="pool.url"
                type="text"
                placeholder="例如: stratum+ssl://equal.xdagminer.com:13003"
                :disabled="formDisabled"
              />
            </div>

            <div class="form-group">
              <label>钱包地址 *</label>
              <input
                v-model="pool.user"
                type="text"
                placeholder="输入您的钱包地址"
                :disabled="formDisabled"
              />
            </div>

            <div class="form-group">
              <label>密码</label>
              <input v-model="pool.pass" type="text" placeholder="默认: x" :disabled="formDisabled" />
            </div>

            <div class="form-group">
              <label>矿机ID</label>
              <input v-model="pool['rig-id']" type="text" placeholder="可选" :disabled="formDisabled" />
            </div>
          </div>

          <div class="form-row">
            <label class="checkbox">
              <input v-model="pool.enabled" type="checkbox" :disabled="formDisabled" />
              <span>启用此矿池</span>
            </label>
            <label class="checkbox">
              <input v-model="pool.nicehash" type="checkbox" :disabled="formDisabled" />
              <span>NiceHash 模式</span>
            </label>
            <label class="checkbox">
              <input v-model="pool.tls" type="checkbox" :disabled="formDisabled" />
              <span>使用 TLS</span>
            </label>
          </div>
        </div>
      </section>

      <!-- CPU配置 -->
      <section class="config-section">
        <h2>🖥️ CPU 配置</h2>

        <div class="form-grid">
          <div class="form-group">
            <div class="range-header">
              <label>最大线程(%)</label>
              <span class="range-value">{{ config.cpu['max-threads-hint'] }}%</span>
            </div>
            <input
              v-model.number="config.cpu['max-threads-hint']"
              type="range"
              min="1"
              max="100"
              step="1"
              :disabled="formDisabled"
            />
            <div class="range-scale">
              <span>1%</span>
              <span>25%</span>
              <span>50%</span>
              <span>75%</span>
              <span>100%</span>
            </div>
            <small>使用 CPU 核心的百分比 (1-100)</small>
          </div>

          <div class="form-group">
            <label>优先级</label>
            <select v-model="config.cpu.priority" :disabled="formDisabled">
              <option :value="null">自动</option>
              <option :value="1">低</option>
              <option :value="2">普通</option>
              <option :value="3">高</option>
              <option :value="4">实时 (不推荐)</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <label class="checkbox">
            <input v-model="config.cpu.enabled" type="checkbox" :disabled="formDisabled" />
            <span>启用 CPU 挖矿</span>
          </label>
          <label class="checkbox">
            <input v-model="config.cpu['huge-pages']" type="checkbox" :disabled="formDisabled" />
            <span>使用大页内存</span>
          </label>
          <label class="checkbox">
            <input v-model="config.cpu.asm" type="checkbox" :disabled="formDisabled" />
            <span>使用汇编优化</span>
          </label>
        </div>
      </section>

      <!-- HTTP API配置 -->
      <section class="config-section">
        <h2>🌐 HTTP API 配置</h2>

        <div class="form-grid">
          <div class="form-group">
            <label>监听地址</label>
            <input v-model="config.http.host" type="text" placeholder="127.0.0.1" :disabled="formDisabled" />
          </div>

          <div class="form-group">
            <label>端口</label>
            <input v-model.number="config.http.port" type="number" min="1" max="65535" :disabled="formDisabled" />
          </div>

          <div class="form-group">
            <label>访问令牌</label>
            <input
              v-model="config.http['access-token']"
              type="text"
              placeholder="可选，用于API访问认证"
              :disabled="formDisabled"
            />
          </div>
        </div>

        <div class="form-row">
          <label class="checkbox">
            <input v-model="config.http.enabled" type="checkbox" :disabled="formDisabled" />
            <span>启用 HTTP API</span>
          </label>
          <label class="checkbox">
            <input v-model="config.http.restricted" type="checkbox" :disabled="formDisabled" />
            <span>限制模式（只读）</span>
          </label>
        </div>
      </section>

      <!-- 操作按钮 -->
      <div class="actions">
        <button class="btn btn-secondary" :disabled="formDisabled" @click="resetToDefault">
          重置默认配置
        </button>
        <button class="btn btn-primary" :disabled="saving || formDisabled" @click="saveConfig">
          {{ saving ? '保存中...' : '💾 保存配置' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.config-panel {
  max-width: 100%;
  position: relative;
}

.loading {
  text-align: center;
  padding: 3rem;
  font-size: 1.2rem;
  color: rgba(255, 255, 255, 0.6);
}

.config-content {
  padding-bottom: 105px;
}

.lock-banner {
  margin: 0 0 1rem 0;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  background: rgba(255, 152, 0, 0.12);
  border: 1px solid rgba(255, 152, 0, 0.4);
  color: #ffb74d;
  font-size: 0.95rem;
}

.config-section {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.config-section h2 {
  margin: 0 0 1.5rem 0;
  font-size: 1.3rem;
  color: #64b5f6;
}

.pool-item {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1rem;
}

.pool-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.pool-header h4 {
  margin: 0;
  color: #fff;
}

.btn-remove {
  background: rgba(244, 67, 54, 0.2);
  border: 1px solid rgba(244, 67, 54, 0.5);
  color: #ff6b6b;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.btn-remove:hover {
  background: rgba(244, 67, 54, 0.3);
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  font-weight: 500;
}

.form-group input,
.form-group select {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  padding: 0.75rem;
  color: #fff;
  font-size: 0.95rem;
  transition: all 0.3s ease;
}

.form-group select option {
  background: #1a1f35;
  color: #fff;
  padding: 0.5rem;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #4a90e2;
  background: rgba(255, 255, 255, 0.12);
}

.form-group small {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
  margin-top: 0.3rem;
}

.range-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.range-value {
  color: #64b5f6;
  font-weight: 600;
}

.form-group input[type='range'] {
  -webkit-appearance: none;
  appearance: none;
  width: 100%;
  height: 6px;
  border-radius: 6px;
  border: none;
  padding: 0;
  background: linear-gradient(135deg, #4a90e2 0%, #357abd 100%);
  outline: none;
}

.form-group input[type='range']::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #fff;
  border: 3px solid #4a90e2;
  box-shadow: 0 2px 8px rgba(74, 144, 226, 0.5);
  cursor: pointer;
}

.form-group input[type='range']::-moz-range-thumb {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #fff;
  border: 3px solid #4a90e2;
  box-shadow: 0 2px 8px rgba(74, 144, 226, 0.5);
  cursor: pointer;
}

.range-scale {
  display: flex;
  justify-content: space-between;
  margin-top: 0.4rem;
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
}

.form-row {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}

.checkbox {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.8);
}

.checkbox input[type='checkbox'] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.actions {
  position: fixed;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 1rem;
  z-index: 100;
  background: rgba(26, 31, 53, 0.9);
  padding: 1rem 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.actions .btn {
  margin: 0;
}

.btn {
  padding: 0.875rem 2rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-small {
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
}

.btn-primary {
  background: linear-gradient(135deg, #4a90e2 0%, #357abd 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(74, 144, 226, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

.btn-secondary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(74, 144, 226, 0.4);
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
