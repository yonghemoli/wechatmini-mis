import React, { useEffect, useState } from 'react'
import { Card, Col, Row, Spin, Statistic, Typography, message } from 'antd'
import {
  HeartOutlined,
  GiftOutlined,
  TeamOutlined,
  TrophyOutlined
} from '@ant-design/icons'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer
} from 'recharts'
import { apiGetCompanionAnalysis } from '@/api'

const { Title } = Typography

const CompanionAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetCompanionAnalysis()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载道侣数据失败')
    } finally {
      setLoading(false)
    }
  }

  const stats = data?.stats || {}
  const giftUsage = data?.gift_usage || []
  const intimacyDist = data?.intimacy_dist || []

  return (
    <Spin spinning={loading}>
      <Title level={4}>道侣系统分析</Title>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="道侣对数"
              value={stats.total_couples || 0}
              prefix={<HeartOutlined />}
              valueStyle={{ color: '#eb2f96' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="普通誓约"
              value={stats.normal_oath_count || 0}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="高级誓约"
              value={stats.advanced_oath_count || 0}
              prefix={<TrophyOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="礼物赠送总数"
              value={stats.total_gifts_sent || 0}
              prefix={<GiftOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Card title="灵契值分布">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={intimacyDist}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="bracket" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="pair_count" name="道侣对数" fill="#eb2f96" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="月度赠礼趋势">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={giftUsage}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="month" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="send_count" name="赠送次数" fill="#52c41a" />
                <Bar dataKey="user_count" name="参与人数" fill="#1890ff" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>
    </Spin>
  )
}

export default CompanionAnalysis
