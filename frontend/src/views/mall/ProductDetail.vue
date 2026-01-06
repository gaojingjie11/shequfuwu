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

      <!-- Comment Section -->
      <div class="comments-section card">
          <h2 class="section-title">商品评价</h2>
          
          <!-- Post Comment Form -->
          <div class="comment-form">
               <el-input 
                  v-model="commentContent" 
                  type="textarea" 
                  :rows="3" 
                  placeholder="写下你的评价..." 
               />
               <div class="form-footer">
                   <el-rate v-model="commentRating" />
                   <el-button type="primary" @click="submitComment">发表评论</el-button>
               </div>
          </div>

          <!-- Comment List -->
          <div class="comment-list">
              <div v-for="item in comments" :key="item.id" class="comment-item">
                  <div class="comment-avatar">
                      <img :src="item.user?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
                  </div>
                  <div class="comment-body">
                      <div class="comment-header">
                          <span class="username">{{ item.user?.real_name || item.user?.username || '匿名用户' }}</span>
                          <el-rate :model-value="item.rating" disabled size="small" />
                          <span class="time">{{ formatDate(item.created_at) }}</span>
                      </div>
                      <div class="comment-content">{{ item.content }}</div>
                  </div>
              </div>
              <el-empty v-if="comments.length === 0" description="暂无评价" />
          </div>

          <!-- Pagination -->
          <div class="pagination-container" v-if="total > 0">
               <el-pagination
                 v-model:current-page="page"
                 v-model:page-size="size"
                 :total="total"
                 layout="prev, pager, next"
                 @current-change="fetchComments"
               />
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
    fetchComments()
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

/* Comment Styles */
.comments-section {
    margin-top: 20px;
    padding: 20px;
}
.section-title {
    font-size: 20px;
    margin-bottom: 20px;
    font-weight: 600;
}
.comment-form {
    margin-bottom: 30px;
    background: #f8f9fa;
    padding: 20px;
    border-radius: 8px;
}
.form-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 10px;
}
.comment-item {
    display: flex;
    gap: 15px;
    padding: 20px 0;
    border-bottom: 1px solid #eee;
}
.comment-avatar img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
}
.comment-body {
    flex: 1;
}
.comment-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 5px;
}
.username {
    font-weight: 600;
    font-size: 14px;
}
.time {
    margin-left: auto;
    color: #999;
    font-size: 12px;
}
.comment-content {
    color: #333;
    line-height: 1.6;
}
.pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
}
</style>
