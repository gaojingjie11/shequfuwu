
import request from '@/utils/request'

export function getCommentList(params) {
    return request({
        url: '/comments',
        method: 'get',
        params
    })
}

export function createComment(data) {
    return request({
        url: '/comment/create',
        method: 'post',
        data
    })
}
