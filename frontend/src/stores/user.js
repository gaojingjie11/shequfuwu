import { defineStore } from 'pinia'
import { login as apiLogin, loginByCode as apiLoginByCode, register as apiRegister, logout as apiLogout, getUserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}'),
        isLoggedIn: !!localStorage.getItem('token')
    }),

    actions: {
        // 登录
        async login(data) {
            try {
                const res = await apiLogin(data)
                this.token = res.token
                this.userInfo = res.user_info
                this.isLoggedIn = true

                localStorage.setItem('token', res.token)
                localStorage.setItem('userInfo', JSON.stringify(res.user_info))

                return res
            } catch (error) {
                throw error
            }
        },

        // 验证码登录
        async loginByCode(data) {
            try {
                const res = await apiLoginByCode(data)
                this.token = res.token
                this.userInfo = res.user_info
                this.isLoggedIn = true

                localStorage.setItem('token', res.token)
                localStorage.setItem('userInfo', JSON.stringify(res.user_info))

                return res
            } catch (error) {
                throw error
            }
        },

        // 注册
        async register(data) {
            return await apiRegister(data)
        },

        // 获取用户信息
        async fetchUserInfo() {
            try {
                const userInfo = await getUserInfo()
                this.userInfo = userInfo
                localStorage.setItem('userInfo', JSON.stringify(userInfo))
            } catch (error) {
                console.error('获取用户信息失败:', error)
            }
        },

        // 退出登录
        async logout() {
            try {
                // 尝试调用后端登出接口，清理 Redis Token
                await apiLogout()
            } catch (error) {
                console.warn('后端登出失败(可能是Token已过期):', error)
            } finally {
                // 无论后端是否成功，前端都要清理掉 local storage
                this.token = ''
                this.userInfo = {}
                this.isLoggedIn = false

                localStorage.removeItem('token')
                localStorage.removeItem('userInfo')
            }
        }
    }
})
