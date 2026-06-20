import React, { useState } from 'react'
import { Avatar, Button, Space, Table, Tag, Typography, message } from 'antd'
import { DownloadOutlined, StopOutlined, UnlockOutlined } from '@ant-design/icons'
import type { UserRecord } from '../misData'
import { exportCsv, users as seedUsers } from '../misData'

const { Title, Text } = Typography

const Users: React.FC = () => {
  const [users, setUsers] = useState<UserRecord[]>(seedUsers)

  const toggleStatus = (id: string) => {
    setUsers(current =>
      current.map(user =>
        user.id === id
          ? { ...user, status: user.status === 'active' ? 'banned' : 'active' }
          : user
      )
    )
    message.success('用户状态已更新')
  }

  const exportUsers = () => {
    exportCsv(
      '家政用户列表.csv',
      users.map(user => ({
        用户ID: user.id,
        昵称: user.nickname,
        注册时间: user.registeredAt,
        累计消费: user.totalSpent,
        最后下单时间: user.lastOrderAt,
        状态: user.status === 'active' ? '正常' : '已封禁'
      }))
    )
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
        <Button icon={<DownloadOutlined />} onClick={exportUsers}>
          导出 Excel
        </Button>
      </div>

      <Table
        rowKey="id"
        dataSource={users}
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
          {
            title: '注册时间',
            dataIndex: 'registeredAt',
            width: 140
          },
          {
            title: '累计消费金额',
            dataIndex: 'totalSpent',
            width: 150,
            render: (value: number) => `¥${value}`
          },
          {
            title: '最后下单时间',
            dataIndex: 'lastOrderAt',
            width: 170
          },
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
                onClick={() => toggleStatus(record.id)}
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
