<template>
  <div class="detail-page">
    <Navbar />
    
    <div class="container">
      <div class="product-detail card" v-if="product">
        <div class="detail-image">
          <img :src="product.image_url || 'https://via.placeholder.com/400'" :alt="product.name" />
        </div>
        
        <div class="detail-info">
          <h1 class="detail-title">{{ product.name }}</h1>
          <p class="detail-desc">{{ product.description }}</p>
          
          <div class="detail-price">
            <span class="price-current">¥{{ product.price }}</span>
            <span class="price-original" v-if="product.original_price">¥{{ product.original_price }}</span>
            <span class="tag tag-success" v-if="product.is_promotion">促销中</span>
          </div>
          
          <div class="detail-meta">
            <div class="meta-item">
              <span class="meta-label">销量:</span>
              <span class="meta-value">{{ product.sales }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-label">库存:</span>
              <span class="meta-value">{{ product.stock }}</span>
            </div>
          </div>
          
          <div class="detail-actions">
            <!-- Fix: Explicitly pass false for silent mode -->
            <button class="btn btn-lg btn-primary" @click="addToCart(false)">加入购物车</button>
            <button 
                class="btn btn-lg" 
                :class="isFavorite ? 'btn-secondary' : 'btn-outline'" 
                @click="toggleFavorite"
            >
                {{ isFavorite ? '已收藏' : '收藏商品' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getProductDetail } from '@/api/product'
import { addToCart as addToCartApi } from '@/api/order'
import { addFavorite, deleteFavorite, checkFavorite } from '@/api/favorite'
import { useCartStore } from '@/stores/cart'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const product = ref(null)
const isFavorite = ref(false)

const addToCart = async (silent) => {
  // Ensure silent is boolean
  if (typeof silent !== 'boolean') silent = false
  
  try {
    await addToCartApi({
      product_id: product.value.id,
      quantity: 1
    })
    if (!silent) {
      ElMessage.success('已添加到购物车！')
    }
    cartStore.fetchCart()
  } catch (error) {
    if (error.response?.status === 401) {
       ElMessage.warning('请先登录')
       router.push('/login')
    } else {
       ElMessage.error('添加失败: ' + error.message)
    }
  }
}

const toggleFavorite = async () => {
    try {
        if (isFavorite.value) {
            await deleteFavorite({ product_id: product.value.id })
            isFavorite.value = false
        } else {
            await addFavorite({ product_id: product.value.id })
            isFavorite.value = true
        }
    } catch (e) {
        if (e.response?.status === 401) {
            router.push('/login')
        } else {
            ElMessage.error('操作失败')
        }
    }
}

const checkFavStatus = async () => {
    try {
        const res = await checkFavorite(product.value.id)
        isFavorite.value = res.is_favorite
    } catch (e) {
        // ignore
    }
}

onMounted(async () => {
  try {
    product.value = await getProductDetail(route.params.id)
    checkFavStatus()
  } catch (error) {
    console.error('获取商品详情失败:', error)
  }
})
</script>

<style scoped>
.detail-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.product-detail {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-xl);
  padding: var(--spacing-xl);
}

.detail-image img {
  width: 100%;
  border-radius: var(--border-radius);
}

.detail-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.detail-title {
  font-size: 32px;
  font-weight: 700;
}

.detail-desc {
  color: var(--text-secondary);
  line-height: 1.8;
}

.detail-price {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.price-current {
  font-size: 36px;
  font-weight: 700;
  color: var(--danger-color);
}

.price-original {
  font-size: var(--font-size-lg);
  color: var(--text-light);
  text-decoration: line-through;
}

.detail-meta {
  display: flex;
  gap: var(--spacing-xl);
}

.meta-item {
  display: flex;
  gap: var(--spacing-sm);
}

.meta-label {
  color: var(--text-secondary);
}

.meta-value {
  font-weight: 600;
}

.detail-actions {
  display: flex;
  gap: var(--spacing-md);
  margin-top: auto;
}

@media (max-width: 768px) {
  .product-detail {
    grid-template-columns: 1fr;
  }
}
</style>
