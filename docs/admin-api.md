# 永和护理 MIS API

Base URL：`/api/v1`。除登录、会话和退出外，均需要内部管理会话 Cookie。

## 账户管理

- `POST /login`
- `GET /session`
- `POST /logout`
- `GET /me`
- `/admin/accounts`：列表、新增、编辑、启停、重置密码。

## 用户管理

- `GET /users`
- `GET /users/export`
- `POST /users/{id}/ban`
- `POST /users/{id}/unban`

## 装修

- `GET /decoration`
- `GET|PUT /decoration/banners`
- `GET|PUT /decoration/customer-service`
- `GET|PUT /decoration/company`
- `/decoration/services`：服务项目列表、新增、编辑、删除。

## 阿姨与申请

- `/caregivers`：列表、详情、新增、编辑、删除。
- `/caregiver-applications`：求职申请列表、分配和状态流转。

## 预约管理

- `GET /appointments`
- `PUT /appointments/{id}/assign`
- `PUT /appointments/{id}/status`
- `GET /status-history/{entityType}/{id}`

## 常见问题

- `/faqs`：列表、新增、编辑、删除、发布和下架。

## 客服在线

- `/chat/sessions`
- `/chat/sessions/{id}/messages`
- `/chat/sessions/{id}/read`
- `/chat/sessions/{id}/close`
- `GET /ws/chat`

旧工作台、订单、服务类型、旧服务、店铺等接口已删除。
