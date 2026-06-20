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
  Typography,
  message
} from 'antd'
import { WarningOutlined, FallOutlined, RiseOutlined } from '@ant-design/icons'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell
} from 'recharts'
import { apiGetRealmBottleneck } from '@/api'

const { Title } = Typography

const COLORS = [
  '#ff4d4f',
  '#fa8c16',
  '#fadb14',
  '#52c41a',
  '#1890ff',
  '#722ed1',
  '#eb2f96',
  '#13c2c2'
]

const RealmBottleneck: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [inactiveDays, setInactiveDays] = useState(7)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [inactiveDays])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetRealmBottleneck(inactiveDays)
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载境界卡点数据失败')
    } finally {
      setLoading(false)
    }
  }

  const bottleneck = data?.bottleneck || []
  const churn = data?.churn || []

  // 找出人数最多的卡点
  const maxStuck =
    bottleneck.length > 0
      ? Math.max(...bottleneck.map((b: any) => b.user_count))
      : 0

  const bottleneckColumns = [
    {
      title: '境界',
      dataIndex: 'realm_name',
      key: 'realm_name',
      render: (v: string) => <Tag color="blue">{v}</Tag>
    },
    {
      title: '停留人数',
      dataIndex: 'user_count',
      key: 'user_count',
      sorter: (a: any, b: any) => a.user_count - b.user_count,
      render: (v: number) => (
        <span
          style={{
            color: v === maxStuck ? '#ff4d4f' : undefined,
            fontWeight: v === maxStuck ? 'bold' : 'normal'
          }}
        >
          {v}
          {v === maxStuck && <WarningOutlined style={{ marginLeft: 4 }} />}
        </span>
      )
    },
    {
      title: '平均停留(天)',
      dataIndex: 'avg_stay_days',
      key: 'avg_stay_days',
      sorter: (a: any, b: any) => a.avg_stay_days - b.avg_stay_days,
      render: (v: number) => (
        <span
          style={{ color: v > 14 ? '#ff4d4f' : v > 7 ? '#fa8c16' : '#52c41a' }}
        >
          {v?.toFixed(1)}
        </span>
      )
    },
    {
      title: '最长停留(天)',
      dataIndex: 'max_stay_days',
      key: 'max_stay_days',
      sorter: (a: any, b: any) => a.max_stay_days - b.max_stay_days
    },
    {
      title: '流失人数',
      dataIndex: 'churned_count',
      key: 'churned_count',
      sorter: (a: any, b: any) => a.churned_count - b.churned_count,
      render: (v: number) => <span style={{ color: '#ff4d4f' }}>{v}</span>
    },
    {
      title: '通过率',
      dataIndex: 'pass_rate',
      key: 'pass_rate',
      sorter: (a: any, b: any) => a.pass_rate - b.pass_rate,
      render: (v: number) => {
        const pct = (v ?? 0).toFixed(1)
        return (
          <Tag color={v > 80 ? 'green' : v > 50 ? 'orange' : 'red'}>{pct}%</Tag>
        )
      }
    }
  ]

  const churnColumns = [
    {
      title: '境界',
      dataIndex: 'realm_name',
      key: 'realm_name',
      render: (v: string) => <Tag color="purple">{v}</Tag>
    },
    {
      title: '总人数',
      dataIndex: 'total',
      key: 'total',
      sorter: (a: any, b: any) => a.total - b.total
    },
    {
      title: '流失人数',
      dataIndex: 'churned',
      key: 'churned',
      sorter: (a: any, b: any) => a.churned - b.churned,
      render: (v: number) => <span style={{ color: '#ff4d4f' }}>{v}</span>
    },
    {
      title: '流失率',
      dataIndex: 'churn_rate',
      key: 'churn_rate',
      sorter: (a: any, b: any) => a.churn_rate - b.churn_rate,
      render: (v: number) => {
        const pct = (v * 100).toFixed(1)
        return (
          <Tag color={v > 0.3 ? 'red' : v > 0.15 ? 'orange' : 'green'}>
            {pct}%
          </Tag>
        )
      }
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
          境界卡点分析
        </Title>
        <Select
          value={inactiveDays}
          onChange={setInactiveDays}
          style={{ width: 160 }}
        >
          <Select.Option value={3}>3天未活跃</Select.Option>
          <Select.Option value={7}>7天未活跃</Select.Option>
          <Select.Option value={14}>14天未活跃</Select.Option>
          <Select.Option value={30}>30天未活跃</Select.Option>
        </Select>
      </div>

      <Spin spinning={loading}>
        {bottleneck.length > 0 && (
          <Row gutter={16} style={{ marginBottom: 16 }}>
            <Col span={8}>
              <Card>
                <Statistic
                  title="卡点境界数"
                  value={bottleneck.length}
                  prefix={<WarningOutlined />}
                />
              </Card>
            </Col>
            <Col span={8}>
              <Card>
                <Statistic
                  title="最大卡点"
                  value={
                    bottleneck.reduce(
                      (a: any, b: any) => (a.user_count > b.user_count ? a : b),
                      bottleneck[0]
                    )?.realm_name || '-'
                  }
                  prefix={<FallOutlined />}
                  valueStyle={{ color: '#ff4d4f' }}
                />
              </Card>
            </Col>
            <Col span={8}>
              <Card>
                <Statistic
                  title="总流失人数"
                  value={churn.reduce(
                    (s: number, c: any) => s + (c.churned || 0),
                    0
                  )}
                  prefix={<RiseOutlined />}
                  valueStyle={{ color: '#fa8c16' }}
                />
              </Card>
            </Col>
          </Row>
        )}

        <Card title="卡点境界分布" style={{ marginBottom: 16 }}>
          <ResponsiveContainer width="100%" height={350}>
            <BarChart data={bottleneck}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="realm_name"
                angle={-20}
                textAnchor="end"
                height={60}
              />
              <YAxis />
              <Tooltip />
              <Bar dataKey="user_count" name="停留人数">
                {bottleneck.map((_: any, i: number) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </Card>

        <Row gutter={16}>
          <Col span={12}>
            <Card title="卡点详情">
              <Table
                columns={bottleneckColumns}
                dataSource={bottleneck}
                rowKey="realm_name"
                pagination={false}
                size="small"
              />
            </Card>
          </Col>
          <Col span={12}>
            <Card title="各境界流失情况">
              <Table
                columns={churnColumns}
                dataSource={churn}
                rowKey="realm_name"
                pagination={false}
                size="small"
              />
            </Card>
          </Col>
        </Row>
      </Spin>
    </div>
  )
}

export default RealmBottleneck
