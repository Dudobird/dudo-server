# dudo-server
A private cloud file storage server based on minio



### 1. 开发环境部署

开发环境下，使用docker-compose，将自动创建一个MySQL数据库和一个Minio存储服务器。

```shell
docker-compose up
```

创建额外的测试数据库，用于测试

```shell
docker exec -it $(docker ps --filter "name=dudo-server_mysql_1" -q) bash
(win10下面需要使用winpty docker exec ...)

# 默认密码为root
mysql -u root -p
CREATE DATABASE `dudo-test` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON `dudo-test`.* TO 'dudouser'@'%' IDENTIFIED BY 'dudoadmin' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

执行测试：
```shell
go test ./...
```

### 2. 启动服务器

