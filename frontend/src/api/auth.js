import request from '@/utils/request'

export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/logout',
    method: 'post'
  })
}

export function sendCode(data) {
  return request({
    url: '/send_code',
    method: 'post',
    data
  })
}

export function loginByCode(data) {
  return request({
    url: '/login_code',
    method: 'post',
    data
  })
}

export function forgetPassword(data) {
  return request({
    url: '/forget_password',
    method: 'post',
    data
  })
}

export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function updateUserInfo(data) {
  return request({
    url: '/user/update',
    method: 'post',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/user/change_password',
    method: 'post',
    data
  })
}
