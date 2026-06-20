import React, { useEffect, useState } from 'react'
import { Card, Col, Row, Spin, Statistic, Typography, message } from 'antd'
import { BookOutlined, UserOutlined, StarOutlined } from '@ant-design/icons'
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
import { apiGetSkillAnalysis } from '@/api'

const { Title } = Typography

const GRADE_COLORS: Record<string, string> = {
  凡物: '#8c8c8c',
  法器: '#52c41a',
  灵器: '#1890ff',
  灵宝: '#722ed1',
  仙器: '#faad14',
  道器: '#ff4d4f'
}

interface SkillGradeItem {
  grade_name: string
  grade_level: number
  skill_count: number
}

interface UserSkillStats {
  total_users: number
  users_with_skill: number
  avg_skill_count: number
  max_skill_level: number
}

interface SkillAnalysisData {
  grade_dist: SkillGradeItem[]
  user_stats: UserSkillStats
}

const SkillAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<SkillAnalysisData | null>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetSkillAnalysis()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载功法数据失败')
    } finally {
      setLoading(false)
    }
  }

  const gradeDist: SkillGradeItem[] = data?.grade_dist ?? []
  const userStats: UserSkillStats = data?.user_stats ?? {
    total_users: 0,
    users_with_skill: 0,
    avg_skill_count: 0,
    max_skill_level: 0
  }

  return (
    <Spin spinning={loading}>
      <Title level={4}>功法 / 技能分析</Title>

      {/* 总览 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={userStats.total_users || 0}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="拥有功法的玩家"
              value={userStats.users_with_skill || 0}
              prefix={<BookOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="人均功法数"
              value={userStats.avg_skill_count?.toFixed(1) || 0}
              prefix={<StarOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="最高功法等级"
              value={userStats.max_skill_level || 0}
              prefix={<StarOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 品级分布 */}
      <Card title="功法品级分布">
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={gradeDist}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="grade_name" />
            <YAxis />
            <Tooltip />
            <Bar dataKey="skill_count" name="功法数量">
              {gradeDist.map((entry, i: number) => (
                <Cell
                  key={i}
                  fill={GRADE_COLORS[entry.grade_name] || '#1890ff'}
                />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </Card>
    </Spin>
  )
}

export default SkillAnalysis
