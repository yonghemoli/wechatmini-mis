import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Table,
  Tag,
  Typography,
  Modal,
  Progress,
  Tooltip,
  Segmented,
  Empty,
  Button,
  message,
  Statistic
} from 'antd'
import {
  FireOutlined,
  ThunderboltOutlined,
  TrophyOutlined,
  RiseOutlined,
  FallOutlined,
  SyncOutlined,
  MessageOutlined,
  TeamOutlined,
  GlobalOutlined,
  NodeIndexOutlined
} from '@ant-design/icons'
import {
  apiGetFeatureScores,
  apiGetFeatureCategories,
  apiGetFeatureTop,
  apiGetFeatureTrend,
  apiRefreshFeatures,
  apiGetFeatureSceneDist,
  apiGetFeatureChannelTop,
  apiGetFeaturePlatformDist,
  apiGetFeatureSceneFeatureTop,
  apiGetFeatureSceneTrend
} from '@/api'

const { Title, Text } = Typography

// ===================== 类型 =====================
interface FeatureScore {
  feature_name: string
  feature_category: string
  total_uses_7d: number
  total_uses_30d: number
  unique_users_7d: number
  unique_users_30d: number
  avg_daily_uses: number
  avg_response_ms: number
  success_rate: number
  usage_growth: number
  user_penetration: number
  quality_score: number
}

interface CategoryStat {
  feature_category: string
  total_uses: number
  unique_users: number
  avg_response_ms: number
  feature_count: number
}

interface TrendItem {
  date: string
  total_uses: number
  unique_users: number
  avg_response_ms: number
  success_rate: number
}

interface SceneDistItem {
  scene: string
  total_uses: number
  unique_users: number
}

interface ChannelTopItem {
  channel_id: string
  scene: string
  total_uses: number
  unique_users: number
}

interface PlatformDistItem {
  platform: string
  total_uses: number
  unique_users: number
}

interface SceneFeatureItem {
  scene: string
  feature_name: string
  feature_category: string
  total_uses: number
  unique_users: number
}

interface SceneTrendItem {
  date: string
  scene: string
  total_uses: number
  unique_users: number
}

// ===================== 颜色 =====================
const categoryColors: Record<string, string> = {
  日常系统: '#52c41a',
  修炼系统: '#722ed1',
  战斗系统: '#f5222d',
  社交系统: '#fa8c16',
  经济系统: '#1890ff',
  成长系统: '#13c2c2',
  充值系统: '#eb2f96',
  信息系统: '#8c8c8c',
  邀请系统: '#2f54eb',
  角色创建: '#faad14'
}

const getCategoryColor = (cat: string) => categoryColors[cat] || '#595959'

const scoreColor = (score: number) => {
  if (score >= 80) return '#52c41a'
  if (score >= 60) return '#1890ff'
  if (score >= 40) return '#faad14'
  return '#f5222d'
}

