import React, { useEffect, useState } from 'react'
import { useDispatch } from 'react-redux'
import { Alert, Button, Card, Form, Input, Spin, Typography, message } from 'antd'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { apiLogin, apiSession } from '@/api'
import { setLoggedIn } from '@/redux'

const { Title, Text } = Typography

const Login: React.FC = () => {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const forbidden = searchParams.get('forbidden') === '1'
  const forbiddenMsg = searchParams.get('msg') || '没有权限，请联系管理员'

  useEffect(() => {
    apiSession()
      .then((res: any) => {
        if (res?.code === 0 && res.data?.user) {
          const u = res.data.user
          dispatch(
            setLoggedIn({
              id: u.id,
              username: u.username,
              name: u.name,
              email: u.email,
              avatar: u.avatar || '',
              status: u.status,
              isSuperAdmin: u.isSuperAdmin,
              roleId: u.roleId ?? null
            })
          )
          navigate('/admin', { replace: true })
          return
        }
        setLoading(false)
      })
      .catch(() => setLoading(false))
  }, [dispatch, navigate])

  const handleLogin = async (values: { username: string; password: string }) => {
    setSubmitting(true)
    try {
      const res = await apiLogin(values)
      if (res?.code === 0 && res.data?.user) {
        const u = res.data.user
        dispatch(
          setLoggedIn({
            id: u.id,
            username: u.username,
            name: u.name,
            email: u.email,
            avatar: u.avatar || '',
            status: u.status,
            isSuperAdmin: u.isSuperAdmin,
            roleId: u.roleId ?? null
          })
        )
        message.success('登录成功')
        navigate('/admin', { replace: true })
        return
      }
      message.error(res?.message || '登录失败')
    } catch (err: any) {
      message.error(err?.message || '登录失败')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="mis-login-page">
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div className="mis-login-page">
      <Card className="mis-login-card">
        <div style={{ textAlign: 'center', marginBottom: 28 }}>
          <Title level={2} style={{ marginBottom: 8 }}>
            永和茉莉
          </Title>
          <Text type="secondary">内部运营管理系统</Text>
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

        <Form layout="vertical" onFinish={handleLogin} initialValues={{ username: 'admin' }}>
          <Form.Item name="username" label="账号" rules={[{ required: true }]}>
            <Input prefix={<UserOutlined />} placeholder="请输入内部账号" size="large" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true }]}>
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请输入密码"
              size="large"
            />
          </Form.Item>
          <Button type="primary" htmlType="submit" size="large" block loading={submitting}>
            登录
          </Button>
        </Form>

        <Text
          type="secondary"
          style={{
            display: 'block',
            textAlign: 'center',
            fontSize: 12,
            marginTop: 16
          }}
        >
          默认开发账号来自 sql/init-seed.sql：admin / admin123
        </Text>
      </Card>
    </div>
  )
}

export default Login
