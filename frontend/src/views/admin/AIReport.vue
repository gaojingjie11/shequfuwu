<template>
  <div class="ai-report-page">
    <Navbar />

    <div class="container">
      <el-card class="report-list-card">
        <template #header>
          <div class="panel-header">
            <div class="panel-title">
              <el-icon><DataAnalysis /></el-icon>
              <span>AI 智能日报汇编</span>
            </div>
            <div class="header-actions">
              <el-button type="success" :icon="Plus" :loading="generating" @click="handleGenerate(true)">
                手动生成报告
              </el-button>
              <el-button type="primary" :icon="Refresh" :loading="listLoading" @click="fetchList">
                刷新列表
              </el-button>
            </div>
          </div>
        </template>

        <el-table :data="list" v-loading="listLoading" stripe>
          <el-table-column label="生成时间" width="220">
            <template #default="{ row }">
              <div class="time-cell">
                <el-icon><Clock /></el-icon>
                <span>{{ formatDateTime(row.created_at) }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="report_summary" label="报告摘要" min-width="620" />
          <el-table-column label="操作" width="140" align="center">
            <template #default="{ row }">
              <el-button type="primary" link @click="openDetail(row.id)">查看详情</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-wrap">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchList"
            @current-change="fetchList"
          />
        </div>
      </el-card>
    </div>

    <el-drawer
      v-model="detailVisible"
      size="70%"
      title="AI 报告详情"
      destroy-on-close
    >
      <div v-loading="detailLoading" class="detail-wrap">
        <template v-if="activeReport">
          <div class="detail-meta">
            <el-tag type="info">生成时间：{{ formatDateTime(activeReport.created_at) }}</el-tag>
            <el-tag type="success">新增报修：{{ activeReport.repair_new_count }}</el-tag>
            <el-tag type="warning">未处理报修：{{ activeReport.repair_pending_count }}</el-tag>
            <el-tag>新增访客：{{ activeReport.visitor_new_count }}</el-tag>
            <el-tag type="success">物业收缴：¥{{ formatAmount(activeReport.property_paid_amount) }}</el-tag>
          </div>

          <article class="markdown-body" v-html="activeReportHtml"></article>
        </template>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { DataAnalysis, Plus, Refresh, Clock } from '@element-plus/icons-vue'
import Navbar from '@/components/layout/Navbar.vue'
import { generateAIReport, getAIReportDetail, getAIReportList } from '@/api/admin'

const listLoading = ref(false)
const generating = ref(false)
const detailLoading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)

const detailVisible = ref(false)
const activeReport = ref(null)
const activeReportHtml = computed(() => markdownToHtml(activeReport.value?.report || ''))

function formatDateTime(value) {
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss')
}

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function isTimeoutError(error) {
  if (!error) return false
  return String(error.code || '') === 'ECONNABORTED' || /timeout/i.test(String(error.message || ''))
}

function isNetworkError(error) {
  if (!error) return false
  if (error.response) return false
  if (isTimeoutError(error)) return false
  return String(error.code || '') === 'ERR_NETWORK' || /Network Error/i.test(String(error.message || ''))
}

async function fetchList() {
  listLoading.value = true
  try {
    const res = await getAIReportList({
      page: page.value,
      size: size.value
    })
    list.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    if (isNetworkError(error)) {
      ElMessage.error('后端服务未连接，请先确认 smartcomunity 已启动在 8080 端口')
      return
    }
    ElMessage.error(error.response?.data?.msg || error.message || '获取报表列表失败')
  } finally {
    listLoading.value = false
  }
}

async function handleGenerate(showMessage) {
  if (generating.value) return
  generating.value = true
  try {
    const data = await generateAIReport()
    if (showMessage) {
      ElMessage.success('AI 报告生成成功')
    }
    await fetchList()
    if (showMessage && data?.id) {
      await openDetail(data.id)
    }
  } catch (error) {
    if (isTimeoutError(error)) {
      ElMessage.warning('请求超时，后台可能仍在生成。请稍后点击“刷新列表”查看结果')
      return
    }
    if (isNetworkError(error)) {
      ElMessage.error('后端服务未连接，请先确认 smartcomunity 已启动在 8080 端口')
      return
    }
    ElMessage.error(error.response?.data?.msg || error.message || 'AI 报告生成失败')
  } finally {
    generating.value = false
  }
}

async function openDetail(id) {
  detailVisible.value = true
  detailLoading.value = true
  try {
    activeReport.value = await getAIReportDetail(id)
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '获取报告详情失败')
  } finally {
    detailLoading.value = false
  }
}

onMounted(async () => {
  await fetchList()
})

