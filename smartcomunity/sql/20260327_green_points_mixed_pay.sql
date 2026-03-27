-- Smart Community migration: green points + mixed payment
-- MySQL 5.7 compatible (uses information_schema + dynamic SQL)

SET @db = DATABASE();

-- sys_user
SET @sql = (
  SELECT IF(COUNT(*)=0,
    'ALTER TABLE sys_user ADD COLUMN green_points BIGINT NOT NULL DEFAULT 0',
    'SELECT 1'
  )
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='sys_user' AND COLUMN_NAME='green_points'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = (
  SELECT IF(COUNT(*)=0,
    'ALTER TABLE sys_user ADD COLUMN balance DECIMAL(10,2) NOT NULL DEFAULT 0.00',
    'ALTER TABLE sys_user MODIFY COLUMN balance DECIMAL(10,2) NOT NULL DEFAULT 0.00'
  )
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='sys_user' AND COLUMN_NAME='balance'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- oms_order
SET @sql = (
  SELECT IF(COUNT(*)=0,
    'ALTER TABLE oms_order ADD COLUMN used_points INT NOT NULL DEFAULT 0',
    'SELECT 1'
  )
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='oms_order' AND COLUMN_NAME='used_points'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = (
  SELECT IF(COUNT(*)=0,
    'ALTER TABLE oms_order ADD COLUMN used_balance DECIMAL(10,2) NOT NULL DEFAULT 0.00',
    'SELECT 1'
  )
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='oms_order' AND COLUMN_NAME='used_balance'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = (
  SELECT IF(COUNT(*)=0,
    'ALTER TABLE oms_order ADD COLUMN paid_at DATETIME NULL',
    'SELECT 1'
  )
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='oms_order' AND COLUMN_NAME='paid_at'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- cms_property_fee
SET @sql = (
  SELECT IF(COUNT(*)=0,
    'ALTER TABLE cms_property_fee ADD COLUMN used_points INT NOT NULL DEFAULT 0',
    'SELECT 1'
  )
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='cms_property_fee' AND COLUMN_NAME='used_points'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = (
  SELECT IF(COUNT(*)=0,
    'ALTER TABLE cms_property_fee ADD COLUMN used_balance DECIMAL(10,2) NOT NULL DEFAULT 0.00',
    'SELECT 1'
  )
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='cms_property_fee' AND COLUMN_NAME='used_balance'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- green_point_record
SET @sql = (
  SELECT IF(COUNT(*)=0,
    'CREATE TABLE green_point_record (
      id BIGINT PRIMARY KEY AUTO_INCREMENT,
      user_id BIGINT NOT NULL,
      action VARCHAR(64) NOT NULL,
      points INT NOT NULL,
      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      INDEX idx_green_point_record_user_id(user_id)
    )',
    'SELECT 1'
  )
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA=@db AND TABLE_NAME='green_point_record'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

