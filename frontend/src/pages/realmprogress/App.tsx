import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Table,
  Tag,
  Tooltip,
  Typography,
  Progress,
  Empty
} from 'antd'
import { RocketOutlined, DollarOutlined } from '@ant-design/icons'
import { apiGetRealmProgress, apiGetRealmPayCorrelation } from '@/api'

const { Title, Text } = Typography

const stageColors: Record<string, string> = {
  凡人: '#8c8c8c',
  练气: '#52c41a',
  筑基: '#1890ff',
  金丹: '#faad14',
  元婴: '#722ed1',
  化神: '#f5222d',
  洞虚: '#eb2f96',
  大乘: '#13c2c2',
  仙: '#fa541c',
  金仙: '#ff4d4f',
  道祖: '#cf1322'
}

// 根据境界名匹配大境界颜色（如 "练气一层" -> 练气 的颜色）
const getRealmColor = (realmName: string): string => {
  if (stageColors[realmName]) return stageColors[realmName]
  // 按 key 长度降序匹配，优先匹配 "金仙" 而非 "金"
  const keys = Object.keys(stageColors).sort((a, b) => b.length - a.length)
  for (const key of keys) {
    if (realmName.startsWith(key)) return stageColors[key]
  }
  return '#595959'
}

interface RealmProgressItem {
  realm_id: number
  realm_name: string
  stage_level: number
  user_count: number
  avg_days_played: number
  avg_recharge: number
  avg_vip_level: number
  active_rate: number
  inactive_count: number
}

interface PayCorrelationItem {
  stage_level: number
  stage_name: string
  user_count: number
  paying_count: number
  pay_rate: number
  avg_recharge: number
  total_recharge: number
}

