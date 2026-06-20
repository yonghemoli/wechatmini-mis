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
  Tooltip,
  Typography,
  Progress,
  message
} from 'antd'
import {
  HomeOutlined,
  TeamOutlined,
  TrophyOutlined,
  FireOutlined
} from '@ant-design/icons'
import { apiGetGuildHealth } from '@/api'

const { Title, Text } = Typography

const GuildHealth: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [inactiveDays, setInactiveDays] = useState(7)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [inactiveDays])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetGuildHealth(inactiveDays)
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载宗门数据失败')
    } finally {
      setLoading(false)
    }
  }

  const overview = data?.overview
  const guilds = data?.guilds || []

  const getHealthLevel = (activeRate: number) => {
    if (activeRate >= 80) return { text: '优秀', color: '#52c41a' }
    if (activeRate >= 60) return { text: '良好', color: '#1890ff' }
    if (activeRate >= 40) return { text: '一般', color: '#faad14' }
    return { text: '低迷', color: '#ff4d4f' }
  }

  const columns = [
    {
      title: '宗门',
      dataIndex: 'guild_name',
      key: 'guild_name',
      render: (v: string, r: any) => (
        <span>
          <Text strong>{v}</Text>
          <Tag style={{ marginLeft: 6 }}>Lv.{r.level}</Tag>
        </span>
      )
    },
    {
      title: '成员',
      key: 'members',
      render: (_: any, r: any) => (
        <Tooltip title={`${r.member_count} / ${r.member_capacity}`}>
          <Progress
            percent={(r.member_count / (r.member_capacity || 1)) * 100}
            size="small"
            format={() => `${r.member_count}/${r.member_capacity}`}
            style={{ width: 100 }}
          />
        </Tooltip>
      ),
      sorter: (a: any, b: any) => a.member_count - b.member_count
    },
    {
      title: '声望',
      dataIndex: 'prestige',
      key: 'prestige',
      render: (v: number) => v?.toLocaleString(),
      sorter: (a: any, b: any) => a.prestige - b.prestige,
      defaultSortOrder: 'descend' as const
    },
    {
      title: '活跃成员',
      dataIndex: 'active_members',
      key: 'active_members',
      sorter: (a: any, b: any) => a.active_members - b.active_members
    },
    {
      title: '活跃率',
      dataIndex: 'active_rate',
      key: 'active_rate',
      render: (v: number) => {
        const health = getHealthLevel(v)
        return (
          <Tag color={health.color}>
            {v?.toFixed(1)}% {health.text}
          </Tag>
        )
      },
      sorter: (a: any, b: any) => a.active_rate - b.active_rate
    },
    {
      title: '平均贡献',
      dataIndex: 'avg_contribution',
      key: 'avg_contribution',
      render: (v: number) => v?.toFixed(0),
      sorter: (a: any, b: any) => a.avg_contribution - b.avg_contribution
    },
    {
      title: '总贡献',
      dataIndex: 'total_contribution',
      key: 'total_contribution',
      render: (v: number) => v?.toLocaleString(),
      sorter: (a: any, b: any) =>
        (a.total_contribution || 0) - (b.total_contribution || 0)
    },
    {
      title: '日均活跃度',
      dataIndex: 'avg_daily_activity',
      key: 'avg_daily_activity',
      render: (v: number) => v?.toFixed(1),
      sorter: (a: any, b: any) => a.avg_daily_activity - b.avg_daily_activity
    }
  ]

  // 健康度分布统计
  const healthDist = {
    excellent: guilds.filter((g: any) => g.active_rate >= 80).length,
    good: guilds.filter((g: any) => g.active_rate >= 60 && g.active_rate < 80)
      .length,
    normal: guilds.filter((g: any) => g.active_rate >= 40 && g.active_rate < 60)
      .length,
    low: guilds.filter((g: any) => g.active_rate < 40).length
  }

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
            <HomeOutlined style={{ marginRight: 8 }} />
            宗门健康度分析
          </Title>
          <Select
            value={inactiveDays}
            onChange={setInactiveDays}
            style={{ width: 120 }}
            options={[
              { label: '3天不活跃', value: 3 },
              { label: '7天不活跃', value: 7 },
              { label: '14天不活跃', value: 14 },
              { label: '30天不活跃', value: 30 }
            ]}
          />
        </div>

        {/* 整体概览 */}
        <Row gutter={[12, 12]} style={{ marginBottom: 16 }}>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="总宗门"
                value={overview?.total_guilds || 0}
                prefix={<HomeOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="总成员"
                value={overview?.total_members || 0}
                prefix={<TeamOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="平均成员"
                value={overview?.avg_member_count || 0}
                precision={1}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="最大宗门"
                value={overview?.max_members || 0}
                prefix={<TrophyOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="平均声望"
                value={overview?.avg_prestige || 0}
                precision={0}
                prefix={<FireOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="加入率"
                value={overview?.coverage_rate || 0}
                suffix="%"
                precision={1}
                valueStyle={{
                  color:
                    (overview?.coverage_rate || 0) > 50 ? '#52c41a' : '#faad14'
                }}
              />
            </Card>
          </Col>
        </Row>

        {/* 健康度分布 */}
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={24}>
            <Card title="宗门健康度分布" size="small">
              <div
                style={{
                  display: 'flex',
                  gap: 24,
                  justifyContent: 'center',
                  padding: '8px 0'
                }}
              >
                {[
                  {
                    label: '优秀(≥80%)',
                    count: healthDist.excellent,
                    color: '#52c41a'
                  },
                  {
                    label: '良好(60-80%)',
                    count: healthDist.good,
                    color: '#1890ff'
                  },
                  {
                    label: '一般(40-60%)',
                    count: healthDist.normal,
                    color: '#faad14'
                  },
                  {
                    label: '低迷(<40%)',
                    count: healthDist.low,
                    color: '#ff4d4f'
                  }
                ].map(item => (
                  <div key={item.label} style={{ textAlign: 'center' }}>
                    <div
                      style={{
                        width: 60,
                        height: 60,
                        borderRadius: '50%',
                        background: item.color,
                        color: '#fff',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 20,
                        fontWeight: 700,
                        margin: '0 auto 6px'
                      }}
                    >
                      {item.count}
                    </div>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      {item.label}
                    </Text>
                  </div>
                ))}
              </div>
            </Card>
          </Col>
        </Row>

        {/* 宗门列表 */}
        <Card title="宗门详情列表" size="small">
          <Table
            dataSource={guilds}
            columns={columns}
            rowKey="guild_id"
            size="small"
            pagination={{ pageSize: 15, showSizeChanger: true }}
            scroll={{ x: 900 }}
          />
        </Card>
      </div>
    </Spin>
  )
}

export default GuildHealth
