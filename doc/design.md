# 设计

## API 网关

> 只处理http post请求，忽略其他类型的请求
> 忽略：直接返回成功

### http请求响应过程

> 加载插件

1. 判断请求方法，忽略全部非`http post`方法的请求
2. 判断**请求来源**，若请求来源**无效**（未注册），则拒绝该请求
    1. 请求来源：http请求头`Unnamed-Plan-Source`字段的值
    2. 无效：启动时注册有效的来源列表，不存在于列表中的请求来源视为无效
3. 插件钩子函数：`before invoke`
4. 执行对应请求的处理函数
    1. 所有处理函数在程序启动时注册，此处根据URI获取
5. 插件钩子函数：`after invoke`
    1. 处理函数执行出错时，不执行`after invoke`函数

### http res

```go 
type ResponseData struct {
	HasError bool   `json:"hasError"`
	Data     string `json:"data"`
}
```

```hasError```表示http请求是否正确执行并按照预期返回：

1. 请求正常执行：```data```是```json string```
2. 请求执行失败：```data```是```string```类型的错误信息

### 多处登录限制 multi-login limit

#### 目的

限制同一用户同时登录多个终端的情况

#### 策略

1. 同一用户，在同一个**限制登录**的来源上，**最多**只能登录一次
2. 当一名用户尝试在同一个来源的**不同终端**上登录时，先登录的终端将被收回登录权限
    1. 不同终端：例如两个浏览器标签页

#### 实现

1. 前端将```token```保存到```sessionStorage```
2. 前端通过自定义的http请求头携带身份验证参数
    1. 来源：```Unnamed-Plan-Source```
    2. 用户ID：```Unnamed-Plan-User```
    3. token：```Unnamed-Plan-Token```
3. 满足以下全部条件的，视为验证通过：
    1. 解析出有效的```user id```与```token```
    2. 程序内保存有该```user id```的登录信息
    3. 该用户在该来源的**上一次请求时间**距**此刻**不超过一定值（可配置）

#### 问题

1. 考虑引入缓存层保存用户登录信息
2. 后端识别**请求头**中的关键参数，然后把关键参数写在**请求体**里返回，  
   这种做法能够最大程度兼容```go http server```和```vue axios interceptor```的原生写法  
   简单来说就是容易实现，后续考虑统一
    1. 考虑使用服务端设置cookie的方式，涉及跨域问题

### 使用ip可能带来的共用session storage问题

> 强调：使用不同的域名（源）不会存在该问题

没有买域名，对于前台界面和后台界面的路径，参考nginx配置：（ip仅供示意）

1. 前台界面：https://117.50.177.201
2. 后台界面：https://117.50.177.201/admin

这就导致两个网站会共用session storage，因为我将登录信息保存在vuex、使用session storage做的vuex防刷新，  
以至于当我登录了前台界面之后，在当前浏览器标签页，通过修改url的方式访问后台界面，后台界面能够读到前台界面的vuex

```ts vue  
// vuex防刷新
export default class App extends Vue {
  private created() {
    if (sessionStorage.getItem(process.env.VUE_APP_axios_source_sign)) {
      this.$store.replaceState(
        Object.assign(
          {},
          this.$store.state,
          JSON.parse(sessionStorage.getItem(process.env.VUE_APP_axios_source_sign))
        )
      );

      sessionStorage.removeItem(process.env.VUE_APP_axios_source_sign);
    }
  }

  private mounted() {
    window.addEventListener("beforeunload", () => {
      sessionStorage.setItem(process.env.VUE_APP_axios_source_sign, JSON.stringify(this.$store.state));
    });
  }
}

// 登录成功，写session storage
sessionStorage.setItem("auth", process.env.VUE_APP_axios_source_sign);
```

又因为我没有为两个网站做代码区分，即**是否登录**字段都叫`isLogin`，所以仅从界面来看，我做到了特殊情况下的单点登录...

对此，我的解决办法是，不同的网站，vuex写session storage时，使用不同的key，以避免与其他网站的session storage混用

效果：（访问其他网站，默认在同一个浏览器标签页内）

1. 访问前台界面，登录，访问后台界面：未登录
2. 访问前台界面，登录，访问后台界面，访问前台界面：登录
3. 访问前台界面，登录，访问后台界面，登录，访问前台界面：登录

## 服务测试(draft)

当前实现：（业务服务）

1. init阶段加载配置
2. init阶段连接数据库

修改：

1. main阶段手动加载配置
2. 临时取消命令行解析功能
3. 取消默认配置文件名，改成允许传入
4. 连接数据库部分，先做检查，如果配置不对，则拒绝连接
5. 然后以服务为单位，编写测试：
    1. 测试不需要启动服务，而是构造一个服务实例，然后针对每一个功能函数构造请求参数
    2. 问题：数据库部分是mock还是新建表？  
       拟新建测试表，每次启动测试，首先清空全部测试表、构造测试数据，然后再进行测试
6. 写完每个服务的测试，再编写统一的测试脚本，直接go test只会在当前目录执行测试
7. 思考：测试是单独写一个powershell，还是引入makefile呢？  
   拟引入makefile

## 配置中心(draft)

> 因未达成设计意图（没想明白），于下一个提交中删除配置中心模块，恢复草稿状态

## 服务注册中心(draft)

目的：解决业务服务实例数量提升所带来的新问题，例如负载均衡、实例状态维护等

1. 服务注册中心实现所有业务服务的客户端代码与转发代码
    1. 调用链调整为：API网关/其他服务 -> 服务注册中心 -> 目标服务
    2. 服务注册中心应实现全部rpc方法的服务端与客户端，即，服务注册中心在收到请求时，仅做转发
    3. 是否要支持**查询服务注册中心当前支持的业务函数**呢？
    4. note：跨服务调用通过服务注册中心，是因为这样调用方就不需要知道目标服务地址了，负载均衡、目标服务状态等都不需要关心

## 性能监控(draft)

pprof

## 调用链追踪(draft)

日志收集，一个请求在各环节拥有同一个请求ID

执行业务的服务中增加日志

rpc方法拟使用自定义错误类型

## 熔断、限流、降级(draft)

熔断、降级拟合并到服务注册中心，限流拟合并到API网关
