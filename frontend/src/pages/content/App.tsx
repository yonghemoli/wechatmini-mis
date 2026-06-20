import React, { useState } from 'react'
import {
  Button,
  Form,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Space,
  Switch,
  Table,
  Typography,
  message
} from 'antd'
import { DeleteOutlined, DownloadOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons'
import type { ContentRecord } from '../misData'
import { contents as seedContents, exportCsv } from '../misData'

const { Title, Text } = Typography

const Content: React.FC = () => {
  const [items, setItems] = useState<ContentRecord[]>(seedContents)
  const [editing, setEditing] = useState<ContentRecord | null>(null)
  const [open, setOpen] = useState(false)
  const [form] = Form.useForm<ContentRecord>()

  const openEditor = (record?: ContentRecord) => {
    setEditing(record || null)
    form.setFieldsValue(
      record || {
        title: '',
        image: '/me.png',
        price: 0,
        description: '',
        visible: true
      }
    )
    setOpen(true)
  }

  const saveItem = async () => {
    const values = await form.validateFields()
    const now = '2026-06-20 12:00'
    if (editing) {
      setItems(current =>
        current.map(item =>
          item.id === editing.id ? { ...editing, ...values, updatedAt: now } : item
        )
      )
    } else {
      setItems(current => [
        {
          ...values,
          id: `SVC${String(current.length + 1).padStart(3, '0')}`,
          updatedAt: now
        },
        ...current
      ])
    }
    setOpen(false)
    message.success('内容已保存')
  }

  const toggleVisible = (id: string) => {
    setItems(current =>
      current.map(item =>
        item.id === id ? { ...item, visible: !item.visible, updatedAt: '2026-06-20 12:00' } : item
      )
    )
  }

  const exportContent = () => {
    exportCsv(
      '家政内容商品.csv',
      items.map(item => ({
        编号: item.id,
        标题: item.title,
        图片: item.image,
        价格: item.price,
        服务介绍: item.description,
        前端可见: item.visible ? '是' : '否',
        更新时间: item.updatedAt
      }))
    )
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            内容/商品发布
          </Title>
          <Text type="secondary">管理小程序首页展示的服务图、价格和介绍文案。</Text>
        </div>
        <Space>
          <Button icon={<DownloadOutlined />} onClick={exportContent}>
            导出 Excel
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openEditor()}>
            新增服务
          </Button>
        </Space>
      </div>

      <Table
        rowKey="id"
        dataSource={items}
        columns={[
          {
            title: '编号',
            dataIndex: 'id',
            width: 100
          },
          {
            title: '服务',
            dataIndex: 'title',
            width: 140
          },
          {
            title: '价格',
            dataIndex: 'price',
            width: 100,
            render: (value: number) => `¥${value}/小时`
          },
          {
            title: '介绍文案',
            dataIndex: 'description'
          },
          {
            title: '前端可见',
            dataIndex: 'visible',
            width: 110,
            render: (_: boolean, record: ContentRecord) => (
              <Switch
                checked={record.visible}
                checkedChildren="上架"
                unCheckedChildren="下架"
                onChange={() => toggleVisible(record.id)}
              />
            )
          },
          {
            title: '更新时间',
            dataIndex: 'updatedAt',
            width: 160
          },
          {
            title: '操作',
            width: 150,
            render: (_: unknown, record: ContentRecord) => (
              <Space>
                <Button icon={<EditOutlined />} onClick={() => openEditor(record)}>
                  编辑
                </Button>
                <Popconfirm
                  title="确认删除该服务？"
                  onConfirm={() =>
                    setItems(current => current.filter(item => item.id !== record.id))
                  }
                >
                  <Button danger icon={<DeleteOutlined />} />
                </Popconfirm>
              </Space>
            )
          }
        ]}
      />

      <Modal
        title={editing ? '编辑服务' : '新增服务'}
        open={open}
        onOk={saveItem}
        onCancel={() => setOpen(false)}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="title" label="服务标题" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="image" label="商品图地址" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="price" label="价格" rules={[{ required: true }]}>
            <InputNumber min={0} precision={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="description" label="服务介绍文案" rules={[{ required: true }]}>
            <Input.TextArea rows={4} />
          </Form.Item>
          <Form.Item name="visible" label="上下架" valuePropName="checked">
            <Switch checkedChildren="上架" unCheckedChildren="下架" />
          </Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}

export default Content
