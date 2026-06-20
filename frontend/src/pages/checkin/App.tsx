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
  CheckCircleOutlined,
  CalendarOutlined,
  TrophyOutlined
} from '@ant-design/icons'
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer
} from 'recharts'
import { apiGetCheckInAnalysis } from '@/api'

const { Title } = Typography

const CheckInAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [days])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetCheckInAnalysis(days)
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载签到数据失败')
    } finally {
      setLoading(false)
    }
  }

  const rate = data?.rate || {}
  const trend = data?.trend || []

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
          签到分析
        </Title>
        <Select value={days} onChange={setDays} style={{ width: 120 }}>
          <Select.Option value={7}>近7天</Select.Option>
          <Select.Option value={14}>近14天</Select.Option>
          <Select.Option value={30}>近30天</Select.Option>
          <Select.Option value={60}>近60天</Select.Option>
        </Select>
      </div>

      {/* 签到率总览 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={rate.total_users || 0}
              prefix={<CalendarOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="今日已签到"
              value={rate.checked_in || 0}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="签到率"
              value={rate.checkin_rate?.toFixed(1) || 0}
              suffix="%"
              prefix={<TrophyOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="月卡使用率"
              value={rate.month_card_rate?.toFixed(1) || 0}
              suffix="%"
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 签到趋势 */}
      <Card title="签到趋势">
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={trend}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip />
            <Area
              type="monotone"
              dataKey="checkin_count"
              name="签到人数"
              stroke="#1890ff"
              fill="#1890ff"
              fillOpacity={0.3}
            />
            <Area
              type="monotone"
              dataKey="month_card_claimed"
              name="月卡领取"
              stroke="#faad14"
              fill="#faad14"
              fillOpacity={0.3}
            />
          </AreaChart>
        </ResponsiveContainer>
      </Card>
    </Spin>
  )
}

export default CheckInAnalysis