const RealmProgress: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [inactiveDays, setInactiveDays] = useState(7)
  const [progressData, setProgressData] = useState<RealmProgressItem[]>([])
  const [payData, setPayData] = useState<PayCorrelationItem[]>([])

  useEffect(() => {
    loadAll()
  }, [inactiveDays])

  const loadAll = async () => {
    setLoading(true)
    try {
      const [pRes, cRes] = await Promise.all([
        apiGetRealmProgress(inactiveDays),
        apiGetRealmPayCorrelation()
      ])
      if (pRes?.data) setProgressData(pRes.data || [])
      if (cRes?.data) setPayData(cRes.data || [])
    } catch {
    } finally {
      setLoading(false)
    }
  }

  const maxUsers = Math.max(...progressData.map(d => d.user_count), 1)
  const maxDays = Math.max(...progressData.map(d => d.avg_days_played), 1)

  const progressColumns = [
    {
      title: '境界',
      dataIndex: 'realm_name',
      key: 'realm_name',
      width: 100,
      render: (v: string) => <Tag color={getRealmColor(v)}>{v}</Tag>
    },
    {
      title: '人数',
      dataIndex: 'user_count',
      key: 'user_count',
      width: 120,
      sorter: (a: RealmProgressItem, b: RealmProgressItem) =>
        a.user_count - b.user_count,
      render: (v: number) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <div
            style={{
              width: `${(v / maxUsers) * 80}px`,
              height: 14,
              minWidth: 2,
              background: '#1890ff',
              borderRadius: 3,
              opacity: 0.7
            }}
          />
          <Text strong>{v}</Text>
        </div>
      )
    },
    {
      title: '平均游龄(天)',
      dataIndex: 'avg_days_played',
      key: 'avg_days',
      width: 130,
      sorter: (a: RealmProgressItem, b: RealmProgressItem) =>
        a.avg_days_played - b.avg_days_played,
      render: (v: number) => (
        <Tooltip title={`平均注册 ${v.toFixed(1)} 天到达此境界`}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <div
              style={{
                width: `${(v / maxDays) * 60}px`,
                height: 14,
                minWidth: 2,
                background: '#722ed1',
                borderRadius: 3,
                opacity: 0.7
              }}
            />
            <span>{v.toFixed(1)}</span>
          </div>
        </Tooltip>
      )
    },
    {
      title: '平均充值(元)',
      dataIndex: 'avg_recharge',
      key: 'avg_recharge',
      width: 110,
      sorter: (a: RealmProgressItem, b: RealmProgressItem) =>
        a.avg_recharge - b.avg_recharge,
      render: (v: number) => (
        <Text style={{ color: v > 0 ? '#faad14' : '#ccc' }}>
          {v.toFixed(2)}
        </Text>
      )
    },
    {
      title: '平均VIP',
      dataIndex: 'avg_vip_level',
      key: 'avg_vip',
      width: 80,
      render: (v: number) =>
        v > 0 ? (
          <Tag color="gold">V{v.toFixed(1)}</Tag>
        ) : (
          <Text type="secondary">-</Text>
        )
    },
    {
      title: '活跃率',
      dataIndex: 'active_rate',
      key: 'active_rate',
      width: 130,
      sorter: (a: RealmProgressItem, b: RealmProgressItem) =>
        a.active_rate - b.active_rate,
      render: (v: number) => (
        <Progress
          percent={Math.round(v)}
          size="small"
          style={{ width: 80, margin: 0 }}
          strokeColor={v > 70 ? '#52c41a' : v > 40 ? '#faad14' : '#f5222d'}
          format={p => `${p}%`}
        />
      )
    },
    {
      title: '不活跃',
      dataIndex: 'inactive_count',
      key: 'inactive',
      width: 80,
      render: (v: number) => (
        <Text style={{ color: v > 0 ? '#f5222d' : '#52c41a' }}>{v}</Text>
      )
    }
  ]

  const payColumns = [
    {
      title: '大境界',
      dataIndex: 'stage_name',
      key: 'stage_name',
      width: 80,
      render: (v: string) => <Tag color={stageColors[v] || '#595959'}>{v}</Tag>
    },
    {
      title: '总人数',
      dataIndex: 'user_count',
      key: 'user_count',
      width: 80,
      render: (v: number) => <Text strong>{v}</Text>
    },
    {
      title: '付费人数',
      dataIndex: 'paying_count',
      key: 'paying_count',
      width: 80,
      render: (v: number) => (
        <Text style={{ color: v > 0 ? '#faad14' : '#ccc' }}>{v}</Text>
      )
    },
    {
      title: '付费率',
      dataIndex: 'pay_rate',
      key: 'pay_rate',
      width: 130,
      sorter: (a: PayCorrelationItem, b: PayCorrelationItem) =>
        a.pay_rate - b.pay_rate,
      render: (v: number) => (
        <Progress
          percent={Math.round(v)}
          size="small"
          style={{ width: 80, margin: 0 }}
          strokeColor={v > 30 ? '#52c41a' : v > 10 ? '#faad14' : '#f5222d'}
          format={p => `${p}%`}
        />
      )
    },
    {
      title: '人均充值(元)',
      dataIndex: 'avg_recharge',
      key: 'avg_recharge',
      width: 110,
      sorter: (a: PayCorrelationItem, b: PayCorrelationItem) =>
        a.avg_recharge - b.avg_recharge,
      render: (v: number) => (
        <Text strong style={{ color: '#faad14' }}>
          {v.toFixed(2)}
        </Text>
      )
    },
    {
      title: '总充值(元)',
      dataIndex: 'total_recharge',
      key: 'total_recharge',
      width: 110,
      render: (v: number) => <Text>{(v / 100).toLocaleString()}</Text>
    }
  ]

  // 计算境界分布热力图数据
  const totalUsers = progressData.reduce((s, d) => s + d.user_count, 0)

  return (
    <Spin spinning={loading}>
      <div>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: 16
          }}
        >
          <Title level={3} style={{ margin: 0 }}>
            <RocketOutlined style={{ marginRight: 8 }} />
            境界进阶分析
          </Title>
          <Select
            value={inactiveDays}
            onChange={setInactiveDays}
            style={{ width: 140 }}
            options={[
              { label: '7天未活跃', value: 7 },
              { label: '14天未活跃', value: 14 },
              { label: '30天未活跃', value: 30 }
            ]}
          />
        </div>

        {/* 境界分布热力条 */}
        <Card size="small" style={{ marginBottom: 16 }}>
          <div style={{ marginBottom: 8 }}>
            <Text type="secondary">境界分布热力图（共 {totalUsers} 人）</Text>
          </div>
          <div
            style={{
              display: 'flex',
              height: 36,
              borderRadius: 6,
              overflow: 'hidden'
            }}
          >
            {progressData.map(d => {
              const pct = totalUsers > 0 ? (d.user_count / totalUsers) * 100 : 0
              if (pct < 0.5) return null
              return (
                <Tooltip
                  key={d.realm_id}
                  title={`${d.realm_name}: ${d.user_count} 人 (${pct.toFixed(1)}%)`}
                >
                  <div
                    style={{
                      width: `${pct}%`,
                      height: '100%',
                      background: getRealmColor(d.realm_name),
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      cursor: 'pointer',
                      transition: 'opacity 0.2s',
                      minWidth: pct > 3 ? undefined : 0
                    }}
                    onMouseEnter={e => (e.currentTarget.style.opacity = '0.8')}
                    onMouseLeave={e => (e.currentTarget.style.opacity = '1')}
                  >
                    {pct > 5 && (
                      <span
                        style={{ color: '#fff', fontSize: 11, fontWeight: 600 }}
                      >
                        {d.realm_name}
                      </span>
                    )}
                  </div>
                </Tooltip>
              )
            })}
          </div>
        </Card>

        {/* 境界明细表 */}
        <Card title="各境界综合统计" size="small" style={{ marginBottom: 16 }}>
          {progressData.length > 0 ? (
            <Table
              dataSource={progressData}
              columns={progressColumns}
              rowKey="realm_id"
              size="small"
              pagination={false}
              scroll={{ x: 800 }}
            />
          ) : (
            <Empty description="暂无数据" />
          )}
        </Card>

        {/* 境界-付费关联 */}
        <Row gutter={16}>
          <Col xs={24} lg={14}>
            <Card
              title={
                <span>
                  <DollarOutlined
                    style={{ color: '#faad14', marginRight: 6 }}
                  />
                  境界-付费关联分析
                </span>
              }
              size="small"
            >
              {payData.length > 0 ? (
                <Table
                  dataSource={payData}
                  columns={payColumns}
                  rowKey="stage_level"
                  size="small"
                  pagination={false}
                  scroll={{ x: 600 }}
                />
              ) : (
                <Empty description="暂无数据" />
              )}
            </Card>
          </Col>
          <Col xs={24} lg={10}>
            <Card title="付费率随境界变化" size="small">
              {payData.length > 0 ? (
                <div style={{ padding: '16px 0' }}>
                  {payData.map(d => {
                    const maxRecharge = Math.max(
                      ...payData.map(p => p.avg_recharge),
                      1
                    )
                    const rechargeWidth = (d.avg_recharge / maxRecharge) * 100
                    return (
                      <div
                        key={d.stage_level}
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          marginBottom: 10,
                          gap: 8
                        }}
                      >
                        <Tag
                          color={stageColors[d.stage_name] || '#595959'}
                          style={{ width: 50, textAlign: 'center', margin: 0 }}
                        >
                          {d.stage_name}
                        </Tag>
                        <div style={{ flex: 1 }}>
                          <div
                            style={{ display: 'flex', gap: 4, marginBottom: 2 }}
                          >
                            <Tooltip title={`付费率 ${d.pay_rate.toFixed(1)}%`}>
                              <div
                                style={{
                                  width: `${d.pay_rate}%`,
                                  height: 10,
                                  minWidth: 2,
                                  background: '#52c41a',
                                  borderRadius: '3px 0 0 3px',
                                  opacity: 0.8
                                }}
                              />
                            </Tooltip>
                          </div>
                          <div style={{ display: 'flex', gap: 4 }}>
                            <Tooltip
                              title={`人均充值 ${d.avg_recharge.toFixed(2)} 元`}
                            >
                              <div
                                style={{
                                  width: `${rechargeWidth}%`,
                                  height: 10,
                                  minWidth: 2,
                                  background: '#faad14',
                                  borderRadius: '3px 0 0 3px',
                                  opacity: 0.8
                                }}
                              />
                            </Tooltip>
                          </div>
                        </div>
                        <div
                          style={{
                            width: 80,
                            textAlign: 'right',
                            fontSize: 11
                          }}
                        >
                          <div style={{ color: '#52c41a' }}>
                            {d.pay_rate.toFixed(1)}%
                          </div>
                          <div style={{ color: '#faad14' }}>
                            ¥{d.avg_recharge.toFixed(1)}
                          </div>
                        </div>
                      </div>
                    )
                  })}
                  <div
                    style={{
                      display: 'flex',
                      gap: 16,
                      justifyContent: 'center',
                      marginTop: 12
                    }}
                  >
                    <span style={{ fontSize: 11 }}>
                      <span
                        style={{
                          display: 'inline-block',
                          width: 10,
                          height: 10,
                          background: '#52c41a',
                          borderRadius: 2,
                          marginRight: 4
                        }}
                      ></span>
                      付费率
                    </span>
                    <span style={{ fontSize: 11 }}>
                      <span
                        style={{
                          display: 'inline-block',
                          width: 10,
                          height: 10,
                          background: '#faad14',
                          borderRadius: 2,
                          marginRight: 4
                        }}
                      ></span>
                      人均充值
                    </span>
                  </div>
                </div>
              ) : (
                <Empty description="暂无数据" />
              )}
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  )
}

export default RealmProgress
