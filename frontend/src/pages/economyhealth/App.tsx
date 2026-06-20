import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Table,
  Tag,
  Tooltip,
  Typography,
  Empty,
  Progress,
  message
} from 'antd'
import {
  BankOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  SwapOutlined
} from '@ant-design/icons'
import { apiGetEconomyHealth, apiGetCurrencyTrend } from '@/api'

const { Title, Text } = Typography

const currencyLabels: Record<string, { text: string; color: string }> = {
  SPIRIT_STONE: { text: '灵石', color: 'gold' },
  SPIRIT_CRYSTAL: { text: '灵晶', color: 'purple' },
  SPIRIT_MILK: { text: '道晶', color: 'cyan' },
  REPUTATION: { text: '声望', color: 'blue' },
  CHECKIN_COIN: { text: '签到币', color: 'green' },
  TOWER_SCORE: { text: '妖塔积分', color: 'red' },
  TONGXIN: { text: '同心值', color: 'magenta' }
}

const EconomyHealth: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<any>(null)
  const [trendCurrency, setTrendCurrency] = useState('SPIRIT_STONE')
  const [trendDays, setTrendDays] = useState(30)
  const [trend, setTrend] = useState<any[]>([])

  useEffect(() => {
    loadData()
  }, [days])
  useEffect(() => {
    loadTrend()
  }, [trendCurrency, trendDays])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetEconomyHealth(days)
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载经济数据失败')
    } finally {
      setLoading(false)
    }
  }

  const loadTrend = async () => {
    try {
      const r = await apiGetCurrencyTrend(trendCurrency, trendDays)
      if (r?.data) setTrend(r.data || [])
    } catch {
      message.error('加载趋势数据失败')
    }
  }

  const flowColumns = [
    {
      title: '币种',
      dataIndex: 'currency_code',
      key: 'currency_code',
      render: (v: string) => {
        const info = currencyLabels[v] || { text: v, color: 'default' }
        return <Tag color={info.color}>{info.text}</Tag>
      }
    },
    {
      title: '总收入',
      dataIndex: 'total_income',
      key: 'total_income',
      render: (v: number) => (
        <Text style={{ color: '#52c41a' }}>+{v?.toLocaleString()}</Text>
      ),
      sorter: (a: any, b: any) => a.total_income - b.total_income
    },
    {
      title: '总支出',
      dataIndex: 'total_expense',
      key: 'total_expense',
      render: (v: number) => (
        <Text style={{ color: '#ff4d4f' }}>-{v?.toLocaleString()}</Text>
      ),
      sorter: (a: any, b: any) => a.total_expense - b.total_expense
    },
    {
      title: '净流量',
      dataIndex: 'net_flow',
      key: 'net_flow',
      render: (v: number) => (
        <Text
          style={{ color: v >= 0 ? '#52c41a' : '#ff4d4f', fontWeight: 600 }}
        >
          {v >= 0 ? '+' : ''}
          {v?.toLocaleString()}
        </Text>
      ),
      sorter: (a: any, b: any) => a.net_flow - b.net_flow
    },
    {
      title: '交易次数',
      dataIndex: 'tx_count',
      key: 'tx_count',
      render: (v: number) => v?.toLocaleString(),
      sorter: (a: any, b: any) => a.tx_count - b.tx_count
    },
    {
      title: '活跃用户',
      dataIndex: 'unique_users',
      key: 'unique_users',
      sorter: (a: any, b: any) => a.unique_users - b.unique_users
    },
    {
      title: '人均收入',
      dataIndex: 'avg_income',
      key: 'avg_income',
      render: (v: number) => v?.toFixed(0),
      sorter: (a: any, b: any) => a.avg_income - b.avg_income
    },
    {
      title: '人均支出',
      dataIndex: 'avg_expense',
      key: 'avg_expense',
      render: (v: number) => v?.toFixed(0),
      sorter: (a: any, b: any) => a.avg_expense - b.avg_expense
    }
  ]

  const wealthBracketOrder = ['0-1K', '1K-10K', '10K-100K', '100K-1M', '1M+']
  const bracketColors = ['#ff4d4f', '#faad14', '#1890ff', '#52c41a', '#722ed1']

  const maxTrendVal = Math.max(
    ...trend.map((d: any) => d.income || 0),
    ...trend.map((d: any) => d.expense || 0),
    1
  )

  return (
    <Spin spinning={loading}>
      <div style={{ padding: 24 }}>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: 16
          }}
        >
          <Title level={4} style={{ margin: 0 }}>
            <BankOutlined style={{ marginRight: 8 }} />
            经济健康度分析
          </Title>
          <Select
            value={days}
            onChange={setDays}
            style={{ width: 100 }}
            options={[
              { label: '7天', value: 7 },
              { label: '14天', value: 14 },
              { label: '30天', value: 30 },
              { label: '90天', value: 90 }
            ]}
          />
        </div>

        {/* 货币流通总览 */}
        {data?.currency_flow && (
          <Row gutter={[12, 12]} style={{ marginBottom: 16 }}>
            {data.currency_flow.map((item: any) => {
              const info = currencyLabels[item.currency_code] || {
                text: item.currency_code,
                color: 'default'
              }
              const healthRatio =
                item.total_expense > 0
                  ? ((item.total_income / item.total_expense) * 100).toFixed(0)
                  : '∞'
              return (
                <Col xs={12} sm={8} md={6} key={item.currency_code}>
                  <Card
                    size="small"
                    style={{
                      borderLeft: `3px solid ${item.net_flow >= 0 ? '#52c41a' : '#ff4d4f'}`
                    }}
                  >
                    <div style={{ marginBottom: 4 }}>
                      <Tag color={info.color}>{info.text}</Tag>
                      <Text type="secondary" style={{ fontSize: 11 }}>
                        收支比 {healthRatio}%
                      </Text>
                    </div>
                    <div style={{ display: 'flex', gap: 12, fontSize: 12 }}>
                      <span style={{ color: '#52c41a' }}>
                        <ArrowUpOutlined />{' '}
                        {(item.total_income || 0).toLocaleString()}
                      </span>
                      <span style={{ color: '#ff4d4f' }}>
                        <ArrowDownOutlined />{' '}
                        {(item.total_expense || 0).toLocaleString()}
                      </span>
                    </div>
                  </Card>
                </Col>
              )
            })}
          </Row>
        )}

        <Row gutter={16}>
          {/* 货币流通明细 */}
          <Col xs={24} lg={14}>
            <Card
              title={
                <span>
                  <SwapOutlined style={{ marginRight: 6 }} />
                  货币流通明细
                </span>
              }
              size="small"
            >
              <Table
                dataSource={data?.currency_flow || []}
                columns={flowColumns}
                rowKey="currency_code"
                size="small"
                pagination={false}
              />
            </Card>
          </Col>

          {/* 灵石财富分布 */}
          <Col xs={24} lg={10}>
            <Card title="灵石财富分布" size="small">
              {data?.wealth_dist?.length > 0 ? (
                <div>
                  {wealthBracketOrder.map((bracket, idx) => {
                    const item = data.wealth_dist.find(
                      (d: any) => d.bracket === bracket
                    )
                    if (!item) return null
                    return (
                      <div key={bracket} style={{ marginBottom: 12 }}>
                        <div
                          style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            marginBottom: 4
                          }}
                        >
                          <Text strong>{bracket}</Text>
                          <Text type="secondary">
                            {item.user_count} 人 ({item.percentage?.toFixed(1)}
                            %)
                          </Text>
                        </div>
                        <Progress
                          percent={item.percentage}
                          showInfo={false}
                          strokeColor={bracketColors[idx]}
                          size="small"
                        />
                      </div>
                    )
                  })}
                </div>
              ) : (
                <Empty description="暂无数据" />
              )}
            </Card>
          </Col>
        </Row>

        {/* 币种日趋势 */}
        <Card
          title="币种日趋势"
          size="small"
          style={{ marginTop: 16 }}
          extra={
            <div style={{ display: 'flex', gap: 8 }}>
              <Select
                value={trendCurrency}
                onChange={setTrendCurrency}
                size="small"
                style={{ width: 100 }}
                options={Object.entries(currencyLabels).map(([k, v]) => ({
                  label: v.text,
                  value: k
                }))}
              />
              <Select
                value={trendDays}
                onChange={setTrendDays}
                size="small"
                style={{ width: 80 }}
                options={[
                  { label: '7天', value: 7 },
                  { label: '14天', value: 14 },
                  { label: '30天', value: 30 }
                ]}
              />
            </div>
          }
        >
          {trend.length > 0 ? (
            <>
              <div
                style={{
                  display: 'flex',
                  alignItems: 'stretch',
                  height: 140,
                  gap: 2,
                  padding: '0 4px',
                  borderBottom: '1px solid #f0f0f0'
                }}
              >
                {trend.map((d: any) => (
                  <Tooltip
                    key={d.date}
                    title={
                      <div>
                        <div>{d.date}</div>
                        <div style={{ color: '#52c41a' }}>
                          收入 {d.income?.toLocaleString()}
                        </div>
                        <div style={{ color: '#ff7875' }}>
                          支出 {d.expense?.toLocaleString()}
                        </div>
                        <div>交易 {d.tx_count} 笔</div>
                      </div>
                    }
                  >
                    <div
                      style={{
                        flex: 1,
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'flex-end',
                        gap: 1,
                        height: '100%'
                      }}
                    >
                      <div
                        style={{
                          height: `${((d.income || 0) / maxTrendVal) * 100}%`,
                          background: '#52c41a',
                          borderRadius: '2px 2px 0 0',
                          minHeight: d.income > 0 ? 2 : 0
                        }}
                      />
                      <div
                        style={{
                          height: `${((d.expense || 0) / maxTrendVal) * 100}%`,
                          background: '#ff4d4f',
                          borderRadius: '0 0 2px 2px',
                          minHeight: d.expense > 0 ? 2 : 0
                        }}
                      />
                    </div>
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
                <span>{trend[0]?.date}</span>
                <span>
                  <span style={{ color: '#52c41a' }}>■ 收入</span>{' '}
                  <span style={{ color: '#ff4d4f' }}>■ 支出</span>
                </span>
                <span>{trend[trend.length - 1]?.date}</span>
              </div>
            </>
          ) : (
            <Empty description="暂无趋势数据" />
          )}
        </Card>
      </div>
    </Spin>
  )
}

export default EconomyHealth
