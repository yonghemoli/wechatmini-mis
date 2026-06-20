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
  Progress,
  Empty
} from 'antd'
import {
  TeamOutlined,
  UserOutlined,
  RiseOutlined,
  FallOutlined,
  CrownOutlined,
  FireOutlined,
  UserAddOutlined,
  BarChartOutlined
} from '@ant-design/icons'
import {
  apiGetGameOverview,
  apiGetPlayerRanking,
  apiGetGuildRanking,
  apiGetNewUsersTrend,
  apiGetRealmStages,
  apiGetRealmChurn
} from '@/api'

const { Title } = Typography

// 大境界颜色
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

const Overview: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [overview, setOverview] = useState<any>(null)
  const [playerSort, setPlayerSort] = useState('realm')
  const [players, setPlayers] = useState<any[]>([])
  const [guildSort, setGuildSort] = useState('prestige')
  const [guilds, setGuilds] = useState<any[]>([])
  const [newUserDays, setNewUserDays] = useState(30)
  const [newUsers, setNewUsers] = useState<any[]>([])
  const [realmStages, setRealmStages] = useState<any[]>([])
  const [realmChurn, setRealmChurn] = useState<any[]>([])
  const [churnDays, setChurnDays] = useState(7)

  useEffect(() => {
    loadAll()
  }, [])

  useEffect(() => {
    loadPlayers()
  }, [playerSort])
  useEffect(() => {
    loadGuilds()
  }, [guildSort])
  useEffect(() => {
    loadNewUsers()
  }, [newUserDays])
  useEffect(() => {
    loadChurn()
  }, [churnDays])

  const loadAll = async () => {
    setLoading(true)
    try {
      await Promise.all([
        loadOverview(),
        loadPlayers(),
        loadGuilds(),
        loadNewUsers(),
        loadRealmStages(),
        loadChurn()
      ])
    } finally {
      setLoading(false)
    }
  }

  const loadOverview = async () => {
    try {
      const r = await apiGetGameOverview()
      if (r?.data) setOverview(r.data)
    } catch {}
  }
  const loadPlayers = async () => {
    try {
      const r = await apiGetPlayerRanking(playerSort, 30)
      if (r?.data) setPlayers(r.data || [])
    } catch {}
  }
  const loadGuilds = async () => {
    try {
      const r = await apiGetGuildRanking(guildSort, 20)
      if (r?.data) setGuilds(r.data || [])
    } catch {}
  }
  const loadNewUsers = async () => {
    try {
      const r = await apiGetNewUsersTrend(newUserDays)
      if (r?.data) setNewUsers(r.data || [])
    } catch {}
  }
  const loadRealmStages = async () => {
    try {
      const r = await apiGetRealmStages()
      if (r?.data) setRealmStages(r.data || [])
    } catch {}
  }
  const loadChurn = async () => {
    try {
      const r = await apiGetRealmChurn(churnDays)
      if (r?.data) setRealmChurn(r.data || [])
    } catch {}
  }

  // 新增玩家柱状图
  const maxNew = Math.max(...newUsers.map((d: any) => d.count), 1)

  // 大境界分布饼状图（纯CSS）
  const totalStage = realmStages.reduce((s: number, d: any) => s + d.count, 0)

  return (
    <Spin spinning={loading}>
      <div style={{ padding: 24 }}>
        <Title level={4}>
          <BarChartOutlined style={{ marginRight: 8 }} />
          游戏总览
        </Title>

        {/* 核心指标卡片 */}
        <Row gutter={[12, 12]}>
          <Col xs={12} sm={6} md={4}>
            <Card size="small">
              <Statistic
                title="总用户"
                value={overview?.total_users || 0}
                prefix={<UserOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={6} md={4}>
            <Card size="small">
              <Statistic
                title="DAU"
                value={overview?.dau || 0}
                prefix={<FireOutlined />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={6} md={4}>
            <Card size="small">
              <Statistic
                title="MAU"
                value={overview?.mau || 0}
                prefix={<TeamOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={6} md={4}>
            <Card size="small">
              <Statistic
                title="DAU/MAU"
                value={overview?.dau_mau_ratio || 0}
                suffix="%"
                precision={1}
                valueStyle={{
                  color:
                    (overview?.dau_mau_ratio || 0) > 20 ? '#52c41a' : '#faad14'
                }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={6} md={4}>
            <Card size="small">
              <Statistic
                title="今日新增"
                value={overview?.new_users_1d || 0}
                prefix={<UserAddOutlined />}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={6} md={4}>
            <Card size="small">
              <Statistic
                title="7日新增"
                value={overview?.new_users_7d || 0}
                prefix={<RiseOutlined />}
              />
            </Card>
          </Col>
        </Row>

        <Row gutter={16} style={{ marginTop: 16 }}>
          {/* 新增玩家趋势 */}
          <Col xs={24} lg={14}>
            <Card
              title={
                <span>
                  <UserAddOutlined
                    style={{ color: '#52c41a', marginRight: 6 }}
                  />
                  新增玩家趋势
                </span>
              }
              size="small"
              extra={
                <Select
                  value={newUserDays}
                  onChange={setNewUserDays}
                  size="small"
                  style={{ width: 90 }}
                  options={[
                    { label: '7天', value: 7 },
                    { label: '14天', value: 14 },
                    { label: '30天', value: 30 },
                    { label: '90天', value: 90 }
                  ]}
                />
              }
            >
              {newUsers.length > 0 ? (
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'flex-end',
                    height: 150,
                    gap: 2,
                    padding: '0 4px',
                    borderBottom: '1px solid #f0f0f0'
                  }}
                >
                  {newUsers.map((d: any) => (
                    <Tooltip key={d.date} title={`${d.date}: ${d.count} 人`}>
                      <div
                        style={{
                          flex: 1,
                          height: `${(d.count / maxNew) * 100}%`,
                          background:
                            'linear-gradient(180deg, #52c41a 0%, #95de64 100%)',
                          borderRadius: '3px 3px 0 0',
                          minHeight: 2,
                          transition: 'height 0.3s'
                        }}
                      />
                    </Tooltip>
                  ))}
                </div>
              ) : (
                <Empty description="暂无数据" />
              )}
              {newUsers.length > 0 && (
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    fontSize: 11,
                    color: '#999',
                    marginTop: 4
                  }}
                >
                  <span>{newUsers[0]?.date}</span>
                  <span>
                    共 {newUsers.reduce((s: number, d: any) => s + d.count, 0)}{' '}
                    人
                  </span>
                  <span>{newUsers[newUsers.length - 1]?.date}</span>
                </div>
              )}
            </Card>
          </Col>

          {/* 大境界分布 */}
          <Col xs={24} lg={10}>
            <Card title="大境界分布" size="small">
              {realmStages.length > 0 ? (
                realmStages.map((s: any) => {
                  const pct = totalStage > 0 ? (s.count / totalStage) * 100 : 0
                  return (
                    <div
                      key={s.stage_level}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        marginBottom: 6
                      }}
                    >
                      <Tag
                        color={stageColors[s.stage_name] || '#595959'}
                        style={{ width: 50, textAlign: 'center', margin: 0 }}
                      >
                        {s.stage_name}
                      </Tag>
                      <div style={{ flex: 1, margin: '0 8px' }}>
                        <div
                          style={{
                            width: `${pct}%`,
                            height: 20,
                            minWidth: 2,
                            background: stageColors[s.stage_name] || '#595959',
                            borderRadius: 3,
                            transition: 'width 0.3s',
                            opacity: 0.8
                          }}
                        />
                      </div>
                      <span
                        style={{
                          width: 50,
                          textAlign: 'right',
                          fontWeight: 600
                        }}
                      >
                        {s.count}
                      </span>
                      <span
                        style={{
                          width: 55,
                          textAlign: 'right',
                          color: '#999',
                          fontSize: 12
                        }}
                      >
                        {pct.toFixed(1)}%
                      </span>
                    </div>
                  )
                })
              ) : (
                <Empty description="暂无数据" />
              )}
            </Card>
          </Col>
        </Row>

        {/* 境界流失率 */}
        <Card
          title={
            <span>
              <FallOutlined style={{ color: '#f5222d', marginRight: 6 }} />
              境界流失率分析
            </span>
          }
          size="small"
          style={{ marginTop: 16 }}
          extra={
            <Select
              value={churnDays}
              onChange={setChurnDays}
              size="small"
              style={{ width: 120 }}
              options={[
                { label: '7天未活跃', value: 7 },
                { label: '14天未活跃', value: 14 },
                { label: '30天未活跃', value: 30 }
              ]}
            />
          }
        >
          <Table
            dataSource={realmChurn}
            rowKey="realm_id"
            size="small"
            pagination={false}
            scroll={{ x: 600 }}
            columns={[
              {
                title: '境界',
                dataIndex: 'realm_name',
                key: 'realm_name',
                width: 100,
                render: (v: string) => (
                  <Tag color={stageColors[v] || '#595959'}>{v}</Tag>
                )
              },
              {
                title: '总人数',
                dataIndex: 'total',
                key: 'total',
                width: 80,
                sorter: (a: any, b: any) => a.total - b.total
              },
              {
                title: '流失人数',
                dataIndex: 'churned',
                key: 'churned',
                width: 80,
                render: (v: number) => (
                  <span style={{ color: v > 0 ? '#f5222d' : '#52c41a' }}>
                    {v}
                  </span>
                )
              },
              {
                title: '流失率',
                dataIndex: 'churn_rate',
                key: 'churn_rate',
                width: 150,
                sorter: (a: any, b: any) => a.churn_rate - b.churn_rate,
                render: (v: number) => {
                  const pct = v * 100
                  return (
                    <div
                      style={{ display: 'flex', alignItems: 'center', gap: 8 }}
                    >
                      <Progress
                        percent={Math.round(pct)}
                        size="small"
                        style={{ width: 80, margin: 0 }}
                        strokeColor={
                          pct > 50
                            ? '#f5222d'
                            : pct > 30
                              ? '#faad14'
                              : '#52c41a'
                        }
                        format={p => `${p}%`}
                      />
                      {pct > 40 && <Tag color="red">瓶颈</Tag>}
                    </div>
                  )
                }
              }
            ]}
          />
        </Card>

        <Row gutter={16} style={{ marginTop: 16 }}>
          <Col xs={24} lg={12}>
            <Card
              title={
                <span>
                  <CrownOutlined style={{ color: '#faad14', marginRight: 6 }} />
                  玩家排行
                </span>
              }
              size="small"
              extra={
                <Select
                  value={playerSort}
                  onChange={setPlayerSort}
                  size="small"
                  style={{ width: 100 }}
                  options={[
                    { label: '境界', value: 'realm' },
                    { label: '充值', value: 'recharge' },
                    { label: '财富', value: 'wealth' },
                    { label: 'VIP', value: 'vip' }
                  ]}
                />
              }
            >
              <Table
                dataSource={players}
                rowKey="id"
                size="small"
                pagination={{
                  pageSize: 10,
                  size: 'small',
                  showSizeChanger: false
                }}
                columns={[
                  {
                    title: '#',
                    key: 'rank',
                    width: 40,
                    render: (_: any, __: any, i: number) => {
                      const icons = [
                        <CrownOutlined
                          style={{ color: '#faad14', fontSize: 16 }}
                        />,
                        <CrownOutlined
                          style={{ color: '#bfbfbf', fontSize: 14 }}
                        />,
                        <CrownOutlined
                          style={{ color: '#d48806', fontSize: 13 }}
                        />
                      ]
                      return i < 3 ? (
                        icons[i]
                      ) : (
                        <span style={{ color: '#999' }}>{i + 1}</span>
                      )
                    }
                  },
                  {
                    title: '名称',
                    dataIndex: 'name',
                    key: 'name',
                    ellipsis: true
                  },
                  {
                    title: '境界',
                    dataIndex: 'realm_name',
                    key: 'realm_name',
                    width: 100,
                    render: (v: string) => <Tag>{v}</Tag>
                  },
                  {
                    title: 'VIP',
                    dataIndex: 'vip_level',
                    key: 'vip_level',
                    width: 55,
                    render: (v: number) =>
                      v > 0 ? (
                        <Tag color="gold">V{v}</Tag>
                      ) : (
                        <span style={{ color: '#ccc' }}>-</span>
                      )
                  },
                  {
                    title: '灵石',
                    dataIndex: 'spirit_stone',
                    key: 'spirit_stone',
                    width: 90,
                    render: (v: number) => v?.toLocaleString()
                  },
                  {
                    title: '充值(元)',
                    dataIndex: 'total_recharge',
                    key: 'total_recharge',
                    width: 90,
                    render: (v: number) => (v / 100).toFixed(2)
                  }
                ]}
              />
            </Card>
          </Col>

          {/* 宗门排行 */}
          <Col xs={24} lg={12}>
            <Card
              title={
                <span>
                  <TeamOutlined style={{ color: '#722ed1', marginRight: 6 }} />
                  宗门排行
                </span>
              }
              size="small"
              extra={
                <Select
                  value={guildSort}
                  onChange={setGuildSort}
                  size="small"
                  style={{ width: 100 }}
                  options={[
                    { label: '声望', value: 'prestige' },
                    { label: '人数', value: 'members' },
                    { label: '财富', value: 'wealth' }
                  ]}
                />
              }
            >
              <Table
                dataSource={guilds}
                rowKey="id"
                size="small"
                pagination={{
                  pageSize: 10,
                  size: 'small',
                  showSizeChanger: false
                }}
                columns={[
                  {
                    title: '#',
                    key: 'rank',
                    width: 40,
                    render: (_: any, __: any, i: number) => {
                      const icons = [
                        <CrownOutlined
                          style={{ color: '#faad14', fontSize: 16 }}
                        />,
                        <CrownOutlined
                          style={{ color: '#bfbfbf', fontSize: 14 }}
                        />,
                        <CrownOutlined
                          style={{ color: '#d48806', fontSize: 13 }}
                        />
                      ]
                      return i < 3 ? (
                        icons[i]
                      ) : (
                        <span style={{ color: '#999' }}>{i + 1}</span>
                      )
                    }
                  },
                  {
                    title: '宗门',
                    dataIndex: 'name',
                    key: 'name',
                    ellipsis: true
                  },
                  {
                    title: '声望',
                    dataIndex: 'prestige',
                    key: 'prestige',
                    width: 80,
                    render: (v: number) => v?.toLocaleString()
                  },
                  {
                    title: '成员',
                    dataIndex: 'member_count',
                    key: 'member_count',
                    width: 60
                  },
                  {
                    title: '灵石',
                    dataIndex: 'spirit_stone',
                    key: 'spirit_stone',
                    width: 90,
                    render: (v: number) => v?.toLocaleString()
                  },
                  {
                    title: '等级',
                    dataIndex: 'level_id',
                    key: 'level_id',
                    width: 55,
                    render: (v: number) => <Tag color="purple">Lv.{v}</Tag>
                  }
                ]}
              />
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  )
}

export default Overview
