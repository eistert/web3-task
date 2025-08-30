CREATE TABLE IF NOT EXISTS books (
  id     BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  title  VARCHAR(128)  NOT NULL,
  author VARCHAR(64)   NOT NULL,
  price  DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  created_at DATETIME NULL,
  updated_at DATETIME NULL
);

-- 可选：一些测试数据
INSERT INTO books (title, author, price, created_at, updated_at) VALUES
('Go 语言实战', 'John', 88.00, NOW(), NOW()),
('数据库系统概念', 'Abraham', 129.00, NOW(), NOW()),
('算法导论', 'CLRS', 99.50, NOW(), NOW()),
('你不知道的JS', 'Kyle', 45.00, NOW(), NOW())
ON DUPLICATE KEY UPDATE price=VALUES(price);
