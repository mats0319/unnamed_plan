# 在公网上部署应用

> 服务器系统：CentOS 8

## 设置使用私钥登录并禁用密码登录

### 生成密钥对

`ssh-keygen`

1. 默认在`/root/.ssh/`目录下，生成`id_rsa`和`id_rsa.pub`文件
2. 相同目录下，新建文件`authorized_keys`，内容与**公钥文件**相同
3. 下载私钥文件到本地

### 启用秘钥登录与禁用密码登录

config file: `/etc/ssh/sshd_config`

1. `PubkeyAuthentication`，设置为`true`
2. `PasswordAuthentication`，设置为`false`

`service sshd restart`

## nginx

### 下载与安装

[reference](https://nginx.org/en/linux_packages.html#RHEL-CentOS)

参考nginx官网提供的下载与安装过程

### 编辑反向代理规则

见同目录下，`nginx`配置文件

### 自签证书

> 使用openssl

`openssl genrsa -out server.key 2048`
`openssl req -new -key server.key -out server.csr`
`openssl x509 -req -in server.csr -out server.crt -signkey server.key -days 3650`

1. `server.crt`证书格式为`pem`
2. 第二步会以问答的形式，要求我们输入一些信息，包括**域名**在内的一些关键信息就是在此处提供的

## postgresql

### 下载与安装

> 使用dnf，centos默认的包管理工具

[reference](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-centos-8)

参考上方链接，下载与安装postgresql数据库，以下为命令与注释：

```shell 
dnf module list postgresql // 查看可用的版本
sudo dnf module enable postgresql:13 // 选择想要安装的版本
sudo dnf install postgresql-server // 安装 优先安装enabled版本
sudo postgresql-setup --initdb // 初始化

sudo systemctl start postgresql // 启动
sudo systemctl enable postgresql // 设置开机自启动
```

### 设置允许外部地址访问

config file: `/var/lib/pgsql/data/postgresql.conf`

add: `listening_address: '*'`

config file: `/var/lib/pgsql/data/pg_hba.conf`

add: `host all all 0.0.0.0/0 md5`

### 创建用户/数据库

```shell 
su - postgres // 切换linux用户
psql

create user [user name] with password '[password]'; // 注意分号
create database [db name] owner [user name];
grant all on database [db name] to [user name];
\q // 退出
```

## 设置开机自启动

在`/etc/rc.d/rc.local`中，编辑想要在开机时执行的代码
