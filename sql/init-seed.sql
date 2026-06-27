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
SELECT 1, 1, '日常保洁', '/me.png', 66, '小时', '适合日常维护，覆盖客厅、卧室、厨房和卫生间基础清洁。', TRUE, 10, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM services WHERE id = 1);

INSERT INTO services (id, type_id, name, image, price, unit, description, visible, sort_order, created_at, updated_at)
SELECT 2, 1, '深度保洁', '/me.png', 82, '小时', '重点处理油污、水垢和卫生死角。', TRUE, 20, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM services WHERE id = 2);

INSERT INTO services (id, type_id, name, image, price, unit, description, visible, sort_order, created_at, updated_at)
SELECT 3, 2, '空调清洗', '/me.png', 118, '台', '拆洗滤网、蒸发器除菌、外壳清洁。', TRUE, 30, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM services WHERE id = 3);

INSERT INTO shops (id, name, contact_name, phone, address, business_hours, status, remark, created_at, updated_at)
SELECT 1, '南宁青秀店', '李店长', '0771-8888888', '广西南宁青秀区民族大道 100 号', '09:00-18:00', 'open', '主城区服务门店', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM shops WHERE id = 1);

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
