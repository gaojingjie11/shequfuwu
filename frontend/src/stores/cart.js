import { defineStore } from 'pinia'
import { getCartList } from '@/api/order'

export const useCartStore = defineStore('cart', {
    state: () => ({
        cartItems: [],
        cartCount: 0
    }),

    actions: {
        // 获取购物车列表
        async fetchCart() {
            try {
                const list = await getCartList()
                this.cartItems = list
                this.cartCount = list.length
            } catch (error) {
                console.error('获取购物车失败:', error)
            }
        },

        // 清空购物车
        clearCart() {
            this.cartItems = []
            this.cartCount = 0
        }
    }
})
