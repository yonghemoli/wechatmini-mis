-- 永和护理 MIS 历史结构清理脚本
-- 警告：本脚本会永久删除旧订单、店铺、地址、家属、配菜和旧服务项目数据。
-- 建议执行前备份：mysqldump yonghemolimis > yonghemolimis-before-cleanup.sql
-- 推荐顺序：1) init-schema.sql  2) cleanup-legacy-tables.sql  3) init-seed.sql

SET @schema_name = DATABASE();

-- 将旧 caregivers.published 平滑迁移为 DRAFT / COMPLETED，并补充数据来源。
SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE caregivers ADD COLUMN application_id VARCHAR(32) DEFAULT NULL AFTER id',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND COLUMN_NAME = 'application_id'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE caregivers ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT ''DRAFT'' AFTER photo_urls',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND COLUMN_NAME = 'status'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE caregivers ADD COLUMN source VARCHAR(20) NOT NULL DEFAULT ''ADMIN'' AFTER status',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND COLUMN_NAME = 'source'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) > 0,
        'UPDATE caregivers SET status = IF(published = TRUE, ''COMPLETED'', ''DRAFT'')',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND COLUMN_NAME = 'published'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) > 0,
        'ALTER TABLE caregivers DROP COLUMN published',
        'SET @schema_noop = 0')
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND COLUMN_NAME = 'published'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE caregivers ADD KEY idx_caregivers_application_id (application_id)',
        'SET @schema_noop = 0')
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND INDEX_NAME = 'idx_caregivers_application_id'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE caregivers ADD KEY idx_caregivers_status_sort (status, sort)',
        'SET @schema_noop = 0')
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND INDEX_NAME = 'idx_caregivers_status_sort'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

SET @migration_sql = (
    SELECT IF(COUNT(*) = 0,
        'ALTER TABLE caregivers ADD KEY idx_caregivers_source (source)',
        'SET @schema_noop = 0')
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'caregivers' AND INDEX_NAME = 'idx_caregivers_source'
);
PREPARE migration_stmt FROM @migration_sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

-- 删除已经退出业务边界的历史表。
DROP TABLE IF EXISTS order_status_histories;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS service_targets;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS meal_packages;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS shops;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS service_types;

-- 删除旧页面遗留配置，不影响新版装修与协议配置。
DELETE FROM app_configs
WHERE `key` IN ('mini.home', 'mini.appointment.tabs', 'mini.business.app', 'mini.business.about');

-- 验收：结果中不应再出现上述历史表。
SELECT TABLE_NAME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
ORDER BY TABLE_NAME;
