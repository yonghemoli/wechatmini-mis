import React, { useEffect, useState, useMemo } from 'react'
import {
  Card,
  Row,
  Col,
  InputNumber,
  Button,
  Tabs,
  Tag,
  Statistic,
  Space,
  Select,
  Typography,
  Alert,
  Tooltip as AntTooltip,
  Spin,
  message,
  Empty
} from 'antd'
import {
  ReloadOutlined,
  SearchOutlined,
  UserOutlined,
  GlobalOutlined,
  RobotOutlined,
  ClockCircleOutlined,
  FireOutlined
} from '@ant-design/icons'
import {
  ResponsiveContainer,
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  Cell
} from 'recharts'
import {
  apiGetPlayerHourly,
  apiGetPlayerDaily,
  apiGetPlayerPeak,
  apiCheckPlayerBot,
  apiGetGlobalHourly,
  apiGetGlobalPeak,
  apiRefreshActivity
} from '@/api'

const { Title, Text } = Typography

// ==================== 类型 ====================
interface HourlyData {
  date: string
  hour: number
  command_count?: number
  active_users?: number
  total_commands?: number
  first_seen_at?: string
  last_seen_at?: string
}

interface DailyData {
  date: string
  total_cmds: number
  active_hours: number
  online_minute: number
  first_seen: string
  last_seen: string
}

interface HourAvg {
  hour: number
  avg_users: number
  avg_commands: number
}

interface WeekdayAvg {
  weekday: number
  weekday_name: string
  avg_users: number
  avg_commands: number
}

interface PeakAnalysis {
  hourly_avg: HourAvg[]
  weekday_avg: WeekdayAvg[]
  peak_hour: number
  peak_weekday: number
  peak_users: number
  total_days: number
  total_cmds?: number
}

interface BotResult {
  uid: number
  is_likely_bot: boolean
  bot_score: number
  reasons: string[]
  avg_daily_cmds: number
  active_days: number
  avg_active_hours: number
  cmd_std_dev: number
  night_activity: number
  regularity: number
}

// ==================== 工具 ====================
const hourLabel = (h: number) => `${String(h).padStart(2, '0')}:00`

const heatColor = (value: number, max: number) => {
  if (max === 0 || value === 0) return '#f5f5f5'
  const ratio = value / max
  if (ratio > 0.8) return '#ff4d4f'
  if (ratio > 0.6) return '#fa8c16'
  if (ratio > 0.4) return '#fadb14'
  if (ratio > 0.2) return '#95de64'
  return '#d9f7be'
}

// ==================== 组件 ====================
const Activity: React.FC = () => {
  return (
    <div>
      <Title level={3}>活跃度分析</Title>
      <Tabs
        defaultActiveKey="global"
        items={[
          {
            key: 'global',
            label: (
              <span>
                <GlobalOutlined /> 全服分析
              </span>
            ),
            children: <GlobalTab />
          },
          {
            key: 'player',
            label: (
              <span>
                <UserOutlined /> 玩家分析
              </span>
            ),
            children: <PlayerTab />
          }
        ]}
      />
    </div>
  )
}

