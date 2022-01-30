# API

一个具体功能的API：

1. uri
2. 描述
3. 输入：输入参数列表
4. 规则：功能逻辑、规则
5. 返回：输出参数列表

## 用户系统

系统中全部密码均使用hash（```sha256```）、不使用明文  
前端：用户输入明文密码，确认后（例如点击登录按钮），程序执行hash，清空输入框  
后端：从前端接收到hash后的密码，与该用户的salt再次执行hash，保存到数据库（创建）或与数据库记录对比（验证）  
数据库：密码字段实际保存的值为```sha256(sha256('明文')+'salt')```

这样一来，只有前端发起的请求被拦截并从中解析出hash后的密码，才会导致密码泄露  
即使数据库字段被获取，因为hash不可逆，理论上无法反推出前端传过来的hash后的密码，所以无法正确执行我们的登录函数

### 前台(1)

#### 登录

/api/login

输入：

1. 用户名 userName
2. 密码 password

规则：

1. 用户名、密码均**不允许为空**

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
    2. 用户名 userName
    3. 昵称 nickname
    4. 锁定状态 isLocked
    5. 用户权限等级 permission

#### 创建

/api/user/create

输入：

1. 创建者ID operatorID
2. 新用户的用户名 userName
3. 新用户的密码 password
4. 新用户的权限等级 permission

规则：

1. 需要创建者权限等级达到**A级管理员权限等级**（可配置）
2. 只能创建比自己权限等级**低**的用户
3. 新用户的昵称与用户名相同

#### 锁定

/api/user/lock

禁用指定用户

输入：

1. 锁定者ID operatorID
2. 待锁定者ID userID

规则：

1. 需要待锁定者当前是**未锁定**状态
2. 需要锁定者权限等级达到**A级管理员权限等级**（可配置）
3. 只能锁定比自己权限等级**低**的用户（隐含：不允许用户锁定自己）

#### 解锁

/api/user/unlock

恢复使用指定账户

输入：

1. 解锁者ID operatorID
2. 待解锁者ID userID

规则：

1. 需要待解锁者当前是**已锁定**状态
2. 需要解锁者的权限等级达到**A级管理员权限等级**（可配置）
3. 只能解锁比自己权限等级**低**的用户（隐含：不允许用户解锁自己）

#### 修改用户信息

/api/user/modifyInfo

修改自己的昵称和密码

输入：

1. 修改者ID operatorID
2. 待修改者ID userID
3. 当前密码 currPwd
4. 新的昵称 nickname
5. 新的密码 password

规则：

1. 只允许修改自己的昵称和密码
2. 若**昵称**与**密码**均无改动，不应允许执行修改
3. 若成功修改密码，则前端退出登录

#### 修改权限

/api/user/modifyPermission

修改指定用户的权限等级

输入：

1. 修改者ID operatorID
2. 待修改者ID userID
3. 待修改者的新权限等级 permission

规则：

1. 需要修改者权限等级达到**S级管理员权限等级**（可配置）
2. 待修改者**当前权限等级**与**新权限等级**均应**小于S级管理员权限等级**

---

## 云文件系统

编写**前台程序移动端界面**，开放查询模块，支持查询自己上传的文件

### 前台(2)

#### 查询

用户可查看自己上传的全部文件，以及其他权限等级**不高于**自身的用户上传的**公开文件**

查询：

1. 当前用户上传的全部文件 /api/cloudFile/listByUploader
2. 当前用户可查看的公开文件 /api/cloudFile/listPublic

规则：

1. 不查询已删除的文件
2. 用户仅可查看其他权限等级**不高于**自己的用户上传的公开文件

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

这样一来，若已知用户ID、文件ID和扩展名，是可以绕过系统的权限验证、从而看到其他人的非公开文件的。  
除了通过后端获取文件，还有什么办法可以解决吗？

### 后台(5)

#### 上传

/api/cloudFile/upload

用户上传文件到云服务器

考虑到系统只是将用户上传的外部文件存储到服务器，而没有参与文件内容编辑，且**上传文件**使用的更多，故api命名未使用```create```

输入：

1. 上传者ID operatorID
2. 文件名 fileName
3. 扩展名 extensionName
4. 最后修改时间 lastModifiedTime
5. 是否公开 isPublic
6. 文件 file

规则：

