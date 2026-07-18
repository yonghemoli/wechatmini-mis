# 永和茉莉小程序接口文档

> 本文档保留历史预约、订单、地址和配菜接口说明。护理小程序 v1.0 的当前规范以
> [API_REQUIREMENTS.md](API_REQUIREMENTS.md) 为准；其中 `GET /services` 已升级为护理服务分类接口，
> 历史“具体服务项目”列表迁移到 `GET /legacy-services`。

## 基础信息

- Base URL：`https://yonghemoli.com/api/mini/`
- 开发 URL：`http://127.0.0.1:8080/api/mini/`
- 请求格式：`Content-Type: application/json`
- 登录后请求头：`Authorization: Bearer <token>`

## 通用响应

成功：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

失败：

```json
{
  "code": -1,
  "message": "错误信息"
}
```

分页：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "page": 1,
    "pageSize": 20,
    "total": 100,
    "list": []
  }
}
```

## 接口分级

公开接口：登录、预约首页、门店、服务分类、服务项目、服务范围、公开菜品和官方套餐。

个人接口：个人信息、常用地址、我的家属、我的配菜、我的订单、客服会话。需要携带 `Authorization`。手机号是小程序用户唯一账号标识，微信授权登录和手机号验证码登录会合并到同一个手机号账户。

平台回调和内部联调接口：微信支付回调、支付预下单 mock、订单状态联调推进。

## 数据结构

### User

```json
{
  "id": "U10081",
  "nickName": "永和护理",
  "avatarUrl": "https://example.com/avatar.png",
  "signature": "认真生活，安心预约",
  "phone": "13800000000",
  "lastLoginAt": "2026-06-29 09:00:00"
}
```

### Address

```json
{
  "id": "addr_1",
  "contactName": "永和护理",
  "phone": "13800000000",
  "district": "广西 南宁 青秀区",
  "detail": "3 栋 2 单元 801",
  "tag": "家",
  "isDefault": true
}
```

### ServiceTarget

```json
{
  "id": "target_1",
  "name": "父亲",
  "category": "eldercare",
  "relation": "父母",
  "age": "70-79 岁",
  "note": "行动较慢",
  "isDefault": true
}
```

### MealPackage

```json
{
  "id": "custom_pkg_1",
  "name": "我家的晚餐组合",
  "scene": "少油少盐",
  "price": 24,
  "dishes": ["番茄炒蛋", "青椒肉丝"]
}
```

### Order

```json
{
  "id": "YH20260629103000123",
  "serviceId": "3",
  "serviceName": "空调清洗",
  "category": "2",
  "status": "pending_confirm",
  "appointmentTime": "2026-06-30 10:00",
  "address": "广西 南宁 青秀区 3 栋 801",
  "contactName": "永和护理",
  "phone": "13800000000",
  "remark": "",
  "serviceTargetId": "",
  "serviceTargetName": "",
  "detailFields": {
    "清洗数量": "挂机 1 台"
  },
  "amount": 0,
  "createdAt": "2026-06-29 10:30:00"
}
```

## 公开接口

### POST `auth/wechat-login`

微信授权登录。`wx.login` 的 `code` 只用于换取 `openid`，不会自动返回手机号。当前业务要求手机号必填，因此首次登录或 `openid` 未绑定手机号账户时，前端必须通过微信手机号授权拿到 `phoneCode` 并随本接口提交；后端用 `phoneCode` 换出手机号后，按手机号查找或创建同一个用户账户，并绑定 `openid`。如果该 `openid` 已绑定到有手机号的账户，后端直接登录并返回用户信息，不再重复校验手机号。

账号规则：

- `phone` 是唯一账号标识，手机号验证码登录和微信手机号授权登录会合并到同一个手机号账户。
- `openid` 是微信身份绑定标识，已绑定且账户手机号存在时，可以免重复手机号授权。
- 如果同一个手机号已经绑定到另一个 `openid`，后端拒绝本次绑定，需要走人工处理或换绑流程。

请求：

```json
{
  "code": "wx.login code",
  "encryptedData": "optional",
  "iv": "optional",
  "phoneCode": "getPhoneNumber 返回的 code，首次绑定/未绑定时必填，已绑定 openid 时可不传",
  "nickName": "永和护理",
  "avatarUrl": "https://example.com/avatar.png"
}
```

响应：

```json
{
  "token": "mini session token",
  "isBoundPhone": true,
  "user": {
    "id": "U20260629103000123",
    "nickName": "永和护理",
    "avatarUrl": "https://example.com/avatar.png",
    "signature": "",
    "phone": "13800000000",
    "lastLoginAt": "2026-06-29 10:30:00"
  }
}
```

### POST `auth/phone-code`

发送手机号验证码。

请求：

```json
{
  "phone": "13800000000"
}
```

响应：

```json
{
  "success": true
}
```

### POST `auth/phone-login`

手机号验证码登录。

请求：

```json
{
  "phone": "13800000000",
  "code": "123456"
}
```

响应同 `auth/wechat-login`。

### GET `appointment/home`

预约页首页数据。

响应：

```json
{
  "store": {
    "id": "1",
    "name": "南宁青秀店",
    "contactName": "李店长",
    "phone": "0771-8888888",
    "address": "广西南宁青秀区民族大道 100 号",
    "businessHours": "09:00-18:00",
    "status": "open",
    "distanceText": "距您约 3.2km · 南宁家庭服务"
  },
  "tabs": [{ "id": "all", "name": "所有" }],
  "groups": [],
  "activities": {}
}
```

### GET `stores/nearest`

获取最近的可服务门店。

响应：

```json
{
  "item": {
    "id": "1",
    "name": "南宁青秀店",
    "phone": "0771-8888888",
    "address": "广西南宁青秀区民族大道 100 号",
    "businessHours": "09:00-18:00",
    "status": "open",
    "distanceText": "距您约 3.2km · 南宁家庭服务"
  }
}
```

### GET `service-categories`

获取服务类型。

响应：

```json
{
  "list": [
    {
      "id": "1",
      "name": "家庭保洁",
      "desc": "日常保洁、深度保洁、新房开荒等家庭清洁服务",
      "icon": "housekeeping",
      "items": []
    }
  ]
}
```

### GET `service-categories/{id}/services`

获取某个服务类型下的预约服务。

响应：

```json
{
  "list": [
    {
      "id": "3",
      "name": "空调清洗",
      "scene": "适合换季前后",
      "summary": "清洗滤网、蒸发器和出风口",
      "priceText": "预约后报价",
      "durationText": "约 1 小时/台",
      "imageText": "设备清洗",
      "category": "2",
      "action": "booking"
    }
  ]
}
```

### GET `services`

服务项目列表。

查询参数：

- `category`：可选，服务类型 ID。
- `keyword`：可选，搜索关键字。

响应：

```json
{
  "list": [
    {
      "id": "3",
      "category": "2",
      "name": "空调清洗",
      "scene": "适合换季前后",
      "summary": "清洗滤网、蒸发器和出风口",
      "priceText": "预约后报价",
      "durationText": "约 1 小时/台",
      "requirementLabel": "清洗数量",
      "requirementOptions": ["挂机 1 台", "挂机 2 台"],
      "suitableFor": ["卧室空调"],
      "scope": ["滤网清洗"],
      "process": ["预约确认"],
      "notes": ["高空外机不在服务范围内"]
    }
  ]
}
```

### GET `services/search`

搜索服务项目。参数和响应同 `services`。

### GET `services/{id}`

服务项目详情。响应为 `services.list[]` 单项。

### GET `service-areas`

获取服务范围。

响应：

```json
{
  "city": "南宁",
  "districts": ["青秀区", "兴宁区", "西乡塘区", "良庆区", "江南区"],
  "notes": ["具体地址以客服确认为准"]
}
```

### GET `meal/pricing`

获取配菜价格配置。

响应：

```json
{
  "dishPrice": 12,
  "deliveryFee": 8
}
```

### GET `meal/dishes`

获取菜品列表。

响应：

```json
{
  "list": [
    {
      "id": "tomato-egg",
      "name": "番茄炒蛋",
      "scene": "酸甜开胃，适合日常晚餐",
      "tag": "家常",
      "price": 12,
      "ingredients": ["番茄 300g", "鸡蛋 3 个"],
      "videoTitle": "番茄炒蛋 8 分钟快手做法",
      "videoUrl": "https://example.com/tomato-egg",
      "comments": ["孩子很爱吃"]
    }
  ]
}
```

### GET `meal/dishes/{nameOrId}`

获取菜品详情。响应为 `meal/dishes.list[]` 单项。

### GET `meal/packages`

获取官方配菜套餐。

响应：

```json
{
  "list": [
    {
      "id": "pkg_1",
      "name": "三菜一汤套餐",
      "scene": "适合 2-3 人工作日晚餐",
      "price": 68,
      "dishes": ["番茄炒蛋", "青椒肉丝", "紫菜蛋花汤"]
    }
  ]
}
```

## 个人接口

### GET `user/profile`

获取个人信息。

响应：`User`。

### PUT `user/profile`

修改个人信息。

说明：手机号是唯一登录账号，不通过该接口修改。手机号变更需要新增独立的手机号验证/换绑流程。

请求：

```json
{
  "nickName": "永和护理",
  "avatarUrl": "https://example.com/avatar.png",
  "signature": "认真生活，安心预约",
  "phone": "13800000000"
}
```

响应：`User`。

### GET `addresses`

获取常用地址列表。

响应：

```json
{
  "list": []
}
```

### GET `addresses/{id}`

获取常用地址详情。响应：`Address`。

### POST `addresses`

新增常用地址。

请求：`Address`，`id` 可不传。

响应：`Address`。

### PUT `addresses/{id}`

修改常用地址。

请求：`Address`。

响应：`Address`。

### PUT `addresses/{id}/default`

设为默认地址。

响应：

```json
{
  "id": "addr_1"
}
```

### DELETE `addresses/{id}`

删除常用地址。

响应：

```json
{
  "id": "addr_1"
}
```

### GET `service-targets`

获取我的家属列表。

查询参数：

- `category`：可选，例如 `eldercare`、`childcare`。

响应：

```json
{
  "list": []
}
```

### GET `service-targets/{id}`

获取我的家属详情。响应：`ServiceTarget`。

### POST `service-targets`

新增家属。

请求：`ServiceTarget`，`id` 可不传。

响应：`ServiceTarget`。

### PUT `service-targets/{id}`

修改家属。

请求：`ServiceTarget`。

响应：`ServiceTarget`。

### PUT `service-targets/{id}/default`

设为默认家属。同一分类下只保留一个默认家属。

响应：

```json
{
  "id": "target_1"
}
```

### DELETE `service-targets/{id}`

删除家属。

响应：

```json
{
  "id": "target_1"
}
```

### GET `meal/custom-packages`

获取我的配菜列表。

响应：

```json
{
  "list": []
}
```

### GET `meal/custom-packages/{id}`

获取我的配菜详情。响应：`MealPackage`。

### POST `meal/custom-packages`

新增我的配菜。

请求：`MealPackage`，`id` 可不传。

响应：`MealPackage`。

### PUT `meal/custom-packages/{id}`

修改我的配菜。

请求：`MealPackage`。

响应：`MealPackage`。

### DELETE `meal/custom-packages/{id}`

删除我的配菜。

响应：

```json
{
  "id": "custom_pkg_1"
}
```

### GET `orders`

获取我的订单列表。

查询参数：

- `status`：可选，默认 `all`。
- `page`：可选，默认 `1`。
- `pageSize`：可选，默认 `20`。

响应：分页结构，`list` 为 `Order[]`。

### GET `orders/{id}`

获取我的订单详情。响应：`Order`。

### POST `orders`

创建订单。服务预约和配菜统一使用该接口。

请求：

```json
{
  "serviceId": "3",
  "serviceName": "空调清洗",
  "category": "2",
  "appointmentTime": "2026-06-30 10:00",
  "address": "广西 南宁 青秀区 3 栋 801",
  "contactName": "永和护理",
  "phone": "13800000000",
  "remark": "",
  "serviceTargetId": "",
  "serviceTargetName": "",
  "detailFields": {
    "清洗数量": "挂机 1 台"
  },
  "amount": 0
}
```

响应：当前请求体。

### POST `orders/{id}/cancel`

取消订单。

响应：`Order`。

### DELETE `orders/{id}`

用户侧删除订单。当前实现为小程序侧软删除，不影响管理端历史数据。

响应：

```json
{
  "id": "YH20260629103000123"
}
```

### GET `chat/session`

获取客服会话。

查询参数：

- `sessionId`：可选。不传时使用当前用户默认会话。

响应：

```json
{
  "item": {}
}
```

### GET `chat/messages`

获取客服消息列表。

查询参数：

- `sessionId`：可选。不传时使用当前用户默认会话。

响应：

```json
{
  "list": []
}
```

### POST `chat/messages`

发送客服消息。

请求：

```json
{
  "sessionId": "chat_U10081",
  "content": "想咨询明天下午保洁",
  "msgType": "text"
}
```

响应：

```json
{
  "item": {}
}
```

## 平台回调和内部联调接口

### POST `payments/wechat/prepay`

微信支付预下单 mock。

请求：

```json
{
  "orderId": "YH20260629103000123",
  "amount": 3200
}
```

响应：

```json
{
  "timeStamp": "1780000000",
  "nonceStr": "dev-mini-nonce",
  "package": "prepay_id=dev_prepay_id",
  "signType": "RSA",
  "paySign": "dev_pay_sign"
}
```

### POST `payments/wechat/notify`

微信支付回调占位。

响应：

```json
{
  "received": true
}
```

### PUT `orders/{id}/status`

内部联调用订单状态推进接口。后续正式业务应由后台或运营端推进。

请求：

```json
{
  "status": "confirmed"
}
```

响应：`Order`。
