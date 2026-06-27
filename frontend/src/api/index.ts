import { axiosInstance, server } from './base'

const download = async (url: string, filename: string) => {
  const res = await axiosInstance({ url, method: 'GET', responseType: 'blob' })
  const blobUrl = URL.createObjectURL(res.data)
  const link = document.createElement('a')
  link.href = blobUrl
  link.download = filename
  link.click()
  URL.revokeObjectURL(blobUrl)
}

// ==================== 内部认证 ====================
export const apiLogin = (data: { username: string; password: string }) =>
  server({ url: '/login', method: 'POST', data })

export const apiSession = () =>
  server({ url: '/session', method: 'GET', silent: true })

export const apiLogout = () => server({ url: '/logout', method: 'POST' })

export const apiGetMe = () => server({ url: '/me', method: 'GET' })

// ==================== 内部账户 ====================
export const apiListAdminAccounts = () =>
  server({ url: '/admin/accounts', method: 'GET' })

export const apiCreateAdminAccount = (data: {
  username: string
  password: string
  name: string
  email: string
  roleId?: number | null
  isSuperAdmin?: boolean
}) => server({ url: '/admin/accounts', method: 'POST', data })

export const apiUpdateAdminAccount = (
  id: number,
  data: {
    name: string
    email: string
    roleId?: number | null
    isSuperAdmin?: boolean
    status?: string
  }
) => server({ url: `/admin/accounts/${id}`, method: 'PUT', data })

export const apiResetAdminPassword = (id: number, password: string) =>
  server({
    url: `/admin/accounts/${id}/reset-password`,
    method: 'POST',
    data: { password }
  })

export const apiDisableAdminAccount = (id: number) =>
  server({ url: `/admin/accounts/${id}/disable`, method: 'POST' })

export const apiEnableAdminAccount = (id: number) =>
  server({ url: `/admin/accounts/${id}/enable`, method: 'POST' })

export type OrderStatus =
  | 'pending_service'
  | 'pending_confirm'
  | 'completed'
  | 'exception'
  | 'refunded'

export const orderStatusText: Record<OrderStatus, string> = {
  pending_service: '待服务',
  pending_confirm: '待核销',
  completed: '已完成',
  exception: '异常待处理',
  refunded: '已退款'
}

export const orderStatusColor: Record<OrderStatus, string> = {
  pending_service: 'processing',
  pending_confirm: 'warning',
  completed: 'success',
  exception: 'error',
  refunded: 'default'
}

export type OrderRecord = {
  id: string
  customer: string
  phone: string
  service: string
  amount: number
  status: OrderStatus
  source: string
  appointmentAt: string
  createdAt: string
  staff: string
  internalNote: string
}

export type UserRecord = {
  id: string
  avatar: string
  nickname: string
  totalSpent: number
  lastOrderAt: string
  status: 'active' | 'banned'
  createdAt: string
}

export const apiDashboardSummary = () =>
  server({ url: '/dashboard/summary', method: 'GET' })

export const apiDashboardExceptions = () =>
  server({ url: '/dashboard/exceptions', method: 'GET' })

export const apiListOrders = (params: {
  status?: string
  keyword?: string
  start?: string
  end?: string
  page?: number
  size?: number
}) => {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, String(value))
  })
  return server({ url: `/orders?${query.toString()}`, method: 'GET' })
}

export const apiConfirmOrder = (id: string) =>
  server({ url: `/orders/${id}/confirm`, method: 'POST' })

export const apiRefundOrder = (id: string, reason = '') =>
  server({ url: `/orders/${id}/refund`, method: 'POST', data: { reason } })

export const apiUpdateOrderNote = (id: string, internalNote: string) =>
  server({ url: `/orders/${id}/note`, method: 'PUT', data: { internalNote } })

export const apiExportOrders = (params: {
  status?: string
  keyword?: string
  start?: string
  end?: string
}) => {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value) query.set(key, String(value))
  })
  return download(`/orders/export?${query.toString()}`, '家政订单列表.csv')
}

export const apiListUsers = () => server({ url: '/users', method: 'GET' })

export const apiBanUser = (id: string) =>
  server({ url: `/users/${id}/ban`, method: 'POST' })

export const apiUnbanUser = (id: string) =>
  server({ url: `/users/${id}/unban`, method: 'POST' })

export const apiExportUsers = () => download('/users/export', '家政用户列表.csv')

export type ServiceTypeRecord = {
  id: number
  name: string
  description: string
  sortOrder: number
  status: 'active' | 'disabled'
}

export type ServiceRecord = {
  id: number
  typeId: number
  typeName: string
  name: string
  image: string
  price: number
  unit: string
  description: string
  visible: boolean
  sortOrder: number
  updatedAt: string
}

export type ShopRecord = {
  id: number
  name: string
  contactName: string
  phone: string
  address: string
  businessHours: string
  status: 'open' | 'closed'
  remark: string
}

export type FAQRecord = {
  id: number
  question: string
  answer: string
  category: string
  sortOrder: number
  visible: boolean
}

