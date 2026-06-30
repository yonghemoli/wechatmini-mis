-- Yonghemoli MIS initial seed data.
-- Run this manually after sql/init-schema.sql has created or upgraded tables.
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

INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
SELECT 'mini.home', '{"appName":"永和护理","banners":["/assets/banners/home-care.png","/assets/banners/meal-prep.png"],"notice":"服务预约后，客服将在 10 分钟内确认上门时间。"}', '小程序首页配置', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM app_configs WHERE `key` = 'mini.home');

INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
SELECT 'mini.appointment.tabs', '[{"id":"all","name":"所有"}]', '预约页 tab 配置', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM app_configs WHERE `key` = 'mini.appointment.tabs');

INSERT INTO faqs (id, question, answer, category, sort_order, visible, created_at, updated_at)
SELECT 1, '下单后多久会有人联系？', '正常情况下客服会在 10 分钟内联系确认服务时间和地址。', '下单', 10, TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM faqs WHERE id = 1);

INSERT INTO faqs (id, question, answer, category, sort_order, visible, created_at, updated_at)
SELECT 2, '服务不满意可以退款吗？', '请在服务完成后及时联系客服，运营会根据情况安排返工或退款。', '售后', 20, TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM faqs WHERE id = 2);
