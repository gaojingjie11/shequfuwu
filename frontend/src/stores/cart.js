import { defineStore } from 'pinia'
import { getCartList } from '@/api/order'

export const useCartStore = defineStore('cart', {
  state: () => ({
    cartItems: [],
    cartCount: 0
  }),

  actions: {
    async fetchCart() {
      try {
        const list = await getCartList()
        this.cartItems = list || []
        this.cartCount = this.cartItems.length
      } catch (error) {
        console.error('fetchCart failed', error)
      }
    },

    clearCart() {
      this.cartItems = []
      this.cartCount = 0
    }
  }
})
