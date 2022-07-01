# 命名

项目中优先使用统一的命名，未使用之处需要特别标注

## http

1. 请求的发起者ID：operatorID
2. 分页信息
    1. 每页条数：pageSize
    2. 当前页码：pageNum （当前页码从1开始）
    3. 符合条件的数据总条数：total
3. 创建时间：createdTime
4. 更新时间：updateTime
5. 是否公开：isPublic

6. 用户列表：users
7. 用户ID：userID
8. 用户名：userName
9. 昵称：nickname
10. 锁定状态：isLocked
11. 密码：password
12. 权限等级：permission
13. 创建人：createdBy

14. 文件列表：files
15. 文件ID：fileID
16. 文件名：fileName
17. 上次修改时间：lastModifiedTime
18. 访问路径：fileURL
19. 扩展名：extensionName
20. 文件大小：fileSize
21. 文件：file

22. 编辑者：writeBy
23. 主题：topic
24. 内容：content
25. 笔记：notes
26. 笔记ID：noteID

27. 任务名：taskName
28. 描述：description
29. 前置任务ID列表：preTaskIDs
30. 任务ID：taskID
31. 任务：tasks

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

## 数据库唯一字段

我们的每一张表都组合了一个`Common`结构体，其中包含`ID`字段，显式定义为数据库主键，其值为uuid

提问：以文件表举例，在有一个主键的前提下，还需要定义一个唯一标识字段`FileID`吗？

回答：尽量使用主键，在后端代码内使用主键作为唯一索引；而在后端代码之外，例如业务需求中，若需要唯一标识，则应定义新的唯一标识字段

举个例子，功能：**允许用户在本地上传文件并保存到云服务器**  
一个用户可能上传多个同名文件，不同用户就更可能了，所以用户上传的文件名无法直接用于在云服务器上保存文件，这时就需要一个新的唯一字段：**文件存储名**  
为了与**用户上传的文件名**做出区分，我们将新的字段命名为`FileID`

那么可不可以用主键作为文件名呢？可以，但不推荐。

因为主键的作用就是唯一标识一条数据库表记录，不应用于后端代码之外，例如上例中的**文件名**；  
如果我们使用主键作为文件名，可能无法很好地应对需求迭代，举个例子，假设有新需求：

1. 云服务器上的文件名要格式统一且有一定意义，采用`[file name] + hash(user id + timestamp)[:4]`格式

此时，新的云文件存储名应该保存到数据库，否则将无法把**数据库记录**与**云文件**一一对应起来；  
无法关联数据库记录和文件，可能会在处理问题时，带来额外的复杂性（例如前端通过url直接访问服务器文件）；  
同时，为了保证**文件存储名**唯一，为其数据库字段`FileID`添加唯一约束

1. 为什么要保证`FileID`唯一：因为该字段同时应用于**服务器保存文件**场景，若不能保证该字段唯一，对应场景（服务器保存文件）将发生**未定义**行为
   （此处使用C语言概念，未定义行为指行为结果不确定，如本例，若`FileID`重复，服务器如何保存文件，将部分依赖操作系统——在后端代码未处理**写入已存在的文件**场景时）

这样一来，我们就通过实现功能的过程，反推，得到了一条编程规则

总结：

1. 后端代码内部需要使用唯一标识时，使用主键
2. 后端代码以外需要使用唯一标识时（例如业务场景），使用自定义唯一标识字段

## rpc方法名称定义

先上例子，以下为云文件服务的全部rpc方法：

```protobuf 
service ICloudFile {
  rpc ListByUploader(CloudFile.ListByUploaderReq) returns (CloudFile.ListByUploaderRes);
  rpc ListPublic(CloudFile.ListPublicReq) returns (CloudFile.ListPublicRes);
  rpc Upload(CloudFile.UploadReq) returns (CloudFile.UploadRes);
  rpc Modify(CloudFile.ModifyReq) returns (CloudFile.ModifyRes);
  rpc Delete(CloudFile.DeleteReq) returns (CloudFile.DeleteRes);
}
```

从例子中可以看到，云文件服务的rpc方法是以**客户端**的角度命名的，即**客户端`upload`文件到云文件服务**

为什么这么设计呢？因为在写代码的时候，是从客户端调用该方法的，即`client.upload()`，  
如果以服务端的角度命名，客户端调用的时候就变成了`client.onUpload()`，略显奇怪

好的，在接受了以上设定之后，我们为一组可以互相调用的服务设计rpc方法，应该把哪个方法写到哪个服务里，就有迹可循了，以下是注册中心的例子。

定义：

1. 服务注册中心-核心模块（core，以下简称c），单独的服务模块，主要负责维护业务服务的全量路由表
2. 服务注册中心-嵌入模块（embedded，以下简称e），嵌入到每个业务服务当中，主要负责与c沟通，以及调用其他业务服务

功能：

1. e可以向c注册
2. e可以向c获取指定服务的实例地址列表
3. c可以向e发送心跳包

代码：

```protobuf 
service IRegistrationCenterCore {
  rpc Register(RegistrationCenterCore.RegisterReq) returns (RegistrationCenterCore.RegisterRes);
  rpc ListServiceTarget(RegistrationCenterCore.ListServiceTargetReq) returns (RegistrationCenterCore.ListServiceTargetRes);
}

service IRegistrationCenterEmbedded {
  rpc CheckHealth(RegistrationCenterEmbedded.CheckHealthReq) returns (RegistrationCenterEmbedded.CheckHealthRes);
}
```
