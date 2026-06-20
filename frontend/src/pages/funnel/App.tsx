import React, { useEffect, useState } from 'react'
import { Card, Select, Spin, Typography, Empty, Tooltip } from 'antd'
import { FunnelPlotOutlined } from '@ant-design/icons'
import { apiGetFunnel } from '@/api'

const { Title, Text } = Typography

interface FunnelData {
  total_registered: number
  first_battle: number
  first_break: number
  first_promotion: number
  joined_guild: number
  first_recharge: number
}

const steps: Array<{
  key: keyof FunnelData
  label: string
  color: string
  icon: string
}> = [
  { key: 'total_registered', label: '注册', color: '#1890ff', icon: '📝' },
  { key: 'first_break', label: '首次突破', color: '#722ed1', icon: '🚀' },
  { key: 'first_promotion', label: '首次晋升', color: '#eb2f96', icon: '🏅' }
]

const Funnel: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(0) // 0 = 全部
  const [data, setData] = useState<FunnelData | null>(null)

  useEffect(() => {
    const loadData = async () => {
      setLoading(true)
      try {
        const res = await apiGetFunnel(days)
        if (res?.data) {
          setData(res.data)
        } else {
          setData(null)
        }
      } catch {
        setData(null)
      } finally {
        setLoading(false)
      }
    }

    void loadData()
  }, [days])

  const getVal = (key: keyof FunnelData) => {
    if (!data) return 0
    return data[key] || 0
  }

  const total = getVal('total_registered') || 1

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
            <FunnelPlotOutlined style={{ marginRight: 8 }} />
            转化漏斗
          </Title>
          <Select
            value={days}
            onChange={setDays}
            style={{ width: 140 }}
            options={[
              { label: '全部时间', value: 0 },
              { label: '近7天注册', value: 7 },
              { label: '近14天注册', value: 14 },
              { label: '近30天注册', value: 30 },
              { label: '近90天注册', value: 90 }
            ]}
          />
        </div>

        {data ? (
          <Card>
            <div style={{ maxWidth: 700, margin: '0 auto', padding: '24px 0' }}>
              {steps.map((step, idx) => {
                const val = getVal(step.key)
                const pct = (val / total) * 100
                const width = Math.max(pct, 8)
                const prevVal = idx > 0 ? getVal(steps[idx - 1].key) : val
                const stepRate = prevVal > 0 ? (val / prevVal) * 100 : 0

                return (
                  <div key={step.key} style={{ marginBottom: 8 }}>
                    {/* 标签行 */}
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        marginBottom: 4
                      }}
                    >
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 6
                        }}
                      >
                        <span style={{ fontSize: 18 }}>{step.icon}</span>
                        <Text strong>{step.label}</Text>
                      </div>
                      <div
                        style={{
                          display: 'flex',
                          gap: 16,
                          alignItems: 'center'
                        }}
                      >
                        <Text
                          style={{
                            fontSize: 20,
                            fontWeight: 700,
                            color: step.color
                          }}
                        >
                          {val.toLocaleString()}
                        </Text>
                        {idx > 0 && (
                          <Tooltip
                            title={`${steps[idx - 1].label} → ${step.label} 的转化率`}
                          >
                            <Text type="secondary" style={{ fontSize: 12 }}>
                              转化 {stepRate.toFixed(1)}%
                            </Text>
                          </Tooltip>
                        )}
                      </div>
                    </div>
                    {/* 漏斗条 */}
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'center'
                      }}
                    >
                      <Tooltip
                        title={`${step.label}: ${val} 人 (${pct.toFixed(1)}%)`}
                      >
                        <div
                          style={{
                            width: `${width}%`,
                            height: 44,
                            background: `linear-gradient(90deg, ${step.color} 0%, ${step.color}99 100%)`,
                            borderRadius: 6,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            transition: 'width 0.5s ease',
                            cursor: 'pointer'
                          }}
                        >
                          <span
                            style={{
                              color: '#fff',
                              fontWeight: 700,
                              fontSize: 14
                            }}
                          >
                            {pct.toFixed(1)}%
                          </span>
                        </div>
                      </Tooltip>
                    </div>
                    {/* 步骤间转化箭头 */}
                    {idx < steps.length - 1 && (
                      <div
                        style={{
                          textAlign: 'center',
                          color: '#d9d9d9',
                          fontSize: 18,
                          lineHeight: '24px'
                        }}
                      >
                        ↓
                      </div>
                    )}
                  </div>
                )
              })}
            </div>

            {/* 总体转化率 */}
            <div
              style={{
                marginTop: 24,
                padding: 16,
                background: '#fafafa',
                borderRadius: 8,
                display: 'flex',
                justifyContent: 'space-around',
                flexWrap: 'wrap',
                gap: 16
              }}
            >
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 12, color: '#999' }}>注册→首次突破</div>
                <div
                  style={{ fontSize: 20, fontWeight: 700, color: '#722ed1' }}
                >
                  {total > 0
                    ? ((getVal('first_break') / total) * 100).toFixed(1)
                    : '0.0'}
                  %
                </div>
              </div>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 12, color: '#999' }}>突破→首次晋升</div>
                <div
                  style={{ fontSize: 20, fontWeight: 700, color: '#eb2f96' }}
                >
                  {getVal('first_break') > 0
                    ? (
                        (getVal('first_promotion') / getVal('first_break')) *
                        100
                      ).toFixed(1)
                    : '0.0'}
                  %
                </div>
              </div>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 12, color: '#999' }}>注册→首次晋升</div>
                <div
                  style={{ fontSize: 20, fontWeight: 700, color: '#eb2f96' }}
                >
                  {total > 0
                    ? ((getVal('first_promotion') / total) * 100).toFixed(1)
                    : '0.0'}
                  %
                </div>
              </div>
            </div>
          </Card>
        ) : (
          <Card>
            <Empty description="暂无数据" />
          </Card>
        )}
      </div>
    </Spin>
  )
}

export default Funnel
