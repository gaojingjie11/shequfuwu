-- Smart Community migration: AI report history table
-- MySQL 5.7 compatible

CREATE TABLE IF NOT EXISTS cms_ai_report (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  repair_new_count BIGINT NOT NULL DEFAULT 0,
  repair_pending_count BIGINT NOT NULL DEFAULT 0,
  visitor_new_count BIGINT NOT NULL DEFAULT 0,
  property_paid_count BIGINT NOT NULL DEFAULT 0,
  property_paid_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  report_summary VARCHAR(255) NOT NULL DEFAULT '',
  report_markdown LONGTEXT,
  generated_by BIGINT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_ai_report_created_at(created_at)
);

