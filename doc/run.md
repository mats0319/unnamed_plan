# 运行

## npm link

[reference](https://docs.npmjs.com/cli/v8/commands/npm-link)

用户界面的公共资源位于`ui/shared`，各用户界面通过`npm link`与之建立联系，执行步骤：

```cmd 
cd ui/shared
npm link // 注册资源库

cd ui/admin_web
npm link shared // 链接资源库
```

`shared`是公共资源库名称，定义于资源库`package.json`，与项目路径无关

然后，在程序中就可以使用`shared`表示`ui/shared`文件夹了

另外一种方法：

```cmd 
cd ui/admin_web
npm link ../shared
```

直接进入主要项目路径，然后通过**路径名**链接另一个库
