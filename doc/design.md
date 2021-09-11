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

1. 符合条件的数据量：total
2. 用户列表 users
    1. 用户ID userID
    2. 昵称 nickname
    3. 用户权限等级 permission

#### 创建

/api/user/create

输入：

1. 创建者ID operatorID
2. 新用户的用户名 userName
3. 新用户的密码 password
4. 新用户的权限等级 permission

规则：

1. 创建者权限等级达到一定值（见配置文件）
2. 只能创建比自己权限等级低的用户
3. 新昵称默认与用户名相同

返回：

1. 创建结果 isSuccess

#### 锁定

/api/user/lock

禁用指定用户

输入：

1. 发起锁定的用户ID operatorID
2. 待锁定用户ID userID

规则：

1. 需要待锁定用户当前是未锁定状态
2. 需要操作员权限等级达到一定值（见配置文件）
3. 只能锁定比自己权限等级低的用户
4. 用户可以锁定自己
    1. 锁定自己不验证权限，锁定成功后退出登录

返回：

1. 锁定结果 isSuccess

#### 解锁

/api/user/unlock

恢复使用指定账户

输入：

1. 发起解锁的用户ID operatorID
2. 待解锁用户ID userID

规则：

1. 需要待解锁用户当前是锁定状态
2. 需要操作员权限等级达到一定值（见配置文件）
3. 只能解锁比自己权限等级低的用户

返回：

1. 解锁结果 isSuccess

#### 修改用户信息（昵称和密码）

/api/user/modifyInfo

修改自己的昵称和密码

输入：

1. 发起修改的用户ID operatorID （当前登录用户）
2. 用户ID userID （想要修改的用户）
3. 新的昵称 nickname
4. 新的密码 password

规则：

1. 只允许修改自己的用户名和密码
2. 修改成功后，退出登录

输出：

1. 修改结果 isSuccess

#### 修改权限

/api/user/modifyPermission

修改指定用户的权限等级

输入：

1. 发起修改的用户ID operatorID
2. 待修改权限的用户ID userID
3. 目标用户的新权限等级 permission

规则：

1. 需要操作员权限等级达到一定值（见配置文件，记为x）
2. 目标用户权限等级应**小于**```x```
3. 目标用户修改后的权限等级应**小于**```x```

返回：

1. 修改结果 isSuccess

## 文档结构

### 功能函数文档结构

1. uri
2. 功能描述
3. 输入参数（参数描述与序列化字段名）
4. 特殊规则（可选）
5. 返回结果（参数描述与序列化字段名）
