# dudo-server
A private cloud file storage server based on minio



### 1. 开发环境部署

开发环境下，使用docker-compose，将自动创建一个MySQL数据库和一个Minio存储服务器。

```
docker-compose up
```

创建额外的测试数据库

```
docker exec -it $(docker ps --filter "name=dudo-mysql" -q) bash

CREATE DATABASE `dudo-test` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON `dudo-test`.* TO 'dudouser'@'%' IDENTIFIED BY 'dudoadmin' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```


### 2. 启动服务器

