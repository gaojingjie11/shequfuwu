import request from '@/utils/request'

// 用户注册
export function register(data) {
    return request({
        url: '/register',
        method: 'post',
        data
    })
}

// 用户登录
export function login(data) {
    return request({
        url: '/login',
        method: 'post',
        data
    })
}

// 退出登录
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

// 忘记密码
export function forgetPassword(data) {
    return request({
        url: '/forget_password',
        method: 'post',
        data
    })
}

// 获取用户信息
export function getUserInfo() {
    return request({
        url: '/user/info',
        method: 'get'
    })
}

// 更新用户信息
export function updateUserInfo(data) {
    return request({
        url: '/user/update',
        method: 'post',
        data
    })
}

// 修改密码
export function changePassword(data) {
    return request({
        url: '/user/change_password',
        method: 'post',
        data
    })
}
