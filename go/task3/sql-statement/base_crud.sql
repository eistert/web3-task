/*
题目1：基本CRUD操作
假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、
grade （学生年级，字符串类型）。
要求 ：
编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/

-- students 建表语句

CREATE TABLE IF NOT EXISTS `students` (
  `id`     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name`   VARCHAR(64)     NOT NULL COMMENT '学生姓名',
  `age`    TINYINT UNSIGNED NOT NULL COMMENT '年龄',
  `grade`  VARCHAR(32)     NOT NULL COMMENT '年级',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci
  COMMENT='学生信息表';



-- 1) 插入一条记录
INSERT INTO students (name, age, grade)
VALUES ('张三', 20, '三年级');

-- 2) 查询年龄 > 18 的学生
SELECT id, name, age, grade
FROM students
WHERE age > 18;

-- 3) 将姓名为“张三”的年级更新为“四年级”
UPDATE students
SET grade = '四年级'
WHERE name = '张三';

-- 4) 删除年龄 < 15 的学生
DELETE FROM students
WHERE age < 15;

