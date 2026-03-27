import request from '@/utils/request'

export function addFavorite(data) {
    return request({
        url: '/favorite/add',
        method: 'post',
        data
    })
}

export function deleteFavorite(data) {
    return request({
        url: '/favorite/delete',
        method: 'post',
        data
    })
}

export function getFavoriteList() {
    return request({
        url: '/favorites',
        method: 'get'
    })
}

export function checkFavorite(productId) {
    return request({
        url: '/favorite/check',
        method: 'get',
        params: { product_id: productId }
    })
}
