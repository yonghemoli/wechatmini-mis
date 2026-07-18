# yonhemoli-mis

## 当前项目状态

- 技术栈：Go 后端 + Vite/React/Ant Design 前端。
- 数据库：MySQL 作为业务读写库。
- 管理端：内部账号密码登录，不再依赖单点登录。
- 前端入口：管理端统一挂载在 `/admin`。
- 本次升级：前端入口、导航和核心页面已收敛为永和茉莉。
- 当前页面：工作台、订单管理、用户管理、服务类型管理、服务管理、账户管理、店铺管理、常见问题管理、客服在线。
- 护理小程序业务：护理服务分类、服务人员、客户需求、求职简历和运营内容均已提供 MIS 接口。
- 数据状态：管理端页面已对接 `/api/v1` 真实接口；开发初始化数据通过 `sql/init-seed.sql` 手动写入。

## 永和茉莉 模块边界

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

4. 服务类型管理
   - 管理服务分类、排序、启用/停用。
   - 给小程序和运营后台统一服务分类口径。

5. 服务管理
   - 管理服务名称、所属类型、图片、价格、单位、介绍、排序。
   - 支持新增、编辑、删除、上下架和导出 CSV。

6. 账户管理
   - 内部管理员多账号管理。
   - 支持创建、编辑、启停用、重置密码，仅超级管理员可见。

7. 店铺管理
   - 管理店铺名称、联系人、电话、地址、营业时间、营业状态和备注。
   - 支持开业/停业。

8. 常见问题管理
   - 管理小程序 FAQ 的问题、答案、分类、排序和上下架。

9. 客服在线
   - 独立全屏聊天页。
   - 支持会话列表、消息查看、文本回复、标记已读和关闭会话。

10. 护理小程序运营
   - `/api/v1/care-service-categories`：稳定服务分类管理。
   - `/api/v1/caregivers`：服务人员档案、发布、推荐和排序。
   - `/api/v1/demands`、`/api/v1/resumes`：咨询需求与求职简历跟进、顾问分配和状态历史。
   - `/api/v1/mini-content/:key`：首页配置、公司介绍和协议内容；`key` 支持 `app-config`、`about`、`agreement-privacy`、`agreement-service`。

后台入口为 `/admin`，管理端接口为 `/api/v1`，小程序接口为 `/api/mini`。小程序完整接口文档见 [docs/mini-api.md](docs/mini-api.md)。

## 初始化数据

默认建议使用独立 SQL 文件初始化表结构和默认数据，不依赖服务启动时自动建表，避免每次启动产生额外 SQL。

`MIS_DB_AUTO_MIGRATE=true` 时，GORM 会对每个模型查询 `information_schema` 检查表、字段、索引和约束，所以启动日志里会出现大量 SQL。生产和日常开发默认保持 `false`。

首次部署或新增模块表不存在时，手动执行：

```sh
# 初始化或补齐表结构
make db-schema

# 初始化默认管理员和开发数据；需要本机有 mysql CLI
make db-seed

# 或一次执行表结构和默认数据
make db-init
```

`sql/init-schema.sql` 可以重复执行：新库会创建完整表结构；旧库会按 `information_schema` 判断并补齐新增字段和索引，例如 `users.phone`、`users.signature`、`users.last_login_at`、`orders.user_id`、`orders.mini_deleted` 和 `idx_orders_user_id`。

护理小程序短信验证码通过通用 HTTP JSON 网关发送，生产环境需要配置：

```text
MIS_MINI_SMS_ENDPOINT=https://sms-gateway.example.com/send
MIS_MINI_SMS_TOKEN=your-token
```

网关收到的 JSON 为 `{"phone":"13800138000","code":"123456","scene":"LOGIN"}`。仅在非 `release` 环境可用 `MIS_MINI_SMS_TEST_CODE=123456` 跳过外部网关进行本地联调；生产环境不会启用固定验证码。

默认管理员账号：

```text
admin / admin123
```

注意：登录不读取 `MIS_ADMIN_USERNAME` / `MIS_ADMIN_PASSWORD`。内部账号必须存在于 MySQL 的 `admins` 表；首次部署或清空数据库后，需要手动执行 `sql/init-seed.sql` 写入默认管理员。

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
