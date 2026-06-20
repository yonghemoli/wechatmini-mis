import React, { useEffect, useState } from 'react'
import { Card, Col, Row, Spin, Statistic, Typography, message } from 'antd'
import {
  ExperimentOutlined,
  ToolOutlined,
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
  Legend
} from 'recharts'
import { apiGetCraftAnalysis } from '@/api'

const { Title } = Typography

const CraftAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetCraftAnalysis()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载炼丹炼器数据失败')
    } finally {
      setLoading(false)
    }
  }

  const stats = data?.stats || {}
  const levelDist = data?.level_dist || []

  // 将 level_dist 按 profession 分组用于图表
  const professions = [
    ...new Set(levelDist.map((d: any) => d.profession_name))
  ] as string[]
  const levelNames = [
    ...new Set(levelDist.map((d: any) => d.level_name))
  ] as string[]
  const chartData = levelNames.map((ln: string) => {
    const row: any = { level_name: ln }
    professions.forEach((pn: string) => {
      const item = levelDist.find(
        (d: any) => d.profession_name === pn && d.level_name === ln
      )
      row[pn] = item?.user_count || 0
    })
    return row
  })

  return (
    <Spin spinning={loading}>
      <Title level={4}>炼丹 / 炼器分析</Title>

      {/* 总览 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic
              title="炼丹配方数"
              value={stats.total_alchemy_recipes || 0}
              prefix={<ExperimentOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="炼器配方数"
              value={stats.total_lianqi_recipes || 0}
              prefix={<ToolOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="参与玩家数"
              value={stats.total_crafters || 0}
              prefix={<TeamOutlined />}
            />
          </Card>
        </Col>
      </Row>

      {/* 职业等级分布 */}
      <Card title="职业等级分布">
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="level_name" />
            <YAxis />
            <Tooltip />
            <Legend />
            {professions.map((pn: string, i: number) => (
              <Bar
                key={pn}
                dataKey={pn}
                name={pn}
                fill={['#52c41a', '#1890ff', '#faad14', '#ff4d4f'][i % 4]}
              />
            ))}
          </BarChart>
        </ResponsiveContainer>
      </Card>
    </Spin>
  )
}

export default CraftAnalysis
