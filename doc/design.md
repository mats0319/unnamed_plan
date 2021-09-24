# 设计

## http res

```go 
type ResponseData struct {
	HasError bool   `json:"hasError"`
	Data     string `json:"data"`
}
```

```hasError```表示http请求是否正确执行并按照预期返回：

1. 请求正常执行：```data```是```json string```
2. 请求执行失败：```data```是```string```类型的错误信息

## 用户系统

用户系统涉及的密码操作：

1. 用户在前端界面输入密码明文```text```
    1. 用户确定后，立刻清空密码框
2. 前端对明文执行```sha256 hash```，将结果转为16进制字符串```input pwd```
3. 后端对前端传过来的```input pwd```加盐```salt```，执行```sha256 hash```，得到```password```
4. ```password```即视作用户输入的密码，保存到数据库或与数据库中的记录进行比较

### 前台

#### 登录

/api/login

输入：

1. 用户名 userName
2. 密码 password

返回：

1. 用户ID userID
2. 昵称 nickname
3. 权限等级 permission

### 后台

#### 登录

> 与前台登录功能相同

#### 查询

/api/user/list

获取用户列表

输入：

1. 查询者ID operatorID
2. 分页信息
    1. 每页条数 pageSize
    2. 当前页码 pageNum

规则：

1. 只能查询权限等级**不超过**自身的用户

返回：

1. 符合条件的数据条数 total
2. 用户列表 users
    1. 用户ID userID
    2. 昵称 nickname
    3. 锁定状态 isLocked
    4. 用户权限等级 permission

#### 修改用户信息（昵称和密码）

/api/user/modifyInfo

修改自己的昵称和密码

输入：

1. 修改者ID operatorID （当前登录用户）
2. 目标用户ID userID （想要修改的用户）
3. 当前密码 currPwd
4. 新的昵称 nickname
5. 新的密码 password

规则：

1. 只允许修改自己的昵称和密码
2. 若成功修改密码，则前端退出登录

输出：

1. 修改结果 isSuccess

#### 创建

/api/user/create

输入：

1. 创建者ID operatorID
2. 新用户的用户名 userName
3. 新用户的密码 password
4. 新用户的权限等级 permission

规则：

1. 需要创建者权限等级达到**A级管理员权限等级**（见配置文件）
2. 只能创建比自己权限等级**低**的用户
3. 新用户的昵称与用户名相同

返回：

1. 创建结果 isSuccess

#### 锁定

/api/user/lock

禁用指定用户

输入：

1. 锁定者ID operatorID
2. 待锁定者ID userID

规则：

1. 需要待锁定者当前是**未锁定**状态
2. 需要锁定者权限等级达到**A级管理员权限等级**（见配置文件）
3. 只能锁定比自己权限等级**低**的用户（隐含：不允许用户锁定自己）

返回：

1. 锁定结果 isSuccess

#### 解锁

/api/user/unlock

恢复使用指定账户

输入：

1. 解锁者ID operatorID
2. 待解锁者ID userID

规则：

1. 需要待解锁者当前是**已锁定**状态
2. 需要解锁者的权限等级达到**A级管理员权限等级**（见配置文件）
3. 只能解锁比自己权限等级**低**的用户（隐含：不允许用户解锁自己）

返回：

1. 解锁结果 isSuccess

#### 修改权限

/api/user/modifyPermission

修改指定用户的权限等级

输入：

1. 修改者ID operatorID
2. 待修改者ID userID
3. 待修改者的新权限等级 permission

规则：

1. 需要修改者权限等级达到**S级管理员权限等级**（见配置文件）
2. 待修改者**当前权限等级**与**新权限等级**均应**小于** **S级管理员权限等级**

返回：

1. 修改结果 isSuccess

## 云文件

### 前台

#### 查询

用户可查看自己上传的全部文件，以及其他权限等级**不高于**自身的用户上传的**公开文件**

查询：

1. 当前用户上传的全部文件 /api/cloudFile/listByUploaderID
2. 当前用户可查看的公开文件 /api/cloudFile/listPublic

输入：

1. 查询者ID operatorID
2. 分页信息
    1. 每页条数 pageSize
    2. 当前页码 pageNum

返回：

1. 符合条件的数据条数 total
2. 文件列表 files
    1. 文件名 fileName
    2. 访问路径 fileURL
    3. 是否公开 isPublic

#### 预览和下载

当前仅支持pdf文件，通过nginx转发请求，让浏览器直接与pdf文件对话，至于浏览器是预览还是下载，都由他去

### 后台

#### 上传

/api/cloudFile/upload

用户上传文件到云服务器

输入：

1. 上传者ID operatorID
2. 文件名 fileName
3. 扩展名 extensionName
5. 是否公开 isPublic
6. 文件 file

规则：

1. 云文件存储结构：
    1. 云文件夹根目录 - 公开文件夹、非公开文件夹（每个用户一个文件夹）

返回：

1. 上传结果 isSuccess

#### 查询

> 仅提供当前用户上传的文件查询，规则与前台相同 /api/cloudFile/listByUploaderID

#### 预览和下载

> 与前台对应模块相同
