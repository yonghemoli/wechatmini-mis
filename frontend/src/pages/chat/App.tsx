import React, { useEffect, useRef, useState } from 'react'
import { Alert, Avatar, Button, Empty, Input, List, Space, Tag, Typography, message, theme } from 'antd'
import { CheckOutlined, CloseCircleOutlined, ReloadOutlined, SendOutlined, UserOutlined } from '@ant-design/icons'
import {
  ChatMessageRecord,
  ChatWsEvent,
  ChatSessionRecord,
  apiCloseChatSession,
  apiListChatMessages,
  apiListChatSessions,
  apiReadChatSession,
  apiSendChatMessage,
  getWsBaseUrl
} from '@/api'

const { Title, Text } = Typography

const Chat: React.FC = () => {
  const { token } = theme.useToken()
  const [sessions, setSessions] = useState<ChatSessionRecord[]>([])
  const [messages, setMessages] = useState<ChatMessageRecord[]>([])
  const [active, setActive] = useState<ChatSessionRecord | null>(null)
  const [text, setText] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const wsRef = useRef<WebSocket | null>(null)
  const [wsStatus, setWsStatus] = useState<'connecting' | 'open' | 'closed'>('closed')

  const activeSessionId = active?.id || ''

  const loadSessions = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await apiListChatSessions()
      const list = res?.data?.list || []
      setSessions(list)
      if (!active && list.length) setActive(list[0])
    } catch {
      setError('客服会话加载失败')
    } finally {
      setLoading(false)
    }
  }

  const loadMessages = async (id: string) => {
    const res = await apiListChatMessages(id)
    setMessages(res?.data?.list || [])
    await apiReadChatSession(id)
  }

  const mergeSession = (next: ChatSessionRecord) => {
    setSessions(prev => {
      const exists = prev.some(item => item.id === next.id)
      const list = exists ? prev.map(item => (item.id === next.id ? { ...item, ...next } : item)) : [next, ...prev]
      return list.sort((a, b) => (b.updatedAt || '').localeCompare(a.updatedAt || ''))
    })
    setActive(prev => (prev?.id === next.id ? { ...prev, ...next } : prev))
  }

  const appendMessage = (next: ChatMessageRecord) => {
    if (next.sessionId !== activeSessionId) return
    setMessages(prev => {
      if (prev.some(item => item.id === next.id)) return prev
      return [...prev, next]
    })
  }

  const connectWs = () => {
    if (wsRef.current && (wsRef.current.readyState === WebSocket.OPEN || wsRef.current.readyState === WebSocket.CONNECTING)) {
      return
    }
    const ws = new WebSocket(`${getWsBaseUrl()}/ws/chat`)
    wsRef.current = ws
    setWsStatus('connecting')
    ws.onopen = () => setWsStatus('open')
    ws.onclose = () => setWsStatus('closed')
    ws.onerror = () => setWsStatus('closed')
    ws.onmessage = ev => {
      let data: ChatWsEvent | null = null
      try {
        data = JSON.parse(ev.data)
      } catch {
        return
      }
      if (!data) return
      if (data.type === 'message') {
        appendMessage(data.message)
        if (activeSessionId === data.sessionId) {
          loadMessages(data.sessionId)
        }
      } else if (data.type === 'session') {
        mergeSession(data.session)
      }
    }
  }

  useEffect(() => {
    loadSessions()
  }, [])

  useEffect(() => {
    if (active) loadMessages(active.id)
  }, [active?.id])

  useEffect(() => {
    connectWs()
    return () => {
      wsRef.current?.close()
      wsRef.current = null
    }
  }, [])

  useEffect(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN && activeSessionId) {
      wsRef.current.send(JSON.stringify({ type: 'ping', sessionId: activeSessionId }))
    }
  }, [activeSessionId])

  const send = async () => {
    if (!active || !text.trim()) return
    try {
      const payload = { content: text.trim(), msgType: 'text' }
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'message', sessionId: active.id, ...payload }))
      } else {
        await apiSendChatMessage(active.id, text.trim())
        await loadMessages(active.id)
        await loadSessions()
      }
      setText('')
    } catch {
      message.error('消息发送失败')
    }
  }

  const closeSession = async () => {
    if (!active) return
    await apiCloseChatSession(active.id)
    message.success('会话已关闭')
    await loadSessions()
  }

  return (
    <div style={{ height: 'calc(100vh - 96px)', display: 'flex', flexDirection: 'column', gap: 12 }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>客服在线</Title>
          <Text type="secondary">独立聊天工作台，处理用户在线咨询。</Text>
        </div>
        <Button icon={<ReloadOutlined />} onClick={loadSessions}>刷新</Button>
      </div>
      {error ? <Alert type="error" showIcon message={error} /> : null}
      <div style={{ flex: 1, minHeight: 0, display: 'grid', gridTemplateColumns: '300px 1fr', border: `1px solid ${token.colorBorderSecondary}`, background: token.colorBgContainer }}>
        <div style={{ borderRight: `1px solid ${token.colorBorderSecondary}`, overflow: 'auto' }}>
          <List
            loading={loading}
            dataSource={sessions}
            locale={{ emptyText: <Empty description="暂无会话" /> }}
            renderItem={item => (
              <List.Item
                onClick={() => setActive(item)}
                style={{
                  cursor: 'pointer',
                  padding: 14,
                  background: active?.id === item.id ? token.colorFillSecondary : undefined
                }}
              >
                <List.Item.Meta
                  avatar={<Avatar src={item.userAvatar} icon={<UserOutlined />} />}
                  title={
                    <Space>
                      {item.userName || item.userId || item.id}
                      {item.unreadCount ? <Tag color="red">{item.unreadCount}</Tag> : null}
                    </Space>
                  }
                  description={item.lastMessage || '暂无消息'}
                />
              </List.Item>
            )}
          />
        </div>
        <div style={{ minWidth: 0, display: 'flex', flexDirection: 'column' }}>
          {active ? (
            <>
              <div style={{ height: 56, padding: '0 16px', display: 'flex', alignItems: 'center', justifyContent: 'space-between', borderBottom: `1px solid ${token.colorBorderSecondary}` }}>
                <Space>
                  <Title level={5} style={{ margin: 0 }}>{active.userName || active.id}</Title>
                  <Tag color={active.status === 'open' ? 'success' : 'default'}>{active.status === 'open' ? '进行中' : '已关闭'}</Tag>
                  <Tag color={wsStatus === 'open' ? 'success' : wsStatus === 'connecting' ? 'processing' : 'default'}>
                    {wsStatus === 'open' ? 'WS在线' : wsStatus === 'connecting' ? 'WS连接中' : 'WS离线'}
                  </Tag>
                </Space>
                <Space>
                  <Button icon={<CheckOutlined />} onClick={() => active && apiReadChatSession(active.id).then(loadSessions)}>标记已读</Button>
                  <Button danger icon={<CloseCircleOutlined />} onClick={closeSession}>关闭会话</Button>
                </Space>
              </div>
              <div style={{ flex: 1, overflow: 'auto', padding: 16 }}>
                {messages.length ? (
                  <Space direction="vertical" style={{ width: '100%' }}>
                    {messages.map(item => (
                      <div key={item.id} style={{ display: 'flex', justifyContent: item.sender === 'admin' ? 'flex-end' : 'flex-start' }}>
                        <div style={{ maxWidth: '70%', padding: '8px 12px', borderRadius: 6, background: item.sender === 'admin' ? token.colorPrimaryBg : token.colorFillSecondary }}>
                          <div>{item.content}</div>
                          <Text type="secondary" style={{ fontSize: 12 }}>{item.createdAt}</Text>
                        </div>
                      </div>
                    ))}
                  </Space>
                ) : (
                  <Empty description="暂无消息" />
                )}
              </div>
              <div style={{ padding: 12, borderTop: `1px solid ${token.colorBorderSecondary}` }}>
                <Space.Compact style={{ width: '100%' }}>
                  <Input value={text} onChange={e => setText(e.target.value)} onPressEnter={send} placeholder="输入回复内容" disabled={active.status !== 'open'} />
                  <Button type="primary" icon={<SendOutlined />} onClick={send} disabled={active.status !== 'open'}>发送</Button>
                </Space.Compact>
              </div>
            </>
          ) : (
            <Empty description="请选择会话" style={{ margin: 'auto' }} />
          )}
        </div>
      </div>
    </div>
  )
}

export default Chat
