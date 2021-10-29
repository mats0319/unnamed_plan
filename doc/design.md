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

多处登录限制：同一用户，**最多**允许同时登录一个```public ui```、一个```admin ui```（```public mobile```不限制）  
实现：

1. 程序内存维护一个**用户**到**登录信息**的```map```，根据上述规则，**登录信息**包括以下内容：
   ```go 
      type loginInfo struct {
         publicUIToken           string
         publicUILastOperateTime int64
         adminUIToken            string
         adminUILastOperateTime  int64
      }
   ```
2. 不同的界面在http请求头添加不同的**标识**与**随机字符串**，若未检测到**有效的标识和随机字符串**，则拒绝请求
    1. 标识：http请求头的```unnamed-plan```字段，其有效值如下，分别对应不同的请求来源界面
        1. ```public ui```：表示该请求来自```public ui```界面
        2. ```admin ui```：表示该请求来自```admin ui```界面
    2. 随机字符串：http请求头的```unnamed-plan-auth-code```字段
    3. 登录请求可以没有**随机字符串**
3. 新的请求需要通过**随机字符串验证**与**上次操作时间验证**才会执行
    1. 随机字符串验证：要求与程序保存的随机字符串一致
    2. 上次操作时间验证：要求本次请求与上次有效操作的时间间隔不多于**一定值**（拟两小时，配置）
    3. 登录请求默认通过两项验证
    4. 成功的登录请求会修改**随机字符串**和**上次操作时间**，返回**新的随机字符串**
    5. 成功的非登录请求会修改**上次操作时间**
4. 成功的登录请求会在请求体里返回新的**随机字符串**，前端接收并在接下来的请求中带上它
    1. 后端识别**请求头**中的关键参数，然后把关键参数写在**请求体**里返回，  
       这种做法能够最大程度兼容```go http server```和```vue axios interceptor```的写法  
       简单来说就是容易，后续考虑统一

问题：

1. 程序内存有限，如果用户总量多了，需要把登录信息保存到其他地方，例如缓存层

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
    2. 用户名 userName
    3. 昵称 nickname
    4. 锁定状态 isLocked
    5. 用户权限等级 permission
    6. 创建人 createdBy

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
2. 若**昵称**与**密码**均无改动，不应允许执行修改
3. 若成功修改密码，则前端退出登录

输出：

1. 修改结果 isSuccess

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

编写**前台程序移动端界面**，开放查询模块，支持查询自己上传的文件

### 前台(2)

#### 查询

用户可查看自己上传的全部文件，以及其他权限等级**不高于**自身的用户上传的**公开文件**

查询：

1. 当前用户上传的全部文件 /api/cloudFile/listByUploader
2. 当前用户可查看的公开文件 /api/cloudFile/listPublic

规则：

1. 默认不查询已删除的文件

输入：

1. 查询者ID operatorID
2. 分页信息
    1. 每页条数 pageSize
    2. 当前页码 pageNum

返回：

1. 符合条件的数据条数 total
2. 文件列表 files
    1. 文件ID fileID （文件存储名称）
    2. 文件名 fileName （文件展示名称）
    3. 文件最后修改时间 lastModifiedTime
    4. 访问路径 fileURL ```/public/ffff.pdf```
    5. 是否公开 isPublic
    6. 上次修改时间 updateTime
    7. 初次上传时间 createdTime

#### 预览

暂时使用浏览器默认解析方式，前端限制只接收pdf类型的文件，文件大小是否要限制、限制多少**待定**

与nginx配合，前端根据url直接定位到服务器文件：

1. 后端返回的```fileURL```举例：```public/ffff.pdf```，
2. 前端拼上**源**和**云文件标识**：```https://mats9693.cn/cloud-file/public/ffff.pdf```
3. nginx识别**云文件标识**：
   ```text 
   location /cloud-file/ {
       alias /home/xxx/cloud_file/;
   }
   ```

##### 问题

这样一来，若已知用户ID、文件ID和扩展名，是可以绕过系统的权限验证，从而看到其他人的非公开文件的。  
除了通过后端获取文件，有什么办法可以解决吗？

### 后台(5)

#### 上传

/api/cloudFile/upload

用户上传文件到云服务器

输入：

