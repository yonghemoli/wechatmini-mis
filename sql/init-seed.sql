-- Yonghemoli MIS initial seed data.
-- Run this manually after the application has created tables with GORM AutoMigrate.
-- Default admin password: admin123

INSERT INTO admins
    (username, password_hash, name, email, role_id, is_super_admin, status, last_login_at, created_at, updated_at)
SELECT
    'admin',
    '$2a$10$H.QR8RkHNSh0tQhBRNkNe.BDYx5AXuQXu7maqomAIKaiGWstQrx..',
    '超级管理员',
    'admin@yonghemoli.local',
    NULL,
    TRUE,
    'active',
    NULL,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM admins WHERE username = 'admin'
);

INSERT INTO users (id, avatar, nickname, total_spent, last_order_at, status, created_at, updated_at)
SELECT 'U10081', 'https://api.dicebear.com/9.x/initials/svg?seed=Lin', '林女士', 1264, '2026-06-20 09:18', 'active', '2026-05-18 09:00:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'U10081');

INSERT INTO users (id, avatar, nickname, total_spent, last_order_at, status, created_at, updated_at)
SELECT 'U10074', 'https://api.dicebear.com/9.x/initials/svg?seed=Zhou', '周先生', 756, '2026-06-19 21:40', 'active', '2026-05-07 10:30:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'U10074');

INSERT INTO users (id, avatar, nickname, total_spent, last_order_at, status, created_at, updated_at)
SELECT 'U10032', 'https://api.dicebear.com/9.x/initials/svg?seed=Zhao', '赵女士', 198, '2026-06-19 08:12', 'active', '2026-04-11 12:00:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'U10032');

INSERT INTO users (id, avatar, nickname, total_spent, last_order_at, status, created_at, updated_at)
SELECT 'U09951', 'https://api.dicebear.com/9.x/initials/svg?seed=Wu', '异常退款用户', 688, '2026-06-17 22:11', 'banned', '2026-03-28 14:00:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'U09951');

INSERT INTO orders (id, customer, phone, service, amount, status, source, appointment_at, staff, internal_note, close_reason, created_at, updated_at)
SELECT 'HS20260620001', '林女士', '138****3201', '深度保洁 4 小时', 328, 'pending_service', '分享', '2026-06-20 14:00', '王阿姨', '客户强调厨房油污重，需带强力清洁剂', '', '2026-06-20 09:18:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM orders WHERE id = 'HS20260620001');

INSERT INTO orders (id, customer, phone, service, amount, status, source, appointment_at, staff, internal_note, close_reason, created_at, updated_at)
SELECT 'HS20260620002', '周先生', '136****7781', '空调清洗 2 台', 236, 'pending_confirm', '搜索', '2026-06-20 10:30', '陈师傅', '师傅已完成，等待用户核销', '', '2026-06-19 21:40:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM orders WHERE id = 'HS20260620002');

INSERT INTO orders (id, customer, phone, service, amount, status, source, appointment_at, staff, internal_note, close_reason, created_at, updated_at)
SELECT 'HS20260619019', '赵女士', '139****8820', '日常保洁 3 小时', 198, 'exception', '扫一扫', '2026-06-19 16:00', '待改派', '服务人员临时请假，需客服回访改期', '', '2026-06-19 08:12:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM orders WHERE id = 'HS20260619019');

INSERT INTO orders (id, customer, phone, service, amount, status, source, appointment_at, staff, internal_note, close_reason, created_at, updated_at)
SELECT 'HS20260619012', '何先生', '135****4910', '玻璃清洁', 168, 'completed', '分享', '2026-06-19 11:00', '刘阿姨', '已评价五星', '', '2026-06-18 19:30:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM orders WHERE id = 'HS20260619012');

INSERT INTO orders (id, customer, phone, service, amount, status, source, appointment_at, staff, internal_note, close_reason, created_at, updated_at)
SELECT 'HS20260618008', '吴女士', '137****6632', '新房开荒', 688, 'refunded', '搜索', '2026-06-18 09:00', '未分配', '客户临时取消，已原路退款', '客户临时取消', '2026-06-17 22:11:00', NOW()
WHERE NOT EXISTS (SELECT 1 FROM orders WHERE id = 'HS20260618008');

