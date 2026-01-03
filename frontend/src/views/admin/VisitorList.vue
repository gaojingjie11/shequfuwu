<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">访客审核</h1>
      </div>

      <div class="table-container card">
        <el-table :data="visitors" stripe border style="width: 100%" v-loading="loadingData">
          <el-table-column prop="name" label="访客姓名" width="120" />
          <el-table-column prop="mobile" label="电话" width="150" />
          <el-table-column prop="reason" label="事由" show-overflow-tooltip />
          <el-table-column prop="visit_time" label="预计时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.visit_time) }}</template>
          </el-table-column>
          <el-table-column prop="created_at" label="提交时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.created_at) }}</template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
             <template #default="scope">
               <el-tag :type="getStatusType(scope.row.status)">
                 {{ getStatusText(scope.row.status) }}
               </el-tag>
             </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
             <template #default="scope">
               <div v-if="scope.row.status === 0">
                 <el-button size="small" type="success" @click="handleAudit(scope.row, 1)">通过</el-button>
                 <el-button size="small" type="danger" @click="handleAudit(scope.row, 2)">拒绝</el-button>
               </div>
               <span v-else class="text-secondary">
                 {{ scope.row.audit_remark || '已处理' }}
               </span>
             </template>
          </el-table-column>
        </el-table>

         <div class="pagination-container">
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
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminVisitorList, auditVisitor } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'

const visitors = ref([])
const loadingData = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusText = (s) => {
    const map = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
    return map[s] || s
}
const getStatusType = (s) => {
    const map = { 0: 'warning', 1: 'success', 2: 'danger' }
    return map[s] || 'info'
}

const fetchList = async () => {
  loadingData.value = true
  try {
    const res = await getAdminVisitorList({ page: page.value, size: size.value })
    if (res.list) {
        visitors.value = res.list
        total.value = res.total
    } else {
        visitors.value = res // Fallback
        total.value = res.length
    }
  } catch (error) {
    console.error(error)
  } finally {
      loadingData.value = false
  }
}

const handleAudit = async (item, status) => {
    if (status === 1) {
        // Approve
         try {
             await auditVisitor({
                id: item.id,
                status: 1,
                audit_remark: '同意'
            })
            ElMessage.success('操作成功')
            fetchList()
         } catch(e) {
             ElMessage.error('操作失败')
         }
    } else {
        // Reject
        ElMessageBox.prompt('请输入拒绝原因', '拒绝审核', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            inputPattern: /.+/,
            inputErrorMessage: '原因不能为空'
        }).then(async ({ value }) => {
            try {
                await auditVisitor({
                    id: item.id,
                    status: 2,
                    audit_remark: value
                })
                ElMessage.success('已拒绝')
                fetchList()
            } catch (e) {
                ElMessage.error('操作失败')
            }
        }).catch(() => {
            // Cancelled
        })
    }
}

onMounted(() => {
    fetchList()
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

.text-secondary {
    color: #909399;
    font-size: 13px;
}
</style>
