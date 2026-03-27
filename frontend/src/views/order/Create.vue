<template>
  <div class="create-order-page">
    <Navbar />

    <div class="container">
      <h1 class="page-title">确认订单</h1>

      <el-alert type="success" show-icon :closable="false" class="tip-alert">
        <template #title>
          当前积分 {{ userStore.userInfo.green_points || 0 }}，余额 ¥{{ formatAmount(userStore.userInfo.balance || 0) }}。
          订单支付时将按 {{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元 自动优先抵扣。
        </template>
      </el-alert>

      <div class="section card">
        <h3>选择服务门店</h3>
        <el-select v-model="selectedStoreId" placeholder="请选择提货门店" style="width: 100%">
          <el-option v-for="store in storeList" :key="store.id" :label="`${store.name} (${store.address})`" :value="store.id" />
        </el-select>
      </div>

      <div class="section card">
        <h3>商品清单</h3>
        <div class="order-items">
          <div class="item" v-for="item in cartItems" :key="item.id">
            <img :src="item.product.image_url" class="item-img" />
            <div class="item-info">
              <div class="name">{{ item.product.name }}</div>
              <div class="price">¥{{ formatAmount(item.product.price) }} x {{ item.quantity }}</div>
            </div>
            <div class="item-total">¥{{ formatAmount(item.product.price * item.quantity) }}</div>
          </div>
        </div>
      </div>

      <div class="section card preview-box">
        <h3>支付预估</h3>
        <div class="preview-line">
          <span>订单总额</span>
          <strong>¥{{ totalPrice }}</strong>
        </div>
        <div class="preview-line">
          <span>预计抵扣积分</span>
          <strong>{{ paymentPreview.points }}</strong>
        </div>
        <div class="preview-line">
          <span>预计扣除余额</span>
          <strong>¥{{ formatAmount(paymentPreview.balance) }}</strong>
        </div>
      </div>

      <div class="footer-bar card">
        <div class="total-info">
          <span>共 {{ totalCount }} 件商品</span>
          <span class="total-price">¥{{ totalPrice }}</span>
        </div>
        <el-button type="primary" size="large" :loading="submitting" @click="submitOrder">
          {{ submitting ? '提交中...' : '提交订单' }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import { useUserStore } from '@/stores/user'
import { getCartList, createOrder } from '@/api/order'
import { getStoreList } from '@/api/service'
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from '@/utils/payment'

const router = useRouter()
const userStore = useUserStore()
const cartItems = ref([])
const storeList = ref([])
const selectedStoreId = ref()
const submitting = ref(false)

const totalCount = computed(() => cartItems.value.reduce((sum, item) => sum + item.quantity, 0))
const totalPriceNumber = computed(() => cartItems.value.reduce((sum, item) => sum + item.product.price * item.quantity, 0))
const totalPrice = computed(() => formatAmount(totalPriceNumber.value))
const paymentPreview = computed(() => getMixedPaymentPreview(totalPriceNumber.value, userStore.userInfo.green_points))

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

async function loadData() {
  const [cartRes, storeRes] = await Promise.all([getCartList(), getStoreList(), userStore.fetchUserInfo()])
  cartItems.value = cartRes || []
  storeList.value = storeRes || []
  if (storeList.value.length > 0) {
    selectedStoreId.value = storeList.value[0].id
  }
  if (cartItems.value.length === 0) {
    ElMessage.warning('购物车为空')
    router.push('/mall')
  }
}

async function submitOrder() {
  if (!selectedStoreId.value) {
    ElMessage.warning('请选择门店')
    return
  }

  submitting.value = true
  try {
    const items = cartItems.value.map((item) => ({
      cart_id: item.id,
      quantity: item.quantity
    }))

    await createOrder({
      store_id: selectedStoreId.value,
      items
    })
    ElMessage.success('订单创建成功')
    router.push('/order')
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '创建订单失败')
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.create-order-page {
  min-height: 100vh;
  padding-bottom: 100px;
}

.tip-alert,
.section {
  margin-bottom: 20px;
}

.order-items {
  display: grid;
  gap: 16px;
}

.item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.item-img {
  width: 64px;
  height: 64px;
  object-fit: cover;
  border-radius: 8px;
}

.item-info {
  flex: 1;
}

.preview-box {
  display: grid;
  gap: 12px;
}

.preview-line {
  display: flex;
  justify-content: space-between;
}

.footer-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 24px;
  border-radius: 0;
  padding: 16px 24px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.08);
}

.total-price {
  margin-left: 12px;
  font-size: 24px;
  font-weight: 700;
  color: var(--danger-color);
}
</style>
