import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Statistic,
  Table,
  Tag,
  Tooltip,
  Typography,
  Empty,
  Progress
} from 'antd'
import {
  DollarOutlined,
  ShoppingOutlined,
  CrownOutlined,
  RiseOutlined
} from '@ant-design/icons'
import {
  apiGetPaymentOverview,
  apiGetRevenueTrend,
  apiGetPackageStats
} from '@/api'

const { Title, Text } = Typography

type RevenueSource = 'all' | 'mis' | 'official'

interface PaymentOverview {
  total_revenue: number
  total_orders: number
  paying_users: number
  total_users: number
  pay_rate: number
  arpu: number
  arppu: number
  revenue_7d: number
  revenue_30d: number
  vip_distribution: Array<{
    vip_level: number
    count: number
  }>
}

interface RevenueItem {
  date: string
  revenue: number
  order_count: number
  paying_users: number
}

interface PackageItem {
  package_id: number
  package_name: string
  package_type: string
  price: number
  sold_count: number
  buyer_count: number
  total_revenue: number
}

const packageTypeLabels: Record<string, { text: string; color: string }> = {
  normal: { text: '普通充值', color: 'blue' },
  first_charge: { text: '首充', color: 'gold' },
  month_card: { text: '月卡', color: 'purple' },
  seclusion_pass: { text: '特殊卡', color: 'cyan' }
}

const revenueSourceOptions: Array<{ label: string; value: RevenueSource }> = [
  { label: '全部', value: 'all' },
  { label: '仅 MIS', value: 'mis' },
  { label: '官网（非 MIS）', value: 'official' }
]

