import React, { useEffect, useState } from 'react'
import {
  Table,
  Card,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Typography,
  Space,
  Popconfirm,
  Tag,
  Divider,
  InputNumber,
  message
} from 'antd'
import {
  PlusOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  EyeOutlined,
  MinusCircleOutlined
} from '@ant-design/icons'
import {
  apiGetSegments,
  apiCreateSegment,
  apiDeleteSegment,
  apiExecuteSegment,
  apiPreviewSegment
} from '@/api'

const { Title, Text } = Typography

const lifecycleLabels: Record<string, string> = {
  NEW: '新手期(0-7天)',
  GROWING: '成长期',
  MATURE: '成熟期(稳定活跃)',
  DECLINING: '衰退期(活跃下降)',
  LOST: '流失期(长期未登录)',
  RETURNED: '回流期(流失后回归)'
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

const lifecycleColors: Record<string, string> = {
  NEW: 'green',
  GROWING: 'blue',
  MATURE: 'purple',
  DECLINING: 'orange',
  LOST: 'red',
  RETURNED: 'cyan'
}

const payTierColors: Record<string, string> = {
  FREE: 'default',
  MINNOW: 'green',
  DOLPHIN: 'blue',
  ORCA: 'purple',
  WHALE: 'gold',
  LEVIATHAN: 'red'
}

// 可选字段
const FIELD_OPTIONS = [
  { label: '生命周期', value: 'lifecycle_stage' },
  { label: '付费等级', value: 'pay_tier' },
  { label: '玩法偏好', value: 'play_style' },
  { label: '社交类型', value: 'social_type' },
  { label: '流失风险', value: 'churn_risk' },
  { label: 'LTV预测', value: 'ltv_predict' },
  { label: '卡关标记', value: 'stuck_flag' },
  { label: '资源预警', value: 'resource_alert' }
]

// 操作符
const OPERATOR_OPTIONS = [
  { label: '等于', value: 'eq' },
  { label: '不等于', value: 'neq' },
  { label: '大于', value: 'gt' },
  { label: '大于等于', value: 'gte' },
  { label: '小于', value: 'lt' },
  { label: '小于等于', value: 'lte' },
  { label: '包含', value: 'in' },
  { label: '模糊匹配', value: 'contains' }
]

// 字段值预设 (value → label 映射)
const FIELD_VALUE_OPTIONS: Record<string, { value: string; label: string }[]> =
  {
    lifecycle_stage: Object.entries(lifecycleLabels).map(([v, l]) => ({
      value: v,
      label: l
    })),
    pay_tier: Object.entries(payTierLabels).map(([v, l]) => ({
      value: v,
      label: l
    })),
    play_style: Object.entries(playStyleLabels).map(([v, l]) => ({
      value: v,
      label: l
    })),
    social_type: [
      '独行侠',
      '双修型',
      '社交达人',
      '宗门活跃',
      '宗门成员',
      '交友型'
    ].map(v => ({ value: v, label: v })),
    stuck_flag: [
      { value: 'true', label: '是' },
      { value: 'false', label: '否' }
    ],
    resource_alert: [
      { value: 'true', label: '是' },
      { value: 'false', label: '否' }
    ]
  }

interface RuleItem {
  field: string
  operator: string
  value: string | number
}

const Segments: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<any[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [resultModal, setResultModal] = useState(false)
  const [resultData, setResultData] = useState<{
    total: number
    profiles: any[]
  }>({
    total: 0,
    profiles: []
  })
  const [form] = Form.useForm()
  const [rules, setRules] = useState<RuleItem[]>([
    { field: 'lifecycle_stage', operator: 'eq', value: '' }
  ])
  const [logic, setLogic] = useState<'AND' | 'OR'>('AND')
  const [previewCount, setPreviewCount] = useState<number | null>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const res = await apiGetSegments()
      if (res?.data) setData(res.data || [])
    } finally {
      setLoading(false)
    }
  }

  const buildRulesJSON = (): string => {
    const validRules = rules.filter(
      r => r.field && r.operator && r.value !== '' && r.value !== undefined
    )
    return JSON.stringify({
      logic,
      rules: validRules.map(r => ({
        field: r.field,
        operator: r.operator,
        value: isNaN(Number(r.value)) ? r.value : Number(r.value)
      }))
    })
  }

  const handlePreview = async () => {
    const rulesJSON = buildRulesJSON()
    try {
      const res = await apiPreviewSegment(rulesJSON)
      if (res?.data) {
        setPreviewCount(res.data.count ?? 0)
      }
    } catch {
      setPreviewCount(null)
    }
  }

  const handleCreate = async () => {
    try {
      const values = await form.validateFields()
      const rulesJSON = buildRulesJSON()
      await apiCreateSegment({
        name: values.name,
        description: values.description,
        rules_json: rulesJSON
      })
      message.success('分群创建成功')
      setModalOpen(false)
      form.resetFields()
      setRules([{ field: 'lifecycle_stage', operator: 'eq', value: '' }])
      setPreviewCount(null)
      loadData()
    } catch {
      message.error('请检查输入项是否完整')
    }
  }

  const handleDelete = async (id: number) => {
    await apiDeleteSegment(id)
    message.success('分群已删除')
    loadData()
  }

  const handleExecute = async (id: number) => {
    try {
      const res = await apiExecuteSegment(id)
      if (res?.data) {
        setResultData({
          total: res.data.total || 0,
          profiles: res.data.profiles || []
        })
        setResultModal(true)
      }
    } catch {
      message.error('执行查询失败')
    }
  }

  const addRule = () => {
    setRules([...rules, { field: 'churn_risk', operator: 'gte', value: '' }])
  }

  const removeRule = (index: number) => {
    setRules(rules.filter((_, i) => i !== index))
  }

  const updateRule = (index: number, key: keyof RuleItem, val: any) => {
    const newRules = [...rules]
    newRules[index] = { ...newRules[index], [key]: val }
    setRules(newRules)
    setPreviewCount(null) // 规则变了，清除预览
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true
    },
    {
      title: '用户数',
      dataIndex: 'user_count',
      key: 'user_count',
      render: (v: number) => <Tag color="blue">{v}</Tag>
    },
    { title: '创建者', dataIndex: 'created_by', key: 'created_by' },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (v: string) => (v ? new Date(v).toLocaleString() : '-')
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: any, row: any) => (
        <Space>
          <Button
            size="small"
            type="primary"
            icon={<PlayCircleOutlined />}
            onClick={() => handleExecute(row.id)}
          >
            执行
          </Button>
          <Popconfirm title="确认删除?" onConfirm={() => handleDelete(row.id)}>
            <Button size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const resultColumns = [
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
      sorter: (a: any, b: any) => a.churn_risk - b.churn_risk,
      defaultSortOrder: 'descend' as const,
      render: (v: number) => (
        <Tag color={v >= 70 ? 'red' : v >= 40 ? 'orange' : 'green'}>{v}%</Tag>
      )
    },
    {
      title: '卡关',
      dataIndex: 'stuck_flag',
      key: 'stuck_flag',
      render: (v: boolean) => (v ? <Tag color="red">是</Tag> : <Tag>否</Tag>)
    },
    {
      title: '资源告急',
      dataIndex: 'resource_alert',
      key: 'resource_alert',
      render: (v: boolean) => (v ? <Tag color="orange">是</Tag> : <Tag>否</Tag>)
    },
    {
      title: 'LTV',
      dataIndex: 'ltv_predict',
      key: 'ltv_predict',
      render: (v: number) => `¥${(v || 0).toFixed(2)}`
    }
  ]

  return (
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
          用户分群
        </Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setModalOpen(true)}
        >
          创建分群
        </Button>
      </div>

      <Card>
        <Table
          columns={columns}
          dataSource={data}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      </Card>

      {/* 创建分群 - 可视化规则构建 */}
      <Modal
        title="创建分群"
        open={modalOpen}
        onOk={handleCreate}
        onCancel={() => {
          setModalOpen(false)
          form.resetFields()
          setRules([{ field: 'lifecycle_stage', operator: 'eq', value: '' }])
          setPreviewCount(null)
        }}
        okText="创建"
        cancelText="取消"
        width={720}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="分群名称" rules={[{ required: true }]}>
            <Input placeholder="例如: 高价值流失预警用户" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input placeholder="分群说明" />
          </Form.Item>
        </Form>

        <Divider>
          筛选规则
          <Select
            value={logic}
            onChange={v => {
              setLogic(v)
              setPreviewCount(null)
            }}
            style={{ width: 80, marginLeft: 12 }}
            size="small"
            options={[
              { label: '且(AND)', value: 'AND' },
              { label: '或(OR)', value: 'OR' }
            ]}
          />
        </Divider>

        {rules.map((rule, idx) => (
          <Space
            key={idx}
            style={{ display: 'flex', marginBottom: 8 }}
            align="start"
          >
            <Select
              value={rule.field}
              onChange={v => updateRule(idx, 'field', v)}
              style={{ width: 140 }}
              options={FIELD_OPTIONS}
              placeholder="字段"
            />
            <Select
              value={rule.operator}
              onChange={v => updateRule(idx, 'operator', v)}
              style={{ width: 120 }}
              options={OPERATOR_OPTIONS}
              placeholder="操作符"
            />
            {FIELD_VALUE_OPTIONS[rule.field] ? (
              <Select
                value={rule.value as string}
                onChange={v => updateRule(idx, 'value', v)}
                style={{ width: 180 }}
                options={FIELD_VALUE_OPTIONS[rule.field]}
                placeholder="选择值"
              />
            ) : (
              <InputNumber
                value={rule.value as number}
                onChange={v => updateRule(idx, 'value', v ?? '')}
                style={{ width: 180 }}
                placeholder={rule.field === 'churn_risk' ? '0-100' : '数值'}
              />
            )}
            {rules.length > 1 && (
              <Button
                type="text"
                danger
                icon={<MinusCircleOutlined />}
                onClick={() => removeRule(idx)}
              />
            )}
          </Space>
        ))}

        <Button
          type="dashed"
          onClick={addRule}
          icon={<PlusOutlined />}
          style={{ width: '100%', marginTop: 8 }}
        >
          添加条件
        </Button>

        <div
          style={{
            marginTop: 16,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center'
          }}
        >
          <Button icon={<EyeOutlined />} onClick={handlePreview}>
            预览人数
          </Button>
          {previewCount !== null && (
            <Text strong>
              匹配用户数: <Tag color="blue">{previewCount}</Tag>
            </Text>
          )}
        </div>
      </Modal>

      {/* 执行结果 */}
      <Modal
        title={`分群查询结果 (${resultData.total} 人)`}
        open={resultModal}
        onCancel={() => setResultModal(false)}
        footer={null}
        width={900}
      >
        <Table
          columns={resultColumns}
          dataSource={resultData.profiles}
          rowKey="uid"
          pagination={{ pageSize: 20 }}
          size="small"
        />
      </Modal>
    </div>
  )
}

export default Segments
