import { server } from './base'

// ==================== SSO 认证 ====================
export const apiSSOConfig = () => server({ url: '/sso/config', method: 'GET' })

export const apiSSOLogin = (data: { code: string }) =>
  server({ url: '/sso/login', method: 'POST', data })

export const apiSSOSession = () =>
  server({ url: '/sso/session', method: 'GET' })

export const apiLogout = () => server({ url: '/logout', method: 'POST' })

export const apiGetMe = () => server({ url: '/me', method: 'GET' })

// ==================== 仪表盘 ====================
export const apiGetOverview = (source = 'official') =>
  server({ url: `/dashboard/overview?source=${source}`, method: 'GET' })

export const apiGetDistribution = (field: string) =>
  server({ url: `/dashboard/distribution?field=${field}`, method: 'GET' })

export const apiGetDAUTrend = (days = 30) =>
  server({ url: `/dashboard/dau-trend?days=${days}`, method: 'GET' })

export const apiGetRealmDistribution = () =>
  server({ url: '/dashboard/realm-distribution', method: 'GET' })

// ==================== 用户画像 ====================
export const apiGetProfiles = (params: {
  page?: number
  page_size?: number
  lifecycle_stage?: string
  pay_tier?: string
  play_style?: string
  min_churn_risk?: number
}) => {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.page_size) query.set('page_size', String(params.page_size))
  if (params.lifecycle_stage)
    query.set('lifecycle_stage', params.lifecycle_stage)
  if (params.pay_tier) query.set('pay_tier', params.pay_tier)
  if (params.play_style) query.set('play_style', params.play_style)
  if (params.min_churn_risk)
    query.set('min_churn_risk', String(params.min_churn_risk))
  return server({ url: `/profiles?${query}`, method: 'GET' })
}

export const apiGetProfile = (uid: number) =>
  server({ url: `/profiles/${uid}`, method: 'GET' })

export const apiRefreshProfile = (uid: number) =>
  server({ url: `/profiles/${uid}/refresh`, method: 'POST' })

export const apiRefreshAllProfiles = () =>
  server({ url: '/profiles/refresh-all', method: 'POST' })

export const apiGetSnapshots = (uid: number, start?: string, end?: string) => {
  const query = new URLSearchParams()
  if (start) query.set('start', start)
  if (end) query.set('end', end)
  return server({ url: `/profiles/${uid}/snapshots?${query}`, method: 'GET' })
}

// ==================== 分群 ====================
export const apiGetSegments = () => server({ url: '/segments', method: 'GET' })

export const apiCreateSegment = (data: {
  name: string
  description?: string
  rules_json: string
}) => server({ url: '/segments', method: 'POST', data })

export const apiDeleteSegment = (id: number) =>
  server({ url: `/segments/${id}`, method: 'DELETE' })

export const apiExecuteSegment = (id: number) =>
  server({ url: `/segments/${id}/execute`, method: 'GET' })

export const apiPreviewSegment = (rules_json: string) =>
  server({ url: '/segments/preview', method: 'POST', data: { rules_json } })

// ==================== 预警 ====================
export const apiGetHighChurnUsers = (threshold = 70, limit = 50) =>
  server({
    url: `/alerts/churn?threshold=${threshold}&limit=${limit}`,
    method: 'GET'
  })

export const apiGetStuckUsers = (limit = 50) =>
  server({ url: `/alerts/stuck?limit=${limit}`, method: 'GET' })

export const apiGetProfileStats = () =>
  server({ url: '/alerts/stats', method: 'GET' })

export const apiGetResourceAlertUsers = (limit = 50) =>
  server({ url: `/alerts/resource?limit=${limit}`, method: 'GET' })

export const apiGetNewbieAtRisk = (threshold = 40, limit = 50) =>
  server({
    url: `/alerts/newbie-risk?threshold=${threshold}&limit=${limit}`,
    method: 'GET'
  })

export const apiGetWhaleChurn = (threshold = 50, limit = 50) =>
  server({
    url: `/alerts/whale-churn?threshold=${threshold}&limit=${limit}`,
    method: 'GET'
  })

export const apiGetBotSuspects = (limit = 50) =>
  server({ url: `/alerts/bot-suspects?limit=${limit}`, method: 'GET' })

// ==================== 功能分析 ====================
export const apiGetFeatureScores = (category?: string) => {
  const q = category ? `?category=${encodeURIComponent(category)}` : ''
  return server({ url: `/features/scores${q}`, method: 'GET' })
}

