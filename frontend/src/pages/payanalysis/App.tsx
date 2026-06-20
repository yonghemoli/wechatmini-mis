import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Statistic,
  Table,
  type TableColumnsType,
  Tabs,
  Typography,
  message
} from 'antd'
import {
  DollarOutlined,
  UserOutlined,
  RiseOutlined,
  CrownOutlined,
  GiftOutlined
} from '@ant-design/icons'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
  Legend,
  PieChart,
  Pie,
  Cell
} from 'recharts'
import { apiGetPayAnalysis } from '@/api'

const { Title } = Typography

const COLORS = [
  '#1890ff',
  '#52c41a',
  '#fa8c16',
  '#ff4d4f',
  '#722ed1',
  '#13c2c2',
  '#eb2f96',
  '#fadb14'
]

interface FirstRechargeItem {
  day_bucket: string
  user_count: number
}

interface RevenueTrendItemRaw {
  date: string
  revenue: number
  order_count: number
  paying_users: number
}

interface RevenueTrendItem {
  date: string
  revenue: number
  order_count: number
  paying_users: number
}

type PackageType =
  | 'normal'
  | 'first_charge'
  | 'month_card'
  | 'seclusion_pass'
  | string

interface PackageStatItem {
  package_id: number
  package_name: string
  package_type: PackageType
  price: number
  sold_count: number
  total_revenue: number
  buyer_count: number
}

interface TierTrendItemRaw {
  date: string
  package_name: string
  order_count: number
  revenue: number
}

interface TierTrendItem {
  date: string
  package_name: string
  order_count: number
  revenue: number
}

interface PassStats {
  active_month_cards: number
  active_seclusion_passes: number
  total_month_card_users: number
  total_seclusion_users: number
}

interface VipDistItem {
  vip_level: number
  count: number
}

interface PayAnalysisData {
  first_recharge_conversion: FirstRechargeItem[]
  revenue_trend: RevenueTrendItemRaw[]
  package_stats: PackageStatItem[]
  tier_trend: TierTrendItemRaw[]
  pass_stats: PassStats | null
  vip_dist: VipDistItem[]
}

interface PieLabelProps {
  package_name?: string
  vip_level?: number
  percent?: number
}

const PayAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<PayAnalysisData | null>(null)

  useEffect(() => {
    const loadData = async () => {
      setLoading(true)
      try {
        const r = await apiGetPayAnalysis(days)
        if (r?.data) {
          setData(r.data as PayAnalysisData)
        } else {
          setData(null)
        }
      } catch {
        setData(null)
        message.error('加载付费分析数据失败')
      } finally {
        setLoading(false)
      }
    }

    void loadData()
  }, [days])

  const firstRecharge = data?.first_recharge_conversion || []
  const revenueTrend: RevenueTrendItem[] = (data?.revenue_trend || []).map(
    r => ({
      ...r,
      date: r.date?.slice(0, 10) || r.date,
      revenue: (r.revenue || 0) / 100
    })
  )
  const packageStats: PackageStatItem[] = data?.package_stats || []
  const tierTrend: TierTrendItem[] = (data?.tier_trend || []).map(r => ({
    ...r,
    date: r.date?.slice(0, 10) || r.date,
    revenue: (r.revenue || 0) / 100
  }))
  const passStats = data?.pass_stats
  const vipDist: VipDistItem[] = data?.vip_dist || []

  // 汇总
  const totalRevenue = revenueTrend.reduce(
    (s: number, r) => s + (r.revenue || 0),
    0
  )
  const totalOrders = revenueTrend.reduce(
    (s: number, r) => s + (r.order_count || 0),
    0
  )

  const packageColumns: TableColumnsType<PackageStatItem> = [
    { title: '礼包名称', dataIndex: 'package_name', key: 'package_name' },
    {
      title: '类型',
      dataIndex: 'package_type',
      key: 'package_type',
      render: (v: string) => {
        if (v === 'month_card') return '月卡'
        if (v === 'seclusion_pass') return '闭关卡'
        return v
      }
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      sorter: (a, b) => a.price - b.price,
      render: (v: number) => `¥${(v / 100)?.toLocaleString()}`
    },
    {
      title: '销售次数',
      dataIndex: 'sold_count',
      key: 'sold_count',
      sorter: (a, b) => a.sold_count - b.sold_count
    },
    {
      title: '购买人数',
      dataIndex: 'buyer_count',
      key: 'buyer_count',
      sorter: (a, b) => a.buyer_count - b.buyer_count
    },
    {
      title: '总收入',
      dataIndex: 'total_revenue',
      key: 'total_revenue',
      sorter: (a, b) => a.total_revenue - b.total_revenue,
      render: (v: number) => (
        <span style={{ color: '#52c41a', fontWeight: 'bold' }}>
          ¥{(v / 100)?.toLocaleString()}
        </span>
      )
    }
  ]

  const tierColumns: TableColumnsType<TierTrendItem> = [
    { title: '日期', dataIndex: 'date', key: 'date' },
    { title: '档位', dataIndex: 'package_name', key: 'package_name' },
    {
      title: '订单数',
      dataIndex: 'order_count',
      key: 'order_count',
      sorter: (a, b) => a.order_count - b.order_count
    },
    {
      title: '收入',
      dataIndex: 'revenue',
      key: 'revenue',
      sorter: (a, b) => a.revenue - b.revenue,
      render: (v: number) => (
        <span style={{ color: '#52c41a', fontWeight: 'bold' }}>
          ¥{v?.toLocaleString()}
        </span>
      )
    }
  ]

  const tabItems = [
    {
      key: 'overview',
      label: (
        <span>
          <DollarOutlined /> 收入趋势
        </span>
      ),
      children: (
        <>
          <Row gutter={16} style={{ marginBottom: 16 }}>
            <Col span={6}>
              <Card>
                <Statistic
                  title="总收入"
                  value={totalRevenue}
                  prefix="¥"
                  valueStyle={{ color: '#52c41a' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="总订单"
                  value={totalOrders}
                  prefix={<DollarOutlined />}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="礼包种类"
                  value={packageStats.length}
                  prefix={<GiftOutlined />}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="月卡/闭关卡活跃"
                  value={
                    (passStats?.active_month_cards || 0) +
                    (passStats?.active_seclusion_passes || 0)
                  }
                  prefix={<CrownOutlined />}
                  valueStyle={{ color: '#722ed1' }}
                />
              </Card>
            </Col>
          </Row>
          <Card title="每日收入趋势">
            <ResponsiveContainer width="100%" height={350}>
              <LineChart data={revenueTrend}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="revenue"
                  name="收入"
                  stroke="#52c41a"
                  strokeWidth={2}
                />
                <Line
                  type="monotone"
                  dataKey="order_count"
                  name="订单数"
                  stroke="#1890ff"
                  strokeWidth={2}
                />
                <Line
                  type="monotone"
                  dataKey="paying_users"
                  name="付费人数"
                  stroke="#ff4d4f"
                  strokeWidth={2}
                />
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </>
      )
    },
    {
      key: 'first',
      label: (
        <span>
          <UserOutlined /> 首充转化
        </span>
      ),
      children: (
        <Card title="注册→首充天数分布">
          <ResponsiveContainer width="100%" height={350}>
            <BarChart data={firstRecharge}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="day_bucket" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="user_count" name="用户数" fill="#1890ff" />
            </BarChart>
          </ResponsiveContainer>
        </Card>
      )
    },
    {
      key: 'packages',
      label: (
        <span>
          <GiftOutlined /> 礼包分析
        </span>
      ),
      children: (
        <Row gutter={16}>
          <Col span={14}>
            <Card title="礼包销售排行">
              <Table
                columns={packageColumns}
                dataSource={packageStats}
                rowKey="package_name"
                size="small"
                pagination={{ pageSize: 10 }}
              />
            </Card>
          </Col>
          <Col span={10}>
            <Card title="礼包收入占比">
              <ResponsiveContainer width="100%" height={350}>
                <PieChart>
                  <Pie
                    data={packageStats.slice(0, 8)}
                    dataKey="total_revenue"
                    nameKey="package_name"
                    cx="50%"
                    cy="50%"
                    outerRadius={120}
                    label={({ package_name, percent }: PieLabelProps) =>
                      `${package_name} ${((percent || 0) * 100).toFixed(0)}%`
                    }
                  >
                    {packageStats.slice(0, 8).map((_, i: number) => (
                      <Cell key={i} fill={COLORS[i % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>
      )
    },
    {
      key: 'tier',
      label: (
        <span>
          <RiseOutlined /> 付费等级趋势
        </span>
      ),
      children: (
        <>
          <Row gutter={16} style={{ marginBottom: 16 }}>
            <Col span={6}>
              <Card>
                <Statistic
                  title="月卡活跃"
                  value={passStats?.active_month_cards || 0}
                  valueStyle={{ color: '#722ed1' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="闭关卡活跃"
                  value={passStats?.active_seclusion_passes || 0}
                  valueStyle={{ color: '#13c2c2' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="月卡累计用户"
                  value={passStats?.total_month_card_users || 0}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="闭关卡累计用户"
                  value={passStats?.total_seclusion_users || 0}
                />
              </Card>
            </Col>
          </Row>
          <Card title="各档位日充值趋势" style={{ marginBottom: 16 }}>
            <Table
              columns={tierColumns}
              dataSource={tierTrend}
              rowKey={r => `${r.date}-${r.package_name}`}
              size="small"
              pagination={{ pageSize: 15 }}
            />
          </Card>
          <Card title="VIP等级分布">
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={vipDist}
                  dataKey="count"
                  nameKey="vip_level"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label={({ vip_level, percent }: PieLabelProps) =>
                    `VIP${vip_level} ${((percent || 0) * 100).toFixed(0)}%`
                  }
                >
                  {vipDist.map((_, i: number) => (
                    <Cell key={i} fill={COLORS[i % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </>
      )
    }
  ]

  return (
    <div style={{ padding: 24 }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 16
        }}
      >
        <Title level={3} style={{ margin: 0 }}>
          付费行为深度分析
        </Title>
        <Select value={days} onChange={setDays} style={{ width: 120 }}>
          <Select.Option value={7}>近7天</Select.Option>
          <Select.Option value={14}>近14天</Select.Option>
          <Select.Option value={30}>近30天</Select.Option>
          <Select.Option value={90}>近90天</Select.Option>
        </Select>
      </div>

      <Spin spinning={loading}>
        <Tabs items={tabItems} />
      </Spin>
    </div>
  )
}

export default PayAnalysis