export type ChatSessionRecord = {
  id: string
  userId: string
  userName: string
  userAvatar: string
  status: 'open' | 'closed'
  lastMessage: string
  unreadCount: number
  updatedAt: string
}

export type ChatMessageRecord = {
  id: number
  sessionId: string
  sender: 'user' | 'admin'
  msgType: string
  content: string
  isRead: boolean
  createdAt: string
}

export type ChatWsEvent =
  | {
      type: 'message'
      sessionId: string
      message: ChatMessageRecord
    }
  | {
      type: 'session'
      sessionId: string
      session: ChatSessionRecord
    }
  | {
      type: 'error'
      error: string
    }
  | {
      type: 'pong'
    }

export const getWsBaseUrl = () => {
  const { protocol, host } = window.location
  return `${protocol === 'https:' ? 'wss:' : 'ws:'}//${host}/api/v1`
}

export const apiListServiceTypes = (keyword = '') =>
  server({ url: `/service-types?keyword=${encodeURIComponent(keyword)}`, method: 'GET' })

export const apiCreateServiceType = (data: Omit<ServiceTypeRecord, 'id'>) =>
  server({ url: '/service-types', method: 'POST', data })

export const apiUpdateServiceType = (id: number, data: Omit<ServiceTypeRecord, 'id'>) =>
  server({ url: `/service-types/${id}`, method: 'PUT', data })

export const apiDeleteServiceType = (id: number) =>
  server({ url: `/service-types/${id}`, method: 'DELETE' })

export const apiEnableServiceType = (id: number) =>
  server({ url: `/service-types/${id}/enable`, method: 'POST' })

export const apiDisableServiceType = (id: number) =>
  server({ url: `/service-types/${id}/disable`, method: 'POST' })

export const apiListServices = (params: { typeId?: number; keyword?: string }) => {
  const query = new URLSearchParams()
  if (params.typeId) query.set('typeId', String(params.typeId))
  if (params.keyword) query.set('keyword', params.keyword)
  return server({ url: `/services?${query.toString()}`, method: 'GET' })
}

export const apiCreateService = (data: Omit<ServiceRecord, 'id' | 'typeName' | 'updatedAt'>) =>
  server({ url: '/services', method: 'POST', data })

export const apiUpdateService = (
  id: number,
  data: Omit<ServiceRecord, 'id' | 'typeName' | 'updatedAt'>
) => server({ url: `/services/${id}`, method: 'PUT', data })

export const apiDeleteService = (id: number) =>
  server({ url: `/services/${id}`, method: 'DELETE' })

export const apiPublishService = (id: number) =>
  server({ url: `/services/${id}/publish`, method: 'POST' })

export const apiUnpublishService = (id: number) =>
  server({ url: `/services/${id}/unpublish`, method: 'POST' })

export const apiExportServices = () => download('/services/export', '家政服务列表.csv')

export const apiListShops = () => server({ url: '/shops', method: 'GET' })

export const apiCreateShop = (data: Omit<ShopRecord, 'id'>) =>
  server({ url: '/shops', method: 'POST', data })

export const apiUpdateShop = (id: number, data: Omit<ShopRecord, 'id'>) =>
  server({ url: `/shops/${id}`, method: 'PUT', data })

export const apiDeleteShop = (id: number) =>
  server({ url: `/shops/${id}`, method: 'DELETE' })

export const apiOpenShop = (id: number) =>
  server({ url: `/shops/${id}/open`, method: 'POST' })

export const apiCloseShop = (id: number) =>
  server({ url: `/shops/${id}/close`, method: 'POST' })

export const apiListFAQs = (category = '') =>
  server({ url: `/faqs?category=${encodeURIComponent(category)}`, method: 'GET' })

export const apiCreateFAQ = (data: Omit<FAQRecord, 'id'>) =>
  server({ url: '/faqs', method: 'POST', data })

export const apiUpdateFAQ = (id: number, data: Omit<FAQRecord, 'id'>) =>
  server({ url: `/faqs/${id}`, method: 'PUT', data })

export const apiDeleteFAQ = (id: number) =>
  server({ url: `/faqs/${id}`, method: 'DELETE' })

export const apiPublishFAQ = (id: number) =>
  server({ url: `/faqs/${id}/publish`, method: 'POST' })

export const apiUnpublishFAQ = (id: number) =>
  server({ url: `/faqs/${id}/unpublish`, method: 'POST' })

export const apiListChatSessions = () =>
  server({ url: '/chat/sessions', method: 'GET' })

export const apiListChatMessages = (id: string) =>
  server({ url: `/chat/sessions/${id}/messages`, method: 'GET' })

export const apiSendChatMessage = (id: string, content: string) =>
  server({ url: `/chat/sessions/${id}/messages`, method: 'POST', data: { content, msgType: 'text' } })

export const apiCloseChatSession = (id: string) =>
  server({ url: `/chat/sessions/${id}/close`, method: 'POST' })

export const apiReadChatSession = (id: string) =>
  server({ url: `/chat/sessions/${id}/read`, method: 'POST' })