export const apiGetFeatureTrend = (name: string, days = 30) =>
  server({
    url: `/features/trend?name=${encodeURIComponent(name)}&days=${days}`,
    method: 'GET'
  })

export const apiGetFeatureCategories = (days = 30) =>
  server({ url: `/features/categories?days=${days}`, method: 'GET' })

export const apiGetFeatureTop = (days = 7, limit = 20) =>
  server({ url: `/features/top?days=${days}&limit=${limit}`, method: 'GET' })

export const apiGetFeatureSceneDist = (days = 30) =>
  server({ url: `/features/scene-dist?days=${days}`, method: 'GET' })

export const apiGetFeatureChannelTop = (days = 7, limit = 20) =>
  server({
    url: `/features/channel-top?days=${days}&limit=${limit}`,
    method: 'GET'
  })

export const apiGetFeaturePlatformDist = (days = 30) =>
  server({ url: `/features/platform-dist?days=${days}`, method: 'GET' })

export const apiGetFeatureSceneFeatureTop = (
  scene: string,
  days = 7,
  limit = 15
) =>
  server({
    url: `/features/scene-feature-top?scene=${encodeURIComponent(scene)}&days=${days}&limit=${limit}`,
    method: 'GET'
  })

export const apiGetFeatureSceneTrend = (days = 30) =>
  server({ url: `/features/scene-trend?days=${days}`, method: 'GET' })

// ==================== 游戏数据统计 ====================
export const apiGetGameOverview = () =>
  server({ url: '/gamestat/overview', method: 'GET' })

export const apiGetPlayerRanking = (sort = 'realm', limit = 50) =>
  server({
    url: `/gamestat/player-ranking?sort=${sort}&limit=${limit}`,
    method: 'GET'
  })

export const apiGetGuildRanking = (sort = 'prestige', limit = 30) =>
  server({
    url: `/gamestat/guild-ranking?sort=${sort}&limit=${limit}`,
    method: 'GET'
  })

export const apiGetNewUsersTrend = (days = 30) =>
  server({ url: `/gamestat/new-users?days=${days}`, method: 'GET' })

export const apiGetRealmStages = () =>
  server({ url: '/gamestat/realm-stages', method: 'GET' })

export const apiGetRealmChurn = (days = 7) =>
  server({ url: `/gamestat/realm-churn?days=${days}`, method: 'GET' })

export const apiGetPaymentOverview = (source = 'official') =>
  server({ url: `/gamestat/payment?source=${source}`, method: 'GET' })

export const apiGetRevenueTrend = (days = 30, source = 'official') =>
  server({
    url: `/gamestat/revenue-trend?days=${days}&source=${source}`,
    method: 'GET'
  })

export const apiGetPackageStats = (days = 0, source = 'official') =>
  server({
    url: `/gamestat/packages?days=${days}&source=${source}`,
    method: 'GET'
  })

// ==================== 留存分析 ====================
export const apiGetRetention = (days = 30) =>
  server({ url: `/retention?days=${days}`, method: 'GET' })

// ==================== 转化漏斗 ====================
export const apiGetFunnel = (days = 0) =>
  server({ url: `/funnel?days=${days}`, method: 'GET' })

// ==================== 境界进阶分析 ====================
export const apiGetRealmProgress = (inactiveDays = 7) =>
  server({
    url: `/realm-progress?inactive_days=${inactiveDays}`,
    method: 'GET'
  })

export const apiGetRealmPayCorrelation = () =>
  server({ url: '/realm-progress/pay-correlation', method: 'GET' })

// ==================== 经济健康度分析 ====================
export const apiGetEconomyHealth = (days = 30) =>
  server({ url: `/economy-health?days=${days}`, method: 'GET' })

export const apiGetCurrencyTrend = (currency: string, days = 30) =>
  server({
    url: `/economy-health/currency-trend?currency=${currency}&days=${days}`,
    method: 'GET'
  })

// ==================== 副本分析 ====================
export const apiGetDungeonAnalysis = (days = 30) =>
  server({ url: `/dungeon-yonghemolimis?days=${days}`, method: 'GET' })

export const apiGetDungeonRealmDist = (days = 30) =>
  server({ url: `/dungeon-yonghemolimis/realm-dist?days=${days}`, method: 'GET' })

