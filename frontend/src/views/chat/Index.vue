<template>
  <div class="chat-page">
    <Navbar />
    <div class="container chat-container">
      <div class="chat-card card">
        <div class="chat-header">
          <h2>🤖 智能社区助手</h2>
          <p>基于通义千问大模型</p>
        </div>
        
        <div class="chat-history" ref="historyRef">
          <div v-for="(msg, index) in messages" :key="index" class="message-item" :class="msg.role">
            <div class="avatar">
              <span v-if="msg.role === 'user'">👤</span>
              <span v-else>🤖</span>
            </div>
            <div class="content">
              <div class="bubble">{{ msg.content }}</div>
              <div class="time">{{ msg.time }}</div>
            </div>
          </div>
          
          <div v-if="loading" class="message-item assistant">
             <div class="avatar">🤖</div>
             <div class="content">
                <div class="bubble loading-bubble">
                  <span>.</span><span>.</span><span>.</span>
                </div>
             </div>
          </div>
        </div>

        <div class="chat-input-area">
          <el-input
            v-model="inputContent"
            type="textarea"
            :rows="3"
            placeholder="请输入您的问题..."
            @keyup.enter.ctrl="handleSend"
          />
          <div class="actions">
            <span class="tip">Ctrl + Enter 发送</span>
            <el-button type="primary" @click="handleSend" :loading="loading">发送</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { sendChat } from '@/api/chat'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const messages = ref([
  { role: 'assistant', content: '您好！我是您的智能社区助手，请问有什么可以帮您？', time: dayjs().format('HH:mm') }
])
const inputContent = ref('')
const loading = ref(false)
const historyRef = ref(null)

const scrollToBottom = () => {
    nextTick(() => {
        if (historyRef.value) {
            historyRef.value.scrollTop = historyRef.value.scrollHeight
        }
    })
}

const handleSend = async () => {
  const content = inputContent.value.trim()
  if (!content) return

  // 1. 添加用户消息
  messages.value.push({
    role: 'user',
    content: content,
    time: dayjs().format('HH:mm')
  })
  inputContent.value = ''
  scrollToBottom()

  // 2. 发送请求
  loading.value = true
  try {
    const res = await sendChat({ content })
    // 3. 添加 AI 回复
    messages.value.push({
      role: 'assistant',
      content: res.reply,
      time: dayjs().format('HH:mm')
    })
  } catch (error) {
    ElMessage.error('AI 响应失败，请稍后重试')
    messages.value.push({
        role: 'system',
        content: '❌ 生成失败: ' + (error.message || '网络错误'),
        time: dayjs().format('HH:mm')
    })
  } finally {
    loading.value = false
    scrollToBottom()
    // 聚焦输入框 (可选)
  }
}
</script>

<style scoped>
.chat-page {
  min-height: 100vh;
  padding-bottom: 20px;
  background-color: #f5f7fa;
}

.chat-container {
    max-width: 800px;
    margin: 20px auto;
}

.chat-card {
    height: 75vh;
    display: flex;
    flex-direction: column;
    padding: 0;
    overflow: hidden;
}

.chat-header {
    padding: 20px;
    border-bottom: 1px solid #eee;
    background: #fff;
    text-align: center;
}
.chat-header h2 { margin: 0; font-size: 18px; }
.chat-header p { margin: 5px 0 0; color: #999; font-size: 12px; }

.chat-history {
    flex: 1;
    padding: 20px;
    overflow-y: auto;
    background: #f9f9f9;
}

.message-item {
    display: flex;
    margin-bottom: 20px;
}

.message-item.user {
    flex-direction: row-reverse;
}

.avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: #e0e0e0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    flex-shrink: 0;
}
.message-item.user .avatar { margin-left: 10px; background: #d9ecff; }
.message-item.assistant .avatar { margin-right: 10px; background: #fff; border: 1px solid #eee; }

.content {
    max-width: 70%;
}

.bubble {
    padding: 10px 15px;
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-all;
}

.message-item.user .bubble {
    background: var(--primary-color);
    color: #fff;
    border-top-right-radius: 0;
}

.message-item.assistant .bubble {
    background: #fff;
    color: #333;
    border: 1px solid #eee;
    border-top-left-radius: 0;
}

.message-item.system .bubble {
    background: #fee;
    color: #f56c6c;
    font-size: 12px;
}

.time {
    font-size: 12px;
    color: #999;
    margin-top: 5px;
    text-align: right;
}
.message-item.assistant .time { text-align: left; }

.chat-input-area {
    padding: 20px;
    background: #fff;
    border-top: 1px solid #eee;
}

.actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 10px;
}

.tip {
    font-size: 12px;
    color: #999;
}

.loading-bubble span {
    display: inline-block;
    animation: bounce 1.4s infinite ease-in-out both;
    margin: 0 2px;
    font-weight: bold;
    color: #999;
}

.loading-bubble span:nth-child(1) { animation-delay: -0.32s; }
.loading-bubble span:nth-child(2) { animation-delay: -0.16s; }

@keyframes bounce {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}
</style>
