import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Statistic,
  Typography,
  message
} from 'antd'
import {
  ClockCircleOutlined,
  UserOutlined,
  CheckCircleOutlined,
  ThunderboltOutlined
} from '@ant-design/icons'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  Legend
} from 'recharts'
import { apiGetSeclusionAnalysis } from '@/api'

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

const SeclusionAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [days])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetSeclusionAnalysis(days)
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载闭关数据失败')
    } finally {
      setLoading(false)
    }
  }

  const stats = data?.stats || {}
  const durationDist = data?.duration_dist || []

  const statusData = [
    { name: '已完成', value: stats.completed_count || 0 },
    { name: '被中断', value: stats.interrupted_count || 0 },
    {
      name: '进行中',
      value: Math.max(
        0,
        (stats.total_sessions || 0) -
          (stats.completed_count || 0) -
          (stats.interrupted_count || 0)
      )
    }
  ].filter(d => d.value > 0)

  const premiumData = [
    { name: '高级闭关', value: stats.premium_sessions || 0 },
    {
      name: '普通闭关',
      value: Math.max(
        0,
        (stats.total_sessions || 0) - (stats.premium_sessions || 0)
      )
    }
  ].filter(d => d.value > 0)

  return (
    <Spin spinning={loading}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 16
        }}
      >
        <Title level={4} style={{ margin: 0 }}>
          闭关修炼分析
        </Title>
        <Select value={days} onChange={setDays} style={{ width: 120 }}>
          <Select.Option value={7}>近7天</Select.Option>
          <Select.Option value={14}>近14天</Select.Option>
          <Select.Option value={30}>近30天</Select.Option>
          <Select.Option value={60}>近60天</Select.Option>
        </Select>
      </div>

      {/* 总览卡片 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={4}>
          <Card>
            <Statistic
              title="总次数"
              value={stats.total_sessions || 0}
              prefix={<ClockCircleOutlined />}
            />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic
              title="参与用户"
              value={stats.unique_users || 0}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic
              title="完成率"
              value={stats.completion_rate?.toFixed(1) || 0}
              suffix="%"
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic
              title="高级闭关率"
              value={stats.premium_rate?.toFixed(1) || 0}
              suffix="%"
              prefix={<ThunderboltOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic
              title="平均时长"
              value={stats.avg_duration?.toFixed(1) || 0}
              suffix="h"
            />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic
              title="被中断"
              value={stats.interrupted_count || 0}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        {/* 状态分布 */}
        <Col span={8}>
          <Card title="闭关状态分布">
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={statusData}
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  dataKey="value"
                  label={({ name, percent }) =>
                    `${name} ${((percent ?? 0) * 100).toFixed(0)}%`
                  }
                >
                  {statusData.map((_, i) => (
                    <Cell key={i} fill={COLORS[i % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Col>

        {/* 高级/普通比例 */}
        <Col span={8}>
          <Card title="普通 vs 高级闭关">
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={premiumData}
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  dataKey="value"
                  label={({ name, percent }) =>
                    `${name} ${((percent ?? 0) * 100).toFixed(0)}%`
                  }
                >
                  {premiumData.map((_, i) => (
                    <Cell key={i} fill={COLORS[i + 3]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Col>

        {/* 时长分布 */}
        <Col span={8}>
          <Card title="闭关时长分布">
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={durationDist}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="duration" unit="h" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="count" name="次数" fill="#1890ff" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>
    </Spin>
  )
}

export default SeclusionAnalysis
