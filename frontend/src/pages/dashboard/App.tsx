import React from 'react'
import { Badge, Button, Col, Row, Space, Statistic, Table, Tag, Typography } from 'antd'
import { AlertOutlined, CheckCircleOutlined, FileTextOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { OrderStatus } from '../misData'
import { orders, statusColor, statusText, users } from '../misData'

const { Title, Text } = Typography

const Dashboard: React.FC = () => {
  const navigate = useNavigate()
  const todayOrders = orders.filter(order => order.createdAt.startsWith('2026-06-20'))
  const todayRevenue = todayOrders.reduce((sum, order) => sum + order.amount, 0)
  const pendingOrders = orders.filter(order =>
    ['pending_service', 'pending_confirm', 'exception'].includes(order.status)
  )
  const exceptionOrders = orders.filter(order => order.status === 'exception')

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Title level={3} style={{ margin: 0 }}>
          工作台
        </Title>
        <Text type="secondary">先处理异常，再看增长。所有关键项控制在一屏内。</Text>
      </div>

      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} lg={6}>
          <div className="mis-metric">
            <Statistic title="今日订单量" value={todayOrders.length} suffix="单" />
          </div>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <div className="mis-metric">
            <Statistic title="今日营收" value={todayRevenue} prefix="¥" precision={0} />
          </div>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <div className="mis-metric mis-metric-danger">
            <Badge count={pendingOrders.length} offset={[10, 0]}>
              <Statistic title="待处理工单" value={pendingOrders.length} />
            </Badge>
          </div>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <div className="mis-metric">
            <Statistic title="用户总数" value={users.length} suffix="人" />
          </div>
        </Col>
      </Row>

      <div className="mis-panel mis-alert-panel">
        <div className="mis-panel-header">
          <Space>
            <AlertOutlined />
            <Title level={5} style={{ margin: 0 }}>
              异常预警
            </Title>
            <Badge count={exceptionOrders.length} />
          </Space>
          <Button type="primary" danger onClick={() => navigate('/orders')}>
            进入订单处理
          </Button>
        </div>
        <Table
          rowKey="id"
          size="middle"
          pagination={false}
          dataSource={pendingOrders}
          columns={[
            {
              title: '订单号',
              dataIndex: 'id',
              width: 160
            },
            {
              title: '客户',
              dataIndex: 'customer',
              width: 100
            },
            {
              title: '服务',
              dataIndex: 'service'
            },
            {
              title: '预约时间',
              dataIndex: 'appointmentAt',
              width: 160
            },
            {
              title: '状态',
              dataIndex: 'status',
              width: 130,
              render: (value: OrderStatus) => (
                <Tag color={statusColor[value]}>{statusText[value]}</Tag>
              )
            },
            {
              title: '内部备注',
              dataIndex: 'internalNote'
            },
            {
              title: '动作',
              width: 120,
              render: () => (
                <Button icon={<FileTextOutlined />} onClick={() => navigate('/orders')}>
                  处理
                </Button>
              )
            }
          ]}
        />
      </div>

      <Row gutter={[12, 12]}>
        <Col xs={24} md={12}>
          <div className="mis-panel">
            <div className="mis-panel-header">
              <Title level={5} style={{ margin: 0 }}>
                今日运营顺序
              </Title>
            </div>
            <Space direction="vertical" size={10}>
              <Text>
                <Badge status="error" /> 先处理异常订单、退款、改派。
              </Text>
              <Text>
                <Badge status="warning" /> 再核销已完成服务，减少财务对账延迟。
              </Text>
              <Text>
                <Badge status="processing" /> 最后检查待服务订单是否已分配服务人员。
              </Text>
            </Space>
          </div>
        </Col>
        <Col xs={24} md={12}>
          <div className="mis-panel">
            <div className="mis-panel-header">
              <Title level={5} style={{ margin: 0 }}>
                系统升级记录
              </Title>
              <CheckCircleOutlined style={{ color: '#237804' }} />
            </div>
            <Text>
              当前已从游戏分析导航收敛为家政 MIS：工作台、订单、用户、内容商品、数据看板。
              初期目标是让运营人员快速处理异常，而不是做展示型大屏。
            </Text>
          </div>
        </Col>
      </Row>
    </Space>
  )
}

export default Dashboard
