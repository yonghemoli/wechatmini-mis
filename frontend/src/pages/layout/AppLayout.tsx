import React, { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Avatar,
  Button,
  Dropdown,
  Layout,
  Menu,
  Typography,
  message,
  theme
} from 'antd'
import {
  DashboardOutlined,
  CommentOutlined,
  HomeOutlined,
  FileTextOutlined,
  QuestionCircleOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  TeamOutlined,
  UserOutlined,
  SafetyOutlined,
  UnorderedListOutlined,
  AppstoreOutlined
} from '@ant-design/icons'
import { useDispatch, useSelector } from 'react-redux'
import { RootState, setLoggedOut } from '@/redux'
import { apiLogout } from '@/api'
import { onSessionExpired, resetSessionExpired } from '@/utils/authEvent'

const { Header, Sider, Content } = Layout

const menuItems = [
  {
    key: '/admin/dashboard',
    icon: <DashboardOutlined />,
    label: '工作台'
  },
  {
    key: '/admin/orders',
    icon: <FileTextOutlined />,
    label: '订单管理'
  },
  {
    key: '/admin/users',
    icon: <TeamOutlined />,
    label: '用户管理'
  },
  {
    key: '/admin/service-types',
    icon: <UnorderedListOutlined />,
    label: '服务类型管理'
  },
  {
    key: '/admin/services',
    icon: <AppstoreOutlined />,
    label: '服务管理'
  },
  {
    key: '/admin/accounts',
    icon: <SafetyOutlined />,
    label: '账户管理'
  },
  {
    key: '/admin/shops',
    icon: <HomeOutlined />,
    label: '店铺管理'
  },
  {
    key: '/admin/faqs',
    icon: <QuestionCircleOutlined />,
    label: '常见问题管理'
  },
  {
    key: '/admin/chat',
    icon: <CommentOutlined />,
    label: '客服在线'
  }
]

const routeKeys = menuItems.map(item => item.key)

const AppLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useDispatch()
  const user = useSelector((s: RootState) => s.auth.user)
  const { token } = theme.useToken()
  const isSuperAdmin = !!user?.isSuperAdmin

  const visibleMenuItems = menuItems.filter(item =>
    item.key === '/admin/accounts' ? isSuperAdmin : true
  )

  useEffect(() => {
    const unsubscribe = onSessionExpired(() => {
      message.warning('会话已过期，请重新登录')
      dispatch(setLoggedOut())
      resetSessionExpired()
      navigate('/admin/login', { replace: true })
    })
    return unsubscribe
  }, [dispatch, navigate])

  const handleLogout = async () => {
    try {
      await apiLogout()
    } catch {
      message.warning('服务端退出请求失败，已清理本地登录态')
    }
    dispatch(setLoggedOut())
    navigate('/admin/login')
  }

  const selectedKey =
    routeKeys.find(key => location.pathname.startsWith(key)) || '/admin/dashboard'

  return (
    <Layout
      style={{
        height: '100dvh',
        overflow: 'hidden',
        background: token.colorBgLayout
      }}
    >
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme="light"
        width={212}
        style={{
          height: '100%',
          overflow: 'hidden',
          borderRight: `1px solid ${token.colorBorderSecondary}`
        }}
      >
        <div
          style={{
            height: 56,
            display: 'flex',
            alignItems: 'center',
            padding: collapsed ? '0 16px' : '0 20px',
            borderBottom: `1px solid ${token.colorBorderSecondary}`
          }}
        >
          <Typography.Title level={5} style={{ margin: 0, whiteSpace: 'nowrap' }}>
            {collapsed ? 'MIS' : '家政 MIS'}
          </Typography.Title>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[selectedKey]}
          items={visibleMenuItems}
          onClick={({ key }) => navigate(key)}
          style={{ border: 'none', paddingTop: 8 }}
        />
      </Sider>
      <Layout style={{ minWidth: 0, height: '100%', overflow: 'hidden' }}>
        <Header
          style={{
            height: 56,
            flex: '0 0 56px',
            padding: '0 20px',
            background: token.colorBgContainer,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            borderBottom: `1px solid ${token.colorBorderSecondary}`
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <Dropdown
            menu={{
              items: [
                {
                  key: 'logout',
                  icon: <LogoutOutlined />,
                  label: '退出登录',
                  danger: true
                }
              ],
              onClick: ({ key }) => {
                if (key === 'logout') handleLogout()
              }
            }}
            placement="bottomRight"
          >
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                cursor: 'pointer',
                padding: '4px 8px',
                borderRadius: 6
              }}
            >
              <Avatar
                src={user?.avatar}
                icon={!user?.avatar ? <UserOutlined /> : undefined}
                size="small"
              />
              <span>{user?.username || '运营管理员'}</span>
            </div>
          </Dropdown>
        </Header>
        <Content
          style={{
            flex: 1,
            minHeight: 0,
            padding: 20,
            overflow: 'auto'
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default AppLayout
