# API 网关

## 介绍

引入网关层，统一接收、处理外部访问，隐藏内部各业务服务的信息

网关层还可以统一管理外部请求，例如过滤无效请求、多处登录限制，以及熔断、限流、降级等，复杂功能考虑做成插件形式

## 基础功能

网关层的基础功能是转发请求，将接收到的请求转发给相应的业务服务，以下为本项目网关层处理请求的过程：

1. 判断请求方法，忽略全部非`http post`方法的请求
    1. 忽略：直接返回成功
2. 判断**请求来源**，拒绝所有**无效**的请求来源
    1. 请求来源：自定义http请求头`Unnamed-Plan-Source`字段的值
    2. 无效的请求来源：网关层启动时，会注册有效的来源列表，不存在于该列表中的请求来源，视为无效
3. 执行插件函数：`before invoke`
4. 调用指定请求对应的处理函数
    1. 网关层在启动时注册处理函数，此处根据URI获取
    2. 若获取指定请求对应的处理函数失败，则返回失败
5. 执行插件函数：`after invoke`
    1. 处理函数执行出错时，不执行该组函数

### 构造与解析参数

1. 根据参数名，从`request body`中获取参数的值
2. 参考下方返回结构，`has error`字段表示请求是否正确执行，`data`字段的含义则根据`has error`，有所不同：
    1. 请求正常执行：`data`是`json string`类型的执行结果
    2. 请求执行失败：`data`是`string`类型的错误信息

```go 
type ResponseData struct {
	HasError bool   `json:"hasError"`
	Data     string `json:"data"`
}
```

前端构造参数示例代码：

```ts 
function login(userName: string, password: string) {
  const data: FormData = new FormData();
  data.append("userName", userName);
  data.append("password", calcSHA256(password));
  return axiosWrapper.post("/api/login", data);
}
```

后端解析参数示例代码：

```go 
func (p *LoginReqParams) Decode(r *http.Request) string {
	p.UserName = r.PostFormValue(params_UserName)
	p.Password = r.PostFormValue(params_Password)

	return ""
}
```

## 插件

### 多处登录限制 multi-login limit

#### 介绍

限制一个用户到处登录的情况

一个用户，可以同时登录`public web`/`admin web`，但不能同时登录两个`public web`

#### 设计

> 本节默认请求来源有效

1. 前端登录，获得`token`，保存到`session storage`
2. 前端其他请求将携带`token`，后端校验，满足以下全部条件，视为验证通过，继续执行请求：
    1. 后端能够从请求中解析出符合格式的`user id`与`token`
    2. `user id`与`token`是匹配的
    3. 该用户在该来源的**上一次请求时间**距**此刻**不超过一定值（可配置）
3. 前端的每一个登录请求都会生成新的`token`，后端仅记录最后一个
    1. 这样一来，新的登录请求成功时，旧的登录失效（旧`token`验证不通过）

## 记录

### 使用ip可能带来的共用session storage问题

> 强调：使用不同的域名（源）不会存在该问题

没有买域名，对于前台界面和后台界面的路径，参考nginx配置：（ip仅供示意）

1. public：https://117.50.177.201
2. admin ：https://117.50.177.201/admin

这就导致两个网站会共用session storage，因为我将登录信息保存在vuex、使用session storage做的vuex防刷新，  
以至于当我登录了前台界面之后，在当前浏览器标签页，通过修改url的方式访问后台界面，后台界面能够读到前台界面的vuex

又因为我没有为两个网站做代码区分，即vuex的`key`都叫`vuex`、**是否登录**字段都叫`isLogin`，所以仅从界面来看，我做到了特殊情况下的单点登录...

对此，我的解决办法是，不同的网站，vuex写session storage时，使用不同的key，以避免与其他网站的session storage混用

## 优化

1. 考虑引入缓存层保存用户登录信息
2. 后端识别**请求头**中的关键参数，然后把关键参数写在**请求体**里返回，  
   这种做法能够最大程度兼容`go http server`和`vue axios interceptor`的原生写法  
   简单来说就是容易实现，后续考虑统一
   1. 考虑使用服务端设置cookie的方式，涉及跨域问题
