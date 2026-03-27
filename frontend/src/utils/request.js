import axios from 'axios'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})

function createBusinessError(res) {
  const error = new Error(res?.msg || '请求失败')
  error.response = {
    status: 200,
    data: res
  }
  error.code = String(res?.code || '')
  return error
}

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

request.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code === 200) {
      return res.data
    }

    if (res.code === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      window.location.href = '/login'
      return Promise.reject(createBusinessError(res))
    }

    return Promise.reject(createBusinessError(res))
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default request