const Features: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [scores, setScores] = useState<FeatureScore[]>([])
  const [categories, setCategories] = useState<CategoryStat[]>([])
  const [topFeatures, setTopFeatures] = useState<any[]>([])
  const [filterCategory, setFilterCategory] = useState('')
  const [topDays, setTopDays] = useState(7)
  const [trendModal, setTrendModal] = useState(false)
  const [trendData, setTrendData] = useState<TrendItem[]>([])
  const [trendName, setTrendName] = useState('')
  const [trendDays, setTrendDays] = useState(30)
  const [trendLoading, setTrendLoading] = useState(false)

  // 场景/频道/平台维度
  const [sceneDist, setSceneDist] = useState<SceneDistItem[]>([])
  const [channelTop, setChannelTop] = useState<ChannelTopItem[]>([])
  const [platformDist, setPlatformDist] = useState<PlatformDistItem[]>([])
  const [sceneFeatureTop, setSceneFeatureTop] = useState<SceneFeatureItem[]>([])
  const [sceneTrend, setSceneTrend] = useState<SceneTrendItem[]>([])
  const [sceneFilter, setSceneFilter] = useState('')
  const [channelDays, setChannelDays] = useState(7)

  useEffect(() => {
    loadAll()
  }, [])

  useEffect(() => {
    loadScores(filterCategory)
  }, [filterCategory])

  useEffect(() => {
    loadTop(topDays)
  }, [topDays])

  const loadAll = async () => {
    setLoading(true)
    try {
      await Promise.all([
        loadScores(''),
        loadCategories(),
        loadTop(7),
        loadSceneDist(),
        loadChannelTop(7),
        loadPlatformDist(),
        loadSceneFeatureTop('', 7),
        loadSceneTrend()
      ])
    } finally {
      setLoading(false)
    }
  }

  const loadScores = async (cat: string) => {
    try {
      const res = await apiGetFeatureScores(cat || undefined)
      if (res?.data) setScores(res.data || [])
    } catch {}
  }

  const loadCategories = async () => {
    try {
      const res = await apiGetFeatureCategories(30)
      if (res?.data) setCategories(res.data || [])
    } catch {}
  }

  const loadTop = async (days: number) => {
    try {
      const res = await apiGetFeatureTop(days, 20)
      if (res?.data) setTopFeatures(res.data || [])
    } catch {}
  }

  const loadSceneDist = async () => {
    try {
      const res = await apiGetFeatureSceneDist(30)
      if (res?.data) setSceneDist(res.data || [])
    } catch {}
  }

  const loadChannelTop = async (days: number) => {
    try {
      const res = await apiGetFeatureChannelTop(days, 20)
      if (res?.data) setChannelTop(res.data || [])
    } catch {}
  }

  const loadPlatformDist = async () => {
    try {
      const res = await apiGetFeaturePlatformDist(30)
      if (res?.data) setPlatformDist(res.data || [])
    } catch {}
  }

  const loadSceneFeatureTop = async (scene: string, days: number) => {
    try {
      const res = await apiGetFeatureSceneFeatureTop(scene, days, 15)
      if (res?.data) setSceneFeatureTop(res.data || [])
    } catch {}
  }

  const loadSceneTrend = async () => {
    try {
      const res = await apiGetFeatureSceneTrend(30)
      if (res?.data) setSceneTrend(res.data || [])
    } catch {}
  }

  useEffect(() => {
    loadChannelTop(channelDays)
  }, [channelDays])

  useEffect(() => {
    loadSceneFeatureTop(sceneFilter, 7)
  }, [sceneFilter])

  const openTrend = async (name: string) => {
    setTrendName(name)
    setTrendModal(true)
    setTrendLoading(true)
    try {
      const res = await apiGetFeatureTrend(name, trendDays)
      if (res?.data) setTrendData(res.data || [])
    } finally {
      setTrendLoading(false)
    }
  }

  const loadTrend = async (name: string, days: number) => {
    setTrendLoading(true)
    try {
      const res = await apiGetFeatureTrend(name, days)
      if (res?.data) setTrendData(res.data || [])
    } finally {
      setTrendLoading(false)
    }
  }

  const handleRefresh = async () => {
    Modal.confirm({
      title: '确认刷新',
      content: '将重新聚合功能埋点数据并计算评分，确定执行吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        const hide = message.loading('正在刷新功能分析数据...', 0)
        try {
          const res = await apiRefreshFeatures()
          hide()
          if (res?.data) {
            const tasks = res.data.tasks || []
            const failed = tasks.filter((t: any) => !t.success)
            if (failed.length > 0) {
              message.warning(
                `刷新完成，但 ${failed.length} 项失败，耗时 ${res.data.total}`
              )
            } else {
              message.success(`刷新完成，耗时 ${res.data.total}`)
            }
          }
          loadAll()
        } catch {
          hide()
          message.error('刷新失败')
        }
      }
    })
  }

  // ===================== 分类概览卡片 =====================
  const renderCategoryCards = () => {
    if (!categories.length)
      return <Empty description="暂无分类数据" style={{ padding: 40 }} />
    const totalUses = categories.reduce((s, c) => s + c.total_uses, 0)
    return (
      <Row gutter={[12, 12]}>
        {categories.map(cat => {
          const pct =
            totalUses > 0
              ? ((cat.total_uses / totalUses) * 100).toFixed(1)
              : '0'
          return (
            <Col xs={12} sm={8} md={6} key={cat.feature_category}>
              <Card
                size="small"
                hoverable
                onClick={() => setFilterCategory(cat.feature_category)}
                style={{
                  borderLeft: `4px solid ${getCategoryColor(cat.feature_category)}`
                }}
              >
                <div style={{ fontWeight: 600, marginBottom: 4 }}>
                  {cat.feature_category}
                </div>
                <div style={{ fontSize: 20, fontWeight: 700 }}>
                  {cat.total_uses.toLocaleString()}
                  <span
                    style={{
                      fontSize: 12,
                      fontWeight: 400,
                      color: '#999',
                      marginLeft: 4
                    }}
                  >
                    次 ({pct}%)
                  </span>
                </div>
                <div style={{ color: '#666', fontSize: 12, marginTop: 2 }}>
                  {cat.unique_users} 用户 · {cat.feature_count} 功能 ·{' '}
                  {cat.avg_response_ms.toFixed(0)}ms
                </div>
              </Card>
            </Col>
          )
        })}
      </Row>
    )
  }

  // ===================== 热门功能柱状图 =====================
  const renderTopChart = () => {
    if (!topFeatures.length) return <Empty description="暂无数据" />
    const maxUses = Math.max(
      ...topFeatures.map((f: any) => f.total_uses || 0),
      1
    )
    return (
      <div>
        {topFeatures.slice(0, 15).map((f: any, i: number) => (
          <div
            key={f.feature_name}
            style={{
              display: 'flex',
              alignItems: 'center',
              marginBottom: 6,
              cursor: 'pointer'
            }}
            onClick={() => openTrend(f.feature_name)}
          >
            <span
              style={{
                width: 24,
                textAlign: 'center',
                fontWeight: 700,
                color: i < 3 ? '#f5222d' : '#999',
                fontSize: i < 3 ? 14 : 12
              }}
            >
              {i + 1}
            </span>
            <Tag
              color={getCategoryColor(f.feature_category)}
              style={{ margin: '0 6px', minWidth: 60, textAlign: 'center' }}
            >
              {f.feature_category}
            </Tag>
            <span style={{ width: 80, fontSize: 13 }}>{f.feature_name}</span>
            <div style={{ flex: 1, marginLeft: 8 }}>
              <div
                style={{
                  width: `${((f.total_uses || 0) / maxUses) * 100}%`,
                  height: 18,
                  background: `linear-gradient(90deg, ${getCategoryColor(f.feature_category)}88, ${getCategoryColor(f.feature_category)})`,
                  borderRadius: 3,
                  minWidth: 2,
                  transition: 'width 0.3s'
                }}
              />
            </div>
            <span
              style={{
                width: 60,
                textAlign: 'right',
                fontSize: 12,
                color: '#666'
              }}
            >
              {(f.total_uses || 0).toLocaleString()}
            </span>
            <span
              style={{
                width: 50,
                textAlign: 'right',
                fontSize: 12,
                color: '#999'
              }}
            >
              {f.unique_users || 0}人
            </span>
          </div>
        ))}
      </div>
    )
  }

  // ===================== 评分表格 =====================
  const scoreColumns = [
    {
      title: '排名',
      key: 'rank',
      width: 50,
      render: (_: any, __: any, i: number) => {
        const icons = [
          <TrophyOutlined style={{ color: '#faad14', fontSize: 16 }} />,
          <TrophyOutlined style={{ color: '#bfbfbf', fontSize: 14 }} />,
          <TrophyOutlined style={{ color: '#d48806', fontSize: 13 }} />
        ]
        return i < 3 ? icons[i] : <span style={{ color: '#999' }}>{i + 1}</span>
      }
    },
    {
      title: '功能',
      dataIndex: 'feature_name',
      key: 'feature_name',
      render: (name: string) => <a onClick={() => openTrend(name)}>{name}</a>
    },
    {
      title: '分类',
      dataIndex: 'feature_category',
      key: 'feature_category',
      width: 90,
      render: (cat: string) => <Tag color={getCategoryColor(cat)}>{cat}</Tag>
    },
    {
      title: '质量评分',
      dataIndex: 'quality_score',
      key: 'quality_score',
      width: 130,
      sorter: (a: FeatureScore, b: FeatureScore) =>
        a.quality_score - b.quality_score,
      defaultSortOrder: 'descend' as const,
      render: (score: number) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
          <Progress
            percent={Math.round(score)}
            size="small"
            strokeColor={scoreColor(score)}
            style={{ width: 80, margin: 0 }}
            format={p => <span style={{ fontSize: 12 }}>{p}</span>}
          />
        </div>
      )
    },
    {
      title: '7日使用',
      dataIndex: 'total_uses_7d',
      key: 'total_uses_7d',
      width: 90,
      sorter: (a: FeatureScore, b: FeatureScore) =>
        a.total_uses_7d - b.total_uses_7d,
      render: (v: number) => v.toLocaleString()
    },
    {
      title: '30日用户',
      dataIndex: 'unique_users_30d',
      key: 'unique_users_30d',
      width: 90,
      sorter: (a: FeatureScore, b: FeatureScore) =>
        a.unique_users_30d - b.unique_users_30d,
      render: (v: number) => v.toLocaleString()
    },
    {
      title: '渗透率',
      dataIndex: 'user_penetration',
      key: 'user_penetration',
      width: 80,
      sorter: (a: FeatureScore, b: FeatureScore) =>
        a.user_penetration - b.user_penetration,
      render: (v: number) => `${(v * 100).toFixed(1)}%`
    },
    {
      title: '环比增长',
      dataIndex: 'usage_growth',
      key: 'usage_growth',
      width: 90,
      sorter: (a: FeatureScore, b: FeatureScore) =>
        a.usage_growth - b.usage_growth,
      render: (v: number) => {
        const pct = (v * 100).toFixed(1)
        const isUp = v > 0
        return (
          <span
            style={{ color: isUp ? '#52c41a' : v < 0 ? '#f5222d' : '#999' }}
          >
            {isUp ? <RiseOutlined /> : v < 0 ? <FallOutlined /> : null}{' '}
            {isUp ? '+' : ''}
            {pct}%
          </span>
        )
      }
    },
    {
      title: '成功率',
      dataIndex: 'success_rate',
      key: 'success_rate',
      width: 75,
      render: (v: number) => (
        <span
          style={{
            color: v >= 0.95 ? '#52c41a' : v >= 0.8 ? '#faad14' : '#f5222d'
          }}
        >
          {(v * 100).toFixed(1)}%
        </span>
      )
    },
    {
      title: '均耗时',
      dataIndex: 'avg_response_ms',
      key: 'avg_response_ms',
      width: 75,
      render: (v: number) => (
        <span
          style={{
            color: v < 200 ? '#52c41a' : v < 500 ? '#faad14' : '#f5222d'
          }}
        >
          {v.toFixed(0)}ms
        </span>
      )
    }
  ]

  // ===================== 趋势弹窗 =====================
  const renderTrendModal = () => {
    const maxUses = Math.max(...trendData.map(d => d.total_uses), 1)
    return (
      <Modal
        title={`${trendName} — 使用趋势`}
        open={trendModal}
        onCancel={() => setTrendModal(false)}
        footer={null}
        width={720}
      >
        <div style={{ marginBottom: 12 }}>
          <Segmented
            options={[
              { label: '7天', value: 7 },
              { label: '14天', value: 14 },
              { label: '30天', value: 30 },
              { label: '90天', value: 90 }
            ]}
            value={trendDays}
            onChange={v => {
              setTrendDays(v as number)
              loadTrend(trendName, v as number)
            }}
          />
        </div>
        {trendLoading ? (
          <Spin style={{ display: 'block', padding: 40 }} />
        ) : trendData.length === 0 ? (
          <Empty description="暂无趋势数据" />
        ) : (
          <div>
            {/* 简易柱状图 */}
            <div
              style={{
                display: 'flex',
                alignItems: 'flex-end',
                height: 140,
                gap: 2,
                padding: '0 4px',
                borderBottom: '1px solid #f0f0f0'
              }}
            >
              {trendData.map(d => (
                <Tooltip
                  key={d.date}
                  title={
                    <div>
                      <div>{d.date}</div>
                      <div>使用 {d.total_uses} 次</div>
                      <div>用户 {d.unique_users} 人</div>
                      <div>耗时 {d.avg_response_ms.toFixed(0)}ms</div>
                      <div>成功率 {(d.success_rate * 100).toFixed(1)}%</div>
                    </div>
                  }
                >
                  <div
                    style={{
                      flex: 1,
                      height: `${(d.total_uses / maxUses) * 100}%`,
                      background: '#1890ff',
                      borderRadius: '3px 3px 0 0',
                      minHeight: 2,
                      cursor: 'pointer',
                      transition: 'height 0.3s'
                    }}
                  />
                </Tooltip>
              ))}
            </div>
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                fontSize: 11,
                color: '#999',
                marginTop: 4
              }}
            >
              <span>{trendData[0]?.date}</span>
              <span>{trendData[trendData.length - 1]?.date}</span>
            </div>
            {/* 汇总 */}
            <Row gutter={16} style={{ marginTop: 16 }}>
              <Col span={6}>
                <Card size="small">
                  <Tooltip title="选定时间范围内的总使用次数">
                    <Text type="secondary">总使用</Text>
                  </Tooltip>
                  <div style={{ fontSize: 20, fontWeight: 700 }}>
                    {trendData
                      .reduce((s, d) => s + d.total_uses, 0)
                      .toLocaleString()}
                  </div>
                </Card>
              </Col>
              <Col span={6}>
                <Card size="small">
                  <Text type="secondary">最高单日</Text>
                  <div style={{ fontSize: 20, fontWeight: 700 }}>
                    {Math.max(
                      ...trendData.map(d => d.total_uses)
                    ).toLocaleString()}
                  </div>
                </Card>
              </Col>
              <Col span={6}>
                <Card size="small">
                  <Text type="secondary">日均</Text>
                  <div style={{ fontSize: 20, fontWeight: 700 }}>
                    {trendData.length > 0
                      ? Math.round(
                          trendData.reduce((s, d) => s + d.total_uses, 0) /
                            trendData.length
                        ).toLocaleString()
                      : 0}
                  </div>
                </Card>
              </Col>
              <Col span={6}>
                <Card size="small">
                  <Text type="secondary">均耗时</Text>
                  <div style={{ fontSize: 20, fontWeight: 700 }}>
                    {trendData.length > 0
                      ? (
                          trendData.reduce((s, d) => s + d.avg_response_ms, 0) /
                          trendData.length
                        ).toFixed(0)
                      : 0}
                    <span style={{ fontSize: 12, fontWeight: 400 }}>ms</span>
                  </div>
                </Card>
              </Col>
            </Row>
          </div>
        )}
      </Modal>
    )
  }

  // ===================== 场景分布（私聊/群聊）=====================
  const sceneLabel = (s: string) => {
    if (s === 'private') return '私聊'
    if (s === 'group') return '群聊'
    return s || '未知'
  }
  const sceneColor = (s: string) => {
    if (s === 'private') return '#1890ff'
    if (s === 'group') return '#52c41a'
    return '#8c8c8c'
  }

  const renderSceneDistribution = () => {
    if (!sceneDist.length)
      return <Empty description="暂无场景数据" style={{ padding: 20 }} />
    const total = sceneDist.reduce((s, d) => s + d.total_uses, 0)
    return (
      <div>
        <Row gutter={[12, 12]} style={{ marginBottom: 12 }}>
          {sceneDist.map(d => {
            const pct =
              total > 0 ? ((d.total_uses / total) * 100).toFixed(1) : '0'
            return (
              <Col xs={12} key={d.scene}>
                <Card
                  size="small"
                  style={{ borderLeft: `4px solid ${sceneColor(d.scene)}` }}
                >
                  <Statistic
                    title={sceneLabel(d.scene)}
                    value={d.total_uses}
                    suffix={
                      <span style={{ fontSize: 12, color: '#999' }}>
                        次 ({pct}%)
                      </span>
                    }
                    valueStyle={{ fontSize: 22 }}
                  />
                  <div style={{ color: '#666', fontSize: 12, marginTop: 2 }}>
                    {d.unique_users} 用户
                  </div>
                </Card>
              </Col>
            )
          })}
        </Row>
        {/* 占比条 */}
        <div
          style={{
            display: 'flex',
            height: 20,
            borderRadius: 4,
            overflow: 'hidden'
          }}
        >
          {sceneDist.map(d => (
            <Tooltip
              key={d.scene}
              title={`${sceneLabel(d.scene)}: ${d.total_uses} 次 (${total > 0 ? ((d.total_uses / total) * 100).toFixed(1) : 0}%)`}
            >
              <div
                style={{
                  width: `${total > 0 ? (d.total_uses / total) * 100 : 0}%`,
                  background: sceneColor(d.scene),
                  transition: 'width 0.3s'
                }}
              />
            </Tooltip>
          ))}
        </div>
        <div
          style={{
            display: 'flex',
            gap: 16,
            marginTop: 8,
            fontSize: 12,
            color: '#666'
          }}
        >
          {sceneDist.map(d => (
            <span key={d.scene}>
              <span
                style={{
                  display: 'inline-block',
                  width: 10,
                  height: 10,
                  borderRadius: 2,
                  background: sceneColor(d.scene),
                  marginRight: 4
                }}
              />
              {sceneLabel(d.scene)}
            </span>
          ))}
        </div>
      </div>
    )
  }

  // ===================== 频道 TOP =====================
  const renderChannelTop = () => {
    if (!channelTop.length)
      return <Empty description="暂无频道数据" style={{ padding: 20 }} />
    const maxUses = Math.max(...channelTop.map(c => c.total_uses), 1)
    return (
      <div>
        {channelTop.slice(0, 15).map((ch, i) => (
          <div
            key={ch.channel_id}
            style={{ display: 'flex', alignItems: 'center', marginBottom: 5 }}
          >
            <span
              style={{
                width: 24,
                textAlign: 'center',
                fontWeight: 700,
                color: i < 3 ? '#f5222d' : '#999',
                fontSize: i < 3 ? 14 : 12
              }}
            >
              {i + 1}
            </span>
            <Tag
              color={sceneColor(ch.scene)}
              style={{ margin: '0 6px', minWidth: 40, textAlign: 'center' }}
            >
              {sceneLabel(ch.scene)}
            </Tag>
            <span
              style={{
                width: 100,
                fontSize: 12,
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap'
              }}
              title={ch.channel_id}
            >
              {ch.channel_id}
            </span>
            <div style={{ flex: 1, marginLeft: 8 }}>
              <div
                style={{
                  width: `${(ch.total_uses / maxUses) * 100}%`,
                  height: 16,
                  background: `linear-gradient(90deg, ${sceneColor(ch.scene)}88, ${sceneColor(ch.scene)})`,
                  borderRadius: 3,
                  minWidth: 2,
                  transition: 'width 0.3s'
                }}
              />
            </div>
            <span
              style={{
                width: 55,
                textAlign: 'right',
                fontSize: 12,
                color: '#666'
              }}
            >
              {ch.total_uses.toLocaleString()}
            </span>
            <span
              style={{
                width: 45,
                textAlign: 'right',
                fontSize: 12,
                color: '#999'
              }}
            >
              {ch.unique_users}人
            </span>
          </div>
        ))}
      </div>
    )
  }

  // ===================== 平台分布 =====================
  const renderPlatformDist = () => {
    if (!platformDist.length)
      return <Empty description="暂无平台数据" style={{ padding: 20 }} />
    const total = platformDist.reduce((s, d) => s + d.total_uses, 0)
    const platformColors = [
      '#1890ff',
      '#52c41a',
      '#faad14',
      '#f5222d',
      '#722ed1',
      '#13c2c2'
    ]
    return (
      <div>
        {platformDist.map((d, i) => {
          const pct =
            total > 0 ? ((d.total_uses / total) * 100).toFixed(1) : '0'
          return (
            <div
              key={d.platform}
              style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}
            >
              <Tag
                color={platformColors[i % platformColors.length]}
                style={{ minWidth: 60, textAlign: 'center' }}
              >
                {d.platform || '未知'}
              </Tag>
              <div style={{ flex: 1, marginLeft: 8 }}>
                <div
                  style={{
                    width: `${total > 0 ? (d.total_uses / total) * 100 : 0}%`,
                    height: 18,
                    background: platformColors[i % platformColors.length],
                    borderRadius: 3,
                    minWidth: 2,
                    transition: 'width 0.3s'
                  }}
                />
              </div>
              <span
                style={{
                  width: 80,
                  textAlign: 'right',
                  fontSize: 12,
                  color: '#666'
                }}
              >
                {d.total_uses.toLocaleString()} ({pct}%)
              </span>
              <span
                style={{
                  width: 45,
                  textAlign: 'right',
                  fontSize: 12,
                  color: '#999'
                }}
              >
                {d.unique_users}人
              </span>
            </div>
          )
        })}
      </div>
    )
  }

  // ===================== 场景功能 TOP =====================
  const renderSceneFeatureTop = () => {
    if (!sceneFeatureTop.length)
      return <Empty description="暂无数据" style={{ padding: 20 }} />
    const maxUses = Math.max(...sceneFeatureTop.map(f => f.total_uses), 1)
    return (
      <div>
        {sceneFeatureTop.slice(0, 15).map((f, i) => (
          <div
            key={`${f.scene}-${f.feature_name}`}
            style={{ display: 'flex', alignItems: 'center', marginBottom: 5 }}
          >
            <span
              style={{
                width: 24,
                textAlign: 'center',
                fontWeight: 700,
                color: i < 3 ? '#f5222d' : '#999',
                fontSize: i < 3 ? 14 : 12
              }}
            >
              {i + 1}
            </span>
            <Tag
              color={sceneColor(f.scene)}
              style={{
                margin: '0 4px',
                minWidth: 36,
                textAlign: 'center',
                fontSize: 11
              }}
            >
              {sceneLabel(f.scene)}
            </Tag>
            <Tag
              color={getCategoryColor(f.feature_category)}
              style={{
                margin: '0 4px',
                minWidth: 50,
                textAlign: 'center',
                fontSize: 11
              }}
            >
              {f.feature_category}
            </Tag>
            <span style={{ width: 70, fontSize: 12 }}>{f.feature_name}</span>
            <div style={{ flex: 1, marginLeft: 6 }}>
              <div
                style={{
                  width: `${(f.total_uses / maxUses) * 100}%`,
                  height: 16,
                  background: `linear-gradient(90deg, ${getCategoryColor(f.feature_category)}88, ${getCategoryColor(f.feature_category)})`,
                  borderRadius: 3,
                  minWidth: 2
                }}
              />
            </div>
            <span
              style={{
                width: 55,
                textAlign: 'right',
                fontSize: 12,
                color: '#666'
              }}
            >
              {f.total_uses.toLocaleString()}
            </span>
            <span
              style={{
                width: 40,
                textAlign: 'right',
                fontSize: 12,
                color: '#999'
              }}
            >
              {f.unique_users}人
            </span>
          </div>
        ))}
      </div>
    )
  }

  // ===================== 场景趋势图 =====================
  const renderSceneTrend = () => {
    if (!sceneTrend.length)
      return <Empty description="暂无趋势数据" style={{ padding: 20 }} />
    // 按日期分组
    const dateMap = new Map<string, Record<string, number>>()
    sceneTrend.forEach(d => {
      if (!dateMap.has(d.date)) dateMap.set(d.date, {})
      const m = dateMap.get(d.date)!
      m[d.scene] = d.total_uses
    })
    const dates = Array.from(dateMap.keys()).sort()
    const scenes = Array.from(new Set(sceneTrend.map(d => d.scene)))
    const maxVal = Math.max(...sceneTrend.map(d => d.total_uses), 1)

    return (
      <div>
        <div
          style={{
            display: 'flex',
            alignItems: 'flex-end',
            height: 120,
            gap: 1,
            borderBottom: '1px solid #f0f0f0'
          }}
        >
          {dates.map(date => {
            const m = dateMap.get(date)!
            return (
              <div
                key={date}
                style={{
                  flex: 1,
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  gap: 1
                }}
              >
                {scenes.map(scene => {
                  const val = m[scene] || 0
                  return (
                    <Tooltip
                      key={scene}
                      title={`${date} ${sceneLabel(scene)}: ${val}次`}
                    >
                      <div
                        style={{
                          width: '100%',
                          height: `${(val / maxVal) * 50}px`,
                          background: sceneColor(scene),
                          borderRadius: '2px 2px 0 0',
                          minHeight: val > 0 ? 2 : 0,
                          cursor: 'pointer'
                        }}
                      />
                    </Tooltip>
                  )
                })}
              </div>
            )
          })}
        </div>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            fontSize: 11,
            color: '#999',
            marginTop: 4
          }}
        >
          <span>{dates[0]}</span>
          <span>{dates[dates.length - 1]}</span>
        </div>
        <div
          style={{
            display: 'flex',
            gap: 16,
            marginTop: 8,
            fontSize: 12,
            color: '#666'
          }}
        >
          {scenes.map(s => (
            <span key={s}>
              <span
                style={{
                  display: 'inline-block',
                  width: 10,
                  height: 10,
                  borderRadius: 2,
                  background: sceneColor(s),
                  marginRight: 4
                }}
              />
              {sceneLabel(s)}
            </span>
          ))}
        </div>
      </div>
    )
  }

  // ===================== 主渲染 =====================
  return (
    <Spin spinning={loading}>
      <div style={{ padding: 24 }}>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: 16
          }}
        >
          <Title level={4} style={{ margin: 0 }}>
            <ThunderboltOutlined style={{ marginRight: 8 }} />
            功能分析
          </Title>
          <Button icon={<SyncOutlined />} onClick={handleRefresh}>
            重新统计
          </Button>
        </div>

        {/* 分类概览 */}
        <Card
          title="分类概览"
          size="small"
          style={{ marginBottom: 16 }}
          extra={
            filterCategory && (
              <Tag closable onClose={() => setFilterCategory('')}>
                筛选: {filterCategory}
              </Tag>
            )
          }
        >
          {renderCategoryCards()}
        </Card>

        <Row gutter={16}>
          {/* 热门功能 */}
          <Col xs={24} lg={10}>
            <Card
              title={
                <span>
                  <FireOutlined style={{ color: '#f5222d', marginRight: 6 }} />
                  热门功能 TOP 15
                </span>
              }
              size="small"
              style={{ marginBottom: 16 }}
              extra={
                <Select
                  value={topDays}
                  onChange={setTopDays}
                  size="small"
                  style={{ width: 90 }}
                  options={[
                    { label: '近7天', value: 7 },
                    { label: '近14天', value: 14 },
                    { label: '近30天', value: 30 }
                  ]}
                />
              }
            >
              {renderTopChart()}
            </Card>
          </Col>

          {/* 评分排行 */}
          <Col xs={24} lg={14}>
            <Card
              title={
                <span>
                  <TrophyOutlined
                    style={{ color: '#faad14', marginRight: 6 }}
                  />
                  功能质量评分
                </span>
              }
              size="small"
              style={{ marginBottom: 16 }}
              extra={
                <Select
                  value={filterCategory}
                  onChange={setFilterCategory}
                  allowClear
                  placeholder="全部分类"
                  size="small"
                  style={{ width: 120 }}
                  options={categories.map(c => ({
                    label: c.feature_category,
                    value: c.feature_category
                  }))}
                />
              }
            >
              <Table
                dataSource={scores}
                columns={scoreColumns}
                rowKey="feature_name"
                size="small"
                pagination={{
                  pageSize: 15,
                  size: 'small',
                  showSizeChanger: false
                }}
                scroll={{ x: 900 }}
              />
            </Card>
          </Col>
        </Row>

        {/* ==================== 场景/频道/平台维度分析 ==================== */}
        <Card
          title={
            <span>
              <NodeIndexOutlined style={{ marginRight: 8 }} />
              场景 · 频道 · 平台维度分析
            </span>
          }
          size="small"
          style={{ marginBottom: 16 }}
        >
          <Row gutter={16}>
            {/* 场景分布 */}
            <Col xs={24} md={8}>
              <Card
                title={
                  <span>
                    <MessageOutlined style={{ marginRight: 6 }} />
                    场景分布（30天）
                  </span>
                }
                size="small"
                type="inner"
                style={{ marginBottom: 16 }}
              >
                {renderSceneDistribution()}
              </Card>
            </Col>

            {/* 平台分布 */}
            <Col xs={24} md={8}>
              <Card
                title={
                  <span>
                    <GlobalOutlined style={{ marginRight: 6 }} />
                    平台分布（30天）
                  </span>
                }
                size="small"
                type="inner"
                style={{ marginBottom: 16 }}
              >
                {renderPlatformDist()}
              </Card>
            </Col>

            {/* 场景趋势 */}
            <Col xs={24} md={8}>
              <Card
                title={
                  <span>
                    <RiseOutlined style={{ marginRight: 6 }} />
                    场景使用趋势（30天）
                  </span>
                }
                size="small"
                type="inner"
                style={{ marginBottom: 16 }}
              >
                {renderSceneTrend()}
              </Card>
            </Col>
          </Row>

          <Row gutter={16}>
            {/* 频道 TOP */}
            <Col xs={24} md={12}>
              <Card
                title={
                  <span>
                    <TeamOutlined style={{ marginRight: 6 }} />
                    频道/群活跃排行
                  </span>
                }
                size="small"
                type="inner"
                style={{ marginBottom: 16 }}
                extra={
                  <Select
                    value={channelDays}
                    onChange={setChannelDays}
                    size="small"
                    style={{ width: 90 }}
                    options={[
                      { label: '近7天', value: 7 },
                      { label: '近14天', value: 14 },
                      { label: '近30天', value: 30 }
                    ]}
                  />
                }
              >
                {renderChannelTop()}
              </Card>
            </Col>

            {/* 场景功能 TOP */}
            <Col xs={24} md={12}>
              <Card
                title={
                  <span>
                    <FireOutlined style={{ marginRight: 6 }} />
                    各场景功能使用 TOP
                  </span>
                }
                size="small"
                type="inner"
                style={{ marginBottom: 16 }}
                extra={
                  <Select
                    value={sceneFilter}
                    onChange={setSceneFilter}
                    allowClear
                    placeholder="全部场景"
                    size="small"
                    style={{ width: 100 }}
                    options={[
                      { label: '私聊', value: 'private' },
                      { label: '群聊', value: 'group' }
                    ]}
                  />
                }
              >
                {renderSceneFeatureTop()}
              </Card>
            </Col>
          </Row>
        </Card>

        {renderTrendModal()}
      </div>
    </Spin>
  )
}

export default Features
