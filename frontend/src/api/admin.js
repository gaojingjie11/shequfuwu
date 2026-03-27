import request from '@/utils/request'

export function getDashboardStats() {
  return request({
    url: '/dashboard/stats',
    method: 'get'
  })
}

export function getAIReport(params) {
  return request({
    url: '/admin/ai-report',
    method: 'get',
    params
  })
}

export function generateAIReport() {
  return request({
    url: '/admin/ai-report/generate',
    method: 'post',
    timeout: 120000
  })
}

export function getAIReportList(params) {
  return request({
    url: '/admin/ai-report/list',
    method: 'get',
    params
  })
}

export function getAIReportDetail(id) {
  return request({
    url: `/admin/ai-report/${id}`,
    method: 'get'
  })
}

export function getUserList(params) {
  return request({
    url: '/admin/user/list',
    method: 'get',
    params
  })
}

export function freezeUser(data) {
  return request({
    url: '/admin/user/freeze',
    method: 'post',
    data
  })
}

export function getRoleList() {
  return request({
    url: '/admin/role/list',
    method: 'get'
  })
}

export function createRole(data) {
  return request({
    url: '/admin/role/create',
    method: 'post',
    data
  })
}

export function getMenuList() {
  return request({
    url: '/admin/menu/list',
    method: 'get'
  })
}

export function getAdminVisitorList(params) {
  return request({
    url: '/visitor/admin/list',
    method: 'get',
    params
  })
}

export function auditVisitor(data) {
  return request({
    url: '/visitor/audit',
    method: 'post',
    data
  })
}

export function getAdminRepairList(params) {
  return request({
    url: '/repair/admin/list',
    method: 'get',
    params
  })
}

export function processRepair(data) {
  return request({
    url: '/repair/process',
    method: 'post',
    data
  })
}

export function createProduct(data) {
  return request({
    url: '/product/create',
    method: 'post',
    data
  })
}

export function updateProduct(data) {
  return request({
    url: '/product/update',
    method: 'post',
    data
  })
}

export function deleteProduct(id) {
  return request({
    url: `/product/${id}`,
    method: 'delete'
  })
}

export function getAdminOrderList(params) {
  return request({
    url: '/order/admin/list',
    method: 'get',
    params
  })
}

export function shipOrder(data) {
  return request({
    url: '/order/ship',
    method: 'post',
    data
  })
}

export function createStore(data) {
  return request({
    url: '/store/create',
    method: 'post',
    data
  })
}

export function updateStore(data) {
  return request({
    url: '/store/update',
    method: 'post',
    data
  })
}

export function deleteStore(id) {
  return request({
    url: `/store/${id}`,
    method: 'delete'
  })
}

export function createNotice(data) {
  return request({
    url: '/notice/create',
    method: 'post',
    data
  })
}

export function deleteNotice(id) {
  return request({
    url: `/notice/${id}`,
    method: 'delete'
  })
}

export function assignRole(data) {
  return request({
    url: '/admin/user/assign_role',
    method: 'post',
    data
  })
}

export function updateUserBalance(data) {
  return request({
    url: '/admin/user/update_balance',
    method: 'post',
    data
  })
}

export function getAdminParkingList(params) {
  return request({
    url: '/parking/admin/list',
    method: 'get',
    params
  })
}

export function getParkingStats() {
  return request({
    url: '/parking/admin/stats',
    method: 'get'
  })
}

export function assignParking(data) {
  return request({
    url: '/parking/admin/assign',
    method: 'post',
    data
  })
}

export function createParking(data) {
  return request({
    url: '/parking/admin/create',
    method: 'post',
    data
  })
}

export function createPropertyFee(data) {
  return request({
    url: '/property/admin/create',
    method: 'post',
    data
  })
}

export function getAdminPropertyFeeList(params) {
  return request({
    url: '/property/admin/list',
    method: 'get',
    params
  })
}
