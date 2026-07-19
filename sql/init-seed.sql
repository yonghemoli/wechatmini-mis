-- 永和护理 MIS 精简版基础数据
-- 前置条件：已执行 sql/init-schema.sql
-- 默认管理员：admin / admin123（首次登录后请立即修改密码）

INSERT INTO admins
    (username, password_hash, name, email, role_id, is_super_admin, status, last_login_at, created_at, updated_at)
VALUES
    ('admin', '$2a$10$H.QR8RkHNSh0tQhBRNkNe.BDYx5AXuQXu7maqomAIKaiGWstQrx..',
     '超级管理员', 'admin@yonghemoli.local', NULL, TRUE, 'active', NULL, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    name = VALUES(name), email = VALUES(email), is_super_admin = TRUE,
    status = 'active', updated_at = NOW();

INSERT INTO mini_service_categories
    (id, name, subtitle, description, icon_url, tags, enabled, sort, created_at, updated_at)
VALUES
    ('maternity', '月嫂', '母婴专业护理', '产妇照护、月子餐协助、新生儿日常护理与成长记录。', '', '["母婴护理","月子餐","新生儿"]', TRUE, 10, NOW(), NOW()),
    ('childcare', '育儿嫂', '科学育儿陪伴', '婴幼儿生活照料、辅食制作、早期互动与成长陪伴。', '', '["科学育儿","辅食","早教互动"]', TRUE, 20, NOW(), NOW()),
    ('nanny', '保姆', '安心家庭服务', '家庭日常保洁、烹饪和生活协助。', '', '["家庭保洁","烹饪","生活协助"]', TRUE, 30, NOW(), NOW()),
    ('hourly', '钟点工', '灵活按时服务', '按小时提供家庭清洁和临时家务服务。', '', '["小时服务","日常清洁"]', TRUE, 40, NOW(), NOW()),
    ('nursing', '护工', '专业照护陪伴', '老人和病患的生活照护、陪诊及康复协助。', '', '["老人照护","陪诊","康复协助"]', TRUE, 50, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    name = VALUES(name), subtitle = VALUES(subtitle), description = VALUES(description),
    tags = VALUES(tags), enabled = VALUES(enabled), sort = VALUES(sort), updated_at = NOW();

INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
VALUES
    ('mini.decoration.banners', '{"items":[]}', '首页轮播配置', NOW(), NOW()),
    ('mini.decoration.customer_service', '{"name":"小禾顾问","phone":"19994740191","avatarUrl":""}', '客服信息配置', NOW(), NOW()),
    ('mini.decoration.company', '{"logoUrl":"","name":"永和护理","address":"","introduction":"专注家庭护理服务。","serviceGuarantees":[{"icon":"verified","title":"身份核验","sub":"服务人员资料核验"},{"icon":"health","title":"健康信息","sub":"健康资料按授权展示"},{"icon":"support","title":"顾问跟进","sub":"专属顾问全程协助"}],"contactPhone":"19994740191"}', '公司信息配置', NOW(), NOW()),
    ('mini.business.agreement.privacy', '{"title":"隐私政策","version":"1.0","updatedAt":"2026-07-19T00:00:00+08:00","effectiveAt":"2026-07-19T00:00:00+08:00","intro":"","sections":[]}', '隐私政策', NOW(), NOW()),
    ('mini.business.agreement.service', '{"title":"用户服务协议","version":"1.0","updatedAt":"2026-07-19T00:00:00+08:00","effectiveAt":"2026-07-19T00:00:00+08:00","intro":"","sections":[]}', '用户服务协议', NOW(), NOW())
ON DUPLICATE KEY UPDATE value = VALUES(value), note = VALUES(note), updated_at = NOW();

INSERT INTO faqs (id, question, answer, category, sort_order, visible, created_at, updated_at)
VALUES
    (1, '提交预约后多久会有人联系？', '正常情况下客服会尽快联系，确认服务类型、时间和人员需求。', '预约', 10, TRUE, NOW(), NOW()),
    (2, '平台如何审核服务人员？', '平台会核验服务人员提交的身份、技能及必要健康资料。', '服务保障', 20, TRUE, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    question = VALUES(question), answer = VALUES(answer), category = VALUES(category),
    sort_order = VALUES(sort_order), visible = VALUES(visible), updated_at = NOW();
