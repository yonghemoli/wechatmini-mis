import React, { useEffect, useState } from 'react'
import {
  Table,
  Card,
  Tag,
  Typography,
  Tabs,
  Statistic,
  Row,
  Col,
  Spin,
  Space
} from 'antd'
import {
  WarningOutlined,
  StopOutlined,
  AlertOutlined,
  ThunderboltOutlined,
  CrownOutlined,
  RobotOutlined
} from '@ant-design/icons'
import {
  apiGetHighChurnUsers,
  apiGetStuckUsers,
  apiGetProfileStats,
  apiGetResourceAlertUsers,
  apiGetNewbieAtRisk,
  apiGetWhaleChurn,
  apiGetBotSuspects
} from '@/api'

const { Title } = Typography

const lifecycleLabels: Record<string, string> = {
  NEW: '新手期',
  GROWING: '成长期',
  MATURE: '成熟期',
  DECLINING: '衰退期',
  LOST: '流失期',
  RETURNED: '回流期'
}

const payTierLabels: Record<string, string> = {
  FREE: '免费玩家',
  MINNOW: '微氪',
  DOLPHIN: '小氪',
  ORCA: '中氪',
  WHALE: '大氪',
  LEVIATHAN: '巨鲸'
}

const lifecycleColors: Record<string, string> = {
  NEW: 'green',
  GROWING: 'blue',
  MATURE: 'purple',
  DECLINING: 'orange',
  LOST: 'red',
  RETURNED: 'cyan'
}

const payTierColors: Record<string, string> = {
  FREE: 'default',
  MINNOW: 'green',
  DOLPHIN: 'blue',
  ORCA: 'purple',
  WHALE: 'gold',
  LEVIATHAN: 'red'
}

const playStyleLabels: Record<string, string> = {
  COMBAT: '战斗型',
  CRAFT: '制造型',
  SOCIAL: '社交型',
  ECONOMY: '经济型',
  EXPLORER: '探索型',
  BALANCED: '均衡型'
}

