<template>
  <div class="chat-page">
    <Navbar />
    <div class="container chat-container">
      <div class="chat-card card">
        <div class="chat-header">
          <h2>🤖 智慧社区助手</h2>
          <p>支持通知总结、报修创建、商品下单与支付</p>
        </div>

        <div ref="historyRef" class="chat-history">
          <div
            v-for="(msg, index) in messages"
            :key="index"
            class="message-item"
            :class="msg.role"
          >
            <div class="avatar">
              <span v-if="msg.role === 'user'">👤</span>
              <span v-else-if="msg.role === 'assistant'">🤖</span>
              <span v-else>⚠️</span>
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
            <el-button type="primary" :loading="loading" @click="handleSend">发送</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import { getChatHistory, sendChat } from '@/api/chat'

const buildGreetingMessage = () => ({
  role: 'assistant',
  content: '您好，我是智慧社区助手。您可以让我帮您总结通知、创建报修、搜索商品、下单和支付。',
  time: dayjs().format('HH:mm')
})

const formatMessageTime = (value) => {
  const parsed = dayjs(value)
  return parsed.isValid() ? parsed.format('HH:mm') : dayjs().format('HH:mm')
}

const messages = ref([buildGreetingMessage()])
const inputContent = ref('')
const loading = ref(false)
const historyRef = ref(null)

const isPayIntent = (text) => ['支付', '付款', '结算', '确认支付'].some((kw) => text.includes(kw))

const scrollToBottom = () => {
  nextTick(() => {
    if (historyRef.value) {
      historyRef.value.scrollTop = historyRef.value.scrollHeight
    }
  })
}

const loadHistory = async () => {
  try {
    const res = await getChatHistory({ limit: 100 })
    const list = Array.isArray(res?.list) ? res.list : []
    if (list.length === 0) {
      messages.value = [buildGreetingMessage()]
      scrollToBottom()
      return
    }

    messages.value = list.map((item) => ({
      role: item.role,
      content: item.content,
      time: formatMessageTime(item.created_at)
    }))
    scrollToBottom()
  } catch {
    messages.value = [buildGreetingMessage()]
    ElMessage.warning('聊天记录加载失败，已展示默认欢迎语')
    scrollToBottom()
  }
}

const handleSend = async () => {
  const content = inputContent.value.trim()
  if (!content || loading.value) {
    return
  }

  let paymentPassword = ''
  if (isPayIntent(content)) {
    try {
      const { value } = await ElMessageBox.prompt(
        '请输入登录密码完成支付',
        '安全支付验证',
        {
          inputType: 'password',
          inputPlaceholder: '登录密码',
          confirmButtonText: '确认支付',
          cancelButtonText: '取消'
        }
      )
      paymentPassword = (value || '').trim()
      if (!paymentPassword) {
        ElMessage.warning('未输入支付密码，已取消本次支付请求')
        return
      }
    } catch {
      return
    }
  }

  messages.value.push({
    role: 'user',
    content,
    time: dayjs().format('HH:mm')
  })
  inputContent.value = ''
  loading.value = true
  scrollToBottom()

  try {
    const res = await sendChat({
      content,
      payment_password: paymentPassword
    })

    const reply = (res?.reply || '').trim()
    if (!reply) {
      throw new Error('empty AI response')
    }
    messages.value.push({
      role: 'assistant',
      content: reply,
      time: dayjs().format('HH:mm')
    })
  } catch (error) {
    ElMessage.error('AI 响应失败，请稍后重试')
    messages.value.push({
      role: 'system',
      content: `生成失败: ${error?.message || '网络错误'}`,
      time: dayjs().format('HH:mm')
    })
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

onMounted(() => {
  loadHistory()
})
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
  height: 78vh;
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

.chat-header h2 {
  margin: 0;
  font-size: 18px;
}

.chat-header p {
  margin: 5px 0 0;
  color: #999;
  font-size: 12px;
}

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

.message-item.user .avatar {
  margin-left: 10px;
  background: #d9ecff;
}

.message-item.assistant .avatar {
  margin-right: 10px;
  background: #fff;
  border: 1px solid #eee;
}

.message-item.system .avatar {
  margin-right: 10px;
  background: #fff0f0;
  border: 1px solid #ffd6d6;
}

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
  background: #fff0f0;
  color: #f56c6c;
  border: 1px solid #ffd6d6;
  border-top-left-radius: 0;
}

.time {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
  text-align: right;
}

.message-item.assistant .time,
.message-item.system .time {
  text-align: left;
}

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

.loading-bubble span:nth-child(1) {
  animation-delay: -0.32s;
}

.loading-bubble span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%,
  80%,
  100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}
</style>
