import React, { useEffect, useCallback, useState } from 'react'
import { useDispatch } from 'react-redux'
import { Card, Typography, Button, Spin, message, Alert } from 'antd'
import { LoginOutlined } from '@ant-design/icons'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { apiSSOConfig, apiSSOSession } from '@/api'
import { setLoggedIn } from '@/redux'

const { Title, Text } = Typography

const Login: React.FC = () => {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [loading, setLoading] = useState(true)
  const forbidden = searchParams.get('forbidden') === '1'
  const forbiddenMsg = searchParams.get('msg') || '没有权限，请联系管理员'

  // 检查已有会话
  useEffect(() => {
    apiSSOSession()
      .then((res: any) => {
        if (res?.code === 0 && res.data?.user) {
          const u = res.data.user
          dispatch(
            setLoggedIn({
              id: u.id,
              username: u.username,
              email: u.email,
              avatar: u.avatar || '',
              isSuperAdmin: u.isSuperAdmin,
              roleId: u.roleId ?? null
            })
          )
          navigate('/', { replace: true })
          return
        }
        setLoading(false)
      })
      .catch(() => setLoading(false))
  }, [dispatch, navigate])

  const handleSSO = useCallback(async () => {
    try {
      const res = await apiSSOConfig()
      if (res?.data?.authURL) {
        const callbackURL = `${window.location.origin}/sso/callback`
        window.location.href = `${res.data.authURL}/api/sso/authorize?redirect_uri=${encodeURIComponent(callbackURL)}`
      } else {
        message.error('SSO 未配置，请联系管理员')
      }
    } catch {
      message.error('获取 SSO 配置失败')
    }
  }, [])

  if (loading) {
    return (
      <div
        style={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
        }}
      >
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
      }}
    >
      <Card
        style={{
          width: 400,
          borderRadius: 12,
          boxShadow: '0 8px 32px rgba(0,0,0,0.15)'
        }}
      >
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <Title level={2} style={{ marginBottom: 8 }}>
            修仙数据分析
          </Title>
          <Text type="secondary">游戏用户画像与数据洞察平台</Text>
        </div>

        {forbidden && (
          <Alert
            message={forbiddenMsg}
            type="warning"
            showIcon
            closable
            style={{ marginBottom: 16 }}
          />
        )}

        <Button
          type="primary"
          icon={<LoginOutlined />}
          size="large"
          block
          onClick={handleSSO}
        >
          SSO 统一登录
        </Button>

        <Text
          type="secondary"
          style={{
            display: 'block',
            textAlign: 'center',
            fontSize: 12,
            marginTop: 16
          }}
        >
          使用统一认证中心账号登录
        </Text>
      </Card>
    </div>
  )
}

export default Login