const Alerts: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [stats, setStats] = useState<any>({})
  const [churnUsers, setChurnUsers] = useState<any[]>([])
  const [stuckUsers, setStuckUsers] = useState<any[]>([])
  const [resourceUsers, setResourceUsers] = useState<any[]>([])
  const [newbieRisk, setNewbieRisk] = useState<any[]>([])
  const [whaleChurn, setWhaleChurn] = useState<any[]>([])
  const [botSuspects, setBotSuspects] = useState<any[]>([])
  const [activeTab, setActiveTab] = useState('churn')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const [
        statsRes,
        churnRes,
        stuckRes,
        resourceRes,
        newbieRes,
        whaleRes,
        botRes
      ] = await Promise.all([
        apiGetProfileStats(),
        apiGetHighChurnUsers(60, 100),
        apiGetStuckUsers(100),
        apiGetResourceAlertUsers(100),
        apiGetNewbieAtRisk(40, 100),
        apiGetWhaleChurn(50, 100),
        apiGetBotSuspects(100)
      ])
      if (statsRes?.data) setStats(statsRes.data)
      if (churnRes?.data) setChurnUsers(churnRes.data || [])
      if (stuckRes?.data) setStuckUsers(stuckRes.data || [])
      if (resourceRes?.data) setResourceUsers(resourceRes.data || [])
      if (newbieRes?.data) setNewbieRisk(newbieRes.data || [])
      if (whaleRes?.data) setWhaleChurn(whaleRes.data || [])
      if (botRes?.data) setBotSuspects(botRes.data || [])
    } finally {
      setLoading(false)
    }
  }

  // ==================== 通用列定义 ====================
  const colUID = { title: 'UID', dataIndex: 'uid', key: 'uid', width: 80 }
  const colLifecycle = {
    title: '生命周期',
    dataIndex: 'lifecycle_stage',
    key: 'lifecycle_stage',
    render: (v: string) => (
      <Tag color={lifecycleColors[v] || 'default'}>
        {lifecycleLabels[v] || v}
      </Tag>
    )
  }
  const colPayTier = {
    title: '付费等级',
    dataIndex: 'pay_tier',
    key: 'pay_tier',
    render: (v: string) => (
      <Tag color={payTierColors[v] || 'default'}>{payTierLabels[v] || v}</Tag>
    )
  }
  const colPlayStyle = {
    title: '玩法偏好',
    dataIndex: 'play_style',
    key: 'play_style',
    render: (v: string) => playStyleLabels[v] || v
  }
  const colChurnRisk = {
    title: '流失风险',
    dataIndex: 'churn_risk',
    key: 'churn_risk',
    defaultSortOrder: 'descend' as const,
    sorter: (a: any, b: any) => a.churn_risk - b.churn_risk,
    render: (v: number) => (
      <Tag color={v >= 85 ? 'red' : v >= 60 ? 'orange' : 'green'}>{v}%</Tag>
    )
  }
  const colLTV = {
    title: 'LTV预测',
    dataIndex: 'ltv_predict',
    key: 'ltv_predict',
    render: (v: number) => `¥${(v || 0).toFixed(2)}`
  }

  if (loading)
    return (
      <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />
    )

  return (
    <div>
      <Title level={3}>风险预警</Title>

      {/* 统计概览 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={12} sm={8} md={4}>
          <Card size="small">
            <Statistic
              title="总画像数"
              value={stats.total || 0}
              prefix={<AlertOutlined />}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={4}>
          <Card size="small">
            <Statistic
              title="高流失风险"
              value={stats.high_churn || 0}
              valueStyle={{ color: '#cf1322' }}
              prefix={<WarningOutlined />}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={4}>
          <Card size="small">
            <Statistic
              title="卡关用户"
              value={stats.stuck || 0}
              valueStyle={{ color: '#faad14' }}
              prefix={<StopOutlined />}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={4}>
          <Card size="small">
            <Statistic
              title="资源告急"
              value={stats.resource_alert || 0}
              valueStyle={{ color: '#fa541c' }}
              prefix={<ThunderboltOutlined />}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={4}>
          <Card size="small">
            <Statistic
              title="新手流失"
              value={stats.newbie_at_risk || 0}
              valueStyle={{ color: '#eb2f96' }}
              prefix={<WarningOutlined />}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={4}>
          <Card size="small">
            <Statistic
              title="大氪沉默"
              value={stats.whale_churn || 0}
              valueStyle={{ color: '#722ed1' }}
              prefix={<CrownOutlined />}
            />
          </Card>
        </Col>
      </Row>

      {/* 分 Tab 预警列表 */}
      <Card>
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={[
            {
              key: 'churn',
              label: (
                <Space>
                  <WarningOutlined />
                  {`流失风险 (${churnUsers.length})`}
                </Space>
              ),
              children: (
                <Table
                  columns={[
                    colUID,
                    colLifecycle,
                    colPayTier,
                    colChurnRisk,
                    colLTV
                  ]}
                  dataSource={churnUsers}
                  rowKey="uid"
                  pagination={{ pageSize: 20 }}
                  size="small"
                />
              )
            },
            {
              key: 'stuck',
              label: (
                <Space>
                  <StopOutlined />
                  {`卡关用户 (${stuckUsers.length})`}
                </Space>
              ),
              children: (
                <Table
                  columns={[colUID, colPayTier, colPlayStyle, colChurnRisk]}
                  dataSource={stuckUsers}
                  rowKey="uid"
                  pagination={{ pageSize: 20 }}
                  size="small"
                />
              )
            },
            {
              key: 'resource',
              label: (
                <Space>
                  <ThunderboltOutlined />
                  {`资源告急 (${resourceUsers.length})`}
                </Space>
              ),
              children: (
                <Table
                  columns={[
                    colUID,
                    colLifecycle,
                    colPayTier,
                    colPlayStyle,
                    colChurnRisk
                  ]}
                  dataSource={resourceUsers}
                  rowKey="uid"
                  pagination={{ pageSize: 20 }}
                  size="small"
                />
              )
            },
            {
              key: 'newbie',
              label: (
                <Space>
                  <WarningOutlined style={{ color: '#eb2f96' }} />
                  {`新手流失 (${newbieRisk.length})`}
                </Space>
              ),
              children: (
                <Table
                  columns={[colUID, colPlayStyle, colChurnRisk, colLTV]}
                  dataSource={newbieRisk}
                  rowKey="uid"
                  pagination={{ pageSize: 20 }}
                  size="small"
                />
              )
            },
            {
              key: 'whale',
              label: (
                <Space>
                  <CrownOutlined style={{ color: '#722ed1' }} />
                  {`大氪沉默 (${whaleChurn.length})`}
                </Space>
              ),
              children: (
                <Table
                  columns={[
                    colUID,
                    colLifecycle,
                    { ...colPayTier },
                    colChurnRisk,
                    colLTV
                  ]}
                  dataSource={whaleChurn}
                  rowKey="uid"
                  pagination={{ pageSize: 20 }}
                  size="small"
                />
              )
            },
            {
              key: 'bot',
              label: (
                <Space>
                  <RobotOutlined />
                  {`在线异常 (${botSuspects.length})`}
                </Space>
              ),
              children: (
                <Table
                  columns={[
                    { title: 'UID', dataIndex: 'uid', key: 'uid', width: 80 },
                    {
                      title: '活跃天数',
                      dataIndex: 'active_days',
                      key: 'active_days',
                      render: (v: number) => `${v} 天`
                    },
                    {
                      title: '日均在线时长',
                      dataIndex: 'avg_hours',
                      key: 'avg_hours',
                      defaultSortOrder: 'descend' as const,
                      sorter: (a: any, b: any) => a.avg_hours - b.avg_hours,
                      render: (v: number) => (
                        <Tag
                          color={v >= 22 ? 'red' : v >= 18 ? 'orange' : 'green'}
                        >
                          {v} 小时
                        </Tag>
                      )
                    },
                    {
                      title: '总记录数',
                      dataIndex: 'total_records',
                      key: 'total_records'
                    }
                  ]}
                  dataSource={botSuspects}
                  rowKey="uid"
                  pagination={{ pageSize: 20 }}
                  size="small"
                />
              )
            }
          ]}
        />
      </Card>
    </div>
  )
}

export default Alerts
