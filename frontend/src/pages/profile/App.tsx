import React from 'react'
import { Avatar, Button, Descriptions, Space, Typography, message } from 'antd'
import { LogoutOutlined, UserOutlined } from '@ant-design/icons'
import { useDispatch, useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { RootState, setLoggedOut } from '@/redux'
import { apiLogout } from '@/api'

const { Title, Text } = Typography

const Profile: React.FC = () => {
  const user = useSelector((s: RootState) => s.auth.user)
  const dispatch = useDispatch()
  const navigate = useNavigate()

  const handleLogout = async () => {
    try {
      await apiLogout()
    } catch {
      message.warning('服务端退出请求失败，已清理本地登录态')
    }
    dispatch(setLoggedOut())
    navigate('/admin/login', { replace: true })
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            个人中心
          </Title>
          <Text type="secondary">查看当前登录信息，并在这里执行退出登录。</Text>
        </div>
        <Button danger icon={<LogoutOutlined />} onClick={handleLogout}>
          退出登录
        </Button>
      </div>

      <div className="mis-panel">
        <Space align="start" size={16}>
          <Avatar
            size={64}
            src={user?.avatar}
            icon={!user?.avatar ? <UserOutlined /> : undefined}
          />
          <div>
            <Title level={4} style={{ margin: '0 0 4px' }}>
              {user?.name || user?.username || '运营管理员'}
            </Title>
            <Text type="secondary">{user?.username || '-'}</Text>
          </div>
        </Space>

        <div style={{ marginTop: 20 }}>
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="账号">{user?.username || '-'}</Descriptions.Item>
            <Descriptions.Item label="姓名">{user?.name || '-'}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{user?.email || '-'}</Descriptions.Item>
            <Descriptions.Item label="角色">
              {user?.isSuperAdmin ? '超级管理员' : '普通管理员'}
            </Descriptions.Item>
            <Descriptions.Item label="状态">{user?.status || '-'}</Descriptions.Item>
          </Descriptions>
        </div>
      </div>
    </Space>
  )
}

export default Profile