// ==================== 宗门健康度 ====================
export const apiGetGuildHealth = (inactiveDays = 7) =>
  server({ url: `/guild-health?inactive_days=${inactiveDays}`, method: 'GET' })

// ==================== 社交网络 ====================
export const apiGetSocialNetwork = () =>
  server({ url: '/social-network', method: 'GET' })

// ==================== 手动刷新 ====================
export const apiRefreshDashboard = () =>
  server({ url: '/refresh/dashboard', method: 'POST' })

export const apiRefreshFeatures = () =>
  server({ url: '/refresh/features', method: 'POST' })

export const apiRefreshAll = () =>
  server({ url: '/refresh/all', method: 'POST' })

// ==================== 活跃度分析 ====================
export const apiGetPlayerHourly = (uid: number, days = 7) =>
  server({ url: `/activity/player/${uid}/hourly?days=${days}`, method: 'GET' })

export const apiGetPlayerDaily = (uid: number, days = 30) =>
  server({ url: `/activity/player/${uid}/daily?days=${days}`, method: 'GET' })

export const apiGetPlayerPeak = (uid: number, days = 30) =>
  server({ url: `/activity/player/${uid}/peak?days=${days}`, method: 'GET' })

export const apiCheckPlayerBot = (uid: number, days = 14) =>
  server({
    url: `/activity/player/${uid}/bot-check?days=${days}`,
    method: 'GET'
  })

export const apiGetGlobalHourly = (days = 7) =>
  server({ url: `/activity/global/hourly?days=${days}`, method: 'GET' })

export const apiGetGlobalPeak = (days = 30) =>
  server({ url: `/activity/global/peak?days=${days}`, method: 'GET' })

export const apiRefreshActivity = (days = 3) =>
  server({ url: `/activity/refresh?days=${days}`, method: 'POST' })

// ==================== 境界卡点分析 ====================
export const apiGetRealmBottleneck = (inactiveDays = 7) =>
  server({
    url: `/realm-bottleneck?inactive_days=${inactiveDays}`,
    method: 'GET'
  })

// ==================== 通天塔与主线进度 ====================
export const apiGetTowerStory = (days = 30) =>
  server({ url: `/tower-story?days=${days}`, method: 'GET' })

export const apiGetStageDifficulty = (chapter = 0) =>
  server({
    url: `/tower-story/stage-difficulty?chapter=${chapter}`,
    method: 'GET'
  })

// ==================== 付费行为深度分析 ====================
export const apiGetPayAnalysis = (days = 30) =>
  server({ url: `/pay-yonghemolimis?days=${days}`, method: 'GET' })

// ==================== 装备与战力分析 ====================
export const apiGetEquipAnalysis = () =>
  server({ url: '/equip-yonghemolimis', method: 'GET' })

// ==================== 交易市场健康度 ====================
export const apiGetMarketHealth = (days = 30) =>
  server({ url: `/market-health?days=${days}`, method: 'GET' })

// ==================== 签到与活跃任务 ====================
export const apiGetCheckInAnalysis = (days = 30) =>
  server({ url: `/checkin-yonghemolimis?days=${days}`, method: 'GET' })

// ==================== 闭关修炼分析 ====================
export const apiGetSeclusionAnalysis = (days = 30) =>
  server({ url: `/seclusion-yonghemolimis?days=${days}`, method: 'GET' })

// ==================== 功法/技能分析 ====================
export const apiGetSkillAnalysis = () =>
  server({ url: '/skill-yonghemolimis', method: 'GET' })

// ==================== 炼丹/炼器分析 ====================
export const apiGetCraftAnalysis = () =>
  server({ url: '/craft-yonghemolimis', method: 'GET' })

// ==================== 新人引导漏斗 ====================
export const apiGetNewcomerFunnel = () =>
  server({ url: '/newcomer-funnel', method: 'GET' })

// ==================== 道侣系统分析 ====================
export const apiGetCompanionAnalysis = () =>
  server({ url: '/companion-yonghemolimis', method: 'GET' })

// ==================== 本命物系统分析 ====================
export const apiGetSoulArtifactAnalysis = () =>
  server({ url: '/soul-artifact-yonghemolimis', method: 'GET' })

// ==================== 商店消费分析 ====================
export const apiGetShopAnalysis = () =>
  server({ url: '/shop-yonghemolimis', method: 'GET' })
