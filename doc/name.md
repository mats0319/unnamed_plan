# 命名

项目中优先使用统一的命名，未使用之处需要特别标注

## http

1. 用户ID：userID
2. 用户名：userName
3. 昵称：nickname
4. 锁定状态：isLocked
5. 密码：password
6. 权限等级：permission
7. 分页信息
    1. 每页条数：pageSize
    2. 当前页码：pageNum （当前页码从1开始）
    3. 符合条件的数据总数：total
8. 请求的发起者ID：operatorID （一般用于需要区分调用者和被操作者的场景，例如锁定用户）
9. 执行结果：isSuccess

## 数据库

### dao

1. 插入：insert
2. 更新：update
3. 查询：query
4. 删除：delete

### 前、后端

1. 新增：create
2. 修改：modify
3. 查询：list
4. 删除：delete