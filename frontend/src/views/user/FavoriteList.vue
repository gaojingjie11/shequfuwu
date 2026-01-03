<template>
  <div class="favorite-page">
    <Navbar />
    <div class="container">
      <h1 class="page-title">我的收藏</h1>

      <div class="product-grid" v-if="list.length > 0">
        <div class="product-card card" v-for="item in list" :key="item.id">
            <!-- 假设后端返回结构是 item.product -->
          <div class="card-img" @click="goDetail(item.product.id)">
            <img :src="item.product.image_url || 'https://via.placeholder.com/200'" :alt="item.product.name" />
          </div>
          <div class="card-body">
            <h3 class="product-name" @click="goDetail(item.product.id)">{{ item.product.name }}</h3>
            <div class="price-row">
              <span class="price">¥{{ item.product.price }}</span>
              <span class="original-price" v-if="item.product.original_price">¥{{ item.product.original_price }}</span>
            </div>
            <div class="action-row">
                <button class="btn btn-sm btn-outline-danger" @click="removeFavorite(item.product.id)">取消收藏</button>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="empty-state">
        <div class="empty-icon">❤️</div>
        <p>暂无收藏商品</p>
        <button class="btn btn-primary" @click="$router.push('/mall')">去逛逛</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getFavoriteList, deleteFavorite } from '@/api/favorite'
import { ElMessage } from 'element-plus'


const router = useRouter()
const list = ref([])

const fetchList = async () => {
    try {
        const res = await getFavoriteList()
        list.value = res.list || []
    } catch (e) {
        console.error(e)
    }
}

const removeFavorite = async (pid) => {
    try {
        await deleteFavorite({ product_id: pid })
        fetchList() // 刷新
    } catch (e) {
        ElMessage.error('操作失败')
    }
}

const goDetail = (id) => {
    router.push(`/product/${id}`)
}

onMounted(fetchList)
</script>

<style scoped>
.favorite-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: var(--spacing-lg);
}

.product-card {
  transition: transform 0.2s;
  overflow: hidden;
}

.product-card:hover {
  transform: translateY(-4px);
}

.card-img {
  height: 200px;
  overflow: hidden;
  cursor: pointer;
}

.card-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-body {
  padding: var(--spacing-md);
}

.product-name {
  font-size: var(--font-size-md);
  margin-bottom: var(--spacing-sm);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
}

.price-row {
  margin-bottom: var(--spacing-sm);
}

.price {
  color: var(--danger-color);
  font-weight: 600;
  font-size: var(--font-size-lg);
  margin-right: var(--spacing-sm);
}

.original-price {
  color: var(--text-light);
  text-decoration: line-through;
  font-size: var(--font-size-sm);
}

.action-row {
    text-align: right;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-xxl) 0;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: var(--spacing-md);
}
</style>
