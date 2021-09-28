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

系统中全部密码均使用hash（```sha256```）、不使用明文  
前端：用户输入明文密码，确认后（例如点击登录按钮），程序执行hash，清空输入框  
后端：从前端接收到hash后的密码，与该用户的salt再次执行hash，保存到数据库（创建）或与数据库记录对比（验证）  
数据库：密码字段实际保存的值为```sha256(sha256('明文')+'salt')```

这样一来，只有前端发起的请求被拦截并从中解析出hash后的密码，才会导致密码泄露  
即使数据库字段被获取，因为hash不可逆，理论上无法反推出前端传过来的hash后的密码，所以无法正确执行我们后端的登录函数

### 前台(1)

#### 登录

/api/login

输入：

1. 用户名 userName
2. 密码 password

返回：

1. 用户ID userID
2. 昵称 nickname
3. 权限等级 permission

### 后台(7)

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
    5. 创建人 createdBy

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

## 云文件系统

用户在后台上传文件到云服务器，可以在前台查看  
文件一经上传，不可更改  
允许用户删除自己上传的文件，采用软删除，保留数据库记录与云文件，只是查询时不显示

编写后台程序移动端界面，拟仅开放云文件系统，且仅支持查询自己上传的文件

### 前台(2)

#### 查询

用户可查看自己上传的全部文件，以及其他权限等级**不高于**自身的用户上传的**公开文件**

查询：

1. 当前用户上传的全部文件 /api/cloudFile/listByUploader
2. 当前用户可查看的公开文件 /api/cloudFile/listPublic
3. 查询已删除的文件（仅后台） /api/cloudFile/listDelete
   1. 规则1：需要**S级管理员权限**
   2. 规则2：仅可查询权限等级**不高于**自身的上传者上传的文件

输入：

1. 查询者ID operatorID
2. 分页信息
    1. 每页条数 pageSize
    2. 当前页码 pageNum

返回：

1. 符合条件的数据条数 total
2. 文件列表 files
    1. 文件ID fileID
    2. 文件名 fileName
    3. 访问路径 fileURL ```/public/ffff.pdf```
    4. 是否公开 isPublic
    5. 更新时间 updateTime
    6. 创建时间 createdTime

#### 预览

暂时使用浏览器默认解析方式，前端限制只接收pdf类型的文件，文件大小是否要限制、限制多少**待定**

与nginx配合，前端根据url直接定位到服务器文件：
1. 后端返回的```fileURL```举例：```/public/ffff.pdf```，
2. 前端拼上**源**和**云文件标识**：```https://mats9693.cn/cloud-file/public/ffff.pdf```
3. nginx识别**云文件标识**：
```text 
location /cloud-file/ {
    alias /home/xxx/cloud_file/;
}
```

### 后台(5)

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

#### 删除

/api/cloudFile/delete

用户删除自己上传的文件

输入：

1. 操作员 operatorID
2. 密码 password
3. 文件ID fileID

规则：

1. 仅允许删除自己上传的文件
2. 软删除，保留数据库记录与服务器文件，修改数据库记录**更新时间**与**是否已删除**字段
3. 需要指定文件当前状态为**未删除**

输出：

1. 删除结果 isSuccess

#### 查询

> 允许根据上传者查询 /api/cloudFile/listByUploader  
> 允许S级管理员查询已删除的文件 /api/cloudFile/listDelete

#### 预览

> 与前台对应模块相同
