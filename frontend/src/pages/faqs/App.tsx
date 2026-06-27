import React, { useEffect, useState } from 'react'
import { Alert, Button, Empty, Form, Input, InputNumber, Modal, Popconfirm, Select, Space, Switch, Table, Typography, message } from 'antd'
import { DeleteOutlined, EditOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons'
import { FAQRecord, apiCreateFAQ, apiDeleteFAQ, apiListFAQs, apiPublishFAQ, apiUnpublishFAQ, apiUpdateFAQ } from '@/api'

const { Title, Text } = Typography

const FAQs: React.FC = () => {
  const [items, setItems] = useState<FAQRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [category, setCategory] = useState('')
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<FAQRecord | null>(null)
  const [form] = Form.useForm<FAQRecord>()

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiListFAQs(category)
      setItems(res?.data?.list || [])
    } catch {
      setError('常见问题加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [category])

  const openEditor = (record?: FAQRecord) => {
    setEditing(record || null)
    form.setFieldsValue(record || { question: '', answer: '', category: '下单', sortOrder: 0, visible: true })
    setOpen(true)
  }

  const save = async () => {
    const values = await form.validateFields()
    try {
      if (editing) await apiUpdateFAQ(editing.id, values)
      else await apiCreateFAQ(values)
      message.success('常见问题已保存')
      setOpen(false)
      load()
    } catch {
      message.error('常见问题保存失败')
    }
  }

  const toggle = async (record: FAQRecord) => {
    try {
      if (record.visible) await apiUnpublishFAQ(record.id)
      else await apiPublishFAQ(record.id)
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
          <Title level={3} style={{ margin: 0 }}>常见问题管理</Title>
          <Text type="secondary">维护小程序和客服常用问答。</Text>
        </div>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>刷新</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openEditor()}>新增问题</Button>
        </Space>
      </div>
      {error ? <Alert type="error" showIcon message={error} /> : null}
      <div className="mis-toolbar">
        <Select allowClear placeholder="全部分类" value={category || undefined} onChange={value => setCategory(value || '')} style={{ width: 180 }} options={['下单', '支付', '服务', '售后'].map(v => ({ label: v, value: v }))} />
      </div>
      <Table
        rowKey="id"
        loading={loading}
        dataSource={items}
        locale={{ emptyText: <Empty description="暂无常见问题" /> }}
        columns={[
          { title: '排序', dataIndex: 'sortOrder', width: 80 },
          { title: '分类', dataIndex: 'category', width: 100 },
          { title: '问题', dataIndex: 'question', width: 240 },
          { title: '答案', dataIndex: 'answer', ellipsis: true },
          { title: '展示', width: 100, render: (_, record: FAQRecord) => <Switch checked={record.visible} checkedChildren="上架" unCheckedChildren="下架" onChange={() => toggle(record)} /> },
          {
            title: '操作',
            width: 150,
            render: (_, record: FAQRecord) => (
              <Space>
                <Button icon={<EditOutlined />} onClick={() => openEditor(record)}>编辑</Button>
                <Popconfirm title="确认删除该问题？" onConfirm={() => apiDeleteFAQ(record.id).then(load)}>
                  <Button danger icon={<DeleteOutlined />} />
                </Popconfirm>
              </Space>
            )
          }
        ]}
      />
      <Modal title={editing ? '编辑常见问题' : '新增常见问题'} open={open} onOk={save} onCancel={() => setOpen(false)} okText="保存" cancelText="取消">
        <Form form={form} layout="vertical">
          <Form.Item name="category" label="分类" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="question" label="问题" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="answer" label="答案" rules={[{ required: true }]}><Input.TextArea rows={5} /></Form.Item>
          <Form.Item name="sortOrder" label="排序"><InputNumber min={0} precision={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="visible" label="上下架" valuePropName="checked"><Switch checkedChildren="上架" unCheckedChildren="下架" /></Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}

export default FAQs