1. 上传者ID operatorID
2. 文件名 fileName
3. 扩展名 extensionName
4. 最后修改时间 lastModifiedTime
5. 是否公开 isPublic
6. 文件 file

规则：

1. 云文件存储结构
    1. 云文件夹根目录 - 公开文件夹、非公开文件夹（每个用户一个文件夹）

返回：

1. 上传结果 isSuccess

#### 修改

/api/cloudFile/modify

输入：

1. 修改者ID operatorID
2. 目标文件ID fileID
3. 修改者密码 password
4. 新的文件名 fileName
5. 新的扩展名 extensionName
6. 是否公开 isPublic
7. 新的文件 file
8. 新的文件最后修改时间 lastModifiedTime

规则：

1. 仅允许修改自己上传的文件
2. 若**文件**、**文件名**、**扩展名**、**是否公开**均无改动，不应允许执行修改
3. 不删除旧文件，将旧文件移动至**修改者非公开文件目录**备份

返回：

1. 修改结果 isSuccess

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

> 见前台对应模块

#### 预览

> 见前台对应模块

## 随想系统

记录一些自己的想法（空想）与体会（经一事，长一智）

### 前台(2)

#### 查询

查询：

1. 当前用户编辑的笔记 /api/thinkingNote/listByWriter
2. 当前用户可查看的公开笔记 /api/thinkingNote/listPublic

规则：

1. 默认不查询已删除的笔记

输入：

1. 查询者ID operatorID
2. 分页信息
    1. 每页条数：pageSize
    2. 当前页码：pageNum

返回：

1. 符合条件的数据条数 total
2. 笔记 notes
    1. 笔记ID noteID
    2. 编辑者 writeBy
    3. 主题 topic
    4. 内容 content
    5. 是否公开 isPublic
    6. 上次更新 updateTime
    7. 创建时间 createdTime

### 后台(5)

#### 记录（创建）

/api/thinkingNote/create

输入：

1. 记录者ID operatorID
2. 主题 topic
3. 内容 content
4. 是否公开 isPublic

返回：

1. 记录结果 isSuccess

#### 修改

/api/thinkingNote/modify

输入：

1. 修改者ID operatorID
2. 目标笔记ID noteID
4. 修改者密码 password
5. 新的主题 topic
6. 新的内容 content
7. 是否公开 isPublic

规则：

1. 仅允许修改自己编辑的笔记
2. 若**主题**、**内容**、**是否公开**均无改动，不应允许执行修改

返回：

1. 修改结果 isSuccess

#### 删除

/api/thinkingNote/delete

输入：

1. 删除者ID operatorID
2. 密码 password
3. 笔记ID noteID

规则：

1. 仅允许删除自己上传的文件
2. 软删除，保留数据库记录，修改数据库记录**更新时间**与**是否已删除**字段
3. 需要指定笔记当前状态为**未删除**

返回：

1. 删除结果 isSuccess

#### 查询

> 见前台对应模块

## 小游戏系统

单人游戏：前端体现规则，后端仅记录游戏结果  
多人游戏：拟于两个前端之间建立websocket连接，后端仅记录对局结果

本节首先列举前台、后台的功能，而后针对每款游戏，还会用单独的章节介绍

### 前台(?)

### 后台(?)，仅A级管理员及以上权限可见

#### 新增游戏

/api/game/create

前端用game id与页面绑定

#### 查询游戏

/api/game/list

通过参数控制是否查询已关闭的游戏

#### 开启游戏

/api/game/open

#### 关闭游戏

/api/game/close

### 翻牌游戏（单人）

规则：

1. **16张牌**，相同的牌背，牌面有**8**种、每种**2**张
2. 每次可以翻**2**张牌，称为一组：
    1. 若全部牌的牌面相同，则消除这一组牌
    2. 若存在不同的牌面，则翻转该组牌至牌背朝上的状态
3. 消除全部牌时，游戏结束，记录[游戏结果（翻牌游戏）]

[游戏结果（翻牌游戏）]:本局游戏时长、使用步数

扩展：

1. 总牌数扩展，如修改为**36**张牌
2. 单次翻牌数与牌型扩展，如修改为每次翻**3**张牌，相应的，总牌型当中，每种牌要有**3**张

实现：
