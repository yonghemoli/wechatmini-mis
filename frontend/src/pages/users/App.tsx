import React, { useEffect, useState } from 'react'
import { Alert, Avatar, Button, Empty, Space, Table, Tag, Typography, message } from 'antd'
import { DownloadOutlined, ReloadOutlined, StopOutlined, UnlockOutlined } from '@ant-design/icons'
import { UserRecord, apiBanUser, apiExportUsers, apiListUsers, apiUnbanUser } from '@/api'

const { Title, Text } = Typography

const Users: React.FC = () => {
  const [users, setUsers] = useState<UserRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiListUsers()
      setUsers(res?.data?.list || [])
    } catch {
      setError('用户列表加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const toggleStatus = async (item: UserRecord) => {
    try {
      if (item.status === 'active') {
        await apiBanUser(item.id)
      } else {
        await apiUnbanUser(item.id)
      }
      message.success('用户状态已更新')
      load()
    } catch {
      message.error('用户状态更新失败')
    }
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            用户管理
          </Title>
          <Text type="secondary">极简用户列表，仅保留运营初期真正需要的封禁/解封。</Text>
        </div>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>
            刷新
          </Button>
          <Button icon={<DownloadOutlined />} onClick={apiExportUsers}>
            导出 Excel
          </Button>
        </Space>
      </div>

      {error ? <Alert type="error" showIcon message={error} /> : null}

      <Table
        rowKey="id"
        loading={loading}
        dataSource={users}
        locale={{ emptyText: <Empty description="暂无用户" /> }}
        columns={[
          {
            title: '用户',
            dataIndex: 'nickname',
            render: (_: string, record: UserRecord) => (
              <Space>
                <Avatar src={record.avatar} />
                <div>
                  <div>{record.nickname}</div>
                  <Text type="secondary">{record.id}</Text>
                </div>
              </Space>
            )
          },
          { title: '注册时间', dataIndex: 'createdAt', width: 180 },
          {
            title: '累计消费金额',
            dataIndex: 'totalSpent',
            width: 150,
            render: (value: number) => `¥${value}`
          },
          { title: '最后下单时间', dataIndex: 'lastOrderAt', width: 170 },
          {
            title: '状态',
            dataIndex: 'status',
            width: 100,
            render: (value: UserRecord['status']) => (
              <Tag color={value === 'active' ? 'success' : 'error'}>
                {value === 'active' ? '正常' : '已封禁'}
              </Tag>
            )
          },
          {
            title: '操作',
            width: 140,
            render: (_: unknown, record: UserRecord) => (
              <Button
                danger={record.status === 'active'}
                icon={record.status === 'active' ? <StopOutlined /> : <UnlockOutlined />}
                onClick={() => toggleStatus(record)}
              >
                {record.status === 'active' ? '手动封禁' : '解封'}
              </Button>
            )
          }
        ]}
      />
    </Space>
  )
}

export default Users
