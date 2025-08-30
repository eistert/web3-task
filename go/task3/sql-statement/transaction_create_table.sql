-- 账户表：余额单位为“分”（int64 对应 BIGINT）
CREATE TABLE IF NOT EXISTS accounts (
  id        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  balance   BIGINT          NOT NULL DEFAULT 0 COMMENT '余额（分）',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB;


-- 转账记录表：amount 也是“分”
CREATE TABLE IF NOT EXISTS transactions (
  id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  from_account_id  BIGINT UNSIGNED NOT NULL,
  to_account_id    BIGINT UNSIGNED NOT NULL,
  amount           BIGINT          NOT NULL COMMENT '转账金额（分）',
  created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_from (from_account_id),
  KEY idx_to   (to_account_id),
  CONSTRAINT fk_tx_from FOREIGN KEY (from_account_id) REFERENCES accounts(id),
  CONSTRAINT fk_tx_to   FOREIGN KEY (to_account_id)   REFERENCES accounts(id),
  CONSTRAINT chk_amount CHECK (amount > 0),
  CONSTRAINT chk_not_self CHECK (from_account_id <> to_account_id)
) ENGINE=InnoDB;


// 测试数据（A 有 100 元、B 有 50 元）：

INSERT INTO accounts (balance) VALUES (10000), (5000);
-- id=1: A(¥100.00), id=2: B(¥50.00)


