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
import { ToolOutlined, StarOutlined, CrownOutlined } from '@ant-design/icons'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  PieChart,
  Pie
} from 'recharts'
import { apiGetEquipAnalysis } from '@/api'

const { Title } = Typography

const GRADE_COLORS: Record<string, string> = {
  凡品: '#8c8c8c',
  良品: '#52c41a',
  优品: '#1890ff',
  极品: '#722ed1',
  仙品: '#eb2f96',
  神品: '#fadb14'
}

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

const EquipAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetEquipAnalysis()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载装备分析数据失败')
    } finally {
      setLoading(false)
    }
  }

  const gradeDist = data?.grade_dist || []
  const setPop = data?.set_popularity || []
  const realmDist = data?.realm_dist || []

  // 汇总
  const totalEquips = gradeDist.reduce(
    (s: number, g: any) => s + (g.count || 0),
    0
  )

  const setColumns = [
    {
      title: '套装名',
      dataIndex: 'set_name',
      key: 'set_name',
      render: (v: string) => <Tag color="purple">{v}</Tag>
    },
    {
      title: '品级',
      dataIndex: 'grade_name',
      key: 'grade_name',
      render: (v: string) => <Tag color={GRADE_COLORS[v] || '#8c8c8c'}>{v}</Tag>
    },
    {
      title: '拥有玩家数',
      dataIndex: 'user_count',
      key: 'user_count',
      sorter: (a: any, b: any) => a.user_count - b.user_count
    },
    {
      title: '人均件数',
      dataIndex: 'piece_avg',
      key: 'piece_avg',
      render: (v: number) => v?.toFixed(1),
      sorter: (a: any, b: any) => a.piece_avg - b.piece_avg
    }
  ]

  const realmColumns = [
    {
      title: '境界',
      dataIndex: 'stage_name',
      key: 'stage_name',
      render: (v: string) => <Tag color="blue">{v}</Tag>
    },
    {
      title: '品级',
      dataIndex: 'grade_name',
      key: 'grade_name',
      render: (v: string) => <Tag color={GRADE_COLORS[v] || '#8c8c8c'}>{v}</Tag>
    },
    {
      title: '人均装备数',
      dataIndex: 'avg_count',
      key: 'avg_count',
      render: (v: number) => v?.toFixed(2),
      sorter: (a: any, b: any) => a.avg_count - b.avg_count
    }
  ]

  return (
    <div style={{ padding: 24 }}>
      <Title level={3}>装备与战力分析</Title>

      <Spin spinning={loading}>
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col span={8}>
            <Card>
              <Statistic
                title="总装备数"
                value={totalEquips}
                prefix={<ToolOutlined />}
              />
            </Card>
          </Col>
          <Col span={8}>
            <Card>
              <Statistic
                title="套装种类"
                value={setPop.length}
                prefix={<StarOutlined />}
              />
            </Card>
          </Col>
          <Col span={8}>
            <Card>
              <Statistic
                title="品质分类"
                value={gradeDist.length}
                prefix={<CrownOutlined />}
              />
            </Card>
          </Col>
        </Row>

        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col span={14}>
            <Card title="装备品质分布">
              <ResponsiveContainer width="100%" height={350}>
                <BarChart data={gradeDist}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="grade_name" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="count" name="数量">
                    {gradeDist.map((g: any, i: number) => (
                      <Cell
                        key={i}
                        fill={
                          GRADE_COLORS[g.grade_name] ||
                          COLORS[i % COLORS.length]
                        }
                      />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </Card>
          </Col>
          <Col span={10}>
            <Card title="品质占比">
              <ResponsiveContainer width="100%" height={350}>
                <PieChart>
                  <Pie
                    data={gradeDist}
                    dataKey="count"
                    nameKey="grade_name"
                    cx="50%"
                    cy="50%"
                    outerRadius={110}
                    label={({ grade_name, percent }: any) =>
                      `${grade_name} ${(percent * 100).toFixed(0)}%`
                    }
                  >
                    {gradeDist.map((g: any, i: number) => (
                      <Cell
                        key={i}
                        fill={
                          GRADE_COLORS[g.grade_name] ||
                          COLORS[i % COLORS.length]
                        }
                      />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col span={12}>
            <Card title="热门套装排行">
              <Table
                columns={setColumns}
                dataSource={setPop}
                rowKey={(r: any) => `${r.set_name}-${r.grade_name}`}
                size="small"
                pagination={{ pageSize: 10 }}
              />
            </Card>
          </Col>
          <Col span={12}>
            <Card title="各境界装备水平">
              <Table
                columns={realmColumns}
                dataSource={realmDist}
                rowKey={(r: any) => `${r.stage_name}-${r.grade_name}`}
                size="small"
                pagination={false}
              />
            </Card>
          </Col>
        </Row>
      </Spin>
    </div>
  )
}

export default EquipAnalysis
