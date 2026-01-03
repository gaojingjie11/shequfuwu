import request from '@/utils/request'

export function getDashboardStats() {
    return request({
        url: '/dashboard/stats',
        method: 'get'
    })
}

// 用户管理
export const getUserList = (params) => {
    return request({
        url: '/admin/user/list',
        method: 'get',
        params
    })
}

export const freezeUser = (data) => {
    return request({
        url: '/admin/user/freeze',
        method: 'post',
        data
    })
}

// 角色管理
export const getRoleList = () => {
    return request({
        url: '/admin/role/list',
        method: 'get'
    })
}

export const createRole = (data) => {
    return request({
        url: '/admin/role/create',
        method: 'post',
        data
    })
}

// 菜单管理
export const getMenuList = () => {
    return request({
        url: '/admin/menu/list',
        method: 'get'
    })
}

// 访客审核
export const getAdminVisitorList = (params) => {
    return request({
        url: '/visitor/admin/list',
        method: 'get',
        params
    })
}

export const auditVisitor = (data) => {
    return request({
        url: '/visitor/audit',
        method: 'post',
        data
    })
}

// 报修管理
export const getAdminRepairList = (params) => {
    return request({
        url: '/repair/admin/list',
        method: 'get',
        params
    })
}

export const processRepair = (data) => {
    return request({
        url: '/repair/process',
        method: 'post',
        data
    })
}

// 商品管理 (管理员)
export const createProduct = (data) => {
    return request({
        url: '/product/create',
        method: 'post',
        data
    })
}

export const updateProduct = (data) => {
    return request({
        url: '/product/update',
        method: 'post',
        data
    })
}

export const deleteProduct = (id) => {
    return request({
        url: `/product/${id}`,
        method: 'delete'
    })
}

// 订单管理 (管理员)
// 订单管理 (管理员)
export const getAdminOrderList = (params) => {
    return request({
        url: '/order/admin/list',
        method: 'get',
        params
    })
}

export const shipOrder = (data) => {
    return request({
        url: '/order/ship',
        method: 'post',
        data
    })
}

// 门店管理
export const createStore = (data) => {
    return request({
        url: '/store/create',
        method: 'post',
        data
    })
}

export const updateStore = (data) => {
    return request({
        url: '/store/update',
        method: 'post',
        data
    })
}

export const deleteStore = (id) => {
    return request({
        url: `/store/${id}`,
        method: 'delete'
    })
}

// 公告管理
export const createNotice = (data) => {
    return request({
        url: '/notice/create',
        method: 'post',
        data
    })
}

export const deleteNotice = (id) => {
    return request({
        url: `/notice/${id}`,
        method: 'delete'
    })
}

// User Management Actions
export const assignRole = (data) => {
    return request({
        url: '/admin/user/assign_role',
        method: 'post',
        data
    })
}

export const updateUserBalance = (data) => {
    return request({
        url: '/admin/user/update_balance',
        method: 'post',
        data
    })
}

// Parking Management
export const getAdminParkingList = (params) => {
    return request({
        url: '/parking/admin/list',
        method: 'get',
        params
    })
}

export const getParkingStats = () => {
    return request({
        url: '/parking/admin/stats',
        method: 'get'
    })
}

export const assignParking = (data) => {
    return request({
        url: '/parking/admin/assign',
        method: 'post',
        data
    })
}

// Property Fee Management
export const createPropertyFee = (data) => {
    return request({
        url: '/property/admin/create',
        method: 'post',
        data
    })
}

export const getAdminPropertyFeeList = (params) => {
    return request({
        url: '/property/admin/list',
        method: 'get',
        params
    })
}
