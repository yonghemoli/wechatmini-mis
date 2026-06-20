import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Select,
  Spin,
  Statistic,
  Table,
  Tag,
  Tooltip,
  Typography,
  Empty,
  message
} from 'antd'
import {
  ExperimentOutlined,
  TrophyOutlined,
  TeamOutlined
} from '@ant-design/icons'
import { apiGetDungeonAnalysis, apiGetDungeonRealmDist } from '@/api'

const { Title } = Typography

const dungeonTypeLabels: Record<number, { text: string; color: string }> = {
  1: { text: '普通秘境', color: 'blue' },
  2: { text: '宗门秘境', color: 'purple' }
}

const DungeonAnalysis: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<any>(null)
  const [realmDist, setRealmDist] = useState<any[]>([])

  useEffect(() => {
    loadAll()
  }, [days])

  const loadAll = async () => {
    setLoading(true)
    try {
      const [r1, r2] = await Promise.all([
        apiGetDungeonAnalysis(days),
        apiGetDungeonRealmDist(days)
      ])
      if (r1?.data) setData(r1.data)
      if (r2?.data) setRealmDist(r2.data || [])
    } catch {
      message.error('加载副本数据失败')
    } finally {
      setLoading(false)
    }
  }

  const stats = data?.stats || []
  const dailyTrend = data?.daily_trend || []

  const totalClears = stats.reduce(
    (s: number, d: any) => s + (d.clear_count || 0),
    0
  )
  const totalUsers = new Set(stats.map((d: any) => d.dungeon_id)).size

  const statsColumns = [
    {
      title: '副本',
      dataIndex: 'dungeon_name',
      key: 'dungeon_name',
      render: (v: string, r: any) => (
        <span>
          {v}
          <Tag
            color={dungeonTypeLabels[r.dungeon_type]?.color || 'default'}
            style={{ marginLeft: 6 }}
          >
            {dungeonTypeLabels[r.dungeon_type]?.text || '未知'}
          </Tag>
        </span>
      )
    },
    {
      title: '通关次数',
      dataIndex: 'clear_count',
      key: 'clear_count',
      render: (v: number) => v?.toLocaleString(),
      sorter: (a: any, b: any) => a.clear_count - b.clear_count,
      defaultSortOrder: 'descend' as const
    },
    {
      title: '参与人数',
      dataIndex: 'unique_users',
      key: 'unique_users',
      sorter: (a: any, b: any) => a.unique_users - b.unique_users
    },
    {
      title: '人均通关',
      dataIndex: 'avg_clears',
      key: 'avg_clears',
      render: (v: number) => v?.toFixed(1),
      sorter: (a: any, b: any) => a.avg_clears - b.avg_clears
    }
  ]

  const maxDailyClears = Math.max(
    ...dailyTrend.map((d: any) => d.clear_count || 0),
    1
  )

  // 按副本分组境界分布
  const dungeonNames = [...new Set(realmDist.map((d: any) => d.dungeon_name))]

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
            <ExperimentOutlined style={{ marginRight: 8 }} />
            副本分析
          </Title>
          <Select
            value={days}
            onChange={setDays}
            style={{ width: 100 }}
            options={[
              { label: '7天', value: 7 },
              { label: '14天', value: 14 },
              { label: '30天', value: 30 },
              { label: '90天', value: 90 }
            ]}
          />
        </div>

        {/* 汇总指标 */}
        <Row gutter={[12, 12]} style={{ marginBottom: 16 }}>
          <Col xs={12} sm={8} md={6}>
            <Card size="small">
              <Statistic
                title="总通关次数"
                value={totalClears}
                prefix={<TrophyOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card size="small">
              <Statistic
                title="活跃副本数"
                value={totalUsers}
                prefix={<ExperimentOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card size="small">
              <Statistic
                title="日均通关"
                value={
                  dailyTrend.length > 0
                    ? Math.round(totalClears / dailyTrend.length)
                    : 0
                }
                prefix={<TeamOutlined />}
              />
            </Card>
          </Col>
        </Row>

        <Row gutter={16}>
          {/* 副本排行 */}
          <Col xs={24} lg={14}>
            <Card title="副本通关排行" size="small">
              <Table
                dataSource={stats}
                columns={statsColumns}
                rowKey="dungeon_id"
                size="small"
                pagination={{ pageSize: 10 }}
              />
            </Card>
          </Col>

          {/* 日趋势 */}
          <Col xs={24} lg={10}>
            <Card title="副本活跃日趋势" size="small">
              {dailyTrend.length > 0 ? (
                <>
                  <div
                    style={{
                      display: 'flex',
                      alignItems: 'flex-end',
                      height: 160,
                      gap: 2,
                      padding: '0 4px',
                      borderBottom: '1px solid #f0f0f0'
                    }}
                  >
                    {dailyTrend.map((d: any) => (
                      <Tooltip
                        key={d.date}
                        title={
                          <div>
                            <div>{d.date}</div>
                            <div>通关 {d.clear_count} 次</div>
                            <div>玩家 {d.user_count} 人</div>
                          </div>
                        }
                      >
                        <div
                          style={{
                            flex: 1,
                            height: `${((d.clear_count || 0) / maxDailyClears) * 100}%`,
                            background:
                              'linear-gradient(180deg, #722ed1 0%, #b37feb 100%)',
                            borderRadius: '3px 3px 0 0',
                            minHeight: 2,
                            cursor: 'pointer'
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
                    <span>{dailyTrend[0]?.date}</span>
                    <span>{dailyTrend[dailyTrend.length - 1]?.date}</span>
                  </div>
                </>
              ) : (
                <Empty description="暂无趋势数据" />
              )}
            </Card>
          </Col>
        </Row>

        {/* 境界参与分布 */}
        {realmDist.length > 0 && (
          <Card
            title="副本参与者境界分布"
            size="small"
            style={{ marginTop: 16 }}
          >
            <div style={{ overflowX: 'auto' }}>
              <table
                style={{
                  width: '100%',
                  borderCollapse: 'collapse',
                  fontSize: 13
                }}
              >
                <thead>
                  <tr style={{ borderBottom: '2px solid #f0f0f0' }}>
                    <th style={{ textAlign: 'left', padding: '8px 12px' }}>
                      副本
                    </th>
                    {[...new Set(realmDist.map(d => d.stage_name))].map(
                      name => (
                        <th
                          key={name}
                          style={{ textAlign: 'center', padding: '8px 6px' }}
                        >
                          {name}
                        </th>
                      )
                    )}
                  </tr>
                </thead>
                <tbody>
                  {dungeonNames.map(dName => {
                    const dungeonData = realmDist.filter(
                      d => d.dungeon_name === dName
                    )
                    const maxCount = Math.max(
                      ...dungeonData.map(d => d.user_count || 0),
                      1
                    )
                    return (
                      <tr
                        key={dName}
                        style={{ borderBottom: '1px solid #f5f5f5' }}
                      >
                        <td style={{ padding: '6px 12px', fontWeight: 500 }}>
                          {dName}
                        </td>
                        {[...new Set(realmDist.map(d => d.stage_name))].map(
                          stageName => {
                            const cell = dungeonData.find(
                              d => d.stage_name === stageName
                            )
                            if (!cell)
                              return (
                                <td
                                  key={stageName}
                                  style={{
                                    textAlign: 'center',
                                    padding: '6px'
                                  }}
                                >
                                  -
                                </td>
                              )
                            const intensity = cell.user_count / maxCount
                            return (
                              <td
                                key={stageName}
                                style={{
                                  textAlign: 'center',
                                  padding: '6px',
                                  background: `rgba(114, 46, 209, ${intensity * 0.3})`
                                }}
                              >
                                <Tooltip
                                  title={`${cell.user_count} 人, ${cell.clear_count} 次`}
                                >
                                  <span>{cell.user_count}</span>
                                </Tooltip>
                              </td>
                            )
                          }
                        )}
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          </Card>
        )}
      </div>
    </Spin>
  )
}

export default DungeonAnalysis
