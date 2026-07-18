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

INSERT INTO mini_service_categories (id, name, subtitle, description, icon_url, tags, enabled, sort, created_at, updated_at) VALUES
('maternity', '月嫂', '母婴专业护理', '产妇照护、月子餐协助、新生儿日常护理与成长记录。', '', '["母婴护理","月子餐","新生儿"]', TRUE, 10, NOW(), NOW()),
('childcare', '育儿嫂', '科学育儿陪伴', '婴幼儿生活照料、辅食制作、早期互动与成长陪伴。', '', '["科学育儿","辅食","早教互动"]', TRUE, 20, NOW(), NOW()),
('nanny', '保姆', '安心家庭服务', '家庭日常保洁、烹饪和生活协助。', '', '["家庭保洁","烹饪","生活协助"]', TRUE, 30, NOW(), NOW()),
('hourly', '钟点工', '灵活按时服务', '按小时提供家庭清洁和临时家务服务。', '', '["小时服务","日常清洁"]', TRUE, 40, NOW(), NOW()),
('nursing', '护工', '专业照护陪伴', '老人和病患的生活照护、陪诊及康复协助。', '', '["老人照护","陪诊","康复协助"]', TRUE, 50, NOW(), NOW())
ON DUPLICATE KEY UPDATE name = VALUES(name), subtitle = VALUES(subtitle), description = VALUES(description), tags = VALUES(tags), enabled = VALUES(enabled), sort = VALUES(sort), updated_at = NOW();

INSERT INTO app_configs (`key`, value, note, created_at, updated_at) VALUES
('mini.business.app', '{"consultant":{"name":"小禾顾问","phone":"19994740191","avatarUrl":""},"homeBanners":[],"trustItems":["身份核验","健康信息","顾问跟进"]}', '护理小程序公共配置', NOW(), NOW()),
('mini.business.about', '{"title":"永和护理","introduction":"专注家庭护理服务。","guarantees":["身份核验","健康信息","顾问跟进"]}', '公司介绍', NOW(), NOW()),
('mini.business.agreement.privacy', '{"title":"隐私政策","version":"1.0","updatedAt":"2026-07-18T00:00:00+08:00","effectiveAt":"2026-07-18T00:00:00+08:00","intro":"","sections":[]}', '隐私政策', NOW(), NOW()),
('mini.business.agreement.service', '{"title":"用户服务协议","version":"1.0","updatedAt":"2026-07-18T00:00:00+08:00","effectiveAt":"2026-07-18T00:00:00+08:00","intro":"","sections":[]}', '用户服务协议', NOW(), NOW())
ON DUPLICATE KEY UPDATE value = VALUES(value), note = VALUES(note), updated_at = NOW();
