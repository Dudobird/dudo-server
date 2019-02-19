### Dudo-Server API列表


#### 1. 用户认证

提供涉及用户注册登入以及密码管理等授权相关的API接口描述

##### 1.1 注册用户

POST /api/auth/signup

```
{
    "email":"",
    "password":""
}
```
Response 

```
{
	"status":"",
	"message":"",
	"data" : {
		"email":"",
		"token":"",
	}
}

```

##### 1.2 用户登入

POST /api/auth/signin

```
{
    "email":"",
    "password":""
}
```
Response 

```
{
	"status":"",
	"message":"",
	"data" : {
		"email":"",
		"token":"",
	}
}

```



##### 1.3 用户退出
用户退出必须在登入状态，发送的HTTP请求必须携带必要的JWT Token
GET /api/auth/logout

Response 

```
{
	"status":"",
	"message":""
}

```
/api/auth/logout GET

##### 1.4 密码修改
修改密码必须在登入状态，发送的HTTP请求必须携带必要的JWT Token
UPDATE /api/auth/password 
```
{
    "email":"",
    "password":""
}
```
Response 

```
{
	"status":"",
	"message":"",
	"data" : {
		"email":"",
		"token":"",
	}
}

```