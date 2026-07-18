-- Yonghemoli MIS schema.
-- Safe to run repeatedly when initializing a MySQL database or upgrading old tables.

CREATE TABLE IF NOT EXISTS admins (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(128) NOT NULL,
    role_id BIGINT NULL,
    is_super_admin BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_admins_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(32) NOT NULL,
    openid VARCHAR(128) DEFAULT NULL,
    avatar VARCHAR(255) NOT NULL DEFAULT '',
    nickname VARCHAR(64) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    signature VARCHAR(255) NOT NULL DEFAULT '',
    total_spent INT NOT NULL DEFAULT 0,
    last_order_at VARCHAR(32) DEFAULT NULL,
    last_login_at VARCHAR(32) DEFAULT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_users_openid (openid),
    UNIQUE KEY idx_users_phone (phone),
    KEY idx_users_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(32) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    customer VARCHAR(64) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    service VARCHAR(128) NOT NULL,
    amount INT NOT NULL DEFAULT 0,
    status VARCHAR(32) NOT NULL,
    source VARCHAR(32) NOT NULL,
    appointment_at VARCHAR(32) DEFAULT NULL,
    staff VARCHAR(64) NOT NULL DEFAULT '',
    internal_note TEXT,
    close_reason TEXT,
    mini_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_orders_user_id (user_id),
    KEY idx_orders_status (status),
    KEY idx_orders_source (source),
    KEY idx_orders_appointment_at (appointment_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS service_types (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS services (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    type_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(128) NOT NULL,
    image VARCHAR(255) NOT NULL DEFAULT '',
    price INT NOT NULL DEFAULT 0,
    unit VARCHAR(32) NOT NULL DEFAULT '小时',
    title VARCHAR(128) NOT NULL DEFAULT '',
    scene VARCHAR(255) NOT NULL DEFAULT '',
    summary TEXT,
    price_text VARCHAR(64) NOT NULL DEFAULT '',
    duration_text VARCHAR(64) NOT NULL DEFAULT '',
    requirement_label VARCHAR(64) NOT NULL DEFAULT '',
    requirement_options TEXT NOT NULL,
    suitable_for TEXT NOT NULL,
    scope TEXT NOT NULL,
    process TEXT NOT NULL,
    notes TEXT NOT NULL,
    description TEXT,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_services_type_id (type_id),
    KEY idx_services_visible (visible)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shops (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    contact_name VARCHAR(64) DEFAULT NULL,
    phone VARCHAR(32) DEFAULT NULL,
    address VARCHAR(255) DEFAULT NULL,
    business_hours VARCHAR(128) DEFAULT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    remark TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS addresses (
    id VARCHAR(32) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    contact_name VARCHAR(64) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    district VARCHAR(128) NOT NULL,
    detail VARCHAR(255) NOT NULL,
    tag VARCHAR(32) NOT NULL DEFAULT '',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_addresses_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS service_targets (
    id VARCHAR(32) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    name VARCHAR(64) NOT NULL,
    category VARCHAR(32) NOT NULL,
    relation VARCHAR(32) NOT NULL,
    age VARCHAR(32) NOT NULL,
    note VARCHAR(255) NOT NULL DEFAULT '',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_service_targets_user_id (user_id),
    KEY idx_service_targets_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS dishes (
    id VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    scene VARCHAR(255) NOT NULL,
    tag VARCHAR(32) NOT NULL DEFAULT '',
    price INT NOT NULL DEFAULT 0,
    ingredients TEXT NOT NULL,
    video_title VARCHAR(255) NOT NULL DEFAULT '',
    video_url VARCHAR(255) NOT NULL DEFAULT '',
    comments TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_dishes_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS meal_packages (
    id VARCHAR(64) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    package_type VARCHAR(20) NOT NULL DEFAULT 'official',
    name VARCHAR(128) NOT NULL,
    scene VARCHAR(255) NOT NULL,
    price INT NOT NULL DEFAULT 0,
    dishes TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_meal_packages_user_id (user_id),
    KEY idx_meal_packages_type (package_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS app_configs (
    `key` VARCHAR(128) NOT NULL,
    value LONGTEXT NOT NULL,
    note VARCHAR(255) NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS faqs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    question VARCHAR(255) NOT NULL,
    answer TEXT NOT NULL,
    category VARCHAR(64) DEFAULT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_faqs_category (category),
    KEY idx_faqs_visible (visible)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_sessions (
    id VARCHAR(40) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    user_name VARCHAR(64) DEFAULT NULL,
    user_avatar VARCHAR(255) DEFAULT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    last_message VARCHAR(255) DEFAULT NULL,
    unread_count INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_messages (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    session_id VARCHAR(40) NOT NULL,
    sender VARCHAR(20) NOT NULL,
    msg_type VARCHAR(20) NOT NULL DEFAULT 'text',
    content TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_chat_messages_session_id (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    published BOOLEAN NOT NULL DEFAULT FALSE,
    sort INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_caregivers_name (name),
    KEY idx_caregivers_origin (origin),
    KEY idx_caregivers_availability (availability_status),
    KEY idx_caregivers_recommended (recommended),
    KEY idx_caregivers_published (published),
    KEY idx_caregivers_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS demands (
    id VARCHAR(32) NOT NULL,
    user_id VARCHAR(40) DEFAULT NULL,
    service_id VARCHAR(32) NOT NULL,
    service_name VARCHAR(64) NOT NULL,
    caregiver_id VARCHAR(40) DEFAULT NULL,
    caregiver_name VARCHAR(64) NOT NULL DEFAULT '',
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

-- Idempotent legacy-table upgrades.
-- CREATE TABLE IF NOT EXISTS does not add columns to existing tables, so keep
-- these guards in sync with new fields added to the CREATE TABLE definitions.
SET @schema_name = DATABASE();

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE users ADD COLUMN openid VARCHAR(128) DEFAULT NULL AFTER id',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'users'
      AND COLUMN_NAME = 'openid'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE users ADD COLUMN phone VARCHAR(32) NOT NULL DEFAULT '''' AFTER nickname',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'users'
      AND COLUMN_NAME = 'phone'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE users ADD COLUMN signature VARCHAR(255) NOT NULL DEFAULT '''' AFTER phone',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'users'
      AND COLUMN_NAME = 'signature'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE users ADD COLUMN last_login_at VARCHAR(32) DEFAULT NULL AFTER last_order_at',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'users'
      AND COLUMN_NAME = 'last_login_at'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE users ADD UNIQUE KEY idx_users_openid (openid)',
        'SET @schema_noop = 0')
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'users'
      AND INDEX_NAME = 'idx_users_openid'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE users ADD UNIQUE KEY idx_users_phone (phone)',
        'SET @schema_noop = 0')
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'users'
      AND INDEX_NAME = 'idx_users_phone'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE orders ADD COLUMN user_id VARCHAR(40) DEFAULT NULL AFTER id',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'orders'
      AND COLUMN_NAME = 'user_id'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE orders ADD COLUMN mini_deleted BOOLEAN NOT NULL DEFAULT FALSE AFTER close_reason',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'orders'
      AND COLUMN_NAME = 'mini_deleted'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE orders ADD KEY idx_orders_user_id (user_id)',
        'SET @schema_noop = 0')
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'orders'
      AND INDEX_NAME = 'idx_orders_user_id'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;
