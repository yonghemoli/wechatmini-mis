import React, { useEffect, useState } from 'react'
import { Alert, Button, Empty, Form, Input, InputNumber, Modal, Popconfirm, Select, Space, Switch, Table, Typography, message } from 'antd'
import { DeleteOutlined, DownloadOutlined, EditOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons'
import {
  ServiceRecord,
  ServiceTypeRecord,
  apiCreateService,
  apiDeleteService,
  apiExportServices,
  apiListServiceTypes,
  apiListServices,
  apiPublishService,
  apiUnpublishService,
  apiUpdateService
} from '@/api'

const { Title, Text } = Typography

const Services: React.FC = () => {
  const [items, setItems] = useState<ServiceRecord[]>([])
  const [types, setTypes] = useState<ServiceTypeRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [typeId, setTypeId] = useState<number | undefined>()
  const [keyword, setKeyword] = useState('')
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<ServiceRecord | null>(null)
  const [form] = Form.useForm<ServiceRecord>()

  const loadTypes = async () => {
    const res = await apiListServiceTypes()
    setTypes(res?.data?.list || [])
  }

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiListServices({ typeId, keyword })
      setItems(res?.data?.list || [])
    } catch {
      setError('服务列表加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTypes()
    load()
  }, [])

  const openEditor = (record?: ServiceRecord) => {
    setEditing(record || null)
    form.setFieldsValue(record || { typeId: types[0]?.id, name: '', image: '/me.png', price: 0, unit: '小时', description: '', visible: true, sortOrder: 0 })
    setOpen(true)
  }

  const save = async () => {
    const values = await form.validateFields()
    try {
      if (editing) await apiUpdateService(editing.id, values)
      else await apiCreateService(values)
      message.success('服务已保存')
      setOpen(false)
      load()
    } catch {
      message.error('服务保存失败')
    }
  }

  const toggle = async (record: ServiceRecord) => {
    try {
      if (record.visible) await apiUnpublishService(record.id)
      else await apiPublishService(record.id)
      message.success('上下架状态已更新')
      load()
    } catch {
      message.error('上下架状态更新失败')
    }
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>服务管理</Title>
          <Text type="secondary">维护前端可预约服务、价格、介绍和上下架状态。</Text>
        </div>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>刷新</Button>
          <Button icon={<DownloadOutlined />} onClick={apiExportServices}>导出 Excel</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openEditor()}>新增服务</Button>
        </Space>
      </div>
      {error ? <Alert type="error" showIcon message={error} /> : null}
      <div className="mis-toolbar">
        <Select allowClear placeholder="全部类型" value={typeId} style={{ width: 180 }} onChange={setTypeId} options={types.map(item => ({ label: item.name, value: item.id }))} />
        <Input.Search allowClear placeholder="搜索服务名称" value={keyword} onChange={e => setKeyword(e.target.value)} onSearch={load} style={{ maxWidth: 320 }} />
      </div>
      <Table
        rowKey="id"
        loading={loading}
        dataSource={items}
        locale={{ emptyText: <Empty description="暂无服务" /> }}
        columns={[
          { title: '排序', dataIndex: 'sortOrder', width: 80 },
          { title: '类型', dataIndex: 'typeName', width: 140 },
          { title: '服务名称', dataIndex: 'name', width: 160 },
          { title: '价格', width: 110, render: (_, record: ServiceRecord) => `¥${record.price}/${record.unit}` },
          { title: '介绍', dataIndex: 'description', ellipsis: true },
          { title: '上架', width: 100, render: (_, record: ServiceRecord) => <Switch checked={record.visible} checkedChildren="上架" unCheckedChildren="下架" onChange={() => toggle(record)} /> },
          {
            title: '操作',
            width: 150,
            render: (_, record: ServiceRecord) => (
              <Space>
                <Button icon={<EditOutlined />} onClick={() => openEditor(record)}>编辑</Button>
                <Popconfirm title="确认删除该服务？" onConfirm={() => apiDeleteService(record.id).then(load)}>
                  <Button danger icon={<DeleteOutlined />} />
                </Popconfirm>
              </Space>
            )
          }
        ]}
      />
      <Modal title={editing ? '编辑服务' : '新增服务'} open={open} onOk={save} onCancel={() => setOpen(false)} okText="保存" cancelText="取消">
        <Form form={form} layout="vertical">
          <Form.Item name="typeId" label="服务类型" rules={[{ required: true }]}><Select options={types.map(item => ({ label: item.name, value: item.id }))} /></Form.Item>
          <Form.Item name="name" label="服务名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="image" label="服务图"><Input /></Form.Item>
          <Form.Item name="price" label="价格" rules={[{ required: true }]}><InputNumber min={0} precision={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="unit" label="单位"><Input /></Form.Item>
          <Form.Item name="description" label="服务介绍"><Input.TextArea rows={4} /></Form.Item>
          <Form.Item name="sortOrder" label="排序"><InputNumber min={0} precision={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="visible" label="上下架" valuePropName="checked"><Switch checkedChildren="上架" unCheckedChildren="下架" /></Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}

export default Services
