<script setup>
import { ref } from 'vue'

const showHelp = ref(false)
</script>

<template>
  <div class="config-help">
    <button class="help-button" @click="showHelp = !showHelp">
      ❓ 配置说明
    </button>

    <div v-if="showHelp" class="help-modal">
      <div class="help-content">
        <div class="help-header">
          <h2>⚙️ 配置说明</h2>
          <button class="close-button" @click="showHelp = false">✕</button>
        </div>

        <div class="help-body">
          <!-- 矿池配置 -->
          <section class="help-section">
            <h3>1️⃣ 矿池配置（最重要！）</h3>
            <div class="help-item">
              <strong>矿池地址 (Pool URL)</strong>
              <p>矿池服务器的地址和端口，例如：</p>
              <ul>
                <li><code>stratum+ssl://equal.xdagminer.com:13003</code> xdagminer.com-equal矿池(推荐)</li>
                <li><code>stratum+tcp://stratum.xdag.org:23656</code> xdag.org-equal矿池(推荐)</li>
                <li><code>stratum+ssl://solo.xdagminer.com:23003</code> xdagminer.com-solo矿池</li>
                <li>
                  <p>更多矿池前往 <a href="https://miningpoolstats.stream/dagger" target="_blank">XDAG矿池列表</a></p>
                </li>
              </ul>
            </div>

            <div class="help-item">
              <strong>钱包地址 (Wallet Address)</strong>
              <p>你的 XDAG 钱包地址，挖到的币会发送到这里。</p>
              <p class="warning">⚠️ 这是最重要的配置！没有钱包地址就无法收到挖矿收益。</p>
            </div>

            <div class="help-item">
              <strong>密码 (Password)</strong>
              <p>通常填 "x" 即可，部分矿池可能需要特定密码。</p>
            </div>

            <div class="help-item">
              <strong>矿机ID (Rig ID)</strong>
              <p>可选，用于在矿池中识别不同的矿机，例如：home-pc、laptop 等。</p>
            </div>
          </section>

          <!-- CPU配置 -->
          <section class="help-section">
            <h3>2️⃣ CPU 配置</h3>
            <div class="help-item">
              <strong>最大线程提示 (%)</strong>
              <p>使用 CPU 核心的百分比：</p>
              <ul>
                <li><code>100</code> - 使用所有核心（性能最高，但电脑会卡）</li>
                <li><code>50</code> - 使用一半核心（平衡模式，推荐）</li>
                <li><code>25</code> - 使用四分之一核心（低功耗模式）</li>
              </ul>
            </div>

            <div class="help-item">
              <strong>优先级 (Priority)</strong>
              <p>挖矿程序的运行优先级：</p>
              <ul>
                <li><code>低</code> - 不影响其他程序（推荐日常使用）</li>
                <li><code>普通</code> - 标准优先级</li>
                <li><code>高</code> - 优先分配CPU资源</li>
                <li><code>实时</code> - 最高优先级（不推荐，可能导致系统卡死）</li>
              </ul>
            </div>

            <div class="help-item">
              <strong>大页内存 (Huge Pages)</strong>
              <p>启用后可提升 5-10% 性能，但需要管理员权限。</p>
            </div>
          </section>

          <!-- HTTP API配置 -->
          <section class="help-section">
            <h3>3️⃣ HTTP API 配置</h3>
            <div class="help-item">
              <strong>启用 HTTP API</strong>
              <p>开启后，本程序可以实时读取挖矿数据（算力、温度等）。</p>
              <p class="tip">💡 建议保持开启，否则看不到实时算力！</p>
            </div>

            <div class="help-item">
              <strong>端口 (Port)</strong>
              <p>默认 3649，如果端口被占用可以改成其他（如 8082,3649）。</p>
            </div>
          </section>

          <!-- 使用流程 -->
          <section class="help-section">
            <h3>📋 首次使用步骤</h3>
            <ol class="steps">
              <li>
                <strong>获取 XDAG 钱包地址</strong>
                <p>前往 <a href="https://xdag.io/" target="_blank">xdag.io</a> 下载钱包，或使用交易所地址。</p>
              </li>
              <li>
                <strong>选择矿池</strong>
                <p>推荐新手使用 <code>stratum+ssl://equal.xdagminer.com:13003</code></p>
              </li>
              <li>
                <strong>填写配置</strong>
                <p>在"矿池配置"中填写矿池地址和钱包地址。</p>
              </li>
              <li>
                <strong>调整 CPU 使用率</strong>
                <p>建议先设置为 50%，避免电脑太卡。</p>
              </li>
              <li>
                <strong>保存配置</strong>
                <p>点击"💾 保存配置"按钮。</p>
              </li>
              <li>
                <strong>开始挖矿</strong>
                <p>返回"控制面板"，点击"▶️ 开始挖矿"。</p>
              </li>
            </ol>
          </section>

          <!-- 常见问题 -->
          <section class="help-section">
            <h3>❓ 常见问题</h3>
            <div class="faq">
              <div class="faq-item">
                <strong>Q: 挖矿能赚多少钱？</strong>
                <p>A: 取决于 CPU 性能，以及XDAG币目前市场价格。</p>
              </div>
              <div class="faq-item">
                <strong>Q: 为什么看不到算力？</strong>
                <p>A: 确保"HTTP API"已开启，端口未被占用。</p>
              </div>
              <div class="faq-item">
                <strong>Q: 电脑太卡怎么办？</strong>
                <p>A: 降低"最大线程提示 (%)"的值，比如改成 25% 或 50%。</p>
              </div>
              <div class="faq-item">
                <strong>Q: 多久能收到收益？</strong>
                <p>A: 大多数矿池在达到最低提现额度后会自动转账到你的钱包。</p>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.help-button {
  position: fixed;
  right: 2rem;
  bottom: 8rem;
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #4a90e2 0%, #357abd 100%);
  color: white;
  border: none;
  border-radius: 50px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(74, 144, 226, 0.4);
  transition: all 0.3s ease;
  z-index: 100;
}

