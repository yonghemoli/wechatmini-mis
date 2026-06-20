# yonhemoli-mis

家政业务 MIS 管理端。MIS 是公司的“大脑”和“发动机”，初期目标不是展示好看，而是让运营、客服和管理者高效发现并处理异常。

## 当前项目状态

- 技术栈：Go 后端 + Vite/React/Ant Design 前端。
- 数据库：MySQL 作为业务读写库。
- 现状：仓库原始形态是游戏数据分析后台，后端和大量旧前端页面仍保留历史模块。
- 本次升级：前端入口、导航和核心页面已收敛为家政 MIS，不再暴露游戏分析入口。
- 当前页面：工作台、订单管理、用户管理、内容/商品发布、数据看板。
- 数据状态：前端先使用本地 mock 数据跑通运营流程，后续再替换为真实家政业务 API。

## 家政 MIS 模块边界

1. 工作台
   - 今日订单量、今日营收、待处理工单数、用户总数。
   - 将异常和待处理项前置，使用红点和数字提醒。

2. 订单管理
   - 全量订单列表。
   - 支持按状态筛选、关键字搜索、时间筛选占位。
   - 支持确认/核销、异常关闭/退款、内部备注。
   - 支持导出 CSV，供 Excel 打开和财务对账。

3. 用户管理
   - 展示头像、昵称、注册时间、累计消费金额、最后下单时间。
   - 仅保留手动封禁/解封。

4. 内容/商品发布
   - 管理小程序首页展示的服务图、价格、服务介绍文案。
   - 支持增删改查和上下架开关。

5. 数据看板
   - 近 7 天/30 天营收曲线。
   - 来源分析：扫一扫、搜索、分享。
   - 支持导出 CSV。

## 后续升级记录

要升级为真正独立的家政 MIS，需要继续处理以下事项：

- 后端重命名和清理：移除 `gamestat`、`realm`、`dungeon` 等游戏领域 API，建立 `orders`、`users`、`contents`、`reports` 等家政领域接口。
- 数据模型：新增订单、用户、服务商品、内部备注、退款/关闭记录、来源分析表。
- 权限体系：按老板/运营/客服/财务分配菜单和操作权限。
- 导出能力：当前为前端 CSV，正式版应由后端按筛选条件生成对账文件。
- 小程序联动：内容/商品发布需要对接小程序首页配置与上下架状态。

## 小程序接口

小程序端接口使用独立前缀 `/api/mini`，不走管理端 `/api/v1` 登录态。

当前已提供：

- `POST /api/mini/auth/wechat-login`：微信授权登录 mock。
- `GET /api/mini/user/profile`：获取当前用户信息。
- `GET /api/mini/home`：小程序首页数据。
- `GET /api/mini/services`：服务项目列表。
- `GET /api/mini/service-areas`：服务范围。
- `GET /api/mini/addresses`：常用地址列表。
- `POST /api/mini/addresses`：新增地址。
- `PUT /api/mini/addresses/:id/default`：设为默认地址。
- `DELETE /api/mini/addresses/:id`：删除地址。
- `GET /api/mini/service-targets`：老人/孩子服务对象列表。
- `POST /api/mini/service-targets`：新增服务对象。
- `DELETE /api/mini/service-targets/:id`：删除服务对象。
- `GET /api/mini/meal/dishes`：菜品列表。
- `GET /api/mini/meal/packages`：官方配菜套餐。
- `GET /api/mini/meal/custom-packages`：我的配菜套餐。
- `POST /api/mini/meal/custom-packages`：保存个性化套餐。
- `DELETE /api/mini/meal/custom-packages/:id`：删除个性化套餐。
- `GET /api/mini/orders`：订单分页列表。
- `GET /api/mini/orders/:id`：订单详情。
- `POST /api/mini/orders`：创建服务预约/配菜订单。
- `PUT /api/mini/orders/:id/status`：更新订单状态。
- `POST /api/mini/orders/:id/cancel`：取消订单。
- `POST /api/mini/payments/wechat/prepay`：微信支付预下单 mock。
- `POST /api/mini/payments/wechat/notify`：微信支付回调占位。

当前实现是内存 mock，目的是满足小程序端按字段契约联调。正式上线前需要替换为微信登录、JWT 校验、数据库持久化、订单状态机和微信支付签名验签。

## 后端数据库

业务数据库使用 MySQL，通过环境变量配置：

```dotenv
ANALYTICS_DB_DRIVER=mysql
ANALYTICS_DB_DSN=root:password@tcp(127.0.0.1:3306)/yonhemoli_mis?charset=utf8mb4&parseTime=True&loc=Local
ANALYTICS_DB_AUTO_MIGRATE=true
```

`ANALYTICS_GAME_DB_*` 是历史游戏只读库配置，后续完成家政 MIS 后端清理时应移除。

## 前端开发

```sh
cd frontend
yarn
yarn dev
```

## 构建检查

```sh
cd frontend
yarn build
```
