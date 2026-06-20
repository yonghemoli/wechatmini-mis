import React, { useMemo, useState } from 'react'
import { Button, Radio, Space, Table, Typography } from 'antd'
import { DownloadOutlined } from '@ant-design/icons'
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from 'recharts'
import { exportCsv, revenueTrend30, revenueTrend7, sourceRows } from '../misData'

const { Title, Text } = Typography

const Reports: React.FC = () => {
  const [range, setRange] = useState<'7' | '30'>('7')
  const trend = useMemo(() => (range === '7' ? revenueTrend7 : revenueTrend30), [range])

  const exportReports = () => {
    exportCsv('家政经营看板.csv', [
      ...trend.map(item => ({
        类型: `${range}天营收趋势`,
        日期或来源: item.date,
        用户数: '',
        订单数: '',
        营收: item.revenue
      })),
      ...sourceRows.map(item => ({
        类型: '来源分析',
        日期或来源: item.source,
        用户数: item.users,
        订单数: item.orders,
        营收: item.revenue
      }))
    ])
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div className="mis-page-title">
        <div>
          <Title level={3} style={{ margin: 0 }}>
            数据看板
          </Title>
          <Text type="secondary">经营之眼：交易趋势、来源分析和财务导出。</Text>
        </div>
        <Button icon={<DownloadOutlined />} onClick={exportReports}>
          导出 Excel
        </Button>
      </div>

      <div className="mis-panel">
        <div className="mis-panel-header">
          <Title level={5} style={{ margin: 0 }}>
            交易趋势
          </Title>
          <Radio.Group
            value={range}
            onChange={event => setRange(event.target.value)}
            options={[
              { label: '近 7 天', value: '7' },
              { label: '近 30 天', value: '30' }
            ]}
            optionType="button"
            buttonStyle="solid"
          />
        </div>
        <div style={{ width: '100%', height: 320 }}>
          <ResponsiveContainer>
            <LineChart data={trend} margin={{ top: 16, right: 24, left: 0, bottom: 8 }}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis />
              <Tooltip formatter={value => [`¥${Number(value || 0)}`, '营收']} />
              <Line
                type="monotone"
                dataKey="revenue"
                stroke="#1677ff"
                strokeWidth={2}
                dot={{ r: 3 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="mis-panel">
        <div className="mis-panel-header">
          <Title level={5} style={{ margin: 0 }}>
            来源分析
          </Title>
          <Text type="secondary">来自微信基础来源：扫一扫、搜索、分享。</Text>
        </div>
        <Table
          rowKey="source"
          pagination={false}
          dataSource={sourceRows}
          columns={[
            {
              title: '来源',
              dataIndex: 'source'
            },
            {
              title: '用户数',
              dataIndex: 'users'
            },
            {
              title: '订单数',
              dataIndex: 'orders'
            },
            {
              title: '营收',
              dataIndex: 'revenue',
              render: (value: number) => `¥${value}`
            },
            {
              title: '订单转化率',
              render: (_, record) => `${((record.orders / record.users) * 100).toFixed(1)}%`
            }
          ]}
        />
      </div>
    </Space>
  )
}

export default Reports
