import request from '@/utils/request'

// 获取商品列表
export function getProductList(params) {
    return request({
        url: '/products',
        method: 'get',
        params
    })
}

// 获取商品详情
export function getProductDetail(id) {
    return request({
        url: `/product/${id}`,
        method: 'get'
    })
}

// 创建商品（管理员）
export function createProduct(data) {
    return request({
        url: '/product/create',
        method: 'post',
        data
    })
}

// 更新商品（管理员）
export function updateProduct(data) {
    return request({
        url: '/product/update',
        method: 'post',
        data
    })
}

// 删除商品（管理员）
export function deleteProduct(id) {
    return request({
        url: `/product/${id}`,
        method: 'delete'
    })
}

// 商品收藏/取消收藏
export function toggleFavorite(productId) {
    return request({
        url: '/product/favorite',
        method: 'post',
        data: { product_id: productId }
    })
}

// 获取收藏列表
export function getFavoriteList() {
    return request({
        url: '/product/favorites',
        method: 'get'
    })
}

// 获取销量排行
export function getProductRank() {
    return request({
        url: '/product/rank',
        method: 'get'
    })
}

// 获取分类列表
export function getCategories() {
    return request({
        url: '/categories',
        method: 'get'
    })
}
