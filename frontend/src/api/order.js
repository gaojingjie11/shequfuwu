import request from '@/utils/request'

export function addToCart(data) {
  return request({
    url: '/cart/add',
    method: 'post',
    data
  })
}

export function getCartList() {
  return request({
    url: '/cart/list',
    method: 'get'
  })
}

export function deleteCartItem(id) {
  return request({
    url: `/cart/${id}`,
    method: 'delete'
  })
}

export function updateCartQuantity(id, quantity) {
  return request({
    url: `/cart/${id}`,
    method: 'post',
    data: { quantity }
  })
}

export function createOrder(data) {
  return request({
    url: '/order/create',
    method: 'post',
    data
  })
}

export function getOrderList(params) {
  return request({
    url: '/order/list',
    method: 'get',
    params
  })
}

export function payOrder(data) {
  return request({
    url: '/finance/pay',
    method: 'post',
    data: {
      business_id: data.order_id,
      pay_type: 1,
      password: data.password || ''
    }
  })
}

export function cancelOrder(orderId) {
  return request({
    url: '/order/cancel',
    method: 'post',
    data: { id: orderId }
  })
}

export function receiveOrder(id) {
  return request({
    url: '/order/receive',
    method: 'post',
    data: { id }
  })
}

export function getOrderDetail(id) {
  return request({
    url: '/order/detail',
    method: 'get',
    params: { id }
  })
}
