<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">车位管理</h1>
        <div class="actions">
          <el-button type="primary" @click="openCreateModal">新增车位</el-button>
        </div>
      </div>

      <!-- Stats -->
      <div class="stats-cards">
        <el-row :gutter="20">
          <el-col :span="8">
             <div class="stat-card">
              <div class="stat-label">总车位</div>
              <div class="stat-value">{{ stats.total || 0 }}</div>
            </div>
          </el-col>
           <el-col :span="8">
             <div class="stat-card">
              <div class="stat-label">已占用</div>
              <div class="stat-value highlight">{{ stats.used || 0 }}</div>
            </div>
          </el-col>
           <el-col :span="8">
             <div class="stat-card">
              <div class="stat-label">空闲</div>
              <div class="stat-value success">{{ stats.free || 0 }}</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <div class="table-container card">
         <el-table :data="list" stripe border style="width: 100%" v-loading="loadingData">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="parking_no" label="车位号" />
          <el-table-column prop="status" label="状态">
             <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'warning' : 'success'">
                  {{ scope.row.status === 1 ? '已占用' : '空闲' }}
                </el-tag>
             </template>
          </el-table-column>
          <el-table-column prop="user_id" label="绑定用户ID">
             <template #default="scope">{{ scope.row.user_id || '-' }}</template>
          </el-table-column>
          <el-table-column prop="car_plate" label="绑定车牌">
             <template #default="scope">{{ scope.row.car_plate || '-' }}</template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
             <template #default="scope">
               <el-button size="small" type="primary" @click="openAssignModal(scope.row)">
                  {{ scope.row.status === 1 ? '修改/解绑' : '分配' }}
               </el-button>
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
            @size-change="fetchData"
            @current-change="fetchData"
          />
        </div>
      </div>
      
      <!-- Assign Modal -->
      <el-dialog
        v-model="showModal"
        :title="`分配车位: ${currentItem.parking_no}`"
        width="400px"
      >
        <el-alert
          title="请输入用户ID和车牌号进行绑定。输入ID为0则解绑。"
          type="info"
          show-icon
          :closable="false"
          style="margin-bottom: 20px"
        />
        <el-form label-position="top">
          <el-form-item label="用户ID (User ID)">
            <el-input v-model.number="form.user_id" type="number" placeholder="输入用户ID" />
          </el-form-item>
          <el-form-item label="车牌号 (Car Plate)">
            <el-input v-model="form.car_plate" placeholder="例如: 辽A66666" />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showModal = false">取消</el-button>
            <el-button type="primary" @click="saveAssign" :loading="loading">
              保存
            </el-button>
          </span>
        </template>
      </el-dialog>

      <!-- Create Modal -->
      <el-dialog v-model="showCreateModal" title="新增车位" width="400px">
        <el-form :model="createForm" label-width="80px">
          <el-form-item label="车位号">
            <el-input v-model="createForm.parking_no" placeholder="例如: A-001" />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showCreateModal = false">取消</el-button>
            <el-button type="primary" @click="saveCreate" :loading="creating">保存</el-button>
          </span>
        </template>
      </el-dialog>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminParkingList, getParkingStats, assignParking, createParking } from '@/api/admin'
import { ElMessage } from 'element-plus'

const list = ref([])
const stats = ref({})
const total = ref(0)
const page = ref(1)
const size = ref(20)
const loadingData = ref(false)

const fetchData = async () => {
  loadingData.value = true
  try {
    const res = await getAdminParkingList({ page: page.value, size: size.value })
    list.value = res.list || []
    total.value = res.total
    
    // fetch stats
    const statsRes = await getParkingStats()
    stats.value = statsRes || {}
  } catch (error) {
    console.error(error)
  } finally {
    loadingData.value = false
  }
}

// Modal
const showModal = ref(false)
const currentItem = ref({})
const form = reactive({
  user_id: '',
  car_plate: ''
})
const loading = ref(false)

const openAssignModal = (item) => {
  currentItem.value = item
  form.user_id = item.user_id || ''
  form.car_plate = item.car_plate || ''
  showModal.value = true
}

const saveAssign = async () => {
  loading.value = true
  try {
    await assignParking({
      id: currentItem.value.id,
      user_id: Number(form.user_id),
      car_plate: form.car_plate
    })
    ElMessage.success('操作成功')
    showModal.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('操作失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// Create Logic
const showCreateModal = ref(false)
const createForm = reactive({ parking_no: '' })
const creating = ref(false)

const openCreateModal = () => {
    createForm.parking_no = ''
    showCreateModal.value = true
}

const saveCreate = async () => {
    if (!createForm.parking_no) {
        ElMessage.warning('请输入车位号')
        return
    }
    creating.value = true
    try {
        await createParking({ parking_no: createForm.parking_no })
        ElMessage.success('创建成功')
        showCreateModal.value = false
        fetchData()
    } catch (e) {
        ElMessage.error(e.response?.data?.msg || '创建失败')
    } finally {
        creating.value = false
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

.stats-cards {
  margin-bottom: var(--spacing-lg);
}

.stat-card {
  background: white;
  padding: var(--spacing-lg);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  text-align: center;
}

.stat-label {
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-xs);
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
}
.stat-value.highlight { color: var(--primary-color); }
.stat-value.success { color: var(--success-color); }

.table-container {
  padding: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
