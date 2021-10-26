# 未命名计划——unnamed plan

## 运行

### ui

用户界面的公共资源位于```shared/ui```，各用户界面通过```npm link```与之建立联系，执行步骤：

1. ```cd shared/ui```
2. ```npm link```/```npm run create-link```，注册公共资源库
3. ```cd admin_ui```
4. ```npm link shared_ui```/```npm run link```，链接公共资源库，其中```shared_ui```是公共资源库名称

然后，在程序中就可以使用```shared_ui```表示```shared/ui```路径了
