<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">物业费管理</h1>
        <div class="actions">
          <el-button type="primary" @click="openCreateModal">新增账单</el-button>
        </div>
      </div>

      <div class="table-container card">
        <el-table :data="list" stripe border style="width: 100%" v-loading="loadingData">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="user_id" label="用户ID" />
          <el-table-column prop="month" label="月份" />
          <el-table-column prop="amount" label="金额">
             <template #default="scope">¥{{ scope.row.amount }}</template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
             <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
                  {{ scope.row.status === 1 ? '已缴' : '未缴' }}
                </el-tag>
             </template>
          </el-table-column>
          <el-table-column prop="pay_time" label="支付时间">
            <template #default="scope">{{ scope.row.pay_time ? formatDate(scope.row.pay_time) : '-' }}</template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchData"
            @current-change="fetchData"
          />
        </div>
      </div>

      <!-- Create Dialog -->
      <el-dialog
        v-model="showModal"
        title="新增物业费账单"
        width="400px"
      >
        <el-form label-position="top">
          <el-form-item label="用户ID">
             <el-input v-model.number="form.user_id" type="number" placeholder="输入用户ID" />
          </el-form-item>
          <el-form-item label="月份 (例如: 2024-01)">
             <el-input v-model="form.month" placeholder="YYYY-MM" />
          </el-form-item>
          <el-form-item label="金额">
             <el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showModal = false">取消</el-button>
            <el-button type="primary" @click="saveCreate" :loading="loading">
              保存
            </el-button>
          </span>
        </template>
      </el-dialog>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminPropertyFeeList, createPropertyFee } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const list = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const loadingData = ref(false)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const fetchData = async () => {
  loadingData.value = true
  try {
    const res = await getAdminPropertyFeeList({ page: page.value, size: size.value })
    list.value = res.list || []
    total.value = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loadingData.value = false
  }
}

// Modal
const showModal = ref(false)
const loading = ref(false)
const form = reactive({
  user_id: '',
  month: dayjs().format('YYYY-MM'),
  amount: 200
})

const openCreateModal = () => {
  form.user_id = ''
  form.month = dayjs().format('YYYY-MM')
  form.amount = 200
  showModal.value = true
}

const saveCreate = async () => {
  if (!form.user_id || !form.month || !form.amount) {
    ElMessage.warning('请填写完整信息')
    return
  }
  loading.value = true
  try {
    await createPropertyFee({
      user_id: Number(form.user_id),
      month: form.month,
      amount: Number(form.amount)
    })
    ElMessage.success('创建成功')
    showModal.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-child-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.page-header {
  margin-bottom: var(--spacing-lg);
  display: flex;
  justify-content: space-between;
}

.table-container {
  padding: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
