import React, { useEffect, useState } from 'react'
import {
  Table,
  Card,
  Select,
  Space,
  Tag,
  Button,
  Typography,
  Modal,
  Descriptions,
  InputNumber,
  message
} from 'antd'
import { ReloadOutlined, EyeOutlined } from '@ant-design/icons'
import {
  apiGetProfiles,
  apiGetProfile,
  apiRefreshProfile,
  apiRefreshAllProfiles,
  apiGetSnapshots
} from '@/api'
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
  Legend
} from 'recharts'

const { Title } = Typography

const lifecycleColors: Record<string, string> = {
  NEW: 'green',
  GROWING: 'blue',
  MATURE: 'purple',
  DECLINING: 'orange',
  LOST: 'red',
  RETURNED: 'cyan'
}

const lifecycleLabels: Record<string, string> = {
  NEW: '新手期（0-7天）',
  GROWING: '成长期',
  MATURE: '成熟期（稳定活跃）',
  DECLINING: '衰退期（活跃下降）',
  LOST: '流失期（长期未登录）',
  RETURNED: '回流期（流失后回归）'
}

const payTierColors: Record<string, string> = {
  FREE: 'default',
  MINNOW: 'green',
  DOLPHIN: 'blue',
  ORCA: 'purple',
  WHALE: 'gold',
  LEVIATHAN: 'red'
}

const payTierLabels: Record<string, string> = {
  FREE: '免费玩家',
  MINNOW: '微氪',
  DOLPHIN: '小氪',
  ORCA: '中氪',
  WHALE: '大氪',
  LEVIATHAN: '巨鲸'
}

const playStyleLabels: Record<string, string> = {
  COMBAT: '战斗型',
  CRAFT: '制造型',
  SOCIAL: '社交型',
  ECONOMY: '经济型',
  EXPLORER: '探索型',
  BALANCED: '均衡型'
}