const Economy: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [source, setSource] = useState<RevenueSource>('official')
  const [payment, setPayment] = useState<PaymentOverview | null>(null)
  const [revDays, setRevDays] = useState(30)
  const [revenue, setRevenue] = useState<RevenueItem[]>([])
  const [pkgDays, setPkgDays] = useState(0)
  const [packages, setPackages] = useState<PackageItem[]>([])

  useEffect(() => {
    const loadAll = async () => {
      setLoading(true)
      try {
        await Promise.all([
          loadPayment(source),
          loadRevenue(revDays, source),
          loadPackages(pkgDays, source)
        ])
      } finally {
        setLoading(false)
      }
    }

    void loadAll()
  }, [pkgDays, revDays, source])

  const loadPayment = async (currentSource: RevenueSource) => {
    try {
      const r = await apiGetPaymentOverview(currentSource)
      if (r?.data) setPayment(r.data)
      else setPayment(null)
    } catch {
      setPayment(null)
    }
  }

  const loadRevenue = async (days: number, currentSource: RevenueSource) => {
    try {
      const r = await apiGetRevenueTrend(days, currentSource)
      if (r?.data) setRevenue(r.data || [])
      else setRevenue([])
    } catch {
      setRevenue([])
    }
  }

  const loadPackages = async (days: number, currentSource: RevenueSource) => {
    try {
      const r = await apiGetPackageStats(days, currentSource)
      if (r?.data) setPackages(r.data || [])
      else setPackages([])
    } catch {
      setPackages([])
    }
  }

  const maxRev = Math.max(...revenue.map(d => d.revenue || 0), 1)

  return (
    <Spin spinning={loading}>
      <div style={{ padding: 24 }}>
        <Title level={4}>
          <DollarOutlined style={{ marginRight: 8 }} />
          经济与付费分析
        </Title>
        <div style={{ marginBottom: 16 }}>
          <Text type="secondary" style={{ marginRight: 8 }}>
            订单来源
          </Text>
          <Select
            value={source}
            onChange={setSource}
            style={{ width: 180 }}
            options={revenueSourceOptions}
          />
        </div>

        {/* 付费核心指标 */}
        <Row gutter={[12, 12]}>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="总充值(元)"
                value={payment ? payment.total_revenue / 100 : 0}
                precision={2}
                prefix={<DollarOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="付费用户"
                value={payment?.paying_users || 0}
                prefix={<CrownOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="付费率"
                value={payment?.pay_rate || 0}
                suffix="%"
                precision={1}
                valueStyle={{
                  color: (payment?.pay_rate || 0) > 10 ? '#52c41a' : '#faad14'
                }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="ARPU(元)"
                value={payment?.arpu || 0}
                precision={2}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="ARPPU(元)"
                value={payment?.arppu || 0}
                precision={2}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="近7日(元)"
                value={payment ? payment.revenue_7d / 100 : 0}
                precision={2}
                prefix={<RiseOutlined />}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
        </Row>

        <Row gutter={16} style={{ marginTop: 16 }}>
          {/* 营收趋势 */}
          <Col xs={24} lg={14}>
            <Card
              title={
                <span>
                  <RiseOutlined style={{ color: '#52c41a', marginRight: 6 }} />
                  营收趋势
                </span>
              }
              size="small"
              extra={
                <Select
                  value={revDays}
                  onChange={setRevDays}
                  size="small"
                  style={{ width: 90 }}
                  options={[
                    { label: '7天', value: 7 },
                    { label: '14天', value: 14 },
                    { label: '30天', value: 30 },
                    { label: '90天', value: 90 }
                  ]}
                />
              }
            >
              {revenue.length > 0 ? (
                <>
                  <div
                    style={{
                      display: 'flex',
                      alignItems: 'flex-end',
                      height: 160,
                      gap: 2,
                      padding: '0 4px',
                      borderBottom: '1px solid #f0f0f0'
                    }}
                  >
                    {revenue.map(d => (
                      <Tooltip
                        key={d.date}
                        title={
                          <div>
                            <div>{d.date}</div>
                            <div>营收 ¥{(d.revenue / 100).toFixed(2)}</div>
                            <div>订单 {d.order_count} 笔</div>
                            <div>付费 {d.paying_users} 人</div>
                          </div>
                        }
                      >
                        <div
                          style={{
                            flex: 1,
                            height: `${((d.revenue || 0) / maxRev) * 100}%`,
                            background:
                              'linear-gradient(180deg, #722ed1 0%, #b37feb 100%)',
                            borderRadius: '3px 3px 0 0',
                            minHeight: 2,
                            transition: 'height 0.3s',
                            cursor: 'pointer'
                          }}
                        />
                      </Tooltip>
                    ))}
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      fontSize: 11,
                      color: '#999',
                      marginTop: 4
                    }}
                  >
                    <span>{revenue[0]?.date}</span>
                    <span>
                      总计 ¥
                      {(
                        revenue.reduce(
                          (s: number, d) => s + (d.revenue || 0),
                          0
                        ) / 100
                      ).toFixed(2)}
                    </span>
                    <span>{revenue[revenue.length - 1]?.date}</span>
                  </div>
                </>
              ) : (
                <Empty description="暂无营收数据" />
              )}
            </Card>
          </Col>

          {/* VIP 分布 */}
          <Col xs={24} lg={10}>
            <Card
              title={
                <span>
                  <CrownOutlined style={{ color: '#faad14', marginRight: 6 }} />
                  VIP 等级分布
                </span>
              }
              size="small"
            >
              {payment?.vip_distribution?.length ? (
                (() => {
                  const vipDistribution = payment.vip_distribution
                  const totalVip = vipDistribution.reduce(
                    (s: number, v) => s + v.count,
                    0
                  )
                  return vipDistribution.map(v => {
                    const pct = totalVip > 0 ? (v.count / totalVip) * 100 : 0
                    const isVip = v.vip_level > 0
                    return (
                      <div
                        key={v.vip_level}
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          marginBottom: 6
                        }}
                      >
                        <Tag
                          color={isVip ? 'gold' : 'default'}
                          style={{ width: 50, textAlign: 'center', margin: 0 }}
                        >
                          {isVip ? `V${v.vip_level}` : '非VIP'}
                        </Tag>
                        <div style={{ flex: 1, margin: '0 8px' }}>
                          <div
                            style={{
                              width: `${pct}%`,
                              height: 18,
                              minWidth: 2,
                              background: isVip
                                ? `hsl(${40 + v.vip_level * 15}, 80%, 55%)`
                                : '#d9d9d9',
                              borderRadius: 3,
                              transition: 'width 0.3s'
                            }}
                          />
                        </div>
                        <span
                          style={{
                            width: 50,
                            textAlign: 'right',
                            fontWeight: 600
                          }}
                        >
                          {v.count}
                        </span>
                        <span
                          style={{
                            width: 55,
                            textAlign: 'right',
                            color: '#999',
                            fontSize: 12
                          }}
                        >
                          {pct.toFixed(1)}%
                        </span>
                      </div>
                    )
                  })
                })()
              ) : (
                <Empty description="暂无 VIP 数据" />
              )}
            </Card>
          </Col>
        </Row>

        {/* 充值套餐销售分析 */}
        <Card
          title={
            <span>
              <ShoppingOutlined style={{ color: '#1890ff', marginRight: 6 }} />
              充值套餐销售分析
            </span>
          }
          size="small"
          style={{ marginTop: 16 }}
          extra={
            <Select
              value={pkgDays}
              onChange={setPkgDays}
              size="small"
              style={{ width: 100 }}
              options={[
                { label: '全部', value: 0 },
                { label: '近7天', value: 7 },
                { label: '近30天', value: 30 },
                { label: '近90天', value: 90 }
              ]}
            />
          }
        >
          <Table
            dataSource={packages}
            rowKey="package_id"
            size="small"
            pagination={false}
            scroll={{ x: 700 }}
            columns={[
              {
                title: '套餐',
                dataIndex: 'package_name',
                key: 'package_name',
                ellipsis: true,
                render: (v: string) => <Text strong>{v}</Text>
              },
              {
                title: '类型',
                dataIndex: 'package_type',
                key: 'package_type',
                width: 90,
                render: (v: string) => {
                  const info = packageTypeLabels[v] || {
                    text: v,
                    color: 'default'
                  }
                  return <Tag color={info.color}>{info.text}</Tag>
                }
              },
              {
                title: '价格(元)',
                dataIndex: 'price',
                key: 'price',
                width: 85,
                render: (v: number) => `¥${(v / 100).toFixed(2)}`,
                sorter: (a: PackageItem, b: PackageItem) => a.price - b.price
              },
              {
                title: '销量',
                dataIndex: 'sold_count',
                key: 'sold_count',
                width: 70,
                sorter: (a: PackageItem, b: PackageItem) =>
                  a.sold_count - b.sold_count
              },
              {
                title: '购买人数',
                dataIndex: 'buyer_count',
                key: 'buyer_count',
                width: 80
              },
              {
                title: '总收入(元)',
                dataIndex: 'total_revenue',
                key: 'total_revenue',
                width: 100,
                render: (v: number) => (
                  <Text strong style={{ color: '#722ed1' }}>
                    ¥{(v / 100).toFixed(2)}
                  </Text>
                ),
                sorter: (a: PackageItem, b: PackageItem) =>
                  a.total_revenue - b.total_revenue,
                defaultSortOrder: 'descend' as const
              },
              {
                title: '收入占比',
                key: 'ratio',
                width: 120,
                render: (_: unknown, row: PackageItem) => {
                  const total = packages.reduce(
                    (s: number, p) => s + (p.total_revenue || 0),
                    0
                  )
                  const pct = total > 0 ? (row.total_revenue / total) * 100 : 0
                  return (
                    <Progress
                      percent={Math.round(pct)}
                      size="small"
                      style={{ width: 80, margin: 0 }}
                      strokeColor="#722ed1"
                      format={p => `${p}%`}
                    />
                  )
                }
              }
            ]}
          />
        </Card>
      </div>
    </Spin>
  )
}

export default Economy
