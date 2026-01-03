import request from '@/utils/request'

// 充值
export function recharge(amount) {
    return request({
        url: '/finance/recharge',
        method: 'post',
        data: { amount }
    })
}

// 转账
export function transfer(data) {
    return request({
        url: '/finance/transfer',
        method: 'post',
        data
    })
}

// 获取交易流水
export function getTransactionList(params) {
    return request({
        url: '/finance/transactions',
        method: 'get',
        params
    })
}

// 获取物业费列表 (User)
export function getPropertyFeeList() {
    return request({
        url: '/property/list', // private route in router.go
        method: 'get'
    })
}
