# SimpleCA

[![Build Status](https://travis-ci.com/520MianXiangDuiXiang520/SimpleCA.svg?branch=main)](https://travis-ci.com/520MianXiangDuiXiang520/SimpleCA)

在非对称加密中，为了证明一个公开的密钥属于某个特定的用户，可以由一个双方都信任的机构（CA）为公钥和用户颁发一个证书，该证书中包含一个 CA 对证书内容的签名，如果 A 要与 B 进行加密通信，在通信前 A 会向 B 请求证书，证书中包含 B 的信息，公钥，签名算法， CA 机构， CA 签名等信息，A 拿到证书后使用 CA 机构的公钥对证书内容重新做签名，如果得到的结果与证书中 CA 机构的签名一致，说明证书未被篡改且 B合法拥有该公钥，之后 A 和 B 便可以使用该密钥对进行加密通信。

本系统将实现一个小型的 CA 认证系统，它将以 Web 的形式接受用户的证书申请，管理员在后台验证通过后，将为其自动生成一个证书并通过邮件发送给申请者，同时，用户可以申请撤销自己持有的某个证书，证书被撤销后，将被加入到 CRL 列表中，用户可以访问特定的 API 获取最新的 CRL 列表数据。

##  功能设计

### 认证和授权

用户需要在登录态下进行申请证书，查看证书，撤销证书等操作，因此必须需要认证和鉴权模块，首先认证模块分为登录，注册，认证三部分，登录注册 API  文档如下：

#### 登录 API

认证使用 Token 方案，客户端发送登录请求后，服务端验证用户名和密码，验证通过后使用 UUID 为其生成一个 32 位的 Token并返回，服务端会保存该 Token， 以后每次需要认证的请求客户端都需要携带该 Token, 服务端先检查 Token 是否存在，是否过期，通过后再为其提供服务。

URI： api/ca/login

method: POST

请求格式：

```json
{
    "username": "",
    "password": ""
}
```

响应格式：

```json
{
    "header": {
        "code": 200,
        "msg": "ok"
    },
    "token": ""
}
```

系统所有响应都包含响应头和响应体两部分，响应体可以为空，响应头中两个字段分别表示响应的状态和状态描述，与标准 HTTP 状态码和状态描述保持一致，具体状态码和描述参考：[http 标准状态码文档](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Status)

#### 注册 API

注册需要两次输密码以确认，但此校验会放在前端进行，在请求时只需要传递注册用户名和密码：

URI ： api/ca/register

method: POST

**注册后密码需要使用 SHA256 加盐，保证数据安全性**

#### 认证

登录获得 Token 后，客户端需要把此 Token 存储在 Cookie 中，以后每次访问需要认证的接口时，需要在 Cookie 中携带该 Token，Key Name 为 SESSIONID 认证部分使用认证中间件实现：

```go
func TokenAuth(context *gin.Context) (middleware.UserBase, bool) {
	token, err := context.Cookie("SESSIONID")
	if err != nil {
		return nil, false
	}
	user, ok := dao.GetUserByToken(token)
	if !ok {
		return nil, false
	}
	return user, true
}
```

##### 数据库设计

user 表

| 字段名   | 长度 | 类型   | 备注             |
| -------- | ---- | ------ | ---------------- |
| username | 16   | string | （加密）唯一索引 |
| password | 255  | string | NOT NULL         |
| email    | 255  | string | （加密）         |
| id       |      | int    | 自增，主键索引   |

user_token 表

| 字段名      | 长度 | 类型   | 备注                       |
| ----------- | ---- | ------ | -------------------------- |
| id          |      | int    | 主键索引，自增             |
| token       | 32   | string | 唯一索引                   |
| expire_time |      | int    | 过期时间戳（避免时区转换） |
| user_id     |      | int    | 外键（user.id）外键索引    |



### 证书请求

用户在申请证书时，需要首先有一个密钥对，密钥对可以由用户自己生成，也可以由 CA 机构代为生成，本系统将通过 js 在前端为用户使用 RSA 加密算法生成密钥对，不在后端生成的原因是私钥在从服务端传输给客户端的过程中可能出现泄露和篡改。

自己生成或 CA 代为生成密钥对后，用户自己保存私钥，将公钥和用户信息以 JSON 的形式发送给后端，这些信息会被存库，等管理员审核通过后再颁发证书。

证书请求格式如下：

```json
{
    "public_key" : "",
    "country": "CN",
    "state_or_province": "Beijing",
    "locality": "",
    "organization": "",
    "organizational_unit_name": "",
    "common_name": "",
    "email_address": ""
}
```

证书请求流程图如下：

<img src="https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1605447541792-1605447541780.png" style="zoom:67%;" />

#### 证书请求数据表设计

request 表

| 字段名                   | 长度 | 类型   | 备注                                             |
| :----------------------- | ---- | ------ | ------------------------------------------------ |
| id                       |      | int    | 主键索引，自增                                   |
| user_id                  |      | int    | 外键（user.id) 外键索引                          |
| public_key               |      | text   | 公钥（加密）                                     |
| country                  | 20   | string | 国家（加密）                                     |
| state_or_province        | 255  | string | 州市（加密）                                     |
| locality                 | 255  | string | 地区（加密）                                     |
| organization             | 255  | string | 组织（加密）                                     |
| organizational_unit_name | 255  | string | 部门（加密）                                     |
| common_name              | 255  | string | 姓名（加密）                                     |
| email_address            | 255  | string | 电子邮件（加密）                                 |
| state                    |      | int    | 状态（1：待审核， 2： 审核通过， 3：审核未通过） |

在存储用户申请时提交的资料时，都必须使用服务端单独的私钥对铭感数据进行加密。

#### 申请审核

管理员同样需要登录进入后台，管理员可以单账号多用户登录，表设计与用户表一样，管理员登录后可以查看解密后的用户信息，审查通过后，点击通过会自动生成证书并发送给申请者。

### 证书生成

证书格式遵循 X.509 协议，证书格式如下：

- 版本号(Version Number）：规范的版本号，目前为版本3，值为0x2；
- 序列号（Serial Number）：由CA维护的为它所发的每个证书分配的一的列号，用来追踪和撤销证书。只要拥有签发者信息和序列号，就可以唯一标识一个证书，最大不能过20个字节；
- 签名算法（Signature Algorithm）：数字签名所采用的算法，如：
  - sha256-with-RSA-Encryption
  - ccdsa-with-SHA2S6；
- 颁发者（Issuer）：发证书单位的标识信息，如 ” C=CN，ST=Beijing, L=Beijing, O=org.example.com，CN=ca.org。example.com ”；
- 有效期(Validity): 证书的有效期很，包括起止时间。
- 主体(Subject) : 证书拥有者的标识信息（Distinguished Name），如：" C=CN，ST=Beijing, L=Beijing, CN=person.org.example.com”；
- 主体的公钥信息(SubJect Public Key Info）：所保护的公钥相关的信息：
  - 公钥算法 (Public Key Algorithm）公钥采用的算法；
  - 主体公钥（Subject Unique Identifier）：公钥的内容。
- 颁发者唯一号（Issuer Unique Identifier）：代表颁发者的唯一信息，仅2、3版本支持，可选；
- 主体唯一号（Subject Unique Identifier）：代表拥有证书实体的唯一信息，仅2，3版本支持，可选：
- 扩展（Extensions，可选）: 可选的一些扩展。中可能包括：
  - Subject Key Identifier：实体的秘钥标识符，区分实体的多对秘钥；
  - Basic Constraints：一指明是否属于CA;
  - Authority Key Identifier：证书颁发者的公钥标识符；
  - CRL Distribution Points: 撤销文件的颁发地址；
  - Key Usage：证书的用途或功能信息。

此外，证书的颁发者还需要对证书内容利用自己的私钥添加签名， 以防止别人对证书的内容进行篡改。

#### 数据库设计

Certificat 表：

| 字段名      | 长度 | 类型 | 备注                                      |
| ----------- | ---- | ---- | ----------------------------------------- |
| id          |      | int  | 主键索引，自增                            |
| state       |      | int  | 状态（1 代表在使用中，2代表已撤销或过期） |
| expire_time |      | int  | 过期时间戳                                |
| user_id     |      | int  | 外键（user.id)                            |
| request_id  |      | int  | 外键（request.id)                         |

### 通知

通知模块使用 `gomail` 三方库实现： [github/gomail](https://github.com/go-gomail/gomail)

### 证书查询

可以通过直接上传证书文件查询证书真伪，同时在登录状态下，用户可以查看自己所有的证书

证书真伪查询流程图：

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1605450868687-1605450868671.png)



### 证书撤销

撤销证书请求需要在登录态下进行，请求格式如下：

```json
{
    "certificat_id": 1
}
```

申请通过后，系统需要修改 Certificat  表中对应记录的 state 字段值为 2，并在 crl 表中插入一条数据

#### 数据库设计

crl 表设计 

| 字段名        | 长度 | 类型 | 备注     |
| ------------- | ---- | ---- | -------- |
| certificat_id |      | int  | 证书 ID  |
| time          |      | int  | 加入时间 |

#### crl 列表发布

CRL 列表通过 json 的方式发布。

URI：api/ca/crl

method: GET

响应格式：

```json
{
    "header": {
        "code": 200,
        "msg": "ok"
    },
    "tag": 1,
    "release_time": 1605364857,
    "crl": [
        {
            "certificat_id": 1,
            "expire_time": 1605364857,
        },
        {
            "certificat_id": 2,
            "expire_time": 1605364857,
        },
        {
            "certificat_id": 3,
            "expire_time": 1605364857,
        }
    ]
}
```




## 模块划分和人员分工

| 模块           | 功能                                                         | 人员   |
| -------------- | ------------------------------------------------------------ | ------ |
| 在线密钥对生成 | -----                                                        |    |
| 认证鉴权       | 普通用户登录前后端，<br />普通用户注册前后端<br />鉴权中间件 |    |
| 邮件通知       |                                                              |  |
| 证书请求       | 表单前后端，请求入库                                         |  |
| 请求审核       | 管理员登录，审核前后端                                       |   |
| 证书生成       | 证书生成，数据入库                                           |  |
| 证书查询       | 证书文件上传，证书真伪查询前后端，个人证书查询前后端         |  |
| 证书撤销       | ------                                                       |  |


