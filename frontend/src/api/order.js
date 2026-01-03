import request from '@/utils/request'

// 添加到购物车
export function addToCart(data) {
    return request({
        url: '/cart/add',
        method: 'post',
        data
    })
}

// 获取购物车列表
export function getCartList() {
    return request({
        url: '/cart/list',
        method: 'get'
    })
}

// 删除购物车项
export function deleteCartItem(id) {
    return request({
        url: `/cart/${id}`,
        method: 'delete'
    })
}

// 修改购物车数量
export function updateCartQuantity(id, quantity) {
    return request({
        url: `/cart/${id}`,
        method: 'post', // 或者是 put, 取决于后端 handler.Update 绑定的是哪个 method。 CartHandler.go 里面对应 Update 方法。
        // 等下，CartHandler.go 里面 Update 方法虽然定义了，但是 router 里面有绑定吗？
        // 检查 internal/router/router.go
        data: { quantity }
    })
}

// 创建订单
export function createOrder(data) {
    return request({
        url: '/order/create',
        method: 'post',
        data
    })
}

// 获取订单列表
export function getOrderList(params) {
    return request({
        url: '/order/list',
        method: 'get',
        params
    })
}

// 支付订单
export const payOrder = (data) => {
    return request({
        url: '/finance/pay',
        method: 'post',
        data: {
            business_id: data.order_id,
            pay_type: 1  // 1表示订单支付
        }
    })
}

// 取消订单
export const cancelOrder = (orderId) => {
    return request({
        url: '/order/cancel',
        method: 'post',
        data: { id: orderId }
    })
}

// 确认收货
export function receiveOrder(id) {
    return request({
        url: '/order/receive',
        method: 'post',
        data: { id }
    })
}
