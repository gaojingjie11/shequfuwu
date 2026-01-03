import axios from 'axios'

// 创建 axios 实例
const request = axios.create({
    baseURL: '/api/v1',
    timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        // 从 localStorage 获取 token
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => {
        const res = response.data

        // 后端返回格式: { code: 200, data: {...}, msg: "success" }
        // 后端返回格式: { code: 200, data: {...}, msg: "success" }
        if (res.code === 200) {
            return res.data
        } else if (res.code === 401) {
            // Token 过期或未登录 (业务状态码 401)
            localStorage.removeItem('token')
            localStorage.removeItem('userInfo')
            window.location.href = '/login'
            return Promise.reject(new Error(res.msg || '登录已失效'))
        } else {
            // 处理错误
            console.error(res.msg || '请求失败')
            return Promise.reject(new Error(res.msg || '请求失败'))
        }
    },
    error => {
        if (error.response?.status === 401) {
            // Token 过期或未登录
            localStorage.removeItem('token')
            localStorage.removeItem('userInfo')
            window.location.href = '/login'
        }
        console.error('请求错误:', error.response?.data?.msg || error.message)
        return Promise.reject(error)
    }
)

export default request
