import request from '@/utils/request'

export function sendChat(data) {
    return request({
        url: '/chat/send',
        method: 'post',
        data,
        timeout: 60000 // AI 响应可能较慢，设置 60s 超时
    })
}
