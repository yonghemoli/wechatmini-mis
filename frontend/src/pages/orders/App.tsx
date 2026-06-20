import React, { useMemo, useState } from 'react'
import {
  Button,
  DatePicker,
  Input,
  Modal,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message
} from 'antd'
import { CheckCircleOutlined, CloseCircleOutlined, DownloadOutlined } from '@ant-design/icons'
import type { OrderRecord, OrderStatus } from '../misData'
import { exportCsv, orders as seedOrders, statusColor, statusText } from '../misData'

const { RangePicker } = DatePicker
const { Title, Text } = Typography

const Orders: React.FC = () => {
  const [orders, setOrders] = useState<OrderRecord[]>(seedOrders)
  const [status, setStatus] = useState<OrderStatus | 'all'>('all')
  const [keyword, setKeyword] = useState('')
  const [editingOrder, setEditingOrder] = useState<OrderRecord | null>(null)
  const [note, setNote] = useState('')

  const filteredOrders = useMemo(
    () =>
      orders.filter(order => {
        const matchStatus = status === 'all' || order.status === status
        const matchKeyword =
          !keyword ||
          [order.id, order.customer, order.phone, order.service, order.staff].some(value =>
            value.toLowerCase().includes(keyword.toLowerCase())
          )
        return matchStatus && matchKeyword
      }),
    [orders, status, keyword]
  )

  const updateStatus = (id: string, nextStatus: OrderStatus) => {
    setOrders(current =>
      current.map(order => (order.id === id ? { ...order, status: nextStatus } : order))
    )
    message.success(nextStatus === 'completed' ? '订单已核销' : '订单已异常关闭/退款')
  }

  const openNote = (order: OrderRecord) => {
    setEditingOrder(order)
    setNote(order.internalNote)
  }

  const saveNote = () => {
    if (!editingOrder) return
    setOrders(current =>
      current.map(order =>
        order.id === editingOrder.id ? { ...order, internalNote: note } : order
      )
    )
    setEditingOrder(null)
    message.success('内部备注已保存')
  }

  const exportOrders = () => {
    exportCsv(
      '家政订单列表.csv',
      filteredOrders.map(order => ({
        订单号: order.id,
        客户: order.customer,
        手机: order.phone,
        服务: order.service,
        金额: order.amount,
        状态: statusText[order.status],
        来源: order.source,
        预约时间: order.appointmentAt,
        创建时间: order.createdAt,
        服务人员: order.staff,
        内部备注: order.internalNote
      }))
    )
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            订单管理
          </Title>
          <Text type="secondary">全量订单操作区，优先解决待服务、待核销和异常退款。</Text>
        </div>
        <Button icon={<DownloadOutlined />} onClick={exportOrders}>
          导出 Excel
        </Button>
      </div>

      <div className="mis-toolbar">
        <Select
          value={status}
          style={{ width: 160 }}
          onChange={setStatus}
          options={[
            { label: '全部状态', value: 'all' },
            { label: '待服务', value: 'pending_service' },
            { label: '待核销', value: 'pending_confirm' },
            { label: '异常待处理', value: 'exception' },
            { label: '已完成', value: 'completed' },
            { label: '已退款', value: 'refunded' }
          ]}
        />
        <RangePicker />
        <Input.Search
          allowClear
          placeholder="搜索订单号、客户、服务人员"
          style={{ maxWidth: 320 }}
          value={keyword}
          onChange={event => setKeyword(event.target.value)}
        />
      </div>

      <Table
        rowKey="id"
        dataSource={filteredOrders}
        scroll={{ x: 1180 }}
        columns={[
          {
            title: '订单号',
            dataIndex: 'id',
            fixed: 'left',
            width: 160
          },
          {
            title: '客户',
            dataIndex: 'customer',
            width: 100
          },
          {
            title: '手机',
            dataIndex: 'phone',
            width: 120
          },
          {
            title: '服务',
            dataIndex: 'service',
            width: 150
          },
          {
            title: '金额',
            dataIndex: 'amount',
            width: 90,
            render: (value: number) => `¥${value}`
          },
          {
            title: '状态',
            dataIndex: 'status',
            width: 120,
            render: (value: OrderStatus) => (
              <Tag color={statusColor[value]}>{statusText[value]}</Tag>
            )
          },
          {
            title: '预约时间',
            dataIndex: 'appointmentAt',
            width: 150
          },
          {
            title: '服务人员',
            dataIndex: 'staff',
            width: 100
          },
          {
            title: '内部备注',
            dataIndex: 'internalNote',
            ellipsis: true
          },
          {
            title: '操作',
            fixed: 'right',
            width: 260,
            render: (_: unknown, record: OrderRecord) => (
              <Space>
                {['pending_service', 'pending_confirm'].includes(record.status) ? (
                  <Button
                    icon={<CheckCircleOutlined />}
                    type="primary"
                    onClick={() => updateStatus(record.id, 'completed')}
                  >
                    确认/核销
                  </Button>
                ) : null}
                {record.status !== 'refunded' && record.status !== 'completed' ? (
                  <Button
                    icon={<CloseCircleOutlined />}
                    danger
                    onClick={() => updateStatus(record.id, 'refunded')}
                  >
                    异常关闭/退款
                  </Button>
                ) : null}
                <Button onClick={() => openNote(record)}>备注</Button>
              </Space>
            )
          }
        ]}
      />

      <Modal
        title={`内部备注 - ${editingOrder?.id || ''}`}
        open={Boolean(editingOrder)}
        onOk={saveNote}
        onCancel={() => setEditingOrder(null)}
        okText="保存"
        cancelText="取消"
      >
        <Input.TextArea
          value={note}
          onChange={event => setNote(event.target.value)}
          rows={5}
          placeholder="记录仅内部可见，方便客服和运营交接"
        />
      </Modal>
    </Space>
  )
}

export default Orders
