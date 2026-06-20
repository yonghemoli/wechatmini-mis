import React, { useEffect, useState } from 'react'
import { Card, Col, Row, Spin, Statistic, Typography, message } from 'antd'
import {
  UserAddOutlined,
  GiftOutlined,
  FunnelPlotOutlined
} from '@ant-design/icons'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  LineChart,
  Line
} from 'recharts'
import { apiGetNewcomerFunnel } from '@/api'

const { Title } = Typography

const COLORS = [
  '#1890ff',
  '#52c41a',
  '#faad14',
  '#ff4d4f',
  '#722ed1',
  '#13c2c2',
  '#eb2f96'
]

const NewcomerFunnel: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetNewcomerFunnel()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载新人引导数据失败')
    } finally {
      setLoading(false)
    }
  }

  const funnel = data?.funnel || []

  // 计算漏斗转化率
  const funnelWithRate = funnel.map((item: any, index: number) => ({
    ...item,
    label: `Day ${item.day_index}`,
    rate:
      index === 0
        ? 100
        : funnel[0].user_count > 0
          ? ((item.user_count / funnel[0].user_count) * 100).toFixed(1)
          : 0
  }))

  const totalNewUsers = funnel.length > 0 ? funnel[0].user_count : 0
  const day7Users = funnel.find((f: any) => f.day_index === 7)?.user_count || 0
  const day7RetentionRate =
    totalNewUsers > 0 ? ((day7Users / totalNewUsers) * 100).toFixed(1) : '0'

  return (
    <Spin spinning={loading}>
      <Title level={4}>新人引导漏斗</Title>

      {/* 总览 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic
              title="新人总数 (Day1)"
              value={totalNewUsers}
              prefix={<UserAddOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="Day7 留存"
              value={day7Users}
              prefix={<GiftOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="7日留存率"
              value={day7RetentionRate}
              suffix="%"
              prefix={<FunnelPlotOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 漏斗图 */}
      <Row gutter={16}>
        <Col span={12}>
          <Card title="每日签到人数">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={funnelWithRate}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="label" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="user_count" name="签到人数">
                  {funnelWithRate.map((_: any, i: number) => (
                    <Cell key={i} fill={COLORS[i % COLORS.length]} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="留存率曲线">
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={funnelWithRate}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="label" />
                <YAxis unit="%" />
                <Tooltip
                  formatter={v => (typeof v === 'number' ? `${v}%` : v)}
                />
                <Line
                  type="monotone"
                  dataKey="rate"
                  name="留存率"
                  stroke="#1890ff"
                  strokeWidth={2}
                  dot={{ fill: '#1890ff', r: 4 }}
                />
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>
    </Spin>
  )
}

export default NewcomerFunnel
