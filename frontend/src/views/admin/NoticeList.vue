<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">公告管理</h1>
        <button class="btn btn-primary" @click="openModal()">+ 发布公告</button>
      </div>

      <div class="table-container card">
        <el-table :data="notices" style="width: 100%" stripe border>
          <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
          <el-table-column prop="content" label="内容摘要" min-width="200" show-overflow-tooltip />
          
          <el-table-column label="发布时间" width="160">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <button class="btn btn-sm btn-danger" @click="handleDelete(row.id)">删除</button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container mt-4">
            <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50]"
                layout="total, sizes, prev, pager, next, jumper"
                :total="total"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
            />
        </div>
      </div>


      <div class="modal-overlay" v-if="showModal">
        <div class="modal card">
          <h3>发布公告</h3>
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label>标题</label>
              <input v-model="form.title" class="input" required />
            </div>
            <div class="form-group">
              <label>内容</label>
              <textarea v-model="form.content" class="input textarea" required></textarea>
            </div>
            <div class="modal-actions">
              <button type="button" class="btn btn-secondary" @click="closeModal">取消</button>
              <button type="submit" class="btn btn-primary">发布</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getNoticeList } from '@/api/service'
import { createNotice, deleteNotice } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'

const notices = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const showModal = ref(false)
const form = ref({ title: '', content: '' })

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const fetchNotices = async () => {
  try {
    const res = await getNoticeList({
        page: currentPage.value,
        size: pageSize.value
    })
    // 兼容处理：如果是列表（首页）则 list=res, 如果是分页则 list=res.list
    if (res.list) {
        notices.value = res.list
        total.value = res.total
    } else {
        notices.value = res
        total.value = res.length
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSizeChange = (val) => {
    pageSize.value = val
    fetchNotices()
}

const handleCurrentChange = (val) => {
    currentPage.value = val
    fetchNotices()
}

const openModal = () => {
    form.value = { title: '', content: '' }
    showModal.value = true
}
const closeModal = () => showModal.value = false

const handleSubmit = async () => {
    try {
        await createNotice(form.value)
        ElMessage.success('发布成功')
        closeModal()
        fetchNotices()
    } catch (e) {
        ElMessage.error('发布失败')
    }
}

const handleDelete = async (id) => {
    try {
        await ElMessageBox.confirm('确定删除?', '删除确认', {
            confirmButtonText: '删除',
            cancelButtonText: '取消',
            type: 'warning'
        })
        await deleteNotice(id)
        ElMessage.success('删除成功')
        fetchNotices()
    } catch (e) {
        if (e !== 'cancel') {
            ElMessage.error('删除失败')
        }
    }
}

onMounted(fetchNotices)
</script>

<style scoped>
/* Reuse styles */
.admin-child-page { min-height: 100vh; padding-bottom: var(--spacing-xl); }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--spacing-lg); }
.table { width: 100%; border-collapse: collapse; }
.table th, .table td { padding: 12px; border-bottom: 1px solid #eee; text-align: left; }
.content-cell { max-width: 300px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; z-index: 2000; }
.modal { padding: 24px; width: 400px; max-width: 90%; background: #fff; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.form-group { margin-bottom: 16px; display: flex; flex-direction: column; }
.textarea { height: 100px; resize: vertical; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
.pagination-container { display: flex; justify-content: flex-end; padding-top: 20px; }
</style>
