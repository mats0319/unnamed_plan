# 命名

项目中优先使用统一的命名，未使用之处需要特别标注

## http

1. 请求的发起者ID：operatorID
2. 分页信息
    1. 每页条数：pageSize
    2. 当前页码：pageNum （当前页码从1开始）
    3. 符合条件的数据总条数：total
3. 执行结果：isSuccess
4. 创建时间：createdTime
5. 更新时间：updateTime
6. 是否公开：isPublic

7. 用户ID：userID
8. 用户名：userName
9. 昵称：nickname
10. 锁定状态：isLocked
11. 密码：password
12. 权限等级：permission
13. 创建人：createdBy

14. 文件名：fileName
15. 扩展名：extensionName
16. 文件大小：fileSize
17. 文件：file
18. 文件路径：files
19. 文件ID：fileID
20. 访问路径：fileURL
21. 上次修改时间：lastModifiedTime

22. 编辑者：writeBy
23. 主题：topic
24. 内容：content
25. 笔记：notes
26. 笔记ID：noteID
27. 编辑者ID：writerID

## CRUD

### 数据库

1. 插入：insert
2. 更新：update
3. 查询：query
4. 删除：delete

### 前、后端程序

1. 新增：create
2. 修改：modify
3. 查询：list
4. 删除：delete
