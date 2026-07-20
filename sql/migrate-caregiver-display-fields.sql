-- 阿姨档案“小程序展示控制”字段迁移。
-- 在已有 yonghemolimis 库执行一次；新库执行 init-schema.sql 即可。
-- 先允许 NULL，兼容 MySQL 对 TEXT 默认值的限制；应用读取 NULL 时会按“全部展示”处理。
ALTER TABLE caregivers ADD COLUMN display_fields TEXT NULL AFTER photo_urls;
UPDATE caregivers SET display_fields = '{}' WHERE display_fields IS NULL;
