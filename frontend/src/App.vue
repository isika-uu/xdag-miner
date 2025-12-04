<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { APP_VERSION } from './version'
import Dashboard from './components/Dashboard.vue'
import ConfigPanel from './components/ConfigPanel.vue'
import LogPanel from './components/LogPanel.vue'
import About from './components/About.vue'

const activeTab = ref('dashboard')
let updateTimer = null

const tabs = [
  { id: 'dashboard', name: '控制面板', icon: '📊' },
  { id: 'config', name: '配置管理', icon: '⚙️' },
  { id: 'logs', name: '运行日志', icon: '📝' },
  { id: 'about', name: '关于', icon: 'ℹ️' }
]

const normalize = (v) => {
  if (!v) return '0.0.0'
  return String(v).trim().replace(/^v/i, '')
}

const cmp = (a, b) => {
  const pa = normalize(a).split('.').map(n => parseInt(n || '0'))
  const pb = normalize(b).split('.').map(n => parseInt(n || '0'))
  const len = Math.max(pa.length, pb.length)
  for (let i = 0; i < len; i++) {
    const ai = pa[i] || 0
    const bi = pb[i] || 0
    if (ai > bi) return 1
    if (ai < bi) return -1
  }
  return 0
}

onMounted(() => {

})

onUnmounted(() => {
  if (updateTimer) {
    clearInterval(updateTimer)
    updateTimer = null
  }
})
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <div class="header-content no-select">
        <h1>⛏️ XDAG 矿工管理器-XDAG Miner</h1>
        <p class="subtitle">简单易用的图形化挖矿工具</p>
      </div>
    </header>

    <nav class="app-nav no-select">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['nav-button', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <span class="icon">{{ tab.icon }}</span>
        <span>{{ tab.name }}</span>
      </button>
    </nav>

    <main class="app-main">
      <Dashboard v-if="activeTab === 'dashboard'" />
      <ConfigPanel v-if="activeTab === 'config'" />
      <LogPanel v-if="activeTab === 'logs'" />
      <About v-if="activeTab === 'about'" />
    </main>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1f35 0%, #2d3561 100%);
  color: #fff;
}

.app-header {
  background: rgba(0, 0, 0, 0.3);
  padding: 1.5rem 2rem;
  border-bottom: 2px solid rgba(74, 144, 226, 0.3);
}

.header-content h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #4a90e2 0%, #64b5f6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  margin: 0.5rem 0 0;
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.6);
}

.update-banner {
  margin-top: 0.5rem;
  background: rgba(76, 175, 80, 0.15);
  border: 1px solid rgba(76, 175, 80, 0.5);
  color: #a5d6a7;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
}

.update-banner a {
  color: #81c784;
  margin-left: 0.5rem;
}

.app-nav {
  display: flex;
  gap: 1rem;
  padding: 1.5rem 2rem;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.nav-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-button:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  transform: translateY(-2px);
}

.nav-button.active {
  background: linear-gradient(135deg, #4a90e2 0%, #357abd 100%);
  color: #fff;
  border-color: #4a90e2;
  box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3);
}

.icon {
  font-size: 1.2rem;
}

.app-main {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}
</style>

<style>
.no-select {
  -webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
</style>