// ==================== 全服分析 ====================
const GlobalTab: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [days, setDays] = useState(7)
  const [hourlyData, setHourlyData] = useState<HourlyData[]>([])
  const [peakData, setPeakData] = useState<PeakAnalysis | null>(null)

  const loadData = async (d: number) => {
    setLoading(true)
    try {
      const [hourlyRes, peakRes] = await Promise.all([
        apiGetGlobalHourly(d),
        apiGetGlobalPeak(d)
      ])
      if (hourlyRes?.data) setHourlyData(hourlyRes.data || [])
      if (peakRes?.data) setPeakData(peakRes.data)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData(days)
  }, [days])

  const handleRefresh = async () => {
    await apiRefreshActivity(days)
    message.success('聚合任务已启动，请稍后刷新页面')
    setTimeout(() => loadData(days), 3000)
  }

  // 热力图数据：日期 × 小时
  const heatmapData = useMemo(() => {
    if (!hourlyData.length)
      return { dates: [] as string[], matrix: [] as number[][], max: 0 }
    const dateSet = new Map<string, Map<number, number>>()
    let max = 0
    for (const d of hourlyData) {
      if (!dateSet.has(d.date)) dateSet.set(d.date, new Map())
      const val = d.active_users || 0
      dateSet.get(d.date)!.set(d.hour, val)
      if (val > max) max = val
    }
    const dates = Array.from(dateSet.keys()).sort()
    const matrix = dates.map(date => {
      const hourMap = dateSet.get(date)!
      return Array.from({ length: 24 }, (_, h) => hourMap.get(h) || 0)
    })
    return { dates, matrix, max }
  }, [hourlyData])

  // 全服每日趋势（按天聚合在线人数和指令量）
  const dailyTrend = useMemo(() => {
    const dayMap = new Map<
      string,
      { users: Set<number>; cmds: number; maxUsers: number }
    >()
    for (const d of hourlyData) {
      if (!dayMap.has(d.date))
        dayMap.set(d.date, { users: new Set(), cmds: 0, maxUsers: 0 })
      const entry = dayMap.get(d.date)!
      entry.cmds += d.total_commands || 0
      const au = d.active_users || 0
      if (au > entry.maxUsers) entry.maxUsers = au
    }
    return Array.from(dayMap.entries())
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([date, data]) => ({
        date: date.slice(5),
        指令总量: data.cmds,
        峰值在线: data.maxUsers
      }))
  }, [hourlyData])

  return (
    <Spin spinning={loading}>
      <Space style={{ marginBottom: 16 }}>
        <Select
          value={days}
          onChange={v => {
            setDays(v)
            loadData(v)
          }}
          options={[
            { label: '近7天', value: 7 },
            { label: '近14天', value: 14 },
            { label: '近30天', value: 30 }
          ]}
          style={{ width: 120 }}
        />
        <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
          重新聚合
        </Button>
      </Space>

      {/* 高峰时段概览 */}
      {peakData && peakData.total_days > 0 && (
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={12} sm={6}>
            <Card size="small">
              <Statistic
                title="高峰时段"
                value={hourLabel(peakData.peak_hour)}
                prefix={<FireOutlined style={{ color: '#f5222d' }} />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card size="small">
              <Statistic
                title="高峰平均在线"
                value={peakData.peak_users}
                suffix="人"
              />
            </Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card size="small">
              <Statistic
                title="最活跃星期"
                value={
                  peakData.weekday_avg[peakData.peak_weekday]?.weekday_name ||
                  '-'
                }
                prefix={<ClockCircleOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card size="small">
              <Statistic
                title="统计天数"
                value={peakData.total_days}
                suffix="天"
              />
            </Card>
          </Col>
        </Row>
      )}

      {/* 热力图 */}
      <Card
        title="在线人数热力图（日期 × 小时）"
        size="small"
        style={{ marginBottom: 16 }}
      >
        {heatmapData.dates.length === 0 ? (
          <Empty description="暂无数据，请先聚合" />
        ) : (
          <div style={{ overflowX: 'auto' }}>
            {/* 小时标题行 */}
            <div style={{ display: 'flex', marginLeft: 80, marginBottom: 4 }}>
              {Array.from({ length: 24 }, (_, h) => (
                <div
                  key={h}
                  style={{
                    width: 28,
                    minWidth: 28,
                    fontSize: 10,
                    textAlign: 'center',
                    color: '#999'
                  }}
                >
                  {h}
                </div>
              ))}
            </div>
            {/* 行 */}
            {heatmapData.dates.map((date, di) => (
              <div
                key={date}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  marginBottom: 2
                }}
              >
                <div
                  style={{
                    width: 76,
                    minWidth: 76,
                    fontSize: 11,
                    color: '#666',
                    textAlign: 'right',
                    paddingRight: 4
                  }}
                >
                  {date.slice(5)}
                </div>
                {heatmapData.matrix[di].map((val, h) => (
                  <AntTooltip
                    key={h}
                    title={`${date} ${hourLabel(h)}: ${val} 人在线`}
                  >
                    <div
                      style={{
                        width: 28,
                        minWidth: 28,
                        height: 20,
                        background: heatColor(val, heatmapData.max),
                        borderRadius: 2,
                        margin: '0 1px',
                        cursor: 'pointer',
                        fontSize: 9,
                        lineHeight: '20px',
                        textAlign: 'center',
                        color: val > heatmapData.max * 0.6 ? '#fff' : '#999'
                      }}
                    >
                      {val > 0 ? val : ''}
                    </div>
                  </AntTooltip>
                ))}
              </div>
            ))}
            {/* 图例 */}
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                marginTop: 8,
                marginLeft: 80,
                gap: 4,
                fontSize: 11,
                color: '#999'
              }}
            >
              <span>少</span>
              {['#d9f7be', '#95de64', '#fadb14', '#fa8c16', '#ff4d4f'].map(
                c => (
                  <div
                    key={c}
                    style={{
                      width: 16,
                      height: 12,
                      background: c,
                      borderRadius: 2
                    }}
                  />
                )
              )}
              <span>多</span>
            </div>
          </div>
        )}
      </Card>

      {/* 每小时平均柱状图 */}
      {peakData && peakData.total_days > 0 && (
        <Row gutter={16}>
          <Col xs={24} md={12}>
            <Card
              title="每小时平均在线人数"
              size="small"
              style={{ marginBottom: 16 }}
            >
              <ResponsiveContainer width="100%" height={240}>
                <BarChart data={peakData.hourly_avg || []}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis
                    dataKey="hour"
                    fontSize={11}
                    tickFormatter={h => `${h}时`}
                  />
                  <YAxis fontSize={11} />
                  <Tooltip labelFormatter={h => `${h}:00`} />
                  <Bar
                    dataKey="avg_users"
                    name="平均在线"
                    radius={[4, 4, 0, 0]}
                  >
                    {(peakData.hourly_avg || []).map((entry, i) => (
                      <Cell
                        key={i}
                        fill={
                          entry.hour === peakData.peak_hour
                            ? '#ff4d4f'
                            : '#1677ff'
                        }
                      />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </Card>
          </Col>
          <Col xs={24} md={12}>
            <Card
              title="每星期几平均在线"
              size="small"
              style={{ marginBottom: 16 }}
            >
              <ResponsiveContainer width="100%" height={240}>
                <BarChart data={peakData.weekday_avg || []}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="weekday_name" fontSize={12} />
                  <YAxis fontSize={11} />
                  <Tooltip />
                  <Bar
                    dataKey="avg_users"
                    name="平均在线"
                    fill="#722ed1"
                    radius={[4, 4, 0, 0]}
                  >
                    {(peakData.weekday_avg || []).map((entry, i) => (
                      <Cell
                        key={i}
                        fill={
                          entry.weekday === peakData.peak_weekday
                            ? '#ff4d4f'
                            : '#722ed1'
                        }
                      />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>
      )}

      {/* 每日总量趋势 */}
      {dailyTrend.length > 0 && (
        <Card title="每日趋势" size="small">
          <ResponsiveContainer width="100%" height={250}>
            <LineChart data={dailyTrend}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" fontSize={11} />
              <YAxis yAxisId="left" fontSize={11} />
              <YAxis yAxisId="right" orientation="right" fontSize={11} />
              <Tooltip />
              <Legend />
              <Line
                yAxisId="left"
                type="monotone"
                dataKey="指令总量"
                stroke="#1677ff"
                strokeWidth={2}
                dot={{ r: 2 }}
              />
              <Line
                yAxisId="right"
                type="monotone"
                dataKey="峰值在线"
                stroke="#f5222d"
                strokeWidth={2}
                dot={{ r: 2 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </Card>
      )}
    </Spin>
  )
}

// ==================== 玩家分析 ====================
const PlayerTab: React.FC = () => {
  const [uid, setUid] = useState<number | null>(null)
  const [days, setDays] = useState(14)
  const [loading, setLoading] = useState(false)
  const [hourlyData, setHourlyData] = useState<HourlyData[]>([])
  const [dailyData, setDailyData] = useState<DailyData[]>([])
  const [peakData, setPeakData] = useState<PeakAnalysis | null>(null)
  const [botResult, setBotResult] = useState<BotResult | null>(null)
  const [searched, setSearched] = useState(false)

  const loadPlayer = async () => {
    if (!uid) {
      message.warning('请输入玩家 UID')
      return
    }
    setLoading(true)
    setSearched(true)
    try {
      const [hourlyRes, dailyRes, peakRes, botRes] = await Promise.all([
        apiGetPlayerHourly(uid, days),
        apiGetPlayerDaily(uid, days),
        apiGetPlayerPeak(uid, days),
        apiCheckPlayerBot(uid, days)
      ])
      setHourlyData(hourlyRes?.data || [])
      setDailyData(dailyRes?.data || [])
      setPeakData(peakRes?.data || null)
      setBotResult(botRes?.data || null)
    } finally {
      setLoading(false)
    }
  }

  // 玩家热力图
  const playerHeatmap = useMemo(() => {
    if (!hourlyData.length)
      return { dates: [] as string[], matrix: [] as number[][], max: 0 }
    const dateSet = new Map<string, Map<number, number>>()
    let max = 0
    for (const d of hourlyData) {
      if (!dateSet.has(d.date)) dateSet.set(d.date, new Map())
      const val = d.command_count || 0
      dateSet.get(d.date)!.set(d.hour, val)
      if (val > max) max = val
    }
    const dates = Array.from(dateSet.keys()).sort()
    const matrix = dates.map(date => {
      const hourMap = dateSet.get(date)!
      return Array.from({ length: 24 }, (_, h) => hourMap.get(h) || 0)
    })
    return { dates, matrix, max }
  }, [hourlyData])

  return (
    <div>
      <Space style={{ marginBottom: 16 }} wrap>
        <InputNumber
          placeholder="输入玩家 UID"
          min={1}
          style={{ width: 160 }}
          onChange={v => setUid(v as number)}
          onPressEnter={loadPlayer}
        />
        <Select
          value={days}
          onChange={setDays}
          options={[
            { label: '近7天', value: 7 },
            { label: '近14天', value: 14 },
            { label: '近30天', value: 30 }
          ]}
          style={{ width: 120 }}
        />
        <Button
          type="primary"
          icon={<SearchOutlined />}
          onClick={loadPlayer}
          loading={loading}
        >
          查询
        </Button>
      </Space>

      {/* 真人检测结果 */}
      {botResult && (
        <Alert
          type={botResult.is_likely_bot ? 'error' : 'success'}
          showIcon
          icon={<RobotOutlined />}
          style={{ marginBottom: 16 }}
          message={
            <Space>
              <span>
                {botResult.is_likely_bot
                  ? `疑似脚本/机器人 (评分: ${botResult.bot_score}/100)`
                  : `真实玩家 (评分: ${botResult.bot_score}/100)`}
              </span>
              <Tag
                color={
                  botResult.bot_score >= 60
                    ? 'red'
                    : botResult.bot_score >= 30
                      ? 'orange'
                      : 'green'
                }
              >
                风险 {botResult.bot_score}%
              </Tag>
            </Space>
          }
          description={
            <div>
              <Row gutter={16} style={{ marginTop: 8 }}>
                <Col span={6}>
                  <Text type="secondary">日均指令: </Text>
                  <Text strong>{botResult.avg_daily_cmds}</Text>
                </Col>
                <Col span={6}>
                  <Text type="secondary">活跃天数: </Text>
                  <Text strong>{botResult.active_days}天</Text>
                </Col>
                <Col span={6}>
                  <Text type="secondary">日均在线: </Text>
                  <Text strong>{botResult.avg_active_hours}h</Text>
                </Col>
                <Col span={6}>
                  <Text type="secondary">凌晨占比: </Text>
                  <Text strong>{botResult.night_activity}%</Text>
                </Col>
              </Row>
              <Row style={{ marginTop: 8 }}>
                <Col span={24}>
                  {(botResult.reasons || []).map((r, i) => (
                    <Tag
                      key={i}
                      color={botResult.is_likely_bot ? 'red' : 'green'}
                      style={{ marginBottom: 4 }}
                    >
                      {r}
                    </Tag>
                  ))}
                </Col>
              </Row>
            </div>
          }
        />
      )}

      {searched && !loading && hourlyData.length === 0 && (
        <Empty description="该玩家暂无活跃数据" />
      )}

      {/* 玩家高峰概览 */}
      {peakData && peakData.total_days > 0 && (
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={8} sm={6}>
            <Card size="small">
              <Statistic
                title="个人高峰"
                value={hourLabel(peakData.peak_hour)}
                prefix={<FireOutlined style={{ color: '#f5222d' }} />}
              />
            </Card>
          </Col>
          <Col xs={8} sm={6}>
            <Card size="small">
              <Statistic
                title="活跃天数"
                value={peakData.total_days}
                suffix="天"
              />
            </Card>
          </Col>
          <Col xs={8} sm={6}>
            <Card size="small">
              <Statistic title="总指令数" value={peakData.total_cmds} />
            </Card>
          </Col>
        </Row>
      )}

      {/* 玩家热力图 */}
      {playerHeatmap.dates.length > 0 && (
        <Card
          title="指令数热力图（日期 × 小时）"
          size="small"
          style={{ marginBottom: 16 }}
        >
          <div style={{ overflowX: 'auto' }}>
            <div style={{ display: 'flex', marginLeft: 80, marginBottom: 4 }}>
              {Array.from({ length: 24 }, (_, h) => (
                <div
                  key={h}
                  style={{
                    width: 28,
                    minWidth: 28,
                    fontSize: 10,
                    textAlign: 'center',
                    color: '#999'
                  }}
                >
                  {h}
                </div>
              ))}
            </div>
            {playerHeatmap.dates.map((date, di) => (
              <div
                key={date}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  marginBottom: 2
                }}
              >
                <div
                  style={{
                    width: 76,
                    minWidth: 76,
                    fontSize: 11,
                    color: '#666',
                    textAlign: 'right',
                    paddingRight: 4
                  }}
                >
                  {date.slice(5)}
                </div>
                {playerHeatmap.matrix[di].map((val, h) => (
                  <AntTooltip
                    key={h}
                    title={`${date} ${hourLabel(h)}: ${val} 次指令`}
                  >
                    <div
                      style={{
                        width: 28,
                        minWidth: 28,
                        height: 20,
                        background: heatColor(val, playerHeatmap.max),
                        borderRadius: 2,
                        margin: '0 1px',
                        cursor: 'pointer',
                        fontSize: 9,
                        lineHeight: '20px',
                        textAlign: 'center',
                        color: val > playerHeatmap.max * 0.6 ? '#fff' : '#999'
                      }}
                    >
                      {val > 0 ? val : ''}
                    </div>
                  </AntTooltip>
                ))}
              </div>
            ))}
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                marginTop: 8,
                marginLeft: 80,
                gap: 4,
                fontSize: 11,
                color: '#999'
              }}
            >
              <span>少</span>
              {['#d9f7be', '#95de64', '#fadb14', '#fa8c16', '#ff4d4f'].map(
                c => (
                  <div
                    key={c}
                    style={{
                      width: 16,
                      height: 12,
                      background: c,
                      borderRadius: 2
                    }}
                  />
                )
              )}
              <span>多</span>
            </div>
          </div>
        </Card>
      )}

      {/* 玩家每日趋势 */}
      {dailyData.length > 0 && (
        <Row gutter={16}>
          <Col xs={24} md={12}>
            <Card
              title="每日指令数 & 在线时长"
              size="small"
              style={{ marginBottom: 16 }}
            >
              <ResponsiveContainer width="100%" height={240}>
                <LineChart
                  data={dailyData.map(d => ({
                    date: d.date.slice(5),
                    指令数: d.total_cmds,
                    在线分钟: Math.round(d.online_minute)
                  }))}
                >
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" fontSize={11} />
                  <YAxis yAxisId="left" fontSize={11} />
                  <YAxis yAxisId="right" orientation="right" fontSize={11} />
                  <Tooltip />
                  <Legend />
                  <Line
                    yAxisId="left"
                    type="monotone"
                    dataKey="指令数"
                    stroke="#1677ff"
                    strokeWidth={2}
                    dot={{ r: 2 }}
                  />
                  <Line
                    yAxisId="right"
                    type="monotone"
                    dataKey="在线分钟"
                    stroke="#52c41a"
                    strokeWidth={2}
                    dot={{ r: 2 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </Card>
          </Col>
          <Col xs={24} md={12}>
            <Card
              title="每日活跃小时数"
              size="small"
              style={{ marginBottom: 16 }}
            >
              <ResponsiveContainer width="100%" height={240}>
                <BarChart
                  data={dailyData.map(d => ({
                    date: d.date.slice(5),
                    活跃小时: d.active_hours
                  }))}
                >
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" fontSize={11} />
                  <YAxis fontSize={11} domain={[0, 24]} />
                  <Tooltip />
                  <Bar dataKey="活跃小时" fill="#722ed1" radius={[4, 4, 0, 0]}>
                    {dailyData.map((d, i) => (
                      <Cell
                        key={i}
                        fill={
                          d.active_hours >= 18
                            ? '#ff4d4f'
                            : d.active_hours >= 12
                              ? '#fa8c16'
                              : '#722ed1'
                        }
                      />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>
      )}

      {/* 玩家个人高峰小时图 */}
      {peakData && peakData.total_days > 0 && (
        <Card title="个人各时段平均指令数" size="small">
          <ResponsiveContainer width="100%" height={200}>
            <BarChart data={peakData.hourly_avg || []}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="hour"
                fontSize={11}
                tickFormatter={h => `${h}时`}
              />
              <YAxis fontSize={11} />
              <Tooltip labelFormatter={h => `${h}:00`} />
              <Bar
                dataKey="avg_commands"
                name="平均指令数"
                radius={[4, 4, 0, 0]}
              >
                {((peakData.hourly_avg || []) as HourAvg[]).map(
                  (entry: HourAvg, i: number) => (
                    <Cell
                      key={i}
                      fill={
                        entry.hour === peakData.peak_hour
                          ? '#ff4d4f'
                          : '#1677ff'
                      }
                    />
                  )
                )}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </Card>
      )}
    </div>
  )
}

export default Activity
