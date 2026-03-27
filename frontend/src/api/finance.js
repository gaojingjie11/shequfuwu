import request from '@/utils/request'

export function recharge(amount) {
  return request({
    url: '/finance/recharge',
    method: 'post',
    data: { amount }
  })
}

export function transfer(data) {
  return request({
    url: '/finance/transfer',
    method: 'post',
    data
  })
}

export function getTransactionList(params) {
  return request({
    url: '/finance/transactions',
    method: 'get',
    params
  })
}

export function getPropertyFeeList(params) {
  return request({
    url: '/property/list',
    method: 'get',
    params
  })
}
