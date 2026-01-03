<template>
  <div class="create-order-page">
    <Navbar />
    <div class="container">
      <h1 class="page-title">确认订单</h1>

      <div class="section card">
        <h3>选择服务门店</h3>
        <div class="store-selector">
          <select v-model="selectedStoreId" class="input select">
            <option :value="0" disabled>请选择提货门店</option>
            <option v-for="store in storeList" :key="store.id" :value="store.id">
              {{ store.name }} ({{ store.address }})
            </option>
          </select>
        </div>
      </div>

      <div class="section card">
        <h3>商品清单</h3>
        <div class="order-items">
          <div class="item" v-for="item in cartItems" :key="item.id">
            <img :src="item.product.image_url" class="item-img">
            <div class="item-info">
              <div class="name">{{ item.product.name }}</div>
              <div class="price">¥{{ item.product.price }} x {{ item.quantity }}</div>
            </div>
            <div class="item-total">¥{{ (item.product.price * item.quantity).toFixed(2) }}</div>
          </div>
        </div>
      </div>

      <div class="footer-bar card">
        <div class="total-info">
          <span>共 {{ totalCount }} 件, 合计:</span>
          <span class="total-price">¥{{ totalPrice }}</span>
        </div>
        <button class="btn btn-primary btn-lg" :disabled="submitting" @click="submitOrder">
          {{ submitting ? '提交中...' : '提交订单' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { useUserStore } from '@/stores/user'
import { getCartList, createOrder } from '@/api/order'
import { getStoreList } from '@/api/service'
import { ElMessage } from 'element-plus'


const router = useRouter()
const userStore = useUserStore()
const cartItems = ref([])
const storeList = ref([])
const selectedStoreId = ref(0)
const submitting = ref(false)

const totalCount = computed(() => cartItems.value.reduce((sum, i) => sum + i.quantity, 0))
const totalPrice = computed(() => cartItems.value.reduce((sum, i) => sum + i.product.price * i.quantity, 0).toFixed(2))

const loadData = async () => {
  try {
    const [cartRes, storeRes] = await Promise.all([getCartList(), getStoreList()])
    cartItems.value = cartRes
    storeList.value = storeRes || []
    if (storeList.value.length > 0) {
        selectedStoreId.value = storeList.value[0].id // 默认选中第一个
    }
    
    if (cartItems.value.length === 0) {
      ElMessage.warning('购物车为空')
      router.push('/mall')
    }
  } catch (e) {
    console.error(e)
  }
}

const submitOrder = async () => {
  if (!selectedStoreId.value) {
    ElMessage.warning('请选择服务门店')
    return
  }
  
  submitting.value = true
  try {
    // 构造后端需要的 items 参数
    const items = cartItems.value.map(item => ({
        cart_id: item.id,
        quantity: item.quantity
    }))

    await createOrder({
        store_id: selectedStoreId.value,
        items: items
    })
    ElMessage.success('订单创建成功')
    router.push('/order') 
  } catch (e) {
    ElMessage.error('创建失败: ' + (e.response?.data?.msg || e.message))
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.create-order-page { min-height: 100vh; padding-bottom: 80px; }
.section { margin-bottom: 24px; padding: 24px; }
.section h3 { margin-bottom: 16px; border-bottom: 1px solid #eee; padding-bottom: 12px;}

.order-items { display: flex; flex-direction: column; gap: 16px; }
.item { display: flex; gap: 16px; align-items: center; }
.item-img { width: 60px; height: 60px; object-fit: cover; border-radius: 4px; }
.item-info { flex: 1; }
.item-total { font-weight: bold; }

.footer-bar {
  position: fixed; bottom: 0; left: 0; right: 0;
  padding: 16px 24px;
  display: flex; justify-content: flex-end; align-items: center; gap: 24px;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
  border-radius: 0;
}
.total-price { font-size: 24px; color: var(--danger-color); font-weight: bold; margin-left: 8px; }
</style>
