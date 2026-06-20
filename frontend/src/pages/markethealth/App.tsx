import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Statistic,
  Table,
  Tabs,
  Tag,
  Typography,
  message
} from 'antd'
import {
  ShopOutlined,
  FireOutlined,
  WarningOutlined,
  DollarOutlined,
  SwapOutlined,
  TeamOutlined
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
  Cell
} from 'recharts'
import { apiGetMarketHealth } from '@/api'

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

const MarketHealth: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [days])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetMarketHealth(days)
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载交易市场数据失败')
    } finally {
      setLoading(false)
    }
  }

  const stats = data?.trade_stats
  const hotItems = data?.hot_items || []
  const dailyTrend = data?.daily_trend || []
  const anomalies = data?.price_anomalies || []
  const largeTrades = data?.large_trades || []
  const playerProfiles = data?.player_profiles || []

  const suspiciousPlayers = playerProfiles.filter(
    (item: any) => item.suspicious_score >= 40
  )
  const mostlySellers = playerProfiles.filter(
    (item: any) => item.trade_role === 'mostly_seller'
  )
  const mostlyBuyers = playerProfiles.filter(
    (item: any) => item.trade_role === 'mostly_buyer'
  )
  const concentratedPlayers = playerProfiles.filter(
    (item: any) =>
      (item.top_buyer?.concentration || 0) >= 65 ||
      (item.top_seller?.concentration || 0) >= 65
  )

  const formatRole = (role: string) => {
    if (role === 'mostly_seller') return <Tag color="volcano">偏卖方</Tag>
    if (role === 'mostly_buyer') return <Tag color="geekblue">偏买方</Tag>
    return <Tag color="green">相对平衡</Tag>
  }

  const formatRisk = (record: any) => {
    if (record.suspicious_level === 'high') {
      return <Tag color="red">高风险 {record.suspicious_score}</Tag>
    }
    if (record.suspicious_level === 'medium') {
      return <Tag color="orange">中风险 {record.suspicious_score}</Tag>
    }
    return <Tag color="default">低风险 {record.suspicious_score}</Tag>
  }

  const renderCounterparty = (party: any, direction: 'sell' | 'buy') => {
    if (!party?.user_id) return '-'
    return (
      <div>
        <div>
          {direction === 'sell' ? '长期卖给' : '长期买自'} #{party.user_id}{' '}
          {party.user_name || '未知玩家'}
        </div>
        <div style={{ color: '#8c8c8c', fontSize: 12 }}>
          {party.trade_count} 笔 / {party.total_volume?.toLocaleString()} / 占比{' '}
          {party.concentration?.toFixed(1)}%
        </div>
      </div>
    )
  }

  const hotColumns = [
    {
      title: '分类',
      dataIndex: 'item_category',
      key: 'item_category',
      render: (v: string) => <Tag color="cyan">{v}</Tag>
    },
    {
      title: '物品',
      dataIndex: 'item_name',
      key: 'item_name',
      render: (v: string) => <Tag color="blue">{v}</Tag>
    },
    {
      title: '交易次数',
      dataIndex: 'trade_count',
      key: 'trade_count',
      sorter: (a: any, b: any) => a.trade_count - b.trade_count
    },
    {
      title: '总成交额',
      dataIndex: 'total_volume',
      key: 'total_volume',
      sorter: (a: any, b: any) => a.total_volume - b.total_volume,
      render: (v: number) => (
        <span style={{ color: '#52c41a' }}>{v?.toLocaleString()}</span>
      )
    },
    {
      title: '买家数',
      dataIndex: 'buyer_count',
      key: 'buyer_count',
      sorter: (a: any, b: any) => a.buyer_count - b.buyer_count
    }
  ]

  const anomalyColumns = [
    {
      title: '分类',
      dataIndex: 'item_category',
      key: 'item_category',
      render: (v: string) => <Tag color="cyan">{v}</Tag>
    },
    { title: '物品', dataIndex: 'item_name', key: 'item_name' },
    {
      title: '基准价',
      dataIndex: 'base_price',
      key: 'base_price',
      render: (v: number) => v?.toLocaleString()
    },
    {
      title: '均价',
      dataIndex: 'avg_price',
      key: 'avg_price',
      render: (v: number) => v?.toFixed(0)
    },
    {
      title: '最低价',
      dataIndex: 'min_price',
      key: 'min_price',
      render: (v: number) => v?.toLocaleString()
    },
    {
      title: '最高价',
      dataIndex: 'max_price',
      key: 'max_price',
      render: (v: number) => v?.toLocaleString()
    },
    {
      title: '偏差(%)',
      dataIndex: 'deviation',
      key: 'deviation',
      sorter: (a: any, b: any) => a.deviation - b.deviation,
      render: (v: number) => (
        <Tag
          color={
            Math.abs(v) > 50 ? 'red' : Math.abs(v) > 20 ? 'orange' : 'green'
          }
        >
          {v?.toFixed(1)}%
          {Math.abs(v) > 50 && <WarningOutlined style={{ marginLeft: 4 }} />}
        </Tag>
      )
    },
    {
      title: '交易次数',
      dataIndex: 'trade_count',
      key: 'trade_count',
      sorter: (a: any, b: any) => a.trade_count - b.trade_count
    }
  ]

  const largeColumns = [
    {
      title: '分类',
      dataIndex: 'item_category',
      key: 'item_category',
      render: (v: string) => <Tag color="cyan">{v}</Tag>
    },
    { title: '物品', dataIndex: 'item_name', key: 'item_name' },
    { title: '卖家ID', dataIndex: 'seller_id', key: 'seller_id' },
    { title: '买家ID', dataIndex: 'buyer_id', key: 'buyer_id' },
    {
      title: '成交价',
      dataIndex: 'total_price',
      key: 'total_price',
      sorter: (a: any, b: any) => a.total_price - b.total_price,
      render: (v: number) => (
        <span style={{ color: '#ff4d4f', fontWeight: 'bold' }}>
          {v?.toLocaleString()}
        </span>
      )
    },
    { title: '数量', dataIndex: 'quantity', key: 'quantity' },
    {
      title: '时间',
      dataIndex: 'trade_time',
      key: 'trade_time'
    }
  ]

  const suspiciousColumns = [
    {
      title: '玩家',
      dataIndex: 'user_id',
      key: 'user_id',
      render: (_: number, row: any) => (
        <div>
          <div>#{row.user_id}</div>
          <div style={{ color: '#8c8c8c', fontSize: 12 }}>
            {row.user_name || '未知玩家'}
          </div>
        </div>
      )
    },
    {
      title: '交易角色',
      dataIndex: 'trade_role',
      key: 'trade_role',
      render: (v: string) => formatRole(v)
    },
    {
      title: '买/卖比',
      dataIndex: 'buy_sell_ratio',
      key: 'buy_sell_ratio',
      sorter: (a: any, b: any) => a.buy_sell_ratio - b.buy_sell_ratio,
      render: (v: number, row: any) => (
        <span>
          {v?.toFixed(2)}
          <span style={{ color: '#8c8c8c', marginLeft: 8 }}>
            买 {row.buy_trade_count} / 卖 {row.sell_trade_count}
          </span>
        </span>
      )
    },
    {
      title: '固定对手方',
      key: 'counterparty',
      render: (_: any, row: any) => (
        <div>
          <div>{renderCounterparty(row.top_buyer, 'sell')}</div>
          <div style={{ marginTop: 6 }}>
            {renderCounterparty(row.top_seller, 'buy')}
          </div>
        </div>
      )
    },
    {
      title: '低价长期给同一人',
      key: 'low_price_transfer',
      render: (_: any, row: any) =>
        row.low_price_to_user_id ? (
          <div>
            <div>
              #{row.low_price_to_user_id}{' '}
              {row.low_price_to_user_name || '未知玩家'}
            </div>
            <div style={{ color: '#8c8c8c', fontSize: 12 }}>
              {row.low_price_item_name || '-'} / 低于市场均价{' '}
              {row.low_price_discount_rate?.toFixed(1)}% /{' '}
              {row.low_price_trade_count} 笔
            </div>
          </div>
        ) : (
          '-'
        )
    },
    {
      title: '疑似小号概率',
      dataIndex: 'suspicious_score',
      key: 'suspicious_score',
      sorter: (a: any, b: any) => a.suspicious_score - b.suspicious_score,
      render: (_: number, row: any) => (
        <div>
          <div>{formatRisk(row)}</div>
          <div style={{ color: '#8c8c8c', fontSize: 12, marginTop: 4 }}>
            概率 {row.suspicious_alt_likelihood?.toFixed(0)}%
          </div>
        </div>
      )
    },
    {
      title: '判定依据',
      dataIndex: 'suspicious_reasons',
      key: 'suspicious_reasons',
      render: (reasons: string[]) =>
        reasons?.length ? (
          <div>
            {reasons.map(reason => (
              <Tag key={reason} color="red">
                {reason}
              </Tag>
            ))}
          </div>
        ) : (
          '-'
        )
    }
  ]

  const playerColumns = [
    {
      title: '玩家',
      dataIndex: 'user_id',
      key: 'user_id',
      render: (_: number, row: any) => (
        <div>
          <div>
            #{row.user_id} {row.user_name || '未知玩家'}
          </div>
          <div style={{ color: '#8c8c8c', fontSize: 12 }}>
            VIP {row.vip_level} / 账号 {row.account_age_days} 天 / 总充值{' '}
            {row.total_recharge_amount?.toLocaleString()} / 真实充值{' '}
            {row.real_recharge_amount?.toLocaleString()}
          </div>
        </div>
      )
    },
    {
      title: '角色',
      dataIndex: 'trade_role',
      key: 'trade_role',
      render: (v: string) => formatRole(v)
    },
    {
      title: '交易次数',
      dataIndex: 'total_trade_count',
      key: 'total_trade_count',
      sorter: (a: any, b: any) => a.total_trade_count - b.total_trade_count,
      render: (_: number, row: any) => (
        <span>
          {row.total_trade_count}
          <span style={{ color: '#8c8c8c', marginLeft: 8 }}>
            买 {row.buy_trade_count} / 卖 {row.sell_trade_count}
          </span>
        </span>
      )
    },
    {
      title: '交易额',
      dataIndex: 'total_trade_volume',
      key: 'total_trade_volume',
      sorter: (a: any, b: any) => a.total_trade_volume - b.total_trade_volume,
      render: (_: number, row: any) => (
        <span>
          {row.total_trade_volume?.toLocaleString()}
          <span style={{ color: '#8c8c8c', marginLeft: 8 }}>
            买 {row.buy_volume?.toLocaleString()} / 卖{' '}
            {row.sell_volume?.toLocaleString()}
          </span>
        </span>
      )
    },
    {
      title: '买卖比',
      dataIndex: 'buy_sell_ratio',
      key: 'buy_sell_ratio',
      sorter: (a: any, b: any) => a.buy_sell_ratio - b.buy_sell_ratio,
      render: (v: number, row: any) => (
        <span>
          买/卖 {v?.toFixed(2)}
          <span style={{ color: '#8c8c8c', marginLeft: 8 }}>
            卖/买 {row.sell_buy_ratio?.toFixed(2)}
          </span>
        </span>
      )
    },
    {
      title: '对手数',
      dataIndex: 'counterparty_count',
      key: 'counterparty_count',
      sorter: (a: any, b: any) => a.counterparty_count - b.counterparty_count
    },
    {
      title: '长期关系',
      key: 'relationships',
      render: (_: any, row: any) => (
        <div>
          <div>{renderCounterparty(row.top_buyer, 'sell')}</div>
          <div style={{ marginTop: 6 }}>
            {renderCounterparty(row.top_seller, 'buy')}
          </div>
        </div>
      )
    },
    {
      title: '风险',
      dataIndex: 'suspicious_score',
      key: 'suspicious_score',
      sorter: (a: any, b: any) => a.suspicious_score - b.suspicious_score,
      render: (_: number, row: any) => formatRisk(row)
    }
  ]

  const tabItems = [
    {
      key: 'overview',
      label: (
        <span>
          <ShopOutlined /> 市场总览
        </span>
      ),
      children: (
        <>
          {stats && (
            <Row gutter={16} style={{ marginBottom: 16 }}>
              <Col span={5}>
                <Card>
                  <Statistic
                    title="总交易次数"
                    value={stats.total_trades}
                    prefix={<SwapOutlined />}
                  />
                </Card>
              </Col>
              <Col span={5}>
                <Card>
                  <Statistic
                    title="总成交额"
                    value={stats.total_volume}
                    prefix={<DollarOutlined />}
                    valueStyle={{ color: '#52c41a' }}
                  />
                </Card>
              </Col>
              <Col span={5}>
                <Card>
                  <Statistic
                    title="活跃交易者"
                    value={stats.unique_traders}
                    valueStyle={{ color: '#1890ff' }}
                  />
                </Card>
              </Col>
              <Col span={5}>
                <Card>
                  <Statistic
                    title="总手续费"
                    value={stats.total_fees}
                    prefix={<DollarOutlined />}
                    valueStyle={{ color: '#fa8c16' }}
                  />
                </Card>
              </Col>
              <Col span={4}>
                <Card>
                  <Statistic
                    title="均价"
                    value={stats.avg_price?.toFixed(0)}
                    valueStyle={{ color: '#722ed1' }}
                  />
                </Card>
              </Col>
            </Row>
          )}
          <Card title="每日交易趋势">
            <ResponsiveContainer width="100%" height={350}>
              <LineChart data={dailyTrend}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis yAxisId="left" />
                <YAxis yAxisId="right" orientation="right" />
                <Tooltip />
                <Legend />
                <Line
                  yAxisId="left"
                  type="monotone"
                  dataKey="trade_count"
                  name="交易次数"
                  stroke="#1890ff"
                  strokeWidth={2}
                />
                <Line
                  yAxisId="right"
                  type="monotone"
                  dataKey="total_volume"
                  name="成交额"
                  stroke="#52c41a"
                  strokeWidth={2}
                />
                <Line
                  yAxisId="left"
                  type="monotone"
                  dataKey="user_count"
                  name="交易人数"
                  stroke="#722ed1"
                  strokeWidth={2}
                />
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </>
      )
    },
    {
      key: 'hot',
      label: (
        <span>
          <FireOutlined /> 热门物品
        </span>
      ),
      children: (
        <Row gutter={16}>
          <Col span={14}>
            <Card title="热门交易物品">
              <Table
                columns={hotColumns}
                dataSource={hotItems}
                rowKey="item_name"
                size="small"
                pagination={{ pageSize: 10 }}
              />
            </Card>
          </Col>
          <Col span={10}>
            <Card title="交易额TOP10">
              <ResponsiveContainer width="100%" height={350}>
                <BarChart data={hotItems.slice(0, 10)} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" />
                  <YAxis dataKey="item_name" type="category" width={80} />
                  <Tooltip />
                  <Bar dataKey="total_volume" name="总交易额" fill="#1890ff">
                    {hotItems.slice(0, 10).map((_: any, i: number) => (
                      <Cell key={i} fill={COLORS[i % COLORS.length]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>
      )
    },
    {
      key: 'anomaly',
      label: (
        <span>
          <WarningOutlined /> 价格异常
        </span>
      ),
      children: (
        <Card title="价格异常检测（价差比 = 最高价/最低价）">
          <Table
            columns={anomalyColumns}
            dataSource={anomalies}
            rowKey="item_name"
            size="small"
            pagination={{ pageSize: 15 }}
          />
        </Card>
      )
    },
    {
      key: 'large',
      label: (
        <span>
          <DollarOutlined /> 大额交易
        </span>
      ),
      children: (
        <Card title="大额交易记录">
          <Table
            columns={largeColumns}
            dataSource={largeTrades}
            rowKey={(r: any) => `${r.seller_id}-${r.buyer_id}-${r.trade_time}`}
            size="small"
            pagination={{ pageSize: 15 }}
          />
        </Card>
      )
    },
    {
      key: 'players',
      label: (
        <span>
          <TeamOutlined /> 玩家画像
        </span>
      ),
      children: (
        <>
          <Card style={{ marginBottom: 16 }}>
            <Row gutter={16}>
              <Col span={6}>
                <Statistic
                  title="疑似小号玩家"
                  value={suspiciousPlayers.length}
                  valueStyle={{ color: '#ff4d4f' }}
                />
              </Col>
              <Col span={6}>
                <Statistic
                  title="总在卖的玩家"
                  value={mostlySellers.length}
                  valueStyle={{ color: '#fa541c' }}
                />
              </Col>
              <Col span={6}>
                <Statistic
                  title="总在买的玩家"
                  value={mostlyBuyers.length}
                  valueStyle={{ color: '#1677ff' }}
                />
              </Col>
              <Col span={6}>
                <Statistic
                  title="固定对手方过强"
                  value={concentratedPlayers.length}
                  valueStyle={{ color: '#722ed1' }}
                />
              </Col>
            </Row>
          </Card>

          <Tabs
            items={[
              {
                key: 'suspicious-players',
                label: '疑似小号画像',
                children: (
                  <Card
                    title="疑似小号交易画像"
                    extra={
                      <span style={{ color: '#8c8c8c' }}>
                        主题：单向交易、固定对手方、低付费高流转、年轻账号异常活跃
                      </span>
                    }
                  >
                    <Table
                      columns={suspiciousColumns}
                      dataSource={suspiciousPlayers}
                      rowKey="user_id"
                      size="small"
                      pagination={{ pageSize: 8 }}
                    />
                  </Card>
                )
              },
              {
                key: 'all-players',
                label: '交易玩家画像',
                children: (
                  <Card title="交易玩家全量画像 Top100">
                    <Table
                      columns={playerColumns}
                      dataSource={playerProfiles}
                      rowKey="user_id"
                      size="small"
                      scroll={{ x: 1400 }}
                      pagination={{ pageSize: 12 }}
                    />
                  </Card>
                )
              }
            ]}
          />
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
          交易市场健康度
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

export default MarketHealth
