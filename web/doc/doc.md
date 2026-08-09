# web

## 目录

public: 图标、开发阶段的wasm文件

src:
- axios: http请求相关
- components: 组件，专注于解决一个问题的小组件
- pinia: 响应式的全局变量及其配套函数
- ts: ts功能代码
- views: 页面
  - components: 页面组件，页面的一部分，包含多种职能
  - folders: 如果一个页面包含多个子页面，则将它们放在一个文件夹里

## 说明

`wasm`文件打包问题：开发阶段使用iframe嵌入`public/flip.html`，而发布以后使用github pages托管的html，所以游戏相关文件不打包进dist

- 通过环境变量文件(`.env.development`)实现使用不同来源的`wasm`文件
- 为什么使用github pages托管游戏html：因为网络带宽，我购买的服务器网络带宽非常小，下载5k资源需要40多秒
- 如何使用github pages托管资源文件：在github新建仓库并上传想要托管的资源，在setting-pages里开启部署
- 如何将游戏文件从dist移除：在打包配置里编写build hock检查并删除指定文件

## 使用工具

- node: 24 (corepack-pnpm:10)
- vue3 (html ts:5.9 less)
- vite
- vue-router
- pinia
- element plus
- axios
- eslint

node在16.13开始，提供`corepack`工具，这是一个**包管理器的管理器**

```cmd
npm  install --global corepack@latest
corepack enable pnpm
corepack use pnpm@latest // 升级pnpm版本
```

pnpm常用功能：

- `pnpm approve-builds`：允许依赖执行脚本，pnpm默认不允许任何依赖运行脚本
- `pnpm  audit`：检查已安装的包是否存在已知的安全问题
- `pnpm  outdated`：检查依赖是否有更新
- `pnpm why/ls [package name]`：列举一个依赖项在项目中的依赖关系，未指定依赖项时列举全部

## 开发小技巧

- TS7016:找不到指定模块的声明文件
  下载对应的声明文件，例如`import CryptoJs from "crypto-js"`出现该错误，则安装`@types/crypto-js`为开发依赖即可
- 样式优先级：class < style < important class
  在同一个组件上，使用style直接定义的样式会覆盖class中的样式；而使用`class xxx !important`定义的样式又会覆盖style样式
- 在vue代码中使用代码(`src`)以外的值（例如`package.json`中的`version`配置项）：声明和赋值分开
  在`global.d.ts`中声明变量名和类型，在`vite.config.ts`为变量赋值
