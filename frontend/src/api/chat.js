import request from '@/utils/request'

export function sendChat(data) {
  return request({
    url: '/chat/send',
    method: 'post',
    data,
    timeout: 60000
  })
}

export function getChatHistory(params) {
  return request({
    url: '/chat/history',
    method: 'get',
    params
  })
}
