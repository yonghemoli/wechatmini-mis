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
  StarOutlined,
  UserOutlined,
  RiseOutlined,
  CrownOutlined
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
import { apiGetSoulArtifactAnalysis } from '@/api'

const { Title } = Typography
const CATEGORY_COLORS: Record<string, string> = {
  物理: '#ff4d4f',
  法术: '#1890ff',
  防御: '#52c41a',
  速度: '#faad14',
  全能: '#722ed1'
}
const COLORS = [
  '#ff4d4f',
  '#1890ff',
  '#52c41a',
  '#faad14',
  '#722ed1',
  '#13c2c2'
]

const SoulArtifactAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetSoulArtifactAnalysis()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载本命物数据失败')
    } finally {
      setLoading(false)
    }
  }

  const userStats = data?.user_stats || {}
  const gradeDist = data?.grade_dist || []
  const categoryDist = data?.category_dist || []
  const levelDist = data?.level_dist || []
  const popularity = data?.popularity || []

  const popularityColumns = [
    { title: '本命物', dataIndex: 'artifact_name', key: 'artifact_name' },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      render: (v: string) => (
        <Tag color={CATEGORY_COLORS[v] || 'default'}>{v}</Tag>
      )
    },
    {
      title: '拥有者数',
      dataIndex: 'owner_count',
      key: 'owner_count',
      sorter: (a: any, b: any) => a.owner_count - b.owner_count,
      defaultSortOrder: 'descend' as const
    }
  ]

  return (
    <Spin spinning={loading}>
      <Title level={4}>本命物系统分析</Title>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="拥有者数"
              value={userStats.total_owners || 0}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="平均品级"
              value={(userStats.avg_grade_level || 0).toFixed(1)}
              prefix={<StarOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="平均等级"
              value={(userStats.avg_level || 0).toFixed(1)}
              prefix={<RiseOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="最高等级"
              value={userStats.max_level || 0}
              prefix={<CrownOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={12}>
          <Card title="用户品级分布">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={gradeDist}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="grade_name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="user_count" name="用户数" fill="#722ed1" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="本命物受欢迎度 TOP10">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={popularity.slice(0, 10)} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis type="number" />
                <YAxis type="category" dataKey="artifact_name" width={80} />
                <Tooltip />
                <Bar dataKey="owner_count" name="拥有者" fill="#1890ff" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={12}>
          <Card title="分类分布">
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={categoryDist}
                  dataKey="user_count"
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
        <Col span={12}>
          <Card title="等级分布">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={levelDist}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="bracket" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="user_count" name="用户数" fill="#13c2c2" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      <Card title="受欢迎度排行">
        <Table
          dataSource={popularity}
          columns={popularityColumns}
          rowKey="artifact_id"
          size="small"
          pagination={{ pageSize: 10, showTotal: t => `共 ${t} 件` }}
        />
      </Card>
    </Spin>
  )
}

export default SoulArtifactAnalysis
