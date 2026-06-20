import React, { useEffect, useState } from 'react'
import {
  Card,
  Select,
  Spin,
  Table,
  Tag,
  Tooltip,
  Typography,
  Empty
} from 'antd'
import { LineChartOutlined } from '@ant-design/icons'
import { apiGetRetention } from '@/api'

const { Title, Text } = Typography

interface RetentionRow {
  cohort_date: string
  cohort_size: number
  day1: number
  day3: number
  day7: number
  day14: number
  day30: number
}

const retentionColor = (rate: number) => {
  if (rate >= 50) return '#52c41a'
  if (rate >= 30) return '#1890ff'
  if (rate >= 15) return '#faad14'
  return '#f5222d'
}

const Retention: React.FC = () => {
  const [loading, setLoading] = useState(true)
  const [days, setDays] = useState(30)
  const [data, setData] = useState<RetentionRow[]>([])

  useEffect(() => {
    loadData()
  }, [days])

  const loadData = async () => {
    setLoading(true)
    try {
      const res = await apiGetRetention(days)
      if (res?.data) setData(res.data || [])
    } catch {
    } finally {
      setLoading(false)
    }
  }

  const renderRate = (val: number, total: number) => {
    if (total === 0) return <Text type="secondary">-</Text>
    const rate = (val / total) * 100
    return (
      <Tooltip title={`${val} / ${total} 人`}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
          <div
            style={{
              width: 40,
              height: 20,
              borderRadius: 3,
              background: retentionColor(rate),
              opacity: 0.15 + (rate / 100) * 0.85,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}
          >
            <span
              style={{
                fontSize: 11,
                fontWeight: 600,
                color: retentionColor(rate)
              }}
            >
              {rate.toFixed(1)}%
            </span>
          </div>
        </div>
      </Tooltip>
    )
  }

  // 计算平均留存率
  const avgRetention = (field: keyof RetentionRow) => {
    const valid = data.filter(d => d.cohort_size > 0)
    if (valid.length === 0) return 0
    const sum = valid.reduce(
      (s, d) => s + ((d[field] as number) / d.cohort_size) * 100,
      0
    )
    return sum / valid.length
  }

  const columns = [
    {
      title: '注册日期',
      dataIndex: 'cohort_date',
      key: 'cohort_date',
      width: 120,
      render: (v: string) => <Text strong>{v}</Text>
    },
    {
      title: '新增人数',
      dataIndex: 'cohort_size',
      key: 'cohort_size',
      width: 90,
      sorter: (a: RetentionRow, b: RetentionRow) =>
        a.cohort_size - b.cohort_size,
      render: (v: number) => <Tag color="blue">{v}</Tag>
    },
    {
      title: '次日留存',
      key: 'day1',
      width: 100,
      render: (_: any, row: RetentionRow) =>
        renderRate(row.day1, row.cohort_size)
    },
    {
      title: '3日留存',
      key: 'day3',
      width: 100,
      render: (_: any, row: RetentionRow) =>
        renderRate(row.day3, row.cohort_size)
    },
    {
      title: '7日留存',
      key: 'day7',
      width: 100,
      render: (_: any, row: RetentionRow) =>
        renderRate(row.day7, row.cohort_size)
    },
    {
      title: '14日留存',
      key: 'day14',
      width: 100,
      render: (_: any, row: RetentionRow) =>
        renderRate(row.day14, row.cohort_size)
    },
    {
      title: '30日留存',
      key: 'day30',
      width: 100,
      render: (_: any, row: RetentionRow) =>
        renderRate(row.day30, row.cohort_size)
    }
  ]

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
            <LineChartOutlined style={{ marginRight: 8 }} />
            留存分析
          </Title>
          <Select
            value={days}
            onChange={setDays}
            style={{ width: 120 }}
            options={[
              { label: '近7天', value: 7 },
              { label: '近14天', value: 14 },
              { label: '近30天', value: 30 },
              { label: '近60天', value: 60 },
              { label: '近90天', value: 90 }
            ]}
          />
        </div>

        {/* 平均留存概览 */}
        <Card size="small" style={{ marginBottom: 16 }}>
          <div style={{ display: 'flex', gap: 32, flexWrap: 'wrap' }}>
            {[
              { label: '平均次日留存', value: avgRetention('day1') },
              { label: '平均3日留存', value: avgRetention('day3') },
              { label: '平均7日留存', value: avgRetention('day7') },
              { label: '平均14日留存', value: avgRetention('day14') },
              { label: '平均30日留存', value: avgRetention('day30') }
            ].map(item => (
              <div key={item.label} style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 12, color: '#999' }}>{item.label}</div>
                <div
                  style={{
                    fontSize: 24,
                    fontWeight: 700,
                    color: retentionColor(item.value)
                  }}
                >
                  {item.value.toFixed(1)}%
                </div>
              </div>
            ))}
          </div>
        </Card>

        {/* 留存热力图表格 */}
        <Card title="留存明细" size="small">
          {data.length > 0 ? (
            <Table
              dataSource={data}
              columns={columns}
              rowKey="cohort_date"
              size="small"
              pagination={{
                pageSize: 15,
                size: 'small',
                showSizeChanger: false
              }}
              scroll={{ x: 800 }}
            />
          ) : (
            <Empty description="暂无留存数据" />
          )}
        </Card>
      </div>
    </Spin>
  )
}

export default Retention
