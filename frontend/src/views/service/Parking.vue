<template>
  <div class="parking-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">车位管理</h1>
      
      <div v-if="parkingList && parkingList.length > 0">
          <div class="parking-info card mb-4" v-for="item in parkingList" :key="item.id">
            <h3>我的车位信息 ({{item.parking_no}})</h3>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">车位号：</span>
                <span class="info-value">{{ item.parking_no }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">状态：</span>
                <span class="tag" :class="item.status === 1 ? 'tag-success' : 'tag-warning'">
                  {{ item.status === 1 ? '已占用' : '空闲' }}
                </span>
              </div>
              <div class="info-item">
                <span class="info-label">车牌号：</span>
                <span class="info-value">{{ item.car_plate || '未绑定' }}</span>
              </div>
            </div>
            
            <div class="bind-form">
              <h4>{{ item.car_plate ? '更新车牌' : '绑定车牌' }}</h4>
              <div class="form-row">
                <div class="form-group">
                  <input 
                    v-model="item.editCarPlate" 
                    class="input" 
                    :placeholder="item.car_plate || '请输入车牌号'"
                  />
                </div>
                <button type="button" class="btn btn-primary" :disabled="loading" @click="handleBindCar(item)">
                  {{ loading ? '提交中...' : '更新' }}
                </button>
              </div>
            </div>
          </div>
      </div>
      
      <div class="empty-state card" v-else>
        <div class="empty-state-icon">🅿️</div>
        <p>您还没有分配车位</p>
        <p class="text-secondary">请联系物业管理员分配车位</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getMyParking, bindCar } from '@/api/service'
import { ElMessage } from 'element-plus'

const parkingList = ref([])
const loading = ref(false)

const fetchParking = async () => {
  try {
    const data = await getMyParking()
    // backend now returns list
    if (Array.isArray(data)) {
        parkingList.value = data.map(p => ({
            ...p,
            editCarPlate: p.car_plate // init buffer
        }))
    } else if (data) {
        // Fallback if backend returns object (legacy cache?)
        parkingList.value = [{...data, editCarPlate: data.car_plate}]
    } else {
        parkingList.value = []
    }
  } catch (error) {
    console.log('未分配车位或获取失败', error)
  }
}

const handleBindCar = async (item) => {
  if (!item.editCarPlate || !item.editCarPlate.trim()) {
    ElMessage.warning('请输入车牌号')
    return
  }
  
  loading.value = true
  try {
    await bindCar({ 
        parking_id: item.id, 
        car_plate: item.editCarPlate 
    })
    ElMessage.success('更新成功！')
    await fetchParking()
  } catch (error) {
    ElMessage.error('操作失败：' + (error.response?.data?.msg || error.message))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchParking()
})
</script>

<style scoped>
.parking-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.parking-info {
  max-width: 600px;
  margin: 0 auto;
}

.parking-info h3 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--spacing-lg);
}

.info-grid {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-weight: 500;
  color: var(--text-secondary);
}

.info-value {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--text-primary);
}

.bind-form h4 {
  font-size: var(--font-size-lg);
  margin-bottom: var(--spacing-md);
}

.bind-form form {
  display: flex;
  gap: var(--spacing-md);
  align-items: flex-end;
}

.bind-form .form-group {
  flex: 1;
}

.mb-4 { margin-bottom: 24px; }
.form-row { display: flex; gap: 16px; align-items: center; }
.form-row .form-group { flex: 1; margin: 0; }

.empty-state {
  max-width: 400px;
  margin: var(--spacing-xl) auto;
  text-align: center;
  padding: var(--spacing-xl);
}
</style>
