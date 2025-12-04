<script setup>
import {ref, onMounted} from 'vue'
import {APP_VERSION} from '../version'

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

const checkUpdate = async () => {
  open("https://github.com/isika-uu/xdag-miner/tags")
}

onMounted(() => {

})
</script>

<template>
  <div class="about">
    <div class="card">
      <h3>关于程序</h3>
      <div class="meta">
        <div class="row"><span>github</span><strong>https://github.com/isika-uu/xdag-miner</strong></div>
      </div>
    </div>

    <div class="card">
      <h3>关于社区</h3>
      <div class="meta">
        <div class="row"><span>官网</span><strong>https://xdag.io</strong></div>
        <div class="row"><span>XDAG Community中文群</span><strong>https://t.me/XDAGCommunityChina</strong></div>
      </div>
    </div>

    <div class="card">
      <h3>版本</h3>
      <div class="grid">
        <div class="item"><span>当前版本</span><strong>{{ APP_VERSION }}</strong></div>
      </div>
      <div class="actions">
        <button class="btn" @click="checkUpdate">检查更新</button>
      </div>
    </div>

    <div class="card">
      <h3>支持开发者</h3>
      <div class="meta">
        <div class="row"><span>XDAG地址</span><strong>NNZabJQEhrQGTPabqABWVK9v3rSsNQ7Sy</strong></div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.about {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

@media (max-width: 800px) {
  .about {
    grid-template-columns: 1fr;
  }
}

.card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1.25rem 1.5rem;
}

h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.4rem;
  color: #64b5f6;
}

h3 {
  margin: 0 0 0.75rem 0;
  font-size: 1.15rem;
  color: #64b5f6;
}

.name {
  margin: 0.25rem 0;
  font-weight: 700;
}

.desc {
  margin: 0.25rem 0 0.75rem 0;
  color: rgba(255, 255, 255, 0.7);
}

.meta .row {
  display: flex;
  justify-content: space-between;
  padding: 0.35rem 0;
}

.meta .row span {
  color: rgba(255, 255, 255, 0.6);
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem 1rem;
}

.item {
  display: flex;
  justify-content: space-between;
}

.item span {
  color: rgba(255, 255, 255, 0.6);
}

.loading {
  color: rgba(255, 255, 255, 0.7);
}

.tips {
  margin: 0.25rem 0 0 1rem;
  padding: 0;
  color: rgba(255, 255, 255, 0.8);
}

.tips li {
  margin: 0.25rem 0;
}

.actions {
  margin-top: 0.75rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.donate-card {
  display: flex;
  flex-direction: column;
}

.donate-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.donate-img {
  width: 180px;
  height: 180px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.05);
  object-fit: cover;
}

.donate-text {
  color: rgba(255, 255, 255, 0.8);
}
</style>
