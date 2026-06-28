import React, { useEffect, useState } from 'react'
import { Alert, Badge, Button, Col, Empty, Row, Space, Spin, Statistic, Table, Tag, Typography, message } from 'antd'
import { AlertOutlined, FileTextOutlined, ReloadOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import {
  OrderRecord,
  apiDashboardExceptions,
  apiDashboardSummary,
  orderStatusColor,
  orderStatusText
} from '@/api'

const { Title, Text } = Typography

type Summary = {
  todayOrders: number
  todayRevenue: number
  pendingCount: number
  userCount: number
}

const Dashboard: React.FC = () => {
  const navigate = useNavigate()
  const [summary, setSummary] = useState<Summary>({
    todayOrders: 0,
    todayRevenue: 0,
    pendingCount: 0,
    userCount: 0
  })
  const [orders, setOrders] = useState<OrderRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const [summaryRes, exceptionsRes] = await Promise.all([
        apiDashboardSummary(),
        apiDashboardExceptions()
      ])
      setSummary(summaryRes?.data || summary)
      setOrders(exceptionsRes?.data?.list || [])
    } catch {
      setError('工作台数据加载失败')
      message.error('工作台数据加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            工作台
          </Title>
          <Text type="secondary">先处理异常，再看增长。所有关键项控制在一屏内。</Text>
        </div>
        <Button icon={<ReloadOutlined />} onClick={load}>
          刷新
        </Button>
      </div>

      {error ? <Alert type="error" showIcon message={error} /> : null}

      <Spin spinning={loading}>
        <Row gutter={[12, 12]}>
          <Col xs={24} sm={12} lg={6}>
            <div className="mis-metric">
              <Statistic title="今日订单量" value={summary.todayOrders} suffix="单" />
            </div>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <div className="mis-metric">
              <Statistic title="今日营收" value={summary.todayRevenue} prefix="¥" precision={0} />
            </div>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <div className="mis-metric mis-metric-danger">
              <Badge count={summary.pendingCount} offset={[10, 0]}>
                <Statistic title="待处理工单" value={summary.pendingCount} />
              </Badge>
            </div>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <div className="mis-metric">
              <Statistic title="用户总数" value={summary.userCount} suffix="人" />
            </div>
          </Col>
        </Row>
      </Spin>

      <div className="mis-panel mis-alert-panel">
        <div className="mis-panel-header">
          <Space>
            <AlertOutlined />
            <Title level={5} style={{ margin: 0 }}>
              异常预警
            </Title>
            <Badge count={orders.filter(order => order.status === 'exception').length} />
          </Space>
          <Button type="primary" danger onClick={() => navigate('/admin/orders')}>
            进入订单处理
          </Button>
        </div>
        <Table
          rowKey="id"
          size="middle"
          loading={loading}
          pagination={false}
          dataSource={orders}
          locale={{ emptyText: <Empty description="暂无待处理订单" /> }}
          columns={[
            { title: '订单号', dataIndex: 'id', width: 160 },
            { title: '客户', dataIndex: 'customer', width: 100 },
            { title: '服务', dataIndex: 'service' },
            { title: '预约时间', dataIndex: 'appointmentAt', width: 160 },
            {
              title: '状态',
              dataIndex: 'status',
              width: 130,
              render: (value: OrderRecord['status']) => (
                <Tag color={orderStatusColor[value]}>{orderStatusText[value]}</Tag>
              )
            },
            { title: '内部备注', dataIndex: 'internalNote' },
            {
              title: '动作',
              width: 120,
              render: () => (
                <Button icon={<FileTextOutlined />} onClick={() => navigate('/admin/orders')}>
                  处理
                </Button>
              )
            }
          ]}
        />
      </div>
    </Space>
  )
}

export default Dashboard
