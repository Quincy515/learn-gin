### 01. oauth2基本概念

OAuth2（开放授权）

> 是一个开放标准，允许用户授权第三方移动应用访问他们存储在另外的服务提供者上的信息，而不需要将用户名和密码提供给第三方移动应用或分享他们数据的所有内容，OAuth2.0是OAuth协议的延续版本，但不向后兼容OAuth 1.0即完全废止了OAuth1.0 
>
>  最常见的场景就是 QQ登录、微信登录、github登录等

https://github.com/go-oauth2/oauth2

`go get -u github.com/go-oauth2/oauth2/v4`

四种模式

**授权码模式** (Authorization Code) 一般只使用这个，后面3个基本不用~~~

简化模式 (Implicit)

密码模式 (Resource Owner Password Credentials)

客户端模式 (Client Credentials)

```shell
     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+
```

<img src="../imgs/14_code.png" style="zoom:95%;" />

### 02. 客户端请求授权码：基本参数

请求参数

| 参数名称      | 参数含义                         | 是否必须 |
| ------------- | -------------------------------- | -------- |
| response_type | 授权类型，此处的值为code         | 必须     |
| client_id     | 客户端 ID                        | 必须     |
| redirect_uri  | 重定向 URI                       | 必须     |
| scope         | 申请的权限范围，多个逗号隔开     | 可选     |
| state         | 客户端的当前状态，可以指定任意值 | 可选     |

授权码返回

| 参数名称 | 参数含义                                          | 是否必须 |
| -------- | ------------------------------------------------- | -------- |
| code     | 授权码。认证服务器返回的授权码                    | 必须     |
| state    | 如果A中请求包含这个参数，资源服务器原封不动的返回 | 可选     |

