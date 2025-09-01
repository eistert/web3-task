太棒了！既然你已经用 `.env` 让服务连上了 MySQL（或 SQLite），下面给你一套**一次跑通**的自测流程（用 `curl`），把注册登录/JWT、文章 CRUD、评论都验证一遍。你也可以把这些步骤搬到 Postman 里逐条执行；我把**请求、头、示例请求体、预期返回**都写清楚了。

> 约定：服务监听 `http://localhost:8080`，前缀是 `/api/v1`
> 如果要在终端里自动取字段，建议装 `jq`（可选）：`brew install jq`

---

# 0. 健康检查

```bash
BASE=http://localhost:8080/api/v1
curl -i $BASE/health
```

预期：`200 OK`，响应体 `OK`

---

# 1. 注册两个用户（alice、bob）

```bash
curl -s -X POST $BASE/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123","email":"alice@example.com"}'

curl -s -X POST $BASE/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"bob","password":"pass123","email":"bob@example.com"}'
```

预期：`201 Created` 或 `200 OK`（如果已存在会提示冲突或校验错误）

---

# 2. 登录获取 JWT

```bash
ALICE_TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123"}' | jq -r '.data.token')

BOB_TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"bob","password":"pass123"}' | jq -r '.data.token')

echo "ALICE_TOKEN=${ALICE_TOKEN}"
echo "BOB_TOKEN=${BOB_TOKEN}"
```

预期：拿到两个非空的 JWT。
（如果不用 `jq`，就把返回 JSON 里的 `token` 手动复制出来。）

---

# 3.（鉴权）Alice 创建一篇文章

```bash
POST_ID=$(curl -s -X POST $BASE/posts \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello GORM","content":"This is my first post."}' | jq -r '.data.ID')

echo "POST_ID=${POST_ID}"
```

预期：`201 Created`，返回里有 `id`。未带 Token 应该是 `401`。

---

# 4.（公开）获取文章列表与详情

```bash
curl -s $BASE/posts | jq .
curl -s $BASE/posts/$POST_ID | jq .
```

预期：能看到刚创建的文章，`user_id` 为 Alice。

---

# 5.（鉴权+授权）文章更新 & 权限校验

* Alice 更新（应成功）：

```bash
curl -i -X PUT $BASE/posts/$POST_ID \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello GORM (edited)","content":"Updated content"}'
```

预期：`200 OK`，内容已更新。

* Bob 尝试更新（应禁止）：

```bash
curl -i -X PUT $BASE/posts/$POST_ID \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hacked","content":"Nope"}'
```

预期：`403 Forbidden`。

---

# 6.（鉴权）Bob 给这篇文章发表评论

```bash
COMMENT_ID=$(curl -s -X POST $BASE/posts/$POST_ID/comments \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"content":"Nice post!"}' | jq -r '.id')

echo "COMMENT_ID=${COMMENT_ID}"
```

预期：`201 Created`，返回评论 `id`。

---

# 7.（公开）查看某篇文章的所有评论

```bash
curl -s $BASE/posts/$POST_ID/comments | jq .
```

预期：列表里能看到 Bob 的评论；如果实现了预加载，可能还能看到评论的 `user` 信息。

---

# 8.（授权）删除文章

* Bob 删除（应禁止）：

```bash
curl -i -X DELETE $BASE/posts/$POST_ID \
  -H "Authorization: Bearer $BOB_TOKEN"
```

预期：`403 Forbidden`。

* Alice 删除（应成功）：

```bash
curl -i -X DELETE $BASE/posts/$POST_ID \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

预期：`200 OK` 或 `204 No Content`。

> 如果你在模型里配置了外键级联，评论也会被一起删；否则需要显式删除或做软删 —— 这取决于你的模型定义。

---

# 9. 常见错误用例（可选但建议测）

* 未带 `Authorization` 创建文章：应 `401`
* 带损坏的 JWT（比如在 token 末尾加几个字符）：应 `401`
* 更新/删除不是自己创建的文章：应 `403`
* 查询不存在的文章 / 评论：应 `404`（如果你这样处理）

---

# 10. 数据库侧确认（可选）

* **MySQL**：`USE blog_backend; SHOW TABLES; SELECT * FROM users; SELECT * FROM posts; SELECT * FROM comments;`
* **SQLite**（默认 `blog.db`）：

  ```bash
  sqlite3 blog.db '.tables'
  sqlite3 blog.db 'SELECT * FROM users;'
  ```

---

# 11. Postman 测试清单（给你或同学使用）

1. `GET /api/v1/health` → 200
2. `POST /auth/register` x2（alice、bob）
3. `POST /auth/login`（alice、bob）→ 保存两个 token 到 Postman 的环境变量
4. `POST /posts`（Alice 的 token）→ 得到 `post_id`
5. `GET /posts`、`GET /posts/:id` → 能看到
6. `PUT /posts/:id`（Alice 的 token）→ 200
7. `PUT /posts/:id`（Bob 的 token）→ 403
8. `POST /posts/:id/comments`（Bob 的 token）→ 201
9. `GET /posts/:id/comments` → 列表正确
10. `DELETE /posts/:id`（Alice 的 token）→ 200/204

---

## 如果遇到问题

* 日志里如果出现 `[DB] using SQLite: blog.db`，说明没读到 `MYSQL_DSN`，确认 `.env` 在 `cmd/api`、并已被加载。
* 403/401 问题，先确认 `Authorization: Bearer <token>` 是否带对。
* 字段校验错误，调请求体的 JSON 键名与模型标签一致（如 `title`、`content`）。

按这套流程走一遍，基本就能证明你的**个人博客后端**已完成作业要求 ✅。如果某一步出错，把请求与响应贴出来，我帮你定位。
