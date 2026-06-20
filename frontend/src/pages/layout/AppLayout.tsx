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
  FileTextOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ShopOutlined,
  TeamOutlined,
  UserOutlined,
  LineChartOutlined
} from '@ant-design/icons'
import { useDispatch, useSelector } from 'react-redux'
import { RootState, setLoggedOut } from '@/redux'
import { apiLogout, apiSSOConfig } from '@/api'
import { onSessionExpired, resetSessionExpired } from '@/utils/authEvent'

const { Header, Sider, Content } = Layout

const menuItems = [
  {
    key: '/dashboard',
    icon: <DashboardOutlined />,
    label: '工作台'
  },
  {
    key: '/orders',
    icon: <FileTextOutlined />,
    label: '订单管理'
  },
  {
    key: '/users',
    icon: <TeamOutlined />,
    label: '用户管理'
  },
  {
    key: '/content',
    icon: <ShopOutlined />,
    label: '内容商品'
  },
  {
    key: '/reports',
    icon: <LineChartOutlined />,
    label: '数据看板'
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

  useEffect(() => {
    const unsubscribe = onSessionExpired(() => {
      message.warning('会话已过期，请重新登录')
      dispatch(setLoggedOut())
      resetSessionExpired()
      navigate('/login', { replace: true })
    })
    return unsubscribe
  }, [dispatch, navigate])

  const handleLogout = async () => {
    await apiLogout()
    dispatch(setLoggedOut())
    try {
      const res = await apiSSOConfig()
      const authURL = res?.data?.authURL || res?.authURL
      if (authURL) {
        window.location.href = `${authURL}/api/sso/logout?redirect_uri=${encodeURIComponent(window.location.origin)}`
        return
      }
    } catch {
      /* fallback */
    }
    navigate('/login')
  }

  const selectedKey =
    routeKeys.find(key => location.pathname.startsWith(key)) || '/dashboard'

  return (
    <Layout style={{ minHeight: '100vh', background: token.colorBgLayout }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme="light"
        width={212}
        style={{
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
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{ border: 'none', paddingTop: 8 }}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            height: 56,
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
