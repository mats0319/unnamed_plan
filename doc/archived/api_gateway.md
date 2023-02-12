# API 网关

## 介绍

引入网关层，统一接收、处理外部访问，隐藏内部各业务服务的信息

网关层还可以统一管理外部请求，例如过滤无效请求、多处登录限制，以及熔断、限流、降级等，复杂功能考虑做成插件形式

## 基础功能

网关层的基础功能是转发请求，将接收到的请求转发给相应的业务服务，以下为本项目网关层处理请求的过程：

1. 验证请求有效性
    1. 判断请求方法，忽略全部非`http post`方法的请求
        1. 忽略：直接返回成功
    2. 判断**请求来源**，拒绝所有**无效**的请求来源
        1. 请求来源：自定义http请求头`Unnamed-Plan-Source`字段的值
        2. 无效的请求来源：网关层启动时，会注册有效的来源列表，不存在于该列表中的请求来源，视为无效
2. 执行插件函数：`before invoke`
3. 调用指定请求对应的处理函数
    1. 网关层在启动时注册处理函数，此处根据URI匹配
    2. 若获取指定请求对应的处理函数失败，则返回失败
4. 执行插件函数：`after invoke`
    1. 处理函数执行出错时，不执行该组函数

### 构造与解析参数

前、后端根据proto文件中定义的类型结构进行序列化与反序列化

1. 对应的go类型结构定义：使用protoc-go工具生成
2. 对应的ts类型结构定义：使用自制工具生成
3. 前端序列化使用类型：`FormData`，对应后端反序列化方法：`req.PostFormValue([param_name])`

以下为前、后端序列化、反序列化的示例：

1. 前端反序列化：将http请求返回值（字符串类型）解析为js对象类型，代码中直接使用（见序列化示例代码的返回值部分）
2. 后端序列化：将go的rpc res结构体（根据proto文件生成）通过json marshal序列化成byte数组

```text 
// 前端序列化示例代码：
public login(userName: string, password: string): Promise<User.LoginRes> {
  let request: User.LoginReq = {
    user_name: userName,
    password: calcSHA256(password),
  }

  return axiosWrapper.post("/api/user/login", objectToFormData(request)) // objectToFormData函数定义见下方
}

// 后端反序列化示例代码：
type LoginReqParams rpc_impl.User_LoginReq

func (p *LoginReqParams) Decode(r *http.Request) {
	p.UserName = r.PostFormValue(params_UserName)
	p.Password = r.PostFormValue(params_Password)
}
```

```ts
// objectToFormData 泛型用于解决`obj[key]`报错问题
export function objectToFormData<T extends object>(obj: T): FormData {
  let data = new FormData()
  for (let key in obj) {
    if (typeof obj[key] == "object") { // if field type is another object
      objectToFormData(obj[key] as object).forEach((value: FormDataEntryValue, key: string) => {
        data.append(key, value)
      })
    } else { // normal
      data.append(key, obj[key] as string)
    }

  }

  return data
}
```

## 插件

### 多处登录限制 multi-login limit

限制一个用户到处登录的情况，即一个用户不能同时登录两个`web`

设计：

1. 用户成功登录后，前端保存`token`，并在发起其他请求时携带`token`，供后端验证
2. 后端为`每个用户/每个界面`保存最后一个`token`
    1. 假设用户A在`web`用户界面上登录两次，那么后端会使用第二次的`token`覆盖第一次的`token`，  
       导致携带第一次`token`的请求验证不通过，以实现拒绝请求的目的
