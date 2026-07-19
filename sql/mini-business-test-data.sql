-- 永和护理小程序 / MIS 新业务基础测试数据
-- 特点：单文件可直接执行并可重复导入；固定 ID 仅供开发、联调和验收环境使用，禁止直接用于生产。

-- 先建立新需求使用的业务表，避免测试库尚未执行新版 init-schema.sql 时出现 1146。
CREATE TABLE IF NOT EXISTS mini_service_categories (
    id VARCHAR(32) NOT NULL,
    name VARCHAR(64) NOT NULL,
    subtitle VARCHAR(128) NOT NULL DEFAULT '',
    description TEXT NOT NULL,
    icon_url VARCHAR(512) NOT NULL DEFAULT '',
    tags TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sort INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_mini_service_categories_enabled (enabled),
    KEY idx_mini_service_categories_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS caregivers (
    id VARCHAR(40) NOT NULL,
    application_id VARCHAR(32) DEFAULT NULL,
    contact_phone VARCHAR(32) NOT NULL DEFAULT '',
    avatar_url VARCHAR(512) NOT NULL DEFAULT '',
    name VARCHAR(64) NOT NULL,
    age INT NOT NULL,
    experience_years INT NOT NULL DEFAULT 0,
    origin VARCHAR(128) NOT NULL DEFAULT '',
    service_ids TEXT NOT NULL,
    jobs TEXT NOT NULL,
    availability_status VARCHAR(32) NOT NULL,
    rating DECIMAL(2,1) NOT NULL DEFAULT 0,
    service_count INT NOT NULL DEFAULT 0,
    recommended BOOLEAN NOT NULL DEFAULT FALSE,
    introduction TEXT NOT NULL,
    education VARCHAR(64) NOT NULL DEFAULT '',
    ethnicity VARCHAR(64) NOT NULL DEFAULT '',
    zodiac VARCHAR(32) NOT NULL DEFAULT '',
    skills TEXT NOT NULL,
    certificates TEXT NOT NULL,
    identity_verified BOOLEAN NOT NULL DEFAULT FALSE,
    physical_exam_verified BOOLEAN NOT NULL DEFAULT FALSE,
    medical_report_image_urls TEXT NOT NULL,
    personal_info TEXT NOT NULL,
    work_history TEXT NOT NULL,
    photo_urls TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    source VARCHAR(20) NOT NULL DEFAULT 'ADMIN',
    sort INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_caregivers_application_id (application_id),
    KEY idx_caregivers_name (name),
    KEY idx_caregivers_origin (origin),
    KEY idx_caregivers_availability (availability_status),
    KEY idx_caregivers_recommended (recommended),
    KEY idx_caregivers_status (status),
    KEY idx_caregivers_source (source),
    KEY idx_caregivers_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS demands (
    id VARCHAR(32) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    service_id VARCHAR(32) NOT NULL,
    service_name VARCHAR(64) NOT NULL,
    caregiver_id VARCHAR(40) DEFAULT NULL,
    caregiver_name VARCHAR(64) NOT NULL DEFAULT '',
    contact_name VARCHAR(64) NOT NULL DEFAULT '',
    contact_phone VARCHAR(32) NOT NULL,
    requirements TEXT NOT NULL,
    source VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_CONTACT',
    assigned_admin_id BIGINT UNSIGNED DEFAULT NULL,
    idempotency_key VARCHAR(128) DEFAULT NULL,
    submission_scope VARCHAR(128) DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_demands_user_id (user_id),
    KEY idx_demands_service_id (service_id),
    KEY idx_demands_caregiver_id (caregiver_id),
    KEY idx_demands_phone (contact_phone),
    KEY idx_demands_status (status),
    KEY idx_demands_source (source),
    KEY idx_demands_idempotency (submission_scope, idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS resumes (
    id VARCHAR(32) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    intention_service_id VARCHAR(32) NOT NULL,
    service_name VARCHAR(64) NOT NULL,
    work_status VARCHAR(32) NOT NULL,
    experience_range VARCHAR(32) NOT NULL,
    entry_year INT NOT NULL,
    contact_name VARCHAR(64) NOT NULL DEFAULT '',
    contact_phone VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_CONTACT',
    assigned_admin_id BIGINT UNSIGNED DEFAULT NULL,
    idempotency_key VARCHAR(128) DEFAULT NULL,
    submission_scope VARCHAR(128) DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_resumes_user_id (user_id),
    KEY idx_resumes_service_id (intention_service_id),
    KEY idx_resumes_phone (contact_phone),
    KEY idx_resumes_status (status),
    KEY idx_resumes_idempotency (submission_scope, idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS business_status_histories (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    entity_type VARCHAR(32) NOT NULL,
    entity_id VARCHAR(32) NOT NULL,
    from_status VARCHAR(32) NOT NULL,
    to_status VARCHAR(32) NOT NULL,
    operator_id BIGINT UNSIGNED DEFAULT NULL,
    note TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_business_history_entity (entity_type, entity_id),
    KEY idx_business_history_operator (operator_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

START TRANSACTION;

-- 1. 护理服务分类
INSERT INTO mini_service_categories
    (id, name, subtitle, description, icon_url, tags, enabled, sort, created_at, updated_at)
VALUES
    ('maternity', '月嫂', '母婴专业护理', '产妇照护、月子餐协助、新生儿日常护理与成长记录。',
     'https://cdn.example.com/services/maternity.png', '["母婴护理","月子餐","新生儿"]', TRUE, 10, NOW(), NOW()),
    ('childcare', '育儿嫂', '科学育儿陪伴', '婴幼儿生活照料、辅食制作、早期互动与成长陪伴。',
     'https://cdn.example.com/services/childcare.png', '["科学育儿","辅食","早教互动"]', TRUE, 20, NOW(), NOW()),
    ('nanny', '保姆', '安心家庭服务', '家庭日常保洁、烹饪和生活协助。',
     'https://cdn.example.com/services/nanny.png', '["家庭保洁","烹饪","生活协助"]', TRUE, 30, NOW(), NOW()),
    ('hourly', '钟点工', '灵活按时服务', '按小时提供家庭清洁和临时家务服务。',
     'https://cdn.example.com/services/hourly.png', '["小时服务","日常清洁"]', TRUE, 40, NOW(), NOW()),
    ('nursing', '护工', '专业照护陪伴', '老人和病患的生活照护、陪诊及康复协助。',
     'https://cdn.example.com/services/nursing.png', '["老人照护","陪诊","康复协助"]', TRUE, 50, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    subtitle = VALUES(subtitle),
    description = VALUES(description),
    icon_url = VALUES(icon_url),
    tags = VALUES(tags),
    enabled = VALUES(enabled),
    sort = VALUES(sort),
    updated_at = NOW();

-- 2. 小程序测试用户。手机号为虚构测试号，仅用于本地数据关联。
INSERT INTO users
    (id, openid, avatar, nickname, phone, signature,
     last_login_at, status, created_at, updated_at)
VALUES
    ('usr_test_1001', 'test_openid_1001', 'https://cdn.example.com/users/test-1001.png',
     '永和测试用户', '13800001001', '护理服务联调账号',
     '2026-07-18 10:00:00', 'active', NOW(), NOW()),
    ('usr_test_1002', 'test_openid_1002', '',
     '求职测试用户', '13800001002', '',
     '2026-07-18 10:05:00', 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    avatar = VALUES(avatar),
    nickname = VALUES(nickname),
    signature = VALUES(signature),
    status = VALUES(status),
    updated_at = NOW();

-- 3. 服务人员。包含推荐、非推荐、完成和阿姨自提交草稿，便于验证筛选规则。
INSERT INTO caregivers
    (id, application_id, avatar_url, name, age, experience_years, origin, service_ids, jobs,
     availability_status, rating, service_count, recommended, introduction,
     education, ethnicity, zodiac, skills, certificates, identity_verified,
     physical_exam_verified, medical_report_image_urls, personal_info,
     work_history, photo_urls, status, source, sort, created_at, updated_at)
VALUES
    (
        'auntie-01',
        NULL,
        'https://cdn.example.com/caregivers/auntie-01.jpg',
        '覃阿姨', 40, 10, '广西南宁',
        '["maternity","childcare","nursing"]',
        '["月嫂","育儿嫂","护工"]',
        'AVAILABLE_NOW', 4.9, 128, TRUE,
        '性格温和，做事细致，擅长新生儿护理和产妇照护。',
        '高中', '壮族', '猴',
        '["新生儿护理","月子餐","辅食制作","早期互动"]',
        '[{"id":"cert_1001","name":"母婴护理证","verified":true,"imageUrls":[]},{"id":"cert_1002","name":"育婴师证","verified":true,"imageUrls":[]}]',
        TRUE, TRUE,
        '["https://cdn.example.com/reports/test-auntie-01.jpg"]',
        '{"heightCm":160,"weightKg":55,"bloodType":"O","gender":"FEMALE","maritalStatus":"MARRIED","religion":"无","languages":["普通话","壮语"],"liveInAvailable":true}',
        '[{"startDate":"2021-03","endDate":null,"periodText":"2021.03—至今","role":"母婴护理师 · 南宁家庭","description":"负责新生儿日常护理、产妇照护及月子餐协助。"},{"startDate":"2018-06","endDate":"2021-02","periodText":"2018.06—2021.02","role":"育儿嫂 · 柳州家庭","description":"负责婴幼儿生活照料、辅食制作和成长陪伴。"}]',
        '["https://cdn.example.com/caregivers/auntie-01-work-01.jpg","https://cdn.example.com/caregivers/auntie-01-work-02.jpg"]',
        'COMPLETED', 'ADMIN', 100, NOW(), NOW()
    ),
    (
        'auntie-02',
        NULL,
        'https://cdn.example.com/caregivers/auntie-02.jpg',
        '黄阿姨', 46, 15, '广西玉林',
        '["nanny","hourly"]',
        '["保姆","钟点工"]',
        'AVAILABLE_IN_3_DAYS', 4.8, 96, TRUE,
        '擅长家庭保洁、家常菜和收纳整理，工作认真守时。',
        '初中', '汉族', '鼠',
        '["家庭保洁","家常菜","收纳整理"]',
        '[{"id":"cert_2001","name":"家政服务员证","verified":true,"imageUrls":[]}]',
        TRUE, FALSE, '[]',
        '{"heightCm":158,"weightKg":53,"bloodType":"A","gender":"FEMALE","maritalStatus":"MARRIED","religion":"无","languages":["普通话","粤语"],"liveInAvailable":false}',
        '[{"startDate":"2019-01","endDate":null,"periodText":"2019.01—至今","role":"家庭服务员 · 南宁家庭","description":"负责家庭保洁、三餐制作和日常收纳。"}]',
        '["https://cdn.example.com/caregivers/auntie-02-work-01.jpg"]',
        'COMPLETED', 'ADMIN', 90, NOW(), NOW()
    ),
    (
        'auntie-03',
        NULL,
        'https://cdn.example.com/caregivers/auntie-03.jpg',
        '李阿姨', 52, 12, '广西桂林',
        '["nursing"]',
        '["护工"]',
        'OPEN_TO_OPPORTUNITIES', 4.7, 74, FALSE,
        '具有老人生活照护和陪诊经验，耐心细致。',
        '高中', '汉族', '龙',
        '["老人照护","陪诊","康复协助"]',
        '[{"id":"cert_3001","name":"养老护理员证","verified":true,"imageUrls":[]}]',
        TRUE, TRUE,
        '["https://cdn.example.com/reports/test-auntie-03.jpg"]',
        '{"heightCm":162,"weightKg":58,"bloodType":"B","gender":"FEMALE","maritalStatus":"MARRIED","religion":"无","languages":["普通话","桂柳话"],"liveInAvailable":true}',
        '[{"startDate":"2020-05","endDate":null,"periodText":"2020.05—至今","role":"养老护理员 · 南宁家庭","description":"负责老人日常照护、用药提醒和陪诊。"}]',
        '["https://cdn.example.com/caregivers/auntie-03-work-01.jpg"]',
        'COMPLETED', 'ADMIN', 80, NOW(), NOW()
    ),
    (
        'auntie-draft-01',
        'RTEST202607180001',
        'https://cdn.example.com/caregivers/auntie-draft-01.jpg',
        '测试待发布人员', 35, 5, '广西北海',
        '["childcare"]', '["育儿嫂"]',
        'UNAVAILABLE', 0.0, 0, FALSE,
        '此记录用于验证草稿人员不会出现在小程序公开接口。',
        '高中', '汉族', '兔',
        '["婴幼儿照料"]', '[]',
        FALSE, FALSE, '[]',
        '{"heightCm":159,"weightKg":52,"bloodType":"","gender":"FEMALE","maritalStatus":"MARRIED","religion":"无","languages":["普通话"],"liveInAvailable":false}',
        '[]', '[]',
        'DRAFT', 'SELF_SUBMITTED', 10, NOW(), NOW()
    )
ON DUPLICATE KEY UPDATE
    avatar_url = VALUES(avatar_url),
    application_id = VALUES(application_id),
    name = VALUES(name),
    age = VALUES(age),
    experience_years = VALUES(experience_years),
    origin = VALUES(origin),
    service_ids = VALUES(service_ids),
    jobs = VALUES(jobs),
    availability_status = VALUES(availability_status),
    rating = VALUES(rating),
    service_count = VALUES(service_count),
    recommended = VALUES(recommended),
    introduction = VALUES(introduction),
    education = VALUES(education),
    ethnicity = VALUES(ethnicity),
    zodiac = VALUES(zodiac),
    skills = VALUES(skills),
    certificates = VALUES(certificates),
    identity_verified = VALUES(identity_verified),
    physical_exam_verified = VALUES(physical_exam_verified),
    medical_report_image_urls = VALUES(medical_report_image_urls),
    personal_info = VALUES(personal_info),
    work_history = VALUES(work_history),
    photo_urls = VALUES(photo_urls),
    status = VALUES(status),
    source = VALUES(source),
    sort = VALUES(sort),
    updated_at = NOW();

-- 4. 客户服务需求。覆盖待联系、已联系、匹配中和关闭状态。
INSERT INTO demands
    (id, user_id, service_id, service_name, caregiver_id, caregiver_name, contact_name,
     contact_phone, requirements, source, status, assigned_admin_id,
     idempotency_key, submission_scope, created_at, updated_at)
VALUES
    ('DTEST202607180001', 'usr_test_1001', 'nursing', '护工', 'auntie-03', '李阿姨', '张女士',
     '13800001001', '工作日上午照护老人，需要陪诊和康复协助。',
     'CAREGIVER_DETAIL', 'PENDING_CONTACT', NULL,
     'test-demand-001', 'user:usr_test_1001', '2026-07-18 10:30:00', '2026-07-18 10:30:00'),
    ('DTEST202607180002', '', 'maternity', '月嫂', 'auntie-01', '覃阿姨', '刘女士',
     '13800001003', '预产期在八月，希望提前了解住家月嫂服务。',
     'HOME_SERVICE', 'CONTACTED', NULL,
     'test-demand-002', 'phone:13800001003', '2026-07-18 11:00:00', '2026-07-18 11:20:00'),
    ('DTEST202607180003', '', 'nanny', '保姆', NULL, '', '陈先生',
     '13800001004', '需要白班保姆，负责做饭和日常保洁。',
     'SERVICE_LIST', 'MATCHING', NULL,
     'test-demand-003', 'phone:13800001004', '2026-07-17 09:00:00', '2026-07-18 09:10:00'),
    ('DTEST202607180004', '', 'hourly', '钟点工', 'auntie-02', '黄阿姨', '杨女士',
     '13800001005', '每周两次家庭保洁，每次三小时。',
     'OTHER', 'CLOSED', NULL,
     'test-demand-004', 'phone:13800001005', '2026-07-16 14:00:00', '2026-07-18 09:30:00')
ON DUPLICATE KEY UPDATE
    user_id = VALUES(user_id),
    service_id = VALUES(service_id),
    service_name = VALUES(service_name),
    caregiver_id = VALUES(caregiver_id),
    caregiver_name = VALUES(caregiver_name),
    contact_name = VALUES(contact_name),
    contact_phone = VALUES(contact_phone),
    requirements = VALUES(requirements),
    source = VALUES(source),
    status = VALUES(status),
    idempotency_key = VALUES(idempotency_key),
    submission_scope = VALUES(submission_scope),
    updated_at = VALUES(updated_at);

-- 5. 求职简历。入行年份与从业年限枚举保持一致。
INSERT INTO resumes
    (id, user_id, intention_service_id, service_name, work_status,
     experience_range, entry_year, contact_phone, status, assigned_admin_id,
     idempotency_key, submission_scope, created_at, updated_at)
VALUES
    ('RTEST202607180001', 'usr_test_1002', 'maternity', '月嫂',
     'AVAILABLE_NOW', 'YEAR_5_TO_10', 2018, '13800001002',
     'PENDING_CONTACT', NULL, 'test-resume-001', 'user:usr_test_1002',
     '2026-07-18 10:40:00', '2026-07-18 10:40:00'),
    ('RTEST202607180002', '', 'nursing', '护工',
     'AVAILABLE_IN_3_DAYS', 'MORE_THAN_10_YEARS', 2010, '13800001006',
     'VERIFYING', NULL, 'test-resume-002', 'phone:13800001006',
     '2026-07-17 15:00:00', '2026-07-18 10:00:00'),
    ('RTEST202607180003', '', 'nanny', '保姆',
     'OPEN_TO_OPPORTUNITIES', 'YEAR_3_TO_5', 2022, '13800001007',
     'APPROVED', NULL, 'test-resume-003', 'phone:13800001007',
     '2026-07-15 09:00:00', '2026-07-18 08:30:00')
ON DUPLICATE KEY UPDATE
    user_id = VALUES(user_id),
    intention_service_id = VALUES(intention_service_id),
    service_name = VALUES(service_name),
    work_status = VALUES(work_status),
    experience_range = VALUES(experience_range),
    entry_year = VALUES(entry_year),
    contact_phone = VALUES(contact_phone),
    status = VALUES(status),
    idempotency_key = VALUES(idempotency_key),
    submission_scope = VALUES(submission_scope),
    updated_at = VALUES(updated_at);

-- 6. 状态流转历史。先删除固定测试实体的历史，保证重复执行后数据不累加。
DELETE FROM business_status_histories
WHERE (entity_type = 'DEMAND' AND entity_id IN (
           'DTEST202607180002', 'DTEST202607180003', 'DTEST202607180004'
       ))
   OR (entity_type = 'RESUME' AND entity_id IN (
           'RTEST202607180002', 'RTEST202607180003'
       ));

INSERT INTO business_status_histories
    (entity_type, entity_id, from_status, to_status, operator_id, note, created_at)
VALUES
    ('DEMAND', 'DTEST202607180002', 'PENDING_CONTACT', 'CONTACTED', NULL,
     '测试记录：顾问已电话联系客户。', '2026-07-18 11:20:00'),
    ('DEMAND', 'DTEST202607180003', 'PENDING_CONTACT', 'CONTACTED', NULL,
     '测试记录：已确认客户服务时间。', '2026-07-17 10:00:00'),
    ('DEMAND', 'DTEST202607180003', 'CONTACTED', 'MATCHING', NULL,
     '测试记录：正在匹配合适服务人员。', '2026-07-18 09:10:00'),
    ('DEMAND', 'DTEST202607180004', 'PENDING_CONTACT', 'CONTACTED', NULL,
     '测试记录：已联系。', '2026-07-16 15:00:00'),
    ('DEMAND', 'DTEST202607180004', 'CONTACTED', 'MATCHING', NULL,
     '测试记录：已安排服务人员。', '2026-07-17 09:00:00'),
    ('DEMAND', 'DTEST202607180004', 'MATCHING', 'CLOSED', NULL,
     '测试记录：客户暂不需要服务。', '2026-07-18 09:30:00'),
    ('RESUME', 'RTEST202607180002', 'PENDING_CONTACT', 'CONTACTED', NULL,
     '测试记录：招聘顾问已联系。', '2026-07-17 16:00:00'),
    ('RESUME', 'RTEST202607180002', 'CONTACTED', 'VERIFYING', NULL,
     '测试记录：正在核验证书和工作经历。', '2026-07-18 10:00:00'),
    ('RESUME', 'RTEST202607180003', 'PENDING_CONTACT', 'CONTACTED', NULL,
     '测试记录：已完成首次沟通。', '2026-07-15 10:00:00'),
    ('RESUME', 'RTEST202607180003', 'CONTACTED', 'VERIFYING', NULL,
     '测试记录：资料核验中。', '2026-07-16 09:00:00'),
    ('RESUME', 'RTEST202607180003', 'VERIFYING', 'APPROVED', NULL,
     '测试记录：资料核验通过。', '2026-07-18 08:30:00');

-- 7. 首页运营配置、公司介绍和协议内容
INSERT INTO app_configs (`key`, value, note, created_at, updated_at)
VALUES
    (
        'mini.decoration.banners',
        '{"items":[{"id":"banner_maternity","imageUrl":"https://cdn.example.com/banners/maternity.jpg","kicker":"专业月嫂","title":"安心月子，从专业照护开始","description":"严选经验服务人员，顾问全程跟进","actionType":"DEMAND","actionValue":"maternity","sort":10},{"id":"banner_nursing","imageUrl":"https://cdn.example.com/banners/nursing.jpg","kicker":"老人照护","title":"多一份陪伴，多一份安心","description":"生活照料、陪诊与康复协助","actionType":"DEMAND","actionValue":"nursing","sort":20}]}',
        '宣传图测试数据', NOW(), NOW()
    ),
    (
        'mini.decoration.customer_service',
        '{"name":"小禾顾问","phone":"19900001001","avatarUrl":"https://cdn.example.com/consultants/xiao-he.jpg"}',
        '客服信息测试数据', NOW(), NOW()
    ),
    (
        'mini.decoration.company',
        '{"logoUrl":"https://cdn.example.com/company/logo.png","name":"永和护理","address":"广西南宁市青秀区测试路 1 号","introduction":"永和护理专注于母婴、育儿、家政和老人照护服务，为家庭提供经过资料核验的服务人员和顾问跟进。","serviceGuarantees":[{"icon":"verified","title":"身份资料核验","sub":"服务人员资料已核验"},{"icon":"health","title":"健康信息展示","sub":"按授权范围展示"},{"icon":"support","title":"专属顾问跟进","sub":"持续跟进服务需求"},{"icon":"return","title":"服务过程回访","sub":"服务后持续回访"}],"contactPhone":"19900001001"}',
        '公司信息测试数据', NOW(), NOW()
    ),
    (
        'mini.business.agreement.privacy',
        '{"title":"隐私政策（测试版）","version":"1.0-test","updatedAt":"2026-07-18T00:00:00+08:00","effectiveAt":"2026-07-18T00:00:00+08:00","intro":"本内容仅用于开发联调，不可作为正式法律文本发布。","sections":[{"title":"信息收集","content":"为处理服务咨询或求职申请，我们会收集用户主动提交的联系方式和需求信息。"},{"title":"信息使用","content":"收集的信息仅用于服务匹配、顾问联系和必要的业务跟进。"}]}',
        '隐私政策测试数据', NOW(), NOW()
    ),
    (
        'mini.business.agreement.service',
        '{"title":"用户服务协议（测试版）","version":"1.0-test","updatedAt":"2026-07-18T00:00:00+08:00","effectiveAt":"2026-07-18T00:00:00+08:00","intro":"本内容仅用于开发联调，不可作为正式法律文本发布。","sections":[{"title":"服务说明","content":"用户提交需求仅代表发起服务咨询，不代表已经下单或预约成功。"},{"title":"用户义务","content":"用户应提供真实、准确且必要的联系和需求信息。"}]}',
        '用户服务协议测试数据', NOW(), NOW()
    )
ON DUPLICATE KEY UPDATE
    value = VALUES(value),
    note = VALUES(note),
    updated_at = NOW();

COMMIT;

-- 快速验收查询
SELECT id, name, enabled, sort
FROM mini_service_categories
ORDER BY sort ASC;

SELECT id, name, availability_status, recommended, status, source, sort
FROM caregivers
ORDER BY sort DESC, id ASC;

SELECT id, service_name, caregiver_name, status, created_at
FROM demands
ORDER BY created_at DESC;

SELECT id, service_name, work_status, experience_range, status, created_at
FROM resumes
ORDER BY created_at DESC;
