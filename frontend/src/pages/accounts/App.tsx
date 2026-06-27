import React, { useEffect, useMemo, useState } from 'react'
import {
  Button,
  Form,
  Input,
  Modal,
  Popconfirm,
  Space,
  Switch,
  Table,
  Tag,
  Typography,
  message
} from 'antd'
import {
  DownloadOutlined,
  LockOutlined,
  PlusOutlined,
  ReloadOutlined,
  StopOutlined,
  UserAddOutlined
} from '@ant-design/icons'
import {
  apiCreateAdminAccount,
  apiDisableAdminAccount,
  apiEnableAdminAccount,
  apiListAdminAccounts,
  apiResetAdminPassword,
  apiUpdateAdminAccount
} from '@/api'
import { exportCsv } from '@/utils/exportCsv'

const { Title, Text } = Typography

type AdminAccount = {
  id: number
  username: string
  name: string
  email: string
  roleId?: number | null
  isSuperAdmin: boolean
  status: string
  lastLoginAt?: string | null
}

const emptyAccount = {
  username: '',
  password: '',
  name: '',
  email: '',
  roleId: undefined,
  isSuperAdmin: false
}

const Accounts: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [items, setItems] = useState<AdminAccount[]>([])
  const [editing, setEditing] = useState<AdminAccount | null>(null)
  const [creating, setCreating] = useState(false)
  const [resetTarget, setResetTarget] = useState<AdminAccount | null>(null)
  const [form] = Form.useForm()
  const [passwordForm] = Form.useForm()

  const load = async () => {
    setLoading(true)
    try {
      const res = await apiListAdminAccounts()
      setItems(res?.data?.list || [])
    } catch {
      message.error('加载账户列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const exportAccounts = () => {
    exportCsv(
      '内部账户.csv',
      items.map(item => ({
        账号: item.username,
        姓名: item.name,
        邮箱: item.email,
        角色: item.isSuperAdmin ? '超级管理员' : item.roleId || '',
        状态: item.status,
        最近登录时间: item.lastLoginAt || ''
      }))
    )
  }

  const openCreate = () => {
    setEditing(null)
    form.setFieldsValue(emptyAccount)
    setCreating(true)
  }

  const openEdit = (item: AdminAccount) => {
    setEditing(item)
    form.setFieldsValue(item)
    setCreating(true)
  }

  const saveAccount = async () => {
    const values = await form.validateFields()
    try {
      if (editing) {
        await apiUpdateAdminAccount(editing.id, values)
        message.success('账户已更新')
      } else {
        await apiCreateAdminAccount(values)
        message.success('账户已创建')
      }
      setCreating(false)
      setEditing(null)
      load()
    } catch (err: any) {
      message.error(err?.message || '保存失败')
    }
  }

  const doResetPassword = async () => {
    const values = await passwordForm.validateFields()
    if (!resetTarget) return
    await apiResetAdminPassword(resetTarget.id, values.password)
    message.success('密码已重置')
    setResetTarget(null)
    passwordForm.resetFields()
  }

  const toggleStatus = async (item: AdminAccount) => {
    try {
      if (item.status === 'active') {
        await apiDisableAdminAccount(item.id)
      } else {
        await apiEnableAdminAccount(item.id)
      }
      message.success('状态已更新')
      load()
    } catch {
      message.error('状态更新失败')
    }
  }

  const columns = useMemo(
    () => [
      {
        title: '账号',
        dataIndex: 'username',
        width: 140
      },
      {
        title: '姓名',
        dataIndex: 'name',
        width: 120
      },
      {
        title: '邮箱',
        dataIndex: 'email',
        width: 180
      },
      {
        title: '角色',
        width: 110,
        render: (_: unknown, record: AdminAccount) =>
          record.isSuperAdmin ? <Tag color="red">超级管理员</Tag> : <Tag>普通管理员</Tag>
      },
      {
        title: '状态',
        dataIndex: 'status',
        width: 100,
        render: (status: string) => (
          <Tag color={status === 'active' ? 'success' : 'error'}>
            {status === 'active' ? '正常' : '已禁用'}
          </Tag>
        )
      },
      {
        title: '最近登录时间',
        dataIndex: 'lastLoginAt',
        width: 180
      },
      {
        title: '操作',
        width: 260,
        render: (_: unknown, record: AdminAccount) => (
          <Space>
            <Button icon={<UserAddOutlined />} onClick={() => openEdit(record)}>
              编辑
            </Button>
            <Button icon={<LockOutlined />} onClick={() => setResetTarget(record)}>
              重置密码
            </Button>
            <Popconfirm
              title={record.status === 'active' ? '确认禁用该账户？' : '确认启用该账户？'}
              onConfirm={() => toggleStatus(record)}
            >
              <Button danger={record.status === 'active'} icon={<StopOutlined />}>
                {record.status === 'active' ? '禁用' : '启用'}
              </Button>
            </Popconfirm>
          </Space>
        )
      }
    ],
    [form]
  )

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            账户管理
          </Title>
          <Text type="secondary">内部账号、角色、启停用和重置密码。</Text>
        </div>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>
            刷新
          </Button>
          <Button icon={<DownloadOutlined />} onClick={exportAccounts}>
            导出 Excel
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            新增账户
          </Button>
        </Space>
      </div>

      <Table loading={loading} rowKey="id" dataSource={items} columns={columns as any} />

      <Modal
        title={editing ? '编辑账户' : '新增账户'}
        open={creating}
        onOk={saveAccount}
        onCancel={() => setCreating(false)}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="username" label="账号" rules={[{ required: true }]}>
            <Input disabled={!!editing} />
          </Form.Item>
          {!editing && (
            <Form.Item name="password" label="初始密码" rules={[{ required: true }]}>
              <Input.Password />
            </Form.Item>
          )}
          <Form.Item name="name" label="姓名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ required: true, type: 'email' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="isSuperAdmin" label="超级管理员" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={`重置密码 - ${resetTarget?.username || ''}`}
        open={Boolean(resetTarget)}
        onOk={doResetPassword}
        onCancel={() => setResetTarget(null)}
        okText="确认重置"
        cancelText="取消"
      >
        <Form form={passwordForm} layout="vertical">
          <Form.Item name="password" label="新密码" rules={[{ required: true }]}>
            <Input.Password />
          </Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}

export default Accounts
