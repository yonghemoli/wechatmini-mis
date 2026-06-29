# yonhemoli-mis

家政业务 MIS 管理端。MIS 是公司的“大脑”和“发动机”，初期目标不是展示好看，而是让运营、客服和管理者高效发现并处理异常。

## 当前项目状态

- 技术栈：Go 后端 + Vite/React/Ant Design 前端。
- 数据库：MySQL 作为业务读写库。
- 管理端：内部账号密码登录，不再依赖单点登录。
- 前端入口：管理端统一挂载在 `/admin`。
- 本次升级：前端入口、导航和核心页面已收敛为永和茉莉。
- 当前页面：工作台、订单管理、用户管理、服务类型管理、服务管理、账户管理、店铺管理、常见问题管理、客服在线。
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

## 后续升级记录

后续可以继续处理以下事项：

- 权限体系：按老板/运营/客服/财务分配菜单和操作权限。
- 小程序联动：后续 `/api/mini` 首页服务列表应从 `services` 和 `service_types` 表读取。

## 小程序接口

小程序端接口使用独立前缀 `/api/mini`，不走管理端 `/api/v1` 登录态。

当前已提供：

- `POST /api/mini/auth/wechat-login`：微信授权登录 mock。
- `POST /api/mini/auth/phone-code`：发送手机号验证码 mock。
- `POST /api/mini/auth/phone-login`：手机号验证码登录 mock。
- `GET /api/mini/user/profile`：获取当前用户信息。
- `GET /api/mini/appointment/home`：预约页首页数据。
- `GET /api/mini/stores/nearest`：获取最近的可服务门店。
- `GET /api/mini/service-categories`：服务类型枚举。
- `GET /api/mini/service-categories/:id/services`：获取某个服务类型下的预约服务。
- `GET /api/mini/services`：服务项目列表，支持 `category`、`keyword` 查询。
- `GET /api/mini/services/search`：搜索服务项目，参数同 `services`。
- `GET /api/mini/services/:id`：服务项目详情。
- `GET /api/mini/service-areas`：服务范围。
- `GET /api/mini/addresses`：常用地址列表。
- `POST /api/mini/addresses`：新增地址。
- `PUT /api/mini/addresses/:id/default`：设为默认地址。
- `DELETE /api/mini/addresses/:id`：删除地址。
- `GET /api/mini/service-targets`：老人/孩子服务对象列表。
- `POST /api/mini/service-targets`：新增服务对象。
- `PUT /api/mini/service-targets/:id/default`：设为默认服务对象。
- `DELETE /api/mini/service-targets/:id`：删除服务对象。
- `GET /api/mini/meal/pricing`：配菜价格配置。
- `GET /api/mini/meal/dishes`：菜品列表。
- `GET /api/mini/meal/dishes/:nameOrId`：菜品详情。
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

业务数据库仅支持 MySQL，通过环境变量配置：

```dotenv
MIS_DB_DRIVER=mysql
MIS_DB_DSN=root:password@tcp(127.0.0.1:3306)/yonhemoli_mis?charset=utf8mb4&parseTime=True&loc=Local
MIS_DB_AUTO_MIGRATE=false
```

后台入口为 `/admin`，管理端接口为 `/api/v1`，小程序接口为 `/api/mini`。

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
