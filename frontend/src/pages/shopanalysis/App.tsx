import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Spin,
  Statistic,
  Table,
  Tag,
  Typography,
  message
} from 'antd'
import {
  ShopOutlined,
  DollarOutlined,
  AppstoreOutlined
} from '@ant-design/icons'
import {
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend
} from 'recharts'
import { apiGetShopAnalysis } from '@/api'

const { Title } = Typography
const COLORS = [
  '#1890ff',
  '#52c41a',
  '#faad14',
  '#ff4d4f',
  '#722ed1',
  '#13c2c2'
]

const SHOP_NAMES: Record<string, string> = {
  WANBAOLOU: '万宝楼',
  SHENQI: '神奇商店',
  REPUTATION: '声望商店',
  CHECKIN: '签到商店',
  TOWER_SCORE: '积分商店',
  DAOLV: '道侣商店'
}

const CURRENCY_NAMES: Record<string, string> = {
  SPIRIT_CRYSTAL: '灵晶',
  SPIRIT_STONE: '灵石',
  SPIRIT_MILK: '道晶',
  REPUTATION: '声望',
  CHECKIN_COIN: '签到币',
  TOWER_SCORE: '积分',
  TONGXIN: '同心值'
}

const ShopAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetShopAnalysis()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载商店数据失败')
    } finally {
      setLoading(false)
    }
  }

  const overview = data?.overview || []
  const priceRangeDist = data?.price_range_dist || []
  const currencyDist = data?.currency_dist || []
  const categoryDist = data?.category_dist || []

  const totalProducts = overview.reduce(
    (s: number, o: any) => s + (o.product_count || 0),
    0
  )
  const shopCodes = [
    ...new Set(overview.map((o: any) => o.shop_code))
  ] as string[]

  const currencyColumns = [
    {
      title: '商店',
      dataIndex: 'shop_code',
      key: 'shop_code',
      render: (v: string) => SHOP_NAMES[v] || v
    },
    {
      title: '货币',
      dataIndex: 'price_currency',
      key: 'price_currency',
      render: (v: string) => <Tag color="blue">{CURRENCY_NAMES[v] || v}</Tag>
    },
    {
      title: '商品数',
      dataIndex: 'product_count',
      key: 'product_count',
      sorter: (a: any, b: any) => a.product_count - b.product_count
    },
    {
      title: '均价',
      dataIndex: 'avg_price',
      key: 'avg_price',
      render: (v: number) => (v ?? 0).toFixed(0)
    }
  ]

  const priceRangeColumns = [
    {
      title: '商店',
      dataIndex: 'shop_code',
      key: 'shop_code',
      render: (v: string) => SHOP_NAMES[v] || v
    },
    { title: '价格区间', dataIndex: 'price_range', key: 'price_range' },
    {
      title: '商品数',
      dataIndex: 'product_count',
      key: 'product_count',
      sorter: (a: any, b: any) => a.product_count - b.product_count
    }
  ]

  return (
    <Spin spinning={loading}>
      <Title level={4}>商店消费分析</Title>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic
              title="商店数量"
              value={shopCodes.length}
              prefix={<ShopOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="总商品数"
              value={totalProducts}
              prefix={<AppstoreOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="货币类型数"
              value={
                [...new Set(currencyDist.map((c: any) => c.price_currency))]
                  .length
              }
              prefix={<DollarOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={12}>
          <Card title="各商店商品数">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart
                data={overview.map((o: any) => ({
                  ...o,
                  shop_name: SHOP_NAMES[o.shop_code] || o.shop_code
                }))}
              >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="shop_name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="product_count" name="商品数" fill="#1890ff" />
                <Bar dataKey="category_count" name="分类数" fill="#52c41a" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="商品分类分布">
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={categoryDist.reduce((acc: any[], item: any) => {
                    const existing = acc.find(
                      (a: any) => a.category === item.category
                    )
                    if (existing) existing.product_count += item.product_count
                    else acc.push({ ...item })
                    return acc
                  }, [])}
                  dataKey="product_count"
                  nameKey="category"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label={({ name, percent }: any) =>
                    `${name} ${((percent ?? 0) * 100).toFixed(0)}%`
                  }
                >
                  {categoryDist.map((_: any, i: number) => (
                    <Cell key={i} fill={COLORS[i % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={12}>
          <Card title="货币使用分布">
            <Table
              dataSource={currencyDist}
              columns={currencyColumns}
              rowKey={(r: any) => `${r.shop_code}-${r.price_currency}`}
              size="small"
              pagination={false}
            />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="价格区间分布">
            <Table
              dataSource={priceRangeDist}
              columns={priceRangeColumns}
              rowKey={(r: any) => `${r.shop_code}-${r.price_range}`}
              size="small"
              pagination={false}
            />
          </Card>
        </Col>
      </Row>
    </Spin>
  )
}

export default ShopAnalysis
