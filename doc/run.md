# 运行

## npm link

用户界面的公共资源位于```ui/shared```，各用户界面通过```npm link```与之建立联系，执行步骤：

1. ```cd ui/shared```
2. ```npm link```/```npm run create-link```，注册公共资源库
3. ```cd ui/admin_web```
4. ```npm link shared```/```npm run link```，链接公共资源库
   1. ```shared```是公共资源库名称，定义于资源库```package.json```

然后，在程序中就可以使用```shared```表示```ui/shared```路径了