const Profiles: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<any[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [size] = useState(20)
  const [filters, setFilters] = useState<Record<string, any>>({})
  const [detailVisible, setDetailVisible] = useState(false)
  const [detail, setDetail] = useState<any>(null)
  const [snapshots, setSnapshots] = useState<any[]>([])

  useEffect(() => {
    loadData()
  }, [page, filters])

  const loadData = async () => {
    setLoading(true)
    try {
      const res = await apiGetProfiles({ page, page_size: size, ...filters })
      if (res?.data) {
        setData(res.data.list || [])
        setTotal(res.data.total || 0)
      }
    } finally {
      setLoading(false)
    }
  }

  const showDetail = async (uid: number) => {
    try {
      const [profileRes, snapRes] = await Promise.all([
        apiGetProfile(uid),
        apiGetSnapshots(uid)
      ])
      if (profileRes?.data) {
        setDetail(profileRes.data)
        setDetailVisible(true)
      }
      setSnapshots(snapRes?.data || [])
    } catch {}
  }

  const handleRefresh = async (uid: number) => {
    try {
      await apiRefreshProfile(uid)
      message.success('画像已刷新')
      loadData()
    } catch {}
  }

  const handleRefreshAll = async () => {
    await apiRefreshAllProfiles()
    message.success('全量画像计算任务已启动')
  }

  const columns = [
    { title: 'UID', dataIndex: 'uid', key: 'uid', width: 80 },
    {
      title: '生命周期',
      dataIndex: 'lifecycle_stage',
      key: 'lifecycle_stage',
      render: (v: string) => (
        <Tag color={lifecycleColors[v] || 'default'}>
          {lifecycleLabels[v] || v}
        </Tag>
      )
    },
    {
      title: '付费等级',
      dataIndex: 'pay_tier',
      key: 'pay_tier',
      render: (v: string) => (
        <Tag color={payTierColors[v] || 'default'}>{payTierLabels[v] || v}</Tag>
      )
    },
    {
      title: '玩法偏好',
      dataIndex: 'play_style',
      key: 'play_style',
      render: (v: string) => playStyleLabels[v] || v
    },
    { title: '社交类型', dataIndex: 'social_type', key: 'social_type' },
    {
      title: '流失风险',
      dataIndex: 'churn_risk',
      key: 'churn_risk',
      sorter: true,
      render: (v: number) => (
        <Tag color={v >= 70 ? 'red' : v >= 40 ? 'orange' : 'green'}>{v}%</Tag>
      )
    },
    {
      title: 'LTV预测',
      dataIndex: 'ltv_predict',
      key: 'ltv_predict',
      render: (v: number) => `¥${(v || 0).toFixed(2)}`
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, row: any) => (
        <Space>
          <Button
            size="small"
            icon={<EyeOutlined />}
            onClick={() => showDetail(row.uid)}
          >
            详情
          </Button>
          <Button
            size="small"
            icon={<ReloadOutlined />}
            onClick={() => handleRefresh(row.uid)}
          >
            刷新
          </Button>
        </Space>
      )
    }
  ]

  return (
    <div>
      <Title level={3}>用户画像</Title>
      <Card style={{ marginBottom: 16 }}>
        <Space wrap>
          <Select
            placeholder="生命周期"
            allowClear
            style={{ width: 140 }}
            onChange={v => setFilters(f => ({ ...f, lifecycle_stage: v }))}
            options={[
              'NEW',
              'GROWING',
              'MATURE',
              'DECLINING',
              'LOST',
              'RETURNED'
            ].map(v => ({
              label: lifecycleLabels[v] || v,
              value: v
            }))}
          />
          <Select
            placeholder="付费等级"
            allowClear
            style={{ width: 140 }}
            onChange={v => setFilters(f => ({ ...f, pay_tier: v }))}
            options={[
              'FREE',
              'MINNOW',
              'DOLPHIN',
              'ORCA',
              'WHALE',
              'LEVIATHAN'
            ].map(v => ({
              label: payTierLabels[v] || v,
              value: v
            }))}
          />
          <Select
            placeholder="玩法偏好"
            allowClear
            style={{ width: 140 }}
            onChange={v => setFilters(f => ({ ...f, play_style: v }))}
            options={[
              'COMBAT',
              'CRAFT',
              'SOCIAL',
              'ECONOMY',
              'EXPLORER',
              'BALANCED'
            ].map(v => ({
              label: playStyleLabels[v] || v,
              value: v
            }))}
          />
          <InputNumber
            placeholder="最低流失风险%"
            min={0}
            max={100}
            onChange={v => setFilters(f => ({ ...f, min_churn_risk: v }))}
            style={{ width: 160 }}
          />
          <Button
            type="primary"
            onClick={handleRefreshAll}
            icon={<ReloadOutlined />}
          >
            全量重算
          </Button>
        </Space>
      </Card>

      <Table
        columns={columns}
        dataSource={data}
        rowKey="uid"
        loading={loading}
        pagination={{
          current: page,
          pageSize: size,
          total,
          onChange: p => setPage(p)
        }}
      />

      <Modal
        title={`用户画像详情 - UID: ${detail?.profile?.uid}`}
        open={detailVisible}
        onCancel={() => {
          setDetailVisible(false)
          setSnapshots([])
        }}
        footer={null}
        width={800}
      >
        {detail && (
          <div>
            <Descriptions bordered column={2} size="small">
              <Descriptions.Item label="用户名">
                {detail.game_user?.name || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="境界ID">
                {detail.game_user?.realm_id || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="生命周期">
                <Tag color={lifecycleColors[detail.profile?.lifecycle_stage]}>
                  {lifecycleLabels[detail.profile?.lifecycle_stage] ||
                    detail.profile?.lifecycle_stage}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="付费等级">
                <Tag color={payTierColors[detail.profile?.pay_tier]}>
                  {payTierLabels[detail.profile?.pay_tier] ||
                    detail.profile?.pay_tier}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="玩法偏好">
                {playStyleLabels[detail.profile?.play_style] ||
                  detail.profile?.play_style}
              </Descriptions.Item>
              <Descriptions.Item label="社交类型">
                {detail.profile?.social_type}
              </Descriptions.Item>
              <Descriptions.Item label="流失风险">
                <Tag color={detail.profile?.churn_risk >= 70 ? 'red' : 'green'}>
                  {detail.profile?.churn_risk}%
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="LTV预测">
                ¥{(detail.profile?.ltv_predict || 0).toFixed(2)}
              </Descriptions.Item>
              <Descriptions.Item label="卡关">
                {detail.profile?.stuck_flag ? '是' : '否'}
              </Descriptions.Item>
              <Descriptions.Item label="资源预警">
                {detail.profile?.resource_alert ? '是' : '否'}
              </Descriptions.Item>
              <Descriptions.Item label="灵石" span={2}>
                {detail.game_user?.spirit_stone?.toLocaleString() || 0}
              </Descriptions.Item>
            </Descriptions>

            {detail.tags?.length > 0 && (
              <Card title="用户标签" size="small" style={{ marginTop: 16 }}>
                <Space wrap>
                  {detail.tags.map((t: any, i: number) => (
                    <Tag key={i} color="blue">
                      {t.tag_group}/{t.tag_key}: {t.tag_value}
                    </Tag>
                  ))}
                </Space>
              </Card>
            )}

            {/* 快照趋势图表 */}
            {snapshots.length > 0 &&
              (() => {
                const chartData = [...snapshots].reverse().map(s => ({
                  date: s.snapshot_date?.slice(5, 10) || '',
                  登录: s.login_count || 0,
                  战斗: s.battle_count || 0,
                  制造: s.craft_count || 0,
                  社交: s.social_count || 0,
                  经济: s.economy_count || 0,
                  灵石收入: s.stone_income || 0,
                  灵石支出: s.stone_expense || 0
                }))

                return (
                  <>
                    {/* 行为趋势折线图 */}
                    <Card
                      title="行为趋势"
                      size="small"
                      style={{ marginTop: 16 }}
                    >
                      <ResponsiveContainer width="100%" height={260}>
                        <LineChart data={chartData}>
                          <CartesianGrid strokeDasharray="3 3" />
                          <XAxis dataKey="date" fontSize={12} />
                          <YAxis fontSize={12} />
                          <Tooltip />
                          <Legend />
                          <Line
                            type="monotone"
                            dataKey="登录"
                            stroke="#1677ff"
                            strokeWidth={2}
                            dot={{ r: 3 }}
                          />
                          <Line
                            type="monotone"
                            dataKey="战斗"
                            stroke="#faad14"
                            strokeWidth={2}
                            dot={{ r: 3 }}
                          />
                          <Line
                            type="monotone"
                            dataKey="制造"
                            stroke="#52c41a"
                            strokeWidth={2}
                            dot={{ r: 3 }}
                          />
                          <Line
                            type="monotone"
                            dataKey="社交"
                            stroke="#eb2f96"
                            strokeWidth={2}
                            dot={{ r: 3 }}
                          />
                          <Line
                            type="monotone"
                            dataKey="经济"
                            stroke="#722ed1"
                            strokeWidth={2}
                            dot={{ r: 3 }}
                          />
                        </LineChart>
                      </ResponsiveContainer>
                    </Card>

                    {/* 灵石收支柱状图 */}
                    <Card
                      title="灵石收支"
                      size="small"
                      style={{ marginTop: 16 }}
                    >
                      <ResponsiveContainer width="100%" height={220}>
                        <BarChart data={chartData}>
                          <CartesianGrid strokeDasharray="3 3" />
                          <XAxis dataKey="date" fontSize={12} />
                          <YAxis fontSize={12} />
                          <Tooltip
                            formatter={v =>
                              typeof v === 'number' ? v.toLocaleString() : v
                            }
                          />
                          <Legend />
                          <Bar
                            dataKey="灵石收入"
                            fill="#52c41a"
                            radius={[4, 4, 0, 0]}
                          />
                          <Bar
                            dataKey="灵石支出"
                            fill="#ff4d4f"
                            radius={[4, 4, 0, 0]}
                          />
                        </BarChart>
                      </ResponsiveContainer>
                    </Card>

                    {/* 保留数据表格 */}
                    <Card
                      title="快照明细"
                      size="small"
                      style={{ marginTop: 16 }}
                    >
                      <Table
                        dataSource={snapshots}
                        rowKey="id"
                        size="small"
                        pagination={false}
                        scroll={{ x: 600 }}
                        columns={[
                          {
                            title: '日期',
                            dataIndex: 'snapshot_date',
                            key: 'date',
                            width: 100,
                            render: (v: string) => v?.slice(0, 10)
                          },
                          {
                            title: '登录',
                            dataIndex: 'login_count',
                            key: 'login_count',
                            width: 60
                          },
                          {
                            title: '战斗',
                            dataIndex: 'battle_count',
                            key: 'battle_count',
                            width: 60
                          },
                          {
                            title: '制造',
                            dataIndex: 'craft_count',
                            key: 'craft_count',
                            width: 60
                          },
                          {
                            title: '社交',
                            dataIndex: 'social_count',
                            key: 'social_count',
                            width: 60
                          },
                          {
                            title: '经济',
                            dataIndex: 'economy_count',
                            key: 'economy_count',
                            width: 60
                          },
                          {
                            title: '灵石收入',
                            dataIndex: 'stone_income',
                            key: 'stone_income',
                            width: 90,
                            render: (v: number) => (v || 0).toLocaleString()
                          },
                          {
                            title: '灵石支出',
                            dataIndex: 'stone_expense',
                            key: 'stone_expense',
                            width: 90,
                            render: (v: number) => (v || 0).toLocaleString()
                          }
                        ]}
                      />
                    </Card>
                  </>
                )
              })()}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default Profiles
