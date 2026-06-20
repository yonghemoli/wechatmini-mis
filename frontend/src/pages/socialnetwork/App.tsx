import React, { useEffect, useState } from 'react'
import {
  Card,
  Col,
  Row,
  Spin,
  Statistic,
  Tag,
  Typography,
  Empty,
  Progress,
  message
} from 'antd'
import {
  HeartOutlined,
  UserOutlined,
  TeamOutlined,
  SmileOutlined
} from '@ant-design/icons'
import { apiGetSocialNetwork } from '@/api'

const { Title, Text } = Typography

const relationColors: Record<string, string> = {
  道友: '#1890ff',
  道侣: '#ff4d4f',
  师父: '#faad14',
  徒弟: '#52c41a'
}

const SocialNetwork: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const r = await apiGetSocialNetwork()
      if (r?.data) setData(r.data)
    } catch {
      message.error('加载社交数据失败')
    } finally {
      setLoading(false)
    }
  }

  const overview = data?.overview
  const relations = data?.relation_dist || []
  const intimacy = data?.intimacy_dist || []

  const totalRelations = relations.reduce(
    (s: number, r: any) => s + (r.count || 0),
    0
  )

  const intimacyOrder = ['0-10', '10-50', '50-100', '100-200', '200+']
  const intimacyColors = ['#d9d9d9', '#91d5ff', '#69c0ff', '#40a9ff', '#1890ff']

  return (
    <Spin spinning={loading}>
      <div style={{ padding: 24 }}>
        <Title level={4}>
          <HeartOutlined style={{ marginRight: 8, color: '#ff4d4f' }} />
          社交网络分析
        </Title>

        {/* 整体概览 */}
        <Row gutter={[12, 12]} style={{ marginBottom: 16 }}>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="好友关系对"
                value={overview?.total_friendships || 0}
                prefix={<TeamOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="有好友用户"
                value={overview?.users_with_friends || 0}
                prefix={<UserOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="好友覆盖率"
                value={overview?.friend_coverage || 0}
                suffix="%"
                precision={1}
                valueStyle={{
                  color:
                    (overview?.friend_coverage || 0) > 50
                      ? '#52c41a'
                      : '#faad14'
                }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="人均好友"
                value={overview?.avg_friends || 0}
                precision={1}
                prefix={<SmileOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="最多好友"
                value={overview?.max_friends || 0}
                prefix={<UserOutlined />}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8} md={4}>
            <Card size="small">
              <Statistic
                title="道侣对数"
                value={overview?.companion_pairs || 0}
                prefix={<HeartOutlined />}
                valueStyle={{ color: '#ff4d4f' }}
              />
            </Card>
          </Col>
        </Row>

        <Row gutter={16}>
          {/* 关系类型分布 */}
          <Col xs={24} lg={12}>
            <Card title="关系类型分布" size="small">
              {relations.length > 0 ? (
                <div>
                  {relations.map((r: any) => {
                    const pct =
                      totalRelations > 0 ? (r.count / totalRelations) * 100 : 0
                    const color = relationColors[r.type_name] || '#1890ff'
                    return (
                      <div key={r.relation_type} style={{ marginBottom: 16 }}>
                        <div
                          style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            marginBottom: 4
                          }}
                        >
                          <span>
                            <Tag color={color}>{r.type_name}</Tag>
                          </span>
                          <Text type="secondary">
                            {r.count?.toLocaleString()} ({pct.toFixed(1)}%)
                          </Text>
                        </div>
                        <Progress
                          percent={pct}
                          showInfo={false}
                          strokeColor={color}
                          size="small"
                        />
                      </div>
                    )
                  })}
                  {/* 可视化饼图效果 */}
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'center',
                      marginTop: 16
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        gap: 16,
                        flexWrap: 'wrap',
                        justifyContent: 'center'
                      }}
                    >
                      {relations.map((r: any) => {
                        const pct =
                          totalRelations > 0
                            ? (r.count / totalRelations) * 100
                            : 0
                        const color = relationColors[r.type_name] || '#1890ff'
                        const size = Math.max(40, Math.min(100, pct * 1.5))
                        return (
                          <div
                            key={r.relation_type}
                            style={{ textAlign: 'center' }}
                          >
                            <div
                              style={{
                                width: size,
                                height: size,
                                borderRadius: '50%',
                                background: color,
                                opacity: 0.8,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                color: '#fff',
                                fontWeight: 700,
                                fontSize: size > 60 ? 14 : 11
                              }}
                            >
                              {pct.toFixed(0)}%
                            </div>
                            <Text
                              style={{
                                fontSize: 11,
                                display: 'block',
                                marginTop: 4
                              }}
                            >
                              {r.type_name}
                            </Text>
                          </div>
                        )
                      })}
                    </div>
                  </div>
                </div>
              ) : (
                <Empty description="暂无关系数据" />
              )}
            </Card>
          </Col>

          {/* 亲密度分布 */}
          <Col xs={24} lg={12}>
            <Card title="亲密度分布" size="small">
              {intimacy.length > 0 ? (
                <div>
                  {intimacyOrder.map((bracket, idx) => {
                    const item = intimacy.find(
                      (d: any) => d.bracket === bracket
                    )
                    if (!item) return null
                    const totalIntimacy = intimacy.reduce(
                      (s: number, d: any) => s + (d.count || 0),
                      0
                    )
                    const pct =
                      totalIntimacy > 0 ? (item.count / totalIntimacy) * 100 : 0
                    return (
                      <div key={bracket} style={{ marginBottom: 16 }}>
                        <div
                          style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            marginBottom: 4
                          }}
                        >
                          <Text strong>亲密度 {bracket}</Text>
                          <Text type="secondary">
                            {item.count?.toLocaleString()} 对 ({pct.toFixed(1)}
                            %)
                          </Text>
                        </div>
                        <Progress
                          percent={pct}
                          showInfo={false}
                          strokeColor={intimacyColors[idx]}
                          size="small"
                        />
                      </div>
                    )
                  })}

                  {/* 社交活跃度指标 */}
                  <Card
                    size="small"
                    style={{ marginTop: 16, background: '#fafafa' }}
                  >
                    <Title level={5} style={{ marginBottom: 8 }}>
                      社交活跃度评估
                    </Title>
                    <Row gutter={12}>
                      <Col span={12}>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          {'高亲密度(≥100)比例'}
                        </Text>
                        <div
                          style={{
                            fontSize: 20,
                            fontWeight: 700,
                            color: '#1890ff'
                          }}
                        >
                          {(() => {
                            const totalI = intimacy.reduce(
                              (s: number, d: any) => s + (d.count || 0),
                              0
                            )
                            const high = intimacy
                              .filter((d: any) =>
                                ['100-200', '200+'].includes(d.bracket)
                              )
                              .reduce(
                                (s: number, d: any) => s + (d.count || 0),
                                0
                              )
                            return totalI > 0
                              ? ((high / totalI) * 100).toFixed(1)
                              : '0'
                          })()}
                          %
                        </div>
                      </Col>
                      <Col span={12}>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          {'低亲密度(<10)比例'}
                        </Text>
                        <div
                          style={{
                            fontSize: 20,
                            fontWeight: 700,
                            color: '#faad14'
                          }}
                        >
                          {(() => {
                            const totalI = intimacy.reduce(
                              (s: number, d: any) => s + (d.count || 0),
                              0
                            )
                            const low = intimacy
                              .filter((d: any) => d.bracket === '0-10')
                              .reduce(
                                (s: number, d: any) => s + (d.count || 0),
                                0
                              )
                            return totalI > 0
                              ? ((low / totalI) * 100).toFixed(1)
                              : '0'
                          })()}
                          %
                        </div>
                      </Col>
                    </Row>
                  </Card>
                </div>
              ) : (
                <Empty description="暂无亲密度数据" />
              )}
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  )
}

export default SocialNetwork
