import React, { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Avatar,
  Button,
  Layout,
  Menu,
  Space,
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
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  TeamOutlined,
  UserOutlined,
  SafetyOutlined,
  UnorderedListOutlined,
  AppstoreOutlined,
  RightOutlined
} from '@ant-design/icons'
import { useDispatch, useSelector } from 'react-redux'
import { RootState, setLoggedOut } from '@/redux'
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
            {collapsed ? 'MIS' : '永和茉莉'}
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
          <Button
            type="text"
            onClick={() => navigate('/admin/profile')}
            style={{
              height: 'auto',
              padding: '4px 8px'
            }}
          >
            <Space size={8}>
              <Avatar
                src={user?.avatar}
                icon={!user?.avatar ? <UserOutlined /> : undefined}
                size="small"
              />
              <span>{user?.username || '运营管理员'}</span>
              <RightOutlined style={{ fontSize: 12 }} />
            </Space>
          </Button>
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
