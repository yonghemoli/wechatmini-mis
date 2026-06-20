import { useEffect, useRef } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { Spin, message } from 'antd'
import { apiSSOLogin } from '@/api'
import { setLoggedIn } from '@/redux'

export default function SSOCallbackPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const processed = useRef(false)

  useEffect(() => {
    if (processed.current) return
    processed.current = true

    const code = searchParams.get('code')
    if (!code) {
      message.error('授权码缺失')
      navigate('/login', { replace: true })
      return
    }

    apiSSOLogin({ code })
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
          message.success('登录成功')
          navigate('/', { replace: true })
        } else {
          message.error(res?.message || 'SSO 登录失败')
          navigate('/login', { replace: true })
        }
      })
      .catch((err: any) => {
        const status = err?.status || err?.response?.status
        if (status === 403) {
          const msg = err?.message || err?.msg || '没有权限'
          navigate(`/login?forbidden=1&msg=${encodeURIComponent(msg)}`, {
            replace: true
          })
          return
        }
        message.error(err?.msg || err?.message || 'SSO 登录失败')
        navigate('/login', { replace: true })
      })
  }, [searchParams, navigate, dispatch])

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
      <Spin size="large" tip="正在验证登录..." />
    </div>
  )
}
