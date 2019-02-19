# dudo-server
A private cloud file storage server based on minio

### 1. 开发环境部署

开发环境部署MySQL数据库，可以直接下载二进制安装在本地，或者使用启动一个本地MySQL容器
##### Windows 版本下的启动方式
```sh
# windows版本
# 需要在docker desktop的设置中启动共享存储，这里将的e盘共享用于数据存储
$ docker run --name dudo-mysql -p 3306:3306  -v {数据库本地存储目录}:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.6 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
$ winpty docker exec -it $(docker ps --filter "name=dudo-mysql" -q) sh

```
##### Linux 版本下的启动方式

```sh
$ docker run --name dudo-mysql -p 3306:3306  -v {数据库本地存储目录}:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.6 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
$ docker exec -it $(docker ps --filter "name=dudo-mysql" -q) sh

```

进入到数据库中后执行下面的语句来创建必要的数据库和用户,用于程序访问

```sql
create database dudo;
GRANT ALL PRIVILEGES ON dudo.* TO 'dudouser'@'%' IDENTIFIED BY 'dudoadmin' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

### 2. 启动服务器