.help-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(74, 144, 226, 0.6);
}

.help-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 2rem;
  overflow-y: auto;
}

.help-content {
  background: linear-gradient(135deg, #1a1f35 0%, #2d3561 100%);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  max-width: 800px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.help-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  position: sticky;
  top: 0;
  background: rgba(26, 31, 53, 0.95);
  backdrop-filter: blur(10px);
  z-index: 10;
}

.help-header h2 {
  margin: 0;
  color: #64b5f6;
  font-size: 1.5rem;
}

.close-button {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  font-size: 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.close-button:hover {
  background: rgba(244, 67, 54, 0.3);
  transform: rotate(90deg);
}

.help-body {
  padding: 2rem;
}

.help-section {
  margin-bottom: 2.5rem;
}

.help-section h3 {
  color: #4a90e2;
  margin-bottom: 1rem;
  font-size: 1.2rem;
}

.help-item {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.help-item strong {
  display: block;
  color: #64b5f6;
  margin-bottom: 0.5rem;
  font-size: 1rem;
}

.help-item p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0.5rem 0;
  line-height: 1.6;
}

.help-item ul,
.steps {
  margin: 0.5rem 0;
  padding-left: 1.5rem;
  color: rgba(255, 255, 255, 0.8);
}

.help-item li,
.steps li {
  margin: 0.3rem 0;
  line-height: 1.6;
}

.steps li {
  margin-bottom: 1rem;
}

.steps li strong {
  display: inline;
  color: #4caf50;
}

.steps li p {
  margin-left: 0;
  margin-top: 0.3rem;
}

code {
  background: rgba(0, 0, 0, 0.3);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  color: #4caf50;
  font-family: 'Consolas', 'Monaco', monospace;
}

.warning {
  background: rgba(255, 152, 0, 0.1);
  border-left: 3px solid #ff9800;
  padding: 0.5rem 1rem !important;
  border-radius: 4px;
  color: #ffb74d !important;
}

.tip {
  background: rgba(33, 150, 243, 0.1);
  border-left: 3px solid #2196f3;
  padding: 0.5rem 1rem !important;
  border-radius: 4px;
  color: #64b5f6 !important;
}

.faq-item {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.faq-item strong {
  color: #ff9800;
  display: block;
  margin-bottom: 0.5rem;
}

.faq-item p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
}

a {
  color: #4a90e2;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

.help-content::-webkit-scrollbar {
  width: 8px;
}

.help-content::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
}

.help-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}
</style>
