import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Statistic,
  Table,
  Tabs,
  Tag,
  Typography,
  message
} from 'antd'
import {
  ThunderboltOutlined,
  ReadOutlined,
  TrophyOutlined,
  FunnelPlotOutlined
} from '@ant-design/icons'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  LineChart,
  Line
} from 'recharts'
import { apiGetTowerStory, apiGetStageDifficulty } from '@/api'

const { Title } = Typography

const COLORS = [
  '#1890ff',
  '#52c41a',
  '#fa8c16',
  '#ff4d4f',
  '#722ed1',
  '#13c2c2',
  '#eb2f96',
  '#fadb14'
]

const TowerStory: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<any>(null)
  const [chapter, setChapter] = useState(0)
  const [stageData, setStageData] = useState<any[]>([])
  const [stageLoading, setStageLoading] = useState(false)

  useEffect(() => {
    loadData()
  }, [days])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetTowerStory(days)
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载通天塔/主线数据失败')
    } finally {
      setLoading(false)
    }
  }

  const loadStage = async (ch: number) => {
    setStageLoading(true)
    try {
      const r = await apiGetStageDifficulty(ch)
      if (r?.data) setStageData(r.data)
    } catch {
      message.error('加载关卡难度失败')
    } finally {
      setStageLoading(false)
    }
  }

  useEffect(() => {
    loadStage(chapter)
  }, [chapter])

  const towerStats = data?.tower_stats
  const towerLayers = data?.tower_layers || []
  const storyProgress = data?.story_progress || []
  const storyFunnel = data?.story_funnel || []

  const stageColumns = [
    { title: '关卡', dataIndex: 'stage_label', key: 'stage_label' },
    {
      title: '精英',
      dataIndex: 'is_elite',
      key: 'is_elite',
      render: (v: boolean) =>
        v ? <Tag color="red">精英</Tag> : <Tag>普通</Tag>
    },
    {
      title: '通关次数',
      dataIndex: 'clear_count',
      key: 'clear_count',
      sorter: (a: any, b: any) => a.clear_count - b.clear_count
    },
    {
      title: '通关人数',
      dataIndex: 'unique_users',
      key: 'unique_users',
      sorter: (a: any, b: any) => a.unique_users - b.unique_users
    },
    {
      title: '平均星级',
      dataIndex: 'avg_stars',
      key: 'avg_stars',
      sorter: (a: any, b: any) => a.avg_stars - b.avg_stars,
      render: (v: number) => {
        const s = (v ?? 0).toFixed(1)
        return (
          <Tag color={v >= 2.5 ? 'green' : v >= 1.5 ? 'orange' : 'red'}>
            {s} ★
          </Tag>
        )
      }
    }
  ]

  const tabItems = [
    {
      key: 'tower',
      label: (
        <span>
          <ThunderboltOutlined /> 通天塔
        </span>
      ),
      children: (
        <>
          {towerStats && (
            <Row gutter={16} style={{ marginBottom: 16 }}>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="总挑战次数"
                    value={towerStats.total_runs}
                    prefix={<ThunderboltOutlined />}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="挑战人数"
                    value={towerStats.unique_users}
                    prefix={<TrophyOutlined />}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="通关率"
                    value={towerStats.completion_rate?.toFixed(1)}
                    suffix="%"
                    valueStyle={{ color: '#52c41a' }}
                  />
                </Card>
              </Col>
              <Col span={6}>
                <Card>
                  <Statistic
                    title="平均最高层"
                    value={towerStats.avg_max_layer?.toFixed(1)}
                    valueStyle={{ color: '#1890ff' }}
                  />
                </Card>
              </Col>
            </Row>
          )}
          <Card title="层数分布">
            <ResponsiveContainer width="100%" height={350}>
              <BarChart data={towerLayers}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="current_layer" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="user_count" name="玩家人数" fill="#1890ff">
                  {towerLayers.map((_: any, i: number) => (
                    <Cell key={i} fill={COLORS[i % COLORS.length]} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </>
      )
    },
    {
      key: 'story',
      label: (
        <span>
          <ReadOutlined /> 主线剧情
        </span>
      ),
      children: (
        <>
          <Card title="各章节完成人数" style={{ marginBottom: 16 }}>
            <ResponsiveContainer width="100%" height={350}>
              <BarChart data={storyProgress}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="chapter_name"
                  angle={-15}
                  textAnchor="end"
                  height={60}
                />
                <YAxis />
                <Tooltip />
                <Bar dataKey="cleared_users" name="完成人数" fill="#52c41a">
                  {storyProgress.map((_: any, i: number) => (
                    <Cell key={i} fill={COLORS[i % COLORS.length]} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </Card>
          <Card title="剧情推进漏斗">
            <ResponsiveContainer width="100%" height={350}>
              <LineChart data={storyFunnel}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="chapter_name"
                  angle={-15}
                  textAnchor="end"
                  height={60}
                />
                <YAxis />
                <Tooltip />
                <Line
                  type="monotone"
                  dataKey="reached_users"
                  name="到达人数"
                  stroke="#1890ff"
                  strokeWidth={2}
                  dot={{ fill: '#1890ff', r: 4 }}
                />
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </>
      )
    },
    {
      key: 'difficulty',
      label: (
        <span>
          <FunnelPlotOutlined /> 关卡难度
        </span>
      ),
      children: (
        <Spin spinning={stageLoading}>
          <div style={{ marginBottom: 16 }}>
            <span>章节筛选：</span>
            <Select
              value={chapter}
              onChange={setChapter}
              style={{ width: 160 }}
            >
              <Select.Option value={0}>全部章节</Select.Option>
              {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(n => (
                <Select.Option key={n} value={n}>
                  第{n}章
                </Select.Option>
              ))}
            </Select>
          </div>
          <Table
            columns={stageColumns}
            dataSource={stageData}
            rowKey="stage_label"
            size="small"
            pagination={{ pageSize: 15 }}
          />
        </Spin>
      )
    }
  ]

  return (
    <div style={{ padding: 24 }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 16
        }}
      >
        <Title level={3} style={{ margin: 0 }}>
          通天塔与主线进度
        </Title>
        <Select value={days} onChange={setDays} style={{ width: 120 }}>
          <Select.Option value={7}>近7天</Select.Option>
          <Select.Option value={14}>近14天</Select.Option>
          <Select.Option value={30}>近30天</Select.Option>
          <Select.Option value={90}>近90天</Select.Option>
        </Select>
      </div>

      <Spin spinning={loading}>
        <Tabs items={tabItems} />
      </Spin>
    </div>
  )
}

export default TowerStory
