<template>
  <div class="cart-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">购物车</h1>
      
      <div v-if="cartItems.length > 0">
        <div class="cart-list">
          <div class="cart-item card" v-for="item in cartItems" :key="item.id">
            <img :src="item.product.image_url || 'https://via.placeholder.com/100'" class="item-image" />
            <div class="item-info">
              <div class="item-name">{{ item.product.name }}</div>
              <div class="item-price">¥{{ item.product.price }}</div>
            </div>
            
            <div class="quantity-control">
              <button class="btn btn-sm btn-outline" :disabled="item.quantity <= 1" @click="changeQuantity(item, -1)">-</button>
              <input type="number" class="qty-input" v-model.number="item.quantity" @change="manualUpdate(item)" min="1">
              <button class="btn btn-sm btn-outline" @click="changeQuantity(item, 1)">+</button>
            </div>

            <button class="btn btn-sm btn-outline ml-4" @click="removeItem(item.id)">删除</button>
          </div>
        </div>
        
        <div class="cart-footer card">
          <div class="total">
            <span>总计:</span>
            <span class="total-price">¥{{ totalPrice }}</span>
          </div>
          <button class="btn btn-primary btn-lg" @click="checkout">结算</button>
        </div>
      </div>
      
      <div class="empty-state" v-else>
        <div class="empty-state-icon">🛒</div>
        <p>购物车是空的</p>
        <button class="btn btn-primary" @click="$router.push('/mall')">去购物</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getCartList, deleteCartItem, updateCartQuantity } from '@/api/order'
import { useCartStore } from '@/stores/cart'
import { ElMessage } from 'element-plus'

const router = useRouter()
const cartStore = useCartStore()
const cartItems = ref([])

const totalPrice = computed(() => {
  return cartItems.value.reduce((sum, item) => {
    return sum + item.product.price * item.quantity
  }, 0).toFixed(2)
})

const fetchCart = async () => {
  try {
    cartItems.value = await getCartList()
    cartStore.fetchCart() // 同步 store 数量
  } catch (error) {
    console.error('获取购物车失败:', error)
  }
}

const changeQuantity = async (item, delta) => {
  const newQty = item.quantity + delta
  if (newQty < 1) return
  
  // 乐观更新
  item.quantity = newQty
  try {
    await updateCartQuantity(item.id, newQty)
  } catch (e) {
    item.quantity -= delta // 回滚
    console.error(e)
  }
}

const manualUpdate = async (item) => {
  if (item.quantity < 1) item.quantity = 1
  try {
    await updateCartQuantity(item.id, item.quantity)
  } catch (e) {
    console.error(e)
  }
}

const removeItem = async (id) => {
  try {
    await deleteCartItem(id)
    fetchCart()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const checkout = () => {
  router.push('/order/create')
}

onMounted(() => {
  fetchCart()
})
</script>

<style scoped>
.cart-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.cart-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.cart-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-lg);
}

.item-image {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: var(--border-radius-sm);
}

.item-info {
  flex: 1;
}

.item-name {
  font-size: var(--font-size-lg);
  font-weight: 600;
  margin-bottom: var(--spacing-sm);
}

.item-price {
  font-size: var(--font-size-lg);
  color: var(--danger-color);
  font-weight: 600;
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.qty-input {
  width: 50px;
  text-align: center;
  padding: 4px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.ml-4 { margin-left: 16px; }

.cart-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg);
}

.total {
  font-size: var(--font-size-xl);
  font-weight: 600;
}

.total-price {
  color: var(--danger-color);
  margin-left: var(--spacing-md);
}
</style>
