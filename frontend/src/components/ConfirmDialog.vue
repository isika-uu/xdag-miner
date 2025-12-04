<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '确认'
  },
  message: {
    type: String,
    default: ''
  },
  confirmText: {
    type: String,
    default: '确定'
  },
  cancelText: {
    type: String,
    default: '取消'
  },
  type: {
    type: String,
    default: 'warning' // warning, danger, info
  }
})

const emit = defineEmits(['confirm', 'cancel', 'close'])

const visible = ref(props.show)

watch(() => props.show, (newVal) => {
  visible.value = newVal
})

const handleConfirm = () => {
  emit('confirm')
  close()
}

const handleCancel = () => {
  emit('cancel')
  close()
}

const close = () => {
  visible.value = false
  emit('close')
}
</script>

<template>
  <Transition name="modal">
    <div v-if="visible" class="modal-overlay" @click="handleCancel">
      <div class="modal-container" @click.stop>
        <div :class="['modal-header', `modal-${type}`]">
          <h3>{{ title }}</h3>
        </div>
        <div class="modal-body">
          <p>{{ message }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-cancel" @click="handleCancel">
            {{ cancelText }}
          </button>
          <button :class="['btn', `btn-${type}`]" @click="handleConfirm">
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  backdrop-filter: blur(4px);
}

.modal-container {
  background: linear-gradient(135deg, #1a1f35 0%, #2d3561 100%);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  min-width: 400px;
  max-width: 500px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h3 {
  margin: 0;
  font-size: 1.3rem;
  color: white;
}

.modal-warning {
  background: rgba(255, 152, 0, 0.1);
  border-bottom-color: rgba(255, 152, 0, 0.3);
}

.modal-danger {
  background: rgba(244, 67, 54, 0.1);
  border-bottom-color: rgba(244, 67, 54, 0.3);
}

.modal-info {
  background: rgba(33, 150, 243, 0.1);
  border-bottom-color: rgba(33, 150, 243, 0.3);
}

.modal-body {
  padding: 2rem 1.5rem;
}

.modal-body p {
  margin: 0;
  color: rgba(255, 255, 255, 0.9);
  font-size: 1rem;
  line-height: 1.6;
}

.modal-footer {
  padding: 1.5rem;
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  min-width: 100px;
}

.btn-cancel {
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-cancel:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-1px);
}

.btn-warning {
  background: linear-gradient(135deg, #ff9800 0%, #f57c00 100%);
  color: white;
}

.btn-warning:hover {
  box-shadow: 0 4px 12px rgba(255, 152, 0, 0.4);
  transform: translateY(-1px);
}

.btn-danger {
  background: linear-gradient(135deg, #f44336 0%, #d32f2f 100%);
  color: white;
}

.btn-danger:hover {
  box-shadow: 0 4px 12px rgba(244, 67, 54, 0.4);
  transform: translateY(-1px);
}

.btn-info {
  background: linear-gradient(135deg, #2196f3 0%, #1976d2 100%);
  color: white;
}

.btn-info:hover {
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.4);
  transform: translateY(-1px);
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9);
}
</style>