function escapeHtml(text) {
  return String(text || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function inlineMarkdown(text) {
  return text
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.+?)\*/g, '<em>$1</em>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
}

function normalizeMarkdownInput(markdown) {
  let text = String(markdown || '').replace(/\r\n/g, '\n').trim()
  if (!text) return ''

  if (text.startsWith('```')) {
    const lines = text.split('\n')
    lines.shift()
    while (lines.length && !lines[lines.length - 1].trim()) lines.pop()
    if (lines.length && lines[lines.length - 1].trim() === '```') {
      lines.pop()
    }
    text = lines.join('\n').trim()
  }

  const lines = text.split('\n')
  if (lines.length && /^(markdown|md)$/i.test(lines[0].trim())) {
    lines.shift()
    text = lines.join('\n').trim()
  }

  return text
}

function isTableHeader(line) {
  return /\|/.test(line)
}

function isTableDivider(line) {
  return /^\|?[\s:-]+\|[\s|:-]*$/.test(line)
}

function isTableRow(line) {
  return /\|/.test(line)
}

function splitTableRow(line) {
  let content = line.trim()
  if (content.startsWith('|')) content = content.slice(1)
  if (content.endsWith('|')) content = content.slice(0, -1)
  return content.split('|').map((cell) => cell.trim())
}

function renderTable(headers, rows) {
  const head = `<thead><tr>${headers.map((h) => `<th>${inlineMarkdown(h)}</th>`).join('')}</tr></thead>`
  const body = `<tbody>${rows.map((row) => `<tr>${row.map((cell) => `<td>${inlineMarkdown(cell)}</td>`).join('')}</tr>`).join('')}</tbody>`
  return `<table>${head}${body}</table>`
}

function markdownToHtml(markdown) {
  const normalized = normalizeMarkdownInput(markdown)
  if (!normalized) return '<p>暂无内容</p>'

  const lines = escapeHtml(normalized).split('\n')
  const html = []
  let inUl = false
  let inOl = false

  const closeLists = () => {
    if (inUl) {
      html.push('</ul>')
      inUl = false
    }
    if (inOl) {
      html.push('</ol>')
      inOl = false
    }
  }

  for (let i = 0; i < lines.length; i += 1) {
    const line = lines[i].trim()
    if (!line) {
      closeLists()
      continue
    }

    if (isTableHeader(line) && i + 1 < lines.length && isTableDivider(lines[i + 1].trim())) {
      closeLists()
      const headers = splitTableRow(line)
      const rows = []
      i += 2
      for (; i < lines.length; i += 1) {
        const rowLine = lines[i].trim()
        if (!rowLine) {
          i -= 1
          break
        }
        if (!isTableRow(rowLine)) {
          i -= 1
          break
        }
        if (isTableDivider(rowLine)) continue
        rows.push(splitTableRow(rowLine))
      }
      html.push(renderTable(headers, rows))
      continue
    }

    if (/^###\s+/.test(line)) {
      closeLists()
      html.push(`<h3>${inlineMarkdown(line.replace(/^###\s+/, ''))}</h3>`)
      continue
    }
    if (/^##\s+/.test(line)) {
      closeLists()
      html.push(`<h2>${inlineMarkdown(line.replace(/^##\s+/, ''))}</h2>`)
      continue
    }
    if (/^#\s+/.test(line)) {
      closeLists()
      html.push(`<h1>${inlineMarkdown(line.replace(/^#\s+/, ''))}</h1>`)
      continue
    }
    if (/^>\s+/.test(line)) {
      closeLists()
      html.push(`<blockquote>${inlineMarkdown(line.replace(/^>\s+/, ''))}</blockquote>`)
      continue
    }
    if (/^[-*]\s+/.test(line)) {
      if (inOl) {
        html.push('</ol>')
        inOl = false
      }
      if (!inUl) {
        html.push('<ul>')
        inUl = true
      }
      html.push(`<li>${inlineMarkdown(line.replace(/^[-*]\s+/, ''))}</li>`)
      continue
    }
    if (/^\d+\.\s+/.test(line)) {
      if (inUl) {
        html.push('</ul>')
        inUl = false
      }
      if (!inOl) {
        html.push('<ol>')
        inOl = true
      }
      html.push(`<li>${inlineMarkdown(line.replace(/^\d+\.\s+/, ''))}</li>`)
      continue
    }
    if (/^---+$/.test(line)) {
      closeLists()
      html.push('<hr />')
      continue
    }

    closeLists()
    html.push(`<p>${inlineMarkdown(line)}</p>`)
  }

  closeLists()
  return html.join('')
}
</script>

<style scoped>
.ai-report-page {
  min-height: 100vh;
  padding-bottom: 40px;
}

.report-list-card {
  border-radius: 14px;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 22px;
  font-weight: 700;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.time-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-secondary);
}

.pagination-wrap {
  margin-top: 18px;
  display: flex;
  justify-content: flex-end;
}

.detail-wrap {
  height: 100%;
  overflow: auto;
}

.detail-meta {
  margin-bottom: 14px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.markdown-body h1),
:deep(.markdown-body h2),
:deep(.markdown-body h3) {
  margin: 12px 0 8px;
  color: #1f7a4d;
}

:deep(.markdown-body p),
:deep(.markdown-body li),
:deep(.markdown-body blockquote),
:deep(.markdown-body td),
:deep(.markdown-body th) {
  line-height: 1.9;
  font-size: 15px;
}

:deep(.markdown-body blockquote) {
  margin: 10px 0;
  padding: 10px 14px;
  border-left: 4px solid #b7d8c6;
  background: #f8fbf9;
}

:deep(.markdown-body code) {
  padding: 2px 6px;
  border-radius: 6px;
  background: #f2f5f7;
}

:deep(.markdown-body table) {
  width: 100%;
  border-collapse: collapse;
  margin: 10px 0 14px;
}

:deep(.markdown-body th),
:deep(.markdown-body td) {
  border: 1px solid #e7edf3;
  padding: 8px 10px;
  text-align: left;
}

:deep(.markdown-body th) {
  background: #f4f8f7;
}

@media (max-width: 960px) {
  .panel-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .panel-title {
    font-size: 20px;
  }
}
</style>
