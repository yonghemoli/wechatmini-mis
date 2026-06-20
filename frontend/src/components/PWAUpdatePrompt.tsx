import { useRegisterSW } from 'virtual:pwa-register/react'
import { Button, notification } from 'antd'
import { useEffect, useRef } from 'react'
import { SyncOutlined } from '@ant-design/icons'

const PWAUpdatePrompt: React.FC = () => {
  const notifiedRef = useRef(false)
  const [api, contextHolder] = notification.useNotification()

  const {
    needRefresh: [needRefresh, setNeedRefresh],
    updateServiceWorker
  } = useRegisterSW({
    onRegisteredSW(_url, registration) {
      // 每 60 秒检查一次更新
      if (registration) {
        setInterval(() => registration.update(), 60 * 1000)
      }
    }
  })

  useEffect(() => {
    if (needRefresh && !notifiedRef.current) {
      notifiedRef.current = true
      api.info({
        key: 'pwa-update',
        message: '发现新版本',
        description: '应用有新版本可用，点击更新以获取最新功能。',
        icon: <SyncOutlined spin style={{ color: '#1890ff' }} />,
        btn: (
          <Button
            type="primary"
            size="small"
            onClick={() => updateServiceWorker(true)}
          >
            立即更新
          </Button>
        ),
        duration: 0,
        onClose: () => {
          setNeedRefresh(false)
          notifiedRef.current = false
        }
      })
    }
  }, [needRefresh, api, setNeedRefresh, updateServiceWorker])

  return contextHolder
}

export default PWAUpdatePrompt