INSERT INTO service_types (id, name, description, sort_order, status, created_at, updated_at)
SELECT 1, '家庭保洁', '日常保洁、深度保洁、新房开荒等家庭清洁服务', 10, 'active', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM service_types WHERE id = 1);

INSERT INTO service_types (id, name, description, sort_order, status, created_at, updated_at)
SELECT 2, '设备清洗', '空调、油烟机、洗衣机等家电清洗服务', 20, 'active', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM service_types WHERE id = 2);

INSERT INTO services (id, type_id, name, image, price, unit, description, visible, sort_order, created_at, updated_at)
SELECT 1, 1, '日常保洁', '/me.png', 66, '小时', '日常保洁', '适合家庭日常维护', '客厅、卧室、厨房、卫生间基础清洁', '66 元/小时起', '约 3 小时', '服务时长', '["3 小时","4 小时","5 小时"]', '["日常打扫","租房保洁","老人家庭"]', '["地面清洁","台面擦拭","厨卫基础清洁"]', '["预约下单","客服确认","阿姨上门","验收核销"]', '["不含高空作业","不含重油污深度清洁"]', '适合日常维护，覆盖客厅、卧室、厨房和卫生间基础清洁。', TRUE, 10, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM services WHERE id = 1);

INSERT INTO services (id, type_id, name, image, price, unit, description, visible, sort_order, created_at, updated_at)
SELECT 2, 1, '深度保洁', '/me.png', 82, '小时', '深度保洁', '适合深度清洁需求', '重点处理油污、水垢和卫生死角', '82 元/小时起', '约 4 小时', '服务时长', '["4 小时","5 小时","6 小时"]', '["油污处理","水垢处理","卫生死角"]', '["厨房深度清洁","卫浴深度清洁","窗框清洁"]', '["需求沟通","上门评估","深度清洁","验收"]', '["需提前确认清洁范围","部分顽固污渍可能加价"]', '重点处理油污、水垢和卫生死角。', TRUE, 20, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM services WHERE id = 2);

INSERT INTO services (id, type_id, name, image, price, unit, description, visible, sort_order, created_at, updated_at)
SELECT 3, 2, '空调清洗', '/me.png', 118, '台', '空调清洗', '适合换季前后', '清洗滤网、蒸发器和出风口', '预约后报价', '约 1 小时/台', '清洗数量', '["挂机 1 台","挂机 2 台","柜机 1 台"]', '["卧室空调","客厅空调","换季清洁"]', '["滤网清洗","蒸发器清洁","出风口擦洗"]', '["预约确认","师傅上门","现场清洗","验收完成"]', '["高空外机不在服务范围内","具体地址以客服确认为准"]', '拆洗滤网、蒸发器除菌、外壳清洁。', TRUE, 30, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM services WHERE id = 3);

INSERT INTO shops (id, name, contact_name, phone, address, business_hours, status, remark, created_at, updated_at)
SELECT 1, '南宁青秀店', '李店长', '0771-8888888', '广西南宁青秀区民族大道 100 号', '09:00-18:00', 'open', '主城区服务门店', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM shops WHERE id = 1);

INSERT INTO addresses (id, user_id, contact_name, phone, district, detail, tag, is_default, created_at, updated_at)
SELECT 'addr_1', 'U10081', '林女士', '138****3201', '广西 南宁 青秀区', '3 栋 2 单元 801', '家', TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM addresses WHERE id = 'addr_1');

INSERT INTO service_targets (id, user_id, name, category, relation, age, note, is_default, created_at, updated_at)
SELECT 'target_1', 'U10081', '父亲', 'eldercare', '父母', '70-79 岁', '行动较慢', TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM service_targets WHERE id = 'target_1');

INSERT INTO dishes (id, name, scene, tag, price, ingredients, video_title, video_url, comments, created_at, updated_at)
SELECT 'tomato-egg', '番茄炒蛋', '酸甜开胃，适合日常晚餐', '家常', 12, '["番茄 300g","鸡蛋 3 个"]', '番茄炒蛋 8 分钟快手做法', 'https://example.com/tomato-egg', '["孩子很爱吃"]', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM dishes WHERE id = 'tomato-egg');

