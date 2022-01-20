# 配置中心

一个具体功能的API：

1. uri
2. 描述
3. 输入：输入参数列表
4. 规则：功能逻辑、规则
5. 返回：输出参数列表

## 配置中心

### 身份验证

#### 登录

/api/cc/login

输入：

1. 用户名 userName
2. 密码 password

### 配置的版本管理

1. **生成**表示生成新的全部服务的配置版本
2. **应用**表示应用一个版本的配置

#### 查询应用中的版本配置

/api/cc/version/listUsing

查询应用中的版本配置

返回：

1. 版本ID versionID
2. 版本号 versionNum
3. 描述 description
4. 支持的服务id列表 serviceIDs
5. 用到的配置id列表 configIDs
6. 支持的服务对应挂载的配置详情 configurations
7. 是否有更新 hasUpdate
8. 修改时间 updateTime

#### 查询版本配置

/api/cc/version/list

查询版本配置

输入：

1. 分页信息
    1. 每页条数 pageSize
    2. 当前页码 pageNum

返回：

1. 符合条件的数据条数 total
2. 版本配置列表 versions
    1. 版本的数据库记录ID id
    2. 版本号 versionNum
    3. 描述 description
    4. 支持的服务id列表 serviceIDs
    5. 用到的配置id列表 configs
    6. 支持的服务对应挂载的配置 configurations
    7. 是否正在应用 isUsing
    8. 是否有更新 hasUpdate
    9. 修改时间 updateTime

#### 生成

/api/cc/version/create

输入：

1. 版本号 versionNum
2. 描述 description
3. 支持的服务id列表 serviceIDs
4. 支持的服务对应挂载的配置 configurations
5. 是否直接应用该版本的配置 isUsing

#### 修改描述

/api/cc/version/modify

修改指定版本配置的描述信息

输入：

1. 新的描述 description

#### 应用

/api/cc/version/apply

输入：

1. 版本ID：versionID

规则：

1. 同一时刻只能有一个版本的配置处于**应用中**
2. 应用一个新的版本，会取消旧版本的**应用中**状态与**是否有更新**标识
    1. 版本配置是否有更新：在修改服务和修改配置时，若改动到了当前版本支持的服务或用到的配置，则视为当前版本配置有更新

### 维护服务

以用户服务举例：

1. **新增**表示配置中心开始维护用户服务的配置，例如允许用户服务获取配置
2. **修改**表示配置中心调整了用户服务的配置，例如新增数据库配置、删除启动配置
3. **删除**表示配置中心停止维护用户服务的配置，不再为用户服务提供配置

用户服务包含的配置以**配置项ID列表**的形式存储

#### 查询

/api/cc/service/list

输入：

1. 分页信息
    1. 每页条数 pageSize
    2. 当前页码 pageNum

返回：

1. 符合条件的数据条数 total
2. 服务列表 services
    1. 数据库记录ID id
    2. 服务ID serviceID
    3. 服务名 serviceName
    4. 挂载的配置id列表 configIDs
    5. 是否隐藏 isShadow
    6. 修改时间 updateTime

#### 新增

/api/cc/service/create

新增对指定服务的支持

输入：

1. 服务ID serviceID
2. 服务名 serviceName
3. 挂载的配置id列表 configIDs
4. 是否隐藏 isShadow

#### 修改

/api/cc/service/modify

修改服务信息，支持服务隐藏功能、服务挂载功能

输入：

1. 服务的数据库记录ID id
2. 新的服务ID serviceID
3. 新的服务名 serviceName
4. 挂载的新的配置id列表 configIDs
5. 是否隐藏 isShadow

规则：

1. 若输入参数（共4项）**全部为空**或**全部与指定服务当前参数相同**，则不应允许执行修改
2. 若修改了当前应用中的版本配置支持的服务，则应为当前应用中的版本配置设置**有更新**标识

#### 删除

/api/cc/service/delete

不再支持指定服务

输入：

1. 服务的数据库记录ID id

规则：

1. 要求当前应用中的版本配置不支持该服务

### 维护配置详情

举个例子，**数据库**配置可能包含**连接地址**、**数据库名**配置项，以此约定：

1. **数据库**称为一个**配置**
2. **连接地址**等称为一个**配置项**
3. **数据库包含的全部配置项**称为**配置详情**

以数据库连接配置举例：

1. **新增**表示开始维护数据库连接的配置，之前没有
2. **修改**表示调整了数据库连接配置的某些配置项，例如新增DBMSName、修改数据库连接地址
3. **删除**表示停止维护数据库连接的配置，不再允许服务挂载数据库配置

#### 查询

/api/cc/config/list

输入：

1. 分页信息
    1. 每页条数 pageSize
    2. 当前页码 pageNum

返回：

1. 符合条件的数据条数 total
2. 配置列表 configs
    1. 数据库记录ID id
    2. 配置ID configID
    3. 配置名 configName
    4. 配置详情 payload
    5. 是否隐藏 isShadow
    6. 修改时间 updateTime

#### 新增

/api/cc/config/create

输入：

1. 配置ID configID
2. 配置名 configName
3. 配置详情 payload
4. 是否隐藏 isShadow

#### 修改

/api/cc/config/modify

修改配置信息，支持配置隐藏功能

输入：

1. 新的配置ID configID
2. 新的配置名 configName
3. 配置详情 payload
4. 是否隐藏 isShadow

规则：

1. 若输入参数（共4项）**全部为空**或**全部与指定配置当前状态相同**，则不应允许执行修改
2. 若修改了当前应用中的版本配置支持的服务所挂载的配置，则应为当前应用中的版本配置设置**有更新**标识

#### 删除

/api/cc/config/delete

输入：

1. 配置的数据库记录ID id

规则：

1. 要求当前应用的版本未使用该配置
