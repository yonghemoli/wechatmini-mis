import React, { useEffect, useState } from 'react'
import {
  Alert,
  Button,
  DatePicker,
  Empty,
  Input,
  Modal,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message
} from 'antd'
import { CheckCircleOutlined, CloseCircleOutlined, DownloadOutlined, ReloadOutlined } from '@ant-design/icons'
import type { Dayjs } from 'dayjs'
import {
  OrderRecord,
  OrderStatus,
  apiConfirmOrder,
  apiExportOrders,
  apiListOrders,
  apiRefundOrder,
  apiUpdateOrderNote,
  orderStatusColor,
  orderStatusText
} from '@/api'

const { RangePicker } = DatePicker
const { Title, Text } = Typography

const Orders: React.FC = () => {
  const [orders, setOrders] = useState<OrderRecord[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [size, setSize] = useState(20)
  const [status, setStatus] = useState<OrderStatus | 'all'>('all')
  const [keyword, setKeyword] = useState('')
  const [range, setRange] = useState<[Dayjs | null, Dayjs | null] | null>(null)
  const [editingOrder, setEditingOrder] = useState<OrderRecord | null>(null)
  const [note, setNote] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const queryParams = {
    status,
    keyword,
    start: range?.[0]?.format('YYYY-MM-DD'),
    end: range?.[1]?.format('YYYY-MM-DD'),
    page,
    size
  }

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiListOrders(queryParams)
      setOrders(res?.data?.list || [])
      setTotal(res?.data?.total || 0)
    } catch {
      setError('订单列表加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [status, page, size])

  const search = () => {
    setPage(1)
    load()
  }

  const updateStatus = async (id: string, nextStatus: OrderStatus) => {
    try {
      if (nextStatus === 'completed') {
        await apiConfirmOrder(id)
        message.success('订单已核销')
      } else {
        await apiRefundOrder(id, '后台异常关闭/退款')
        message.success('订单已异常关闭/退款')
      }
      load()
    } catch {
      message.error('订单状态更新失败')
    }
  }

  const openNote = (order: OrderRecord) => {
    setEditingOrder(order)
    setNote(order.internalNote)
  }

  const saveNote = async () => {
    if (!editingOrder) return
    try {
      await apiUpdateOrderNote(editingOrder.id, note)
      setEditingOrder(null)
      message.success('内部备注已保存')
      load()
    } catch {
      message.error('内部备注保存失败')
    }
  }

  const exportOrders = () =>
    apiExportOrders({
      status,
      keyword,
      start: range?.[0]?.format('YYYY-MM-DD'),
      end: range?.[1]?.format('YYYY-MM-DD')
    })

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            订单管理
          </Title>
          <Text type="secondary">全量订单操作区，优先解决待服务、待核销和异常退款。</Text>
        </div>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>
            刷新
          </Button>
          <Button icon={<DownloadOutlined />} onClick={exportOrders}>
            导出 Excel
          </Button>
        </Space>
      </div>

      {error ? <Alert type="error" showIcon message={error} /> : null}

      <div className="mis-toolbar">
        <Select
          value={status}
          style={{ width: 160 }}
          onChange={value => {
            setStatus(value)
            setPage(1)
          }}
          options={[
            { label: '全部状态', value: 'all' },
            { label: '待服务', value: 'pending_service' },
            { label: '待核销', value: 'pending_confirm' },
            { label: '异常待处理', value: 'exception' },
            { label: '已完成', value: 'completed' },
            { label: '已退款', value: 'refunded' }
          ]}
        />
        <RangePicker value={range} onChange={value => setRange(value)} />
        <Input.Search
          allowClear
          placeholder="搜索订单号、客户、服务人员"
          style={{ maxWidth: 320 }}
          value={keyword}
          onChange={event => setKeyword(event.target.value)}
          onSearch={search}
        />
      </div>

      <Table
        rowKey="id"
        loading={loading}
        dataSource={orders}
        scroll={{ x: 1180 }}
        locale={{ emptyText: <Empty description="暂无订单" /> }}
        pagination={{
          current: page,
          pageSize: size,
          total,
          showSizeChanger: true,
          onChange: (nextPage, nextSize) => {
            setPage(nextPage)
            setSize(nextSize)
          }
        }}
        columns={[
          { title: '订单号', dataIndex: 'id', fixed: 'left', width: 160 },
          { title: '客户', dataIndex: 'customer', width: 100 },
          { title: '手机', dataIndex: 'phone', width: 120 },
          { title: '服务', dataIndex: 'service', width: 150 },
          { title: '金额', dataIndex: 'amount', width: 90, render: (value: number) => `¥${value}` },
          {
            title: '状态',
            dataIndex: 'status',
            width: 120,
            render: (value: OrderStatus) => <Tag color={orderStatusColor[value]}>{orderStatusText[value]}</Tag>
          },
          { title: '预约时间', dataIndex: 'appointmentAt', width: 150 },
          { title: '服务人员', dataIndex: 'staff', width: 100 },
          { title: '内部备注', dataIndex: 'internalNote', ellipsis: true },
          {
            title: '操作',
            fixed: 'right',
            width: 260,
            render: (_: unknown, record: OrderRecord) => (
              <Space>
                {['pending_service', 'pending_confirm'].includes(record.status) ? (
                  <Button icon={<CheckCircleOutlined />} type="primary" onClick={() => updateStatus(record.id, 'completed')}>
                    确认/核销
                  </Button>
                ) : null}
                {record.status !== 'refunded' && record.status !== 'completed' ? (
                  <Button icon={<CloseCircleOutlined />} danger onClick={() => updateStatus(record.id, 'refunded')}>
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
