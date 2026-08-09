# 在公网上部署应用

> 服务器系统：Ubuntu 24.04

## 数据库(postgresql)

ubuntu系统自带一个pg数据库的指定版本快照，如果想要安装其他版本，参考[官方文档](https://www.postgresql.org/download/linux/ubuntu/)

因为我们想要安装pg 18（这个版本原生支持uuidv7），所以记录我的安装过程：

1. 系统自带pg 16
2. 使apt可以从pg下载：
    - `sudo apt install -y postgresql-common`
    - `sudo /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh`
3. 下载pg 18：
    - `sudo apt update`
    - `sudo apt install postgresql-18`
4. 下载完成之后，终端弹窗说可以直接把16升到18，选择升级，它执行了两个`pg_xxx`工具就好了；
   数据可以正常访问，对外部使用来说没有影响，如果你的升级没有这一步，参考官方文档继续进行后续步骤：
   [文档](https://www.postgresql.org/docs/current/index.html)

### 创建用户/数据库

> 记得在云服务器厂商的控制台，编辑网络策略：允许访问5432端口（如果你使用了非默认的端口号，这里相应修改）

```cmd 
su - postgres // 切换linux用户
psql

create user [user name] with password '[password]'; // 注意分号、单引号
create database [db name] owner [user name];
grant all on database [db name] to [user name];

\du+ // 查看用户
\l   // 查看数据库
```

### 限制访问权限

config file: `/etc/postgresql/16/main/postgresql.conf`(注意版本号)
add: `listen_address: '*'`

config file: `/etc/postgresql/16/main/pg_hba.conf`
add:

```txt
host cloud all 127.0.0.1/32 md5     // 允许本机访问cloud数据库
host cloud all 0.0.0.0/0    reject  // 不允许其他任何人访问cloud
host all   all 0.0.0.0/0    md5     // 允许其他所有访问
// 写在前面的优先级更高，考虑可能是从后往前覆盖式计算规则，例如本例：允许所有人访问所有数据库-禁止所有人访问cloud-允许本机访问cloud
```

数据库可视化工具往往可以通过先登录云服务器再访问的方式，绕过远程访问限制

重启/重载服务：`sudo systemctl restart/reload postgresql`

## nginx

### 下载与安装

可以使用包管理工具下载，或参考[nginx官网](https://nginx.org/en/linux_packages.html#Ubuntu)提供的下载与安装过程

- `sudo apt install nginx`
- `sudo systemctl start nginx`
- `curl http://127.0.0.1` // 验证安装结果

### 编辑反向代理规则

配置文件路径：`/etc/nginx/nginx.conf`
日志文件路径：`/var/log/nginx/access.log`（配置文件内可修改）

默认配置包含`/etc/nginx/conf.d/*.conf`，所以我们将编辑好的`.conf`文件丢到那个文件夹下面即可

详细配置，见本文档同目录下，`*.conf`配置文件

- 如果想要使用https访问，可以使用自签证书（生成方式见下一节）

nginx常用命令：

- `nginx -t`：测试配置文件的语法是否正确
- `nginx -s reload`：重启nginx，更新并应用新的配置文件

### 自签证书

> 目的是允许使用https访问

生成http证书：

- `openssl genrsa -out server.key 2048`
- `openssl req -new -key server.key -out server.csr`
    - 会以问答的形式，要求我们输入一些信息，包括**域名**(Common Name)在内的一些关键信息就是在此处提供的
- `openssl x509 -req -in server.csr -out server.crt -signkey server.key -days 3650`

## 服务器设置

### 设置使用密钥登录

- `ssh-keygen`生成密钥对
    - 默认在`/root/.ssh/`目录生成`id_ed25519`和`id_ed25519.pub`文件
    - 将没有后缀名的私钥下载到本地
- 设置公钥：`cat ./id_ed25519.pub >> ./authorized_keys`(注意切换路径)
    - 将`id_ed25519.pub`写入`authorized_keys`
- 本地验证公钥的有效性：`ssh -i ./id_ed25519 localhost`
    - 修改ssh服务器安全策略：`sudo vim /etc/ssh/sshd_config`
        - 找到以下配置项并修改为对应值
          ```txt
          PermitRootLogin without-password
          PubkeyAuthentication yes
          PasswordAuthentication no // 禁用密码登录
          ```
        - 测试语法是否正确：`sudo ssh -t`
        - 重启/检查ssh服务：`sudo systemctl restart/status ssh`

### 设置开机自启动

在`/etc/rc.local`中，编辑想要在开机时执行的代码

```txt
#!/bin/sh -e
([.exe path] &) 
exit 0
```

&：不阻塞
()：在命令行窗口关闭后仍能持续运行

## 部署web应用

构建ui、服务端程序，参考`scripts/build.sh`

将构建内容上传到云服务器，启动/重新启动服务端程序，参考`scripts/restart_server.sh`
