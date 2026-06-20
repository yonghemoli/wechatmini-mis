export type OrderStatus =
  | 'pending_service'
  | 'pending_confirm'
  | 'completed'
  | 'exception'
  | 'refunded'

export interface OrderRecord {
  id: string
  customer: string
  phone: string
  service: string
  amount: number
  status: OrderStatus
  source: '扫一扫' | '搜索' | '分享'
  appointmentAt: string
  createdAt: string
  staff: string
  internalNote: string
}

export interface UserRecord {
  id: string
  avatar: string
  nickname: string
  registeredAt: string
  totalSpent: number
  lastOrderAt: string
  status: 'active' | 'banned'
}

export interface ContentRecord {
  id: string
  title: string
  image: string
  price: number
  description: string
  visible: boolean
  updatedAt: string
}

export const orders: OrderRecord[] = [
  {
    id: 'HS20260620001',
    customer: '林女士',
    phone: '138****3201',
    service: '深度保洁 4 小时',
    amount: 328,
    status: 'pending_service',
    source: '分享',
    appointmentAt: '2026-06-20 14:00',
    createdAt: '2026-06-20 09:18',
    staff: '王阿姨',
    internalNote: '客户强调厨房油污重，需带强力清洁剂'
  },
  {
    id: 'HS20260620002',
    customer: '周先生',
    phone: '136****7781',
    service: '空调清洗 2 台',
    amount: 236,
    status: 'pending_confirm',
    source: '搜索',
    appointmentAt: '2026-06-20 10:30',
    createdAt: '2026-06-19 21:40',
    staff: '陈师傅',
    internalNote: '师傅已完成，等待用户核销'
  },
  {
    id: 'HS20260619019',
    customer: '赵女士',
    phone: '139****8820',
    service: '日常保洁 3 小时',
    amount: 198,
    status: 'exception',
    source: '扫一扫',
    appointmentAt: '2026-06-19 16:00',
    createdAt: '2026-06-19 08:12',
    staff: '待改派',
    internalNote: '服务人员临时请假，需客服回访改期'
  },
  {
    id: 'HS20260619012',
    customer: '何先生',
    phone: '135****4910',
    service: '玻璃清洁',
    amount: 168,
    status: 'completed',
    source: '分享',
    appointmentAt: '2026-06-19 11:00',
    createdAt: '2026-06-18 19:30',
    staff: '刘阿姨',
    internalNote: '已评价五星'
  },
  {
    id: 'HS20260618008',
    customer: '吴女士',
    phone: '137****6632',
    service: '新房开荒',
    amount: 688,
    status: 'refunded',
    source: '搜索',
    appointmentAt: '2026-06-18 09:00',
    createdAt: '2026-06-17 22:11',
    staff: '未分配',
    internalNote: '客户临时取消，已原路退款'
  }
]

export const users: UserRecord[] = [
  {
    id: 'U10081',
    avatar: 'https://api.dicebear.com/9.x/initials/svg?seed=Lin',
    nickname: '林女士',
    registeredAt: '2026-05-18',
    totalSpent: 1264,
    lastOrderAt: '2026-06-20 09:18',
    status: 'active'
  },
  {
    id: 'U10074',
    avatar: 'https://api.dicebear.com/9.x/initials/svg?seed=Zhou',
    nickname: '周先生',
    registeredAt: '2026-05-07',
    totalSpent: 756,
    lastOrderAt: '2026-06-19 21:40',
    status: 'active'
  },
  {
    id: 'U10032',
    avatar: 'https://api.dicebear.com/9.x/initials/svg?seed=Zhao',
    nickname: '赵女士',
    registeredAt: '2026-04-11',
    totalSpent: 198,
    lastOrderAt: '2026-06-19 08:12',
    status: 'active'
  },
  {
    id: 'U09951',
    avatar: 'https://api.dicebear.com/9.x/initials/svg?seed=Wu',
    nickname: '异常退款用户',
    registeredAt: '2026-03-28',
    totalSpent: 688,
    lastOrderAt: '2026-06-17 22:11',
    status: 'banned'
  }
]

export const contents: ContentRecord[] = [
  {
    id: 'SVC001',
    title: '日常保洁',
    image: '/me.png',
    price: 66,
    description: '适合日常维护，按小时预约，覆盖客厅、卧室、厨房和卫生间基础清洁。',
    visible: true,
    updatedAt: '2026-06-19 18:20'
  },
  {
    id: 'SVC002',
    title: '深度保洁',
    image: '/me.png',
    price: 82,
    description: '重点处理油污、水垢和卫生死角，适合长期未系统清洁的家庭。',
    visible: true,
    updatedAt: '2026-06-18 16:10'
  },
  {
    id: 'SVC003',
    title: '空调清洗',
    image: '/me.png',
    price: 118,
    description: '拆洗滤网、蒸发器除菌、外壳清洁，按台计费。',
    visible: false,
    updatedAt: '2026-06-17 10:45'
  }
]

export const statusText: Record<OrderStatus, string> = {
  pending_service: '待服务',
  pending_confirm: '待核销',
  completed: '已完成',
  exception: '异常待处理',
  refunded: '已退款'
}

export const statusColor: Record<OrderStatus, string> = {
  pending_service: 'processing',
  pending_confirm: 'warning',
  completed: 'success',
  exception: 'error',
  refunded: 'default'
}

export const revenueTrend7 = [
  { date: '06-14', revenue: 1680 },
  { date: '06-15', revenue: 2210 },
  { date: '06-16', revenue: 1890 },
  { date: '06-17', revenue: 2520 },
  { date: '06-18', revenue: 1978 },
  { date: '06-19', revenue: 3064 },
  { date: '06-20', revenue: 1618 }
]

export const revenueTrend30 = Array.from({ length: 30 }, (_, index) => ({
  date: `D-${29 - index}`,
  revenue: 1200 + ((index * 237) % 1800) + (index % 5) * 120
}))

export const sourceRows = [
  { source: '分享', users: 486, orders: 128, revenue: 24360 },
  { source: '搜索', users: 352, orders: 96, revenue: 18840 },
  { source: '扫一扫', users: 191, orders: 44, revenue: 7920 }
]

export const exportCsv = (
  filename: string,
  rows: Array<Record<string, string | number | boolean>>
) => {
  if (!rows.length) return
  const headers = Object.keys(rows[0])
  const escape = (value: string | number | boolean) =>
    `"${String(value).replace(/"/g, '""')}"`
  const csv = [
    headers.join(','),
    ...rows.map(row => headers.map(header => escape(row[header])).join(','))
  ].join('\n')
  const blob = new Blob([`\uFEFF${csv}`], {
    type: 'text/csv;charset=utf-8;'
  })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}
