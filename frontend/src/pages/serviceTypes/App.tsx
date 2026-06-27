import React, { useEffect, useState } from 'react'
import { Alert, Button, Empty, Form, Input, InputNumber, Modal, Popconfirm, Space, Table, Tag, Typography, message } from 'antd'
import { DeleteOutlined, EditOutlined, PlusOutlined, ReloadOutlined, StopOutlined, UnlockOutlined } from '@ant-design/icons'
import {
  ServiceTypeRecord,
  apiCreateServiceType,
  apiDeleteServiceType,
  apiDisableServiceType,
  apiEnableServiceType,
  apiListServiceTypes,
  apiUpdateServiceType
} from '@/api'

const { Title, Text } = Typography

const ServiceTypes: React.FC = () => {
  const [items, setItems] = useState<ServiceTypeRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [keyword, setKeyword] = useState('')
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<ServiceTypeRecord | null>(null)
  const [form] = Form.useForm<ServiceTypeRecord>()

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiListServiceTypes(keyword)
      setItems(res?.data?.list || [])
    } catch {
      setError('服务类型加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const openEditor = (record?: ServiceTypeRecord) => {
    setEditing(record || null)
    form.setFieldsValue(record || { name: '', description: '', sortOrder: 0, status: 'active' })
    setOpen(true)
  }

  const save = async () => {
    const values = await form.validateFields()
    try {
      if (editing) await apiUpdateServiceType(editing.id, values)
      else await apiCreateServiceType(values)
      message.success('服务类型已保存')
      setOpen(false)
      load()
    } catch {
      message.error('服务类型保存失败')
    }
  }

  const toggle = async (record: ServiceTypeRecord) => {
    try {
      if (record.status === 'active') await apiDisableServiceType(record.id)
      else await apiEnableServiceType(record.id)
      message.success('状态已更新')
      load()
    } catch {
      message.error('状态更新失败')
    }
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>服务类型管理</Title>
          <Text type="secondary">维护服务分类、排序和启停用状态。</Text>
        </div>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>刷新</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openEditor()}>新增类型</Button>
        </Space>
      </div>
      {error ? <Alert type="error" showIcon message={error} /> : null}
      <div className="mis-toolbar">
        <Input.Search allowClear placeholder="搜索类型名称" value={keyword} onChange={e => setKeyword(e.target.value)} onSearch={load} style={{ maxWidth: 320 }} />
      </div>
      <Table
        rowKey="id"
        loading={loading}
        dataSource={items}
        locale={{ emptyText: <Empty description="暂无服务类型" /> }}
        columns={[
          { title: '排序', dataIndex: 'sortOrder', width: 80 },
          { title: '类型名称', dataIndex: 'name', width: 180 },
          { title: '描述', dataIndex: 'description' },
          {
            title: '状态',
            dataIndex: 'status',
            width: 100,
            render: value => <Tag color={value === 'active' ? 'success' : 'default'}>{value === 'active' ? '启用' : '停用'}</Tag>
          },
          {
            title: '操作',
            width: 230,
            render: (_, record: ServiceTypeRecord) => (
              <Space>
                <Button icon={<EditOutlined />} onClick={() => openEditor(record)}>编辑</Button>
                <Button icon={record.status === 'active' ? <StopOutlined /> : <UnlockOutlined />} onClick={() => toggle(record)}>
                  {record.status === 'active' ? '停用' : '启用'}
                </Button>
                <Popconfirm title="确认删除该类型？" onConfirm={() => apiDeleteServiceType(record.id).then(load)}>
                  <Button danger icon={<DeleteOutlined />} />
                </Popconfirm>
              </Space>
            )
          }
        ]}
      />
      <Modal title={editing ? '编辑服务类型' : '新增服务类型'} open={open} onOk={save} onCancel={() => setOpen(false)} okText="保存" cancelText="取消">
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="类型名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={3} /></Form.Item>
          <Form.Item name="sortOrder" label="排序"><InputNumber min={0} precision={0} style={{ width: '100%' }} /></Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}

export default ServiceTypes
