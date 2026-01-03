<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">用户管理</h1>
        <div class="actions">
          <el-input 
            v-model="keyword" 
            placeholder="搜索用户名/手机号/姓名" 
            style="width: 300px; margin-right: 10px"
            @keyup.enter="handleSearch"
            clearable
          >
            <template #append>
              <el-button @click="handleSearch"><el-icon><Search /></el-icon></el-button>
            </template>
          </el-input>
        </div>
      </div>

      <div class="table-container card">
        <el-table :data="users" stripe border style="width: 100%" v-loading="loadingData">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="real_name" label="姓名">
            <template #default="scope">{{ scope.row.real_name || '-' }}</template>
          </el-table-column>
          <el-table-column prop="mobile" label="手机号" width="120" />
          <el-table-column prop="balance" label="余额" width="100">
            <template #default="scope">¥{{ scope.row.balance }}</template>
          </el-table-column>
          <el-table-column prop="role" label="角色">
             <template #default="scope">
               <el-tag :type="getRoleType(scope.row.role)">
                 {{ getRoleText(scope.row.role) }}
               </el-tag>
             </template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                {{ scope.row.status === 1 ? '正常' : '冻结' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="注册时间" width="160">
             <template #default="scope">{{ formatDate(scope.row.created_at) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button size="small" @click="openEditModal(scope.row)">编辑</el-button>
              <el-popconfirm
                :title="`确定要${scope.row.status === 1 ? '冻结' : '解冻'}该用户吗？`"
                @confirm="handleFreeze(scope.row)"
              >
                <template #reference>
                  <el-button 
                    size="small" 
                    :type="scope.row.status === 1 ? 'danger' : 'success'"
                  >
                    {{ scope.row.status === 1 ? '冻结' : '解冻' }}
                  </el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSearch"
            @current-change="handleSearch"
          />
        </div>
      </div>

      <!-- Edit Dialog -->
      <el-dialog
        v-model="showModal"
        title="编辑用户"
        width="500px"
      >
        <el-form label-width="120px">
          <el-form-item label="当前用户">
            <strong>{{ currentUser.username }}</strong>
          </el-form-item>
           <el-form-item label="角色">
             <el-select v-model="editForm.role" class="m-2" placeholder="Select">
                <el-option label="普通用户" value="user" />
                <el-option label="超级管理员" value="admin" />
                <el-option label="商场管理员" value="store" />
                <el-option label="物业管理员" value="property" />
             </el-select>
          </el-form-item>
           <el-form-item label="余额调整">
             <el-input-number v-model="editForm.balanceAmount" :precision="2" :step="10" />
             <div class="tips">正数充值，负数扣费 (当前: ¥{{ currentUser.balance }})</div>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showModal = false">取消</el-button>
            <el-button type="primary" @click="saveEdit" :loading="loading">
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
import { getUserList, freezeUser, assignRole, updateUserBalance } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const users = ref([])
const keyword = ref('')
const total = ref(0) 
const page = ref(1)
const size = ref(20)
const loadingData = ref(false)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getRoleText = (role) => {
  const map = { 'admin': '管理员', 'user': '用户', 'property': '物业', 'store': '商场' }
  return map[role] || role
}
const getRoleType = (role) => {
  if (role === 'admin') return ''
  if (role === 'property') return 'info'
  if (role === 'store') return 'warning'
  return 'success'
}

const handleSearch = () => {
    fetchUsers()
}

const fetchUsers = async () => {
  loadingData.value = true
  try {
    const res = await getUserList({ 
      page: page.value, 
      size: size.value, 
      keyword: keyword.value 
    })
    users.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取用户列表失败', error)
  } finally {
    loadingData.value = false
  }
}

const handleFreeze = async (user) => {
  try {
    const newStatus = user.status === 1 ? 0 : 1
    await freezeUser({ id: user.id, status: newStatus })
    ElMessage.success('操作成功')
    await fetchUsers()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// Modal Logic
const showModal = ref(false)
const loading = ref(false)
const currentUser = ref({})
const editForm = reactive({
  role: 'user',
  balanceAmount: 0
})

const openEditModal = (user) => {
  currentUser.value = user
  editForm.role = user.role
  editForm.balanceAmount = 0
  showModal.value = true
}

const saveEdit = async () => {
  loading.value = true
  try {
    // 1. Update Role if changed
    if (editForm.role !== currentUser.value.role) {
      await assignRole({ user_id: currentUser.value.id, role_code: editForm.role })
    }
    // 2. Update Balance if entered
    if (editForm.balanceAmount !== 0) {
      await updateUserBalance({ user_id: currentUser.value.id, amount: editForm.balanceAmount })
    }
    
    ElMessage.success('修改成功')
    showModal.value = false
    fetchUsers()
  } catch (error) {
    console.error(error)
    ElMessage.error('操作失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.admin-child-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.table-container {
  padding: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.tips {
    font-size: 12px;
    color: #909399;
    margin-top: 5px;
}
</style>
