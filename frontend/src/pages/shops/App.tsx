import React, { useEffect, useState } from 'react'
import { Alert, Button, Empty, Form, Input, Modal, Popconfirm, Space, Table, Tag, Typography, message } from 'antd'
import { DeleteOutlined, EditOutlined, PlusOutlined, ReloadOutlined, ShopOutlined } from '@ant-design/icons'
import { ShopRecord, apiCloseShop, apiCreateShop, apiDeleteShop, apiListShops, apiOpenShop, apiUpdateShop } from '@/api'

const { Title, Text } = Typography

const Shops: React.FC = () => {
  const [items, setItems] = useState<ShopRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<ShopRecord | null>(null)
  const [form] = Form.useForm<ShopRecord>()

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiListShops()
      setItems(res?.data?.list || [])
    } catch {
      setError('店铺列表加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const openEditor = (record?: ShopRecord) => {
    setEditing(record || null)
    form.setFieldsValue(record || { name: '', contactName: '', phone: '', address: '', businessHours: '09:00-18:00', status: 'open', remark: '' })
    setOpen(true)
  }

  const save = async () => {
    const values = await form.validateFields()
    try {
      if (editing) await apiUpdateShop(editing.id, values)
      else await apiCreateShop(values)
      message.success('店铺已保存')
      setOpen(false)
      load()
    } catch {
      message.error('店铺保存失败')
    }
  }

  const toggle = async (record: ShopRecord) => {
    try {
      if (record.status === 'open') await apiCloseShop(record.id)
      else await apiOpenShop(record.id)
      message.success('营业状态已更新')
      load()
    } catch {
      message.error('营业状态更新失败')
    }
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>店铺管理</Title>
          <Text type="secondary">维护门店资料、营业时间和营业状态。</Text>
        </div>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>刷新</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openEditor()}>新增店铺</Button>
        </Space>
      </div>
      {error ? <Alert type="error" showIcon message={error} /> : null}
      <Table
        rowKey="id"
        loading={loading}
        dataSource={items}
        locale={{ emptyText: <Empty description="暂无店铺" /> }}
        columns={[
          { title: '店铺', dataIndex: 'name', width: 180 },
          { title: '联系人', dataIndex: 'contactName', width: 100 },
          { title: '电话', dataIndex: 'phone', width: 140 },
          { title: '地址', dataIndex: 'address' },
          { title: '营业时间', dataIndex: 'businessHours', width: 140 },
          { title: '状态', dataIndex: 'status', width: 100, render: value => <Tag color={value === 'open' ? 'success' : 'default'}>{value === 'open' ? '营业中' : '已停业'}</Tag> },
          {
            title: '操作',
            width: 240,
            render: (_, record: ShopRecord) => (
              <Space>
                <Button icon={<EditOutlined />} onClick={() => openEditor(record)}>编辑</Button>
                <Button icon={<ShopOutlined />} onClick={() => toggle(record)}>{record.status === 'open' ? '停业' : '开业'}</Button>
                <Popconfirm title="确认删除该店铺？" onConfirm={() => apiDeleteShop(record.id).then(load)}>
                  <Button danger icon={<DeleteOutlined />} />
                </Popconfirm>
              </Space>
            )
          }
        ]}
      />
      <Modal title={editing ? '编辑店铺' : '新增店铺'} open={open} onOk={save} onCancel={() => setOpen(false)} okText="保存" cancelText="取消">
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="店铺名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="contactName" label="联系人"><Input /></Form.Item>
          <Form.Item name="phone" label="电话"><Input /></Form.Item>
          <Form.Item name="address" label="地址"><Input /></Form.Item>
          <Form.Item name="businessHours" label="营业时间"><Input /></Form.Item>
          <Form.Item name="remark" label="备注"><Input.TextArea rows={3} /></Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}

export default Shops
