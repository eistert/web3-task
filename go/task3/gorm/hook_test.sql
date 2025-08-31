-- 可选：USE blog;
-- USE gorm_demo;

SET NAMES utf8mb4;

-- 清空并重置自增，保证我们能用固定的 id
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE comments;
TRUNCATE TABLE posts;
TRUNCATE TABLE users;
SET FOREIGN_KEY_CHECKS = 1;

-- 1) 用户：Alice(id=1), Bob(id=2)，post_count 先设为 0
INSERT INTO users (id, name, email, post_count, created_at, updated_at) VALUES
(1, 'Alice', 'alice@example.com', 0, NOW(), NOW()),
(2, 'Bob',   'bob@example.com',   0, NOW(), NOW());

/* 2) 用于测试 Comment.AfterDelete 的基线：
      - 创建一篇仅有“1条评论”的文章（post_id=201）
      - 将该文章的 comment_status 设为“有评论”
      - 评论 id=301
      删除这唯一一条评论后，Hook 应把该文章的 comment_status 改为“无评论”
*/
INSERT INTO posts (id, user_id, title, content, comment_status, created_at, updated_at)
VALUES (201, 1, '删除评论Hook测试', '这篇文章只有一条评论', '有评论', NOW(), NOW());

INSERT INTO comments (id, post_id, user_id, content, created_at, updated_at)
VALUES (301, 201, 1, '唯一一条评论', NOW(), NOW());

-- 3) （可选）再建一篇没有评论的文章，用于 AfterCreate 之“对比查询”
INSERT INTO posts (id, user_id, title, content, comment_status, created_at, updated_at)
VALUES (101, 1, '基线文章（无评论）', '用于查询对比', '无评论', NOW(), NOW());