1. 允许上传空文件
2. 根据上传者ID与当前时间戳计算hash，作为文件ID，也是文件存储名
3. 先保存文件、后插入数据库记录，若数据库记录新增失败，应删除文件

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
2. 若**文件名**、**扩展名**、**是否公开**均无改动且**未上传新文件**，不应允许执行修改
3. 若上传有新文件，则**新的文件最后修改时间**要求必填
4. 将旧文件移动到私有目录下，不删除
5. 先保存文件、后插入数据库记录，若数据库记录新增失败，应删除文件

#### 删除

/api/cloudFile/delete

用户删除自己上传的文件

输入：

1. 操作员 operatorID
2. 密码 password
3. 文件ID fileID

规则：

1. 仅允许删除自己上传的文件
2. 需要指定文件当前状态为**未删除**
3. 仅修改数据库表中**是否删除**字段

#### 查询

> 见前台对应模块

#### 预览

> 见前台对应模块

---

## 随想系统

记录自己对一些事情的看法和体会

### 前台(2)

#### 查询

查询：

1. 当前用户编辑的笔记 /api/thinkingNote/listByWriter
2. 当前用户可查看的公开笔记 /api/thinkingNote/listPublic

规则：

1. 不查询已删除的笔记
2. 用户仅可查看其他权限等级**不高于**自己的用户编辑的公开笔记

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

规则：

1. 仅主题允许为空

#### 修改

/api/thinkingNote/modify

输入：

1. 修改者ID operatorID
2. 待修改笔记ID noteID
4. 修改者密码 password
5. 新的主题 topic
6. 新的内容 content
7. 是否公开 isPublic

规则：

1. 仅允许修改自己编辑的笔记
2. 仅**新的主题**允许为空
3. 若**主题**、**内容**、**是否公开**均无改动，不应允许执行修改

#### 删除

/api/thinkingNote/delete

输入：

1. 删除者ID operatorID
2. 密码 password
3. 笔记ID noteID

规则：

1. 仅允许删除自己上传的文件
2. 需要指定笔记当前状态为**未删除**
3. 仅修改数据库表中**是否删除**字段

#### 查询

> 见前台对应模块

## 任务系统

记录和管理自己接下来一段时间要学习的知识

UI设计：

左侧侧边栏主体：任务

主体展开：

1. 发布任务：简单的发布页面，参考其他新增页面
2. 查询任务：想要展示任务之间的关系，但是还没想好怎么做
    1. 第一步，获取所有无前置任务的任务，加入**已获取**组
    2. 第二步，获取所有前置任务均在**已获取**组中的任务，将此类任务加入**已获取**组
    3. 重复上一步，直到所有任务均已加入**已获取**组或无法加入（异常，例如一个前置任务因为状态为**历史任务**而未查出）

### 后台

> 先做后台，不确定要不要做一个前台，感觉没啥必要

#### 查询

/api/task/list

获取指定用户发布的所有任务，接口返回所有任务，前端计算各任务之间的关系

输入：

1. 发布者ID operatorID

规则：

1. 查询不分页，但如果目标数据量过大（拟>200），则返回部分（同上设定值）数据并查询失败，提示将部分任务设置为**历史任务**
2. 查询结果按照**修改时间升序**，即越久未修改过的任务越容易被查询到，前端自行计算与展示

返回：

1. 符合条件的数据条数 total
2. 任务 tasks
    1. 任务ID taskID
    2. 任务名 taskName
    3. 描述 description
    4. 前置任务ID列表 preTaskIDs
    5. 任务状态 status
    6. 上次更新 updateTime
    7. 发布时间 createdTime

#### 发布

/api/task/create

发布一个新的任务

输入：

1. 发布者ID operatorID
2. 任务名 taskName
3. 描述 description
4. 前置任务ID列表 preTaskIDs

规则：

1. 若任务名为空或重复，则不应允许发布任务
    1. 重复的任务：两个拥有相同**发布者**和**任务名**的任务

#### 修改

/api/task/modify

修改一个任务，支持**挂载前置任务**功能、**设置任务状态**功能

输入：

1. 修改者ID operatorID
2. 待修改任务ID taskID
3. 修改者密码 password
4. 新的任务名 taskName
5. 新的描述 description
6. 新的前置任务ID列表 preTaskIDs
7. 新的任务状态 status

规则：

1. 仅允许修改自己发布的任务
2. 前端会将原本的值填入修改框作为默认值，新的任务名不允许为空
3. 若**任务名**、**描述**、**前置任务ID列表**与**状态**均无改动，则不应允许执行修改
