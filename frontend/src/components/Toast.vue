<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'info' // success, error, warning, info
  },
  message: {
    type: String,
    default: ''
  },
  duration: {
    type: Number,
    default: 3000
  }
})

const emit = defineEmits(['close'])

const visible = ref(props.show)

watch(() => props.show, (newVal) => {
  visible.value = newVal
  if (newVal && props.duration > 0) {
    setTimeout(() => {
      close()
    }, props.duration)
  }
})

const close = () => {
  visible.value = false
  emit('close')
}

const getIcon = () => {
  const icons = {
    success: '✓',
    error: '✕',
    warning: '⚠',
    info: 'ℹ'
  }
  return icons[props.type] || icons.info
}
</script>

<template>
  <Transition name="toast">
    <div v-if="visible" :class="['toast', `toast-${type}`]" @click="close">
      <span class="toast-icon">{{ getIcon() }}</span>
      <span class="toast-message">{{ message }}</span>
      <button class="toast-close" @click.stop="close">✕</button>
    </div>
  </Transition>
</template>

<style scoped>
.toast {
  position: fixed;
  top: 2rem;
  right: 2rem;
  min-width: 300px;
  max-width: 500px;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  z-index: 9999;
  cursor: pointer;
  animation: slideIn 0.3s ease;
}

.toast-success {
  background: rgba(76, 175, 80, 0.95);
  border: 1px solid rgba(76, 175, 80, 1);
  color: white;
}

.toast-error {
  background: rgba(244, 67, 54, 0.95);
  border: 1px solid rgba(244, 67, 54, 1);
  color: white;
}

.toast-warning {
  background: rgba(255, 152, 0, 0.95);
  border: 1px solid rgba(255, 152, 0, 1);
  color: white;
}

.toast-info {
  background: rgba(33, 150, 243, 0.95);
  border: 1px solid rgba(33, 150, 243, 1);
  color: white;
}

.toast-icon {
  font-size: 1.5rem;
  font-weight: bold;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  flex-shrink: 0;
}

.toast-message {
  flex: 1;
  font-size: 0.95rem;
  font-weight: 500;
}

.toast-close {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.toast-close:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(90deg);
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(100%);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