INSERT INTO dishes (id, name, scene, tag, price, ingredients, video_title, video_url, comments, created_at, updated_at)
SELECT 'pepper-pork', '青椒肉丝', '下饭快手菜', '家常', 16, '["青椒 250g","猪肉 200g"]', '青椒肉丝家常做法', 'https://example.com/pepper-pork', '["适合工作日晚餐"]', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM dishes WHERE id = 'pepper-pork');

INSERT INTO dishes (id, name, scene, tag, price, ingredients, video_title, video_url, comments, created_at, updated_at)
SELECT 'seaweed-egg-soup', '紫菜蛋花汤', '清淡快手，适合搭配主菜', '汤品', 10, '["紫菜 20g","鸡蛋 2 个"]', '紫菜蛋花汤家常做法', 'https://example.com/seaweed-egg-soup', '["老人孩子都适合"]', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM dishes WHERE id = 'seaweed-egg-soup');

INSERT INTO meal_packages (id, user_id, package_type, name, scene, price, dishes, created_at, updated_at)
SELECT 'pkg_1', NULL, 'official', '三菜一汤套餐', '适合 2-3 人工作日晚餐', 68, '["番茄炒蛋","青椒肉丝","紫菜蛋花汤"]', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM meal_packages WHERE id = 'pkg_1');

INSERT INTO meal_packages (id, user_id, package_type, name, scene, price, dishes, created_at, updated_at)
SELECT 'custom_pkg_1', 'U10081', 'custom', '我家的晚餐组合', '少油少盐', 24, '["番茄炒蛋","青椒肉丝"]', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM meal_packages WHERE id = 'custom_pkg_1');

INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
SELECT 'mini.home', '{"appName":"永和护理","banners":["/assets/banners/home-care.png","/assets/banners/meal-prep.png"],"notice":"服务预约后，客服将在 10 分钟内确认上门时间。"}', '小程序首页配置', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM app_configs WHERE `key` = 'mini.home');

INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
SELECT 'mini.appointment.tabs', '[{"id":"all","name":"所有"}]', '预约页 tab 配置', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM app_configs WHERE `key` = 'mini.appointment.tabs');

INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
SELECT 'mini.service_areas', '{"city":"南宁","districts":["青秀区","兴宁区","西乡塘区","良庆区","江南区"],"notes":["具体地址以客服确认为准"]}', '服务范围配置', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM app_configs WHERE `key` = 'mini.service_areas');

INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
SELECT 'mini.meal_pricing', '{"dishPrice":12,"deliveryFee":8}', '配菜价格配置', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM app_configs WHERE `key` = 'mini.meal_pricing');

INSERT INTO faqs (id, question, answer, category, sort_order, visible, created_at, updated_at)
SELECT 1, '下单后多久会有人联系？', '正常情况下客服会在 10 分钟内联系确认服务时间和地址。', '下单', 10, TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM faqs WHERE id = 1);

INSERT INTO faqs (id, question, answer, category, sort_order, visible, created_at, updated_at)
SELECT 2, '服务不满意可以退款吗？', '请在服务完成后及时联系客服，运营会根据情况安排返工或退款。', '售后', 20, TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM faqs WHERE id = 2);

INSERT INTO chat_sessions (id, user_id, user_name, user_avatar, status, last_message, unread_count, created_at, updated_at)
SELECT 'chat_demo_1', 'U10081', '林女士', 'https://api.dicebear.com/9.x/initials/svg?seed=Lin', 'open', '想确认今天下午的保洁是否能提前到 13 点', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM chat_sessions WHERE id = 'chat_demo_1');

INSERT INTO chat_messages (session_id, sender, msg_type, content, is_read, created_at)
SELECT 'chat_demo_1', 'user', 'text', '想确认今天下午的保洁是否能提前到 13 点', FALSE, NOW()
WHERE NOT EXISTS (SELECT 1 FROM chat_messages WHERE session_id = 'chat_demo_1' AND content = '想确认今天下午的保洁是否能提前到 13 点');
