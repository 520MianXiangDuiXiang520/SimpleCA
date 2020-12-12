#   SimpleCA

[![Build Status](https://travis-ci.com/520MianXiangDuiXiang520/SimpleCA.svg?branch=main)](https://travis-ci.com/520MianXiangDuiXiang520/SimpleCA)


这是软件大型实验周的课设作品，用来实现一个简单的 CA 系统，它包含以下功能：

* 证书生成：用户提供 Certificate Signing Request （CSR）和 公钥后，系统会自动为用户生成证书并通过邮箱发放，支持 用于SSL 和代码签名的两类证书。
* 证书吊销：用户发起证书吊销请求后，系统会为其更新 Certificate Revocation List （CRL），但考虑系统负荷，用户的吊销请求会被暂时记录，然后以天为单位更新，客户端可以通过 CRL Distribution Point（CRL 分发点：本系统分发点为：[CRL Distribution Point](http://39.106.168.39/crl.crl)）获取最新的 CRL 列表。
* 证书审核：用户的证书请求需要由管理员审核后才能颁发，因此我们写了一个简单的管理员审核功能。
* 密钥对生成：基于 Vue 的 jsrsasign 插件，我们支持客户端本地生成 RSA 公私密钥对。
* 用户认证：包括登录登出注册鉴权等常用操作，使用 Token 实现。
* 代码签名：除了 CA 系统自生功能外，本项目还提供了一个简单的代码签名工具，使用 Java 编写，用来验证项目成果。

## 最终成果

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606726448325-1606726448320.png)

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606575852677-1606575852669.png)

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606575905576-1606575905569.png)

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606576339973-1606576339969.png)

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606576362986-1606576362983.png)

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606635923562-1606635923554.png)

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606660019894-1606660019886.png)

## 接口文档

[API 文档](https://documenter.getpostman.com/view/9355097/TVev3jei#ec21108d-d16f-49a6-83c1-4aee91e42ceb)

为了统一前后端请求，所有请求皆使用 POST

### CSR 文件上传

接口：`/api/ca/file`

文件标识：`CSR_FILE`

返回：该接口会返回通过 CSR 解析出的用户公钥及个人信息

```json
{"header":{"code":200,"msg":"ok"},"country":"CN","province":"ShanXi","locality":"TaiYuan","organization":"中北大学","common_name":"junebao.top","email_address":"","organizational_unit":"IT","public_key":"-----BEGIN RSA PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArj/g1UgvJAEnXxlDFHU6\naMwaaZ7BP4zXXWLnfw5kiVTjMPL4kS21do82bygUg4tCPM3pnDt5BpibGwsAZhQ8\nNH527z0Is2yHUT8S/RWT7t7AAJ06NdsEzdyaKAzVHa3xfq6zjVHc11nn0eLB0M0G\nahteIZebHjNhMX3dyVbvUx9e0iAjPDxbCvficbBDhQ1fzZbUoxmS175ENDuoNRY1\nory8+fFAnRhwTJn12mhB/U+QaHiBIfzhC7exMffcUJYK8WYqt0W0+3oS47gyiCt7\nlEyhzQ9UoJVA7O5zsmB39xPLXHIsRLkuAZBIB9YLibZOse5gRCVgd6OYJTyhHbGl\nRwIDAQAB\n-----END RSA PUBLIC KEY-----\n"}
```



##  功能设计

### 系统流程

<img src="https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606738791416-1606738791402.png" style="zoom:67%;" />

### 数据库设计

```sql
CREATE TABLE `ca_requests` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL COMMENT '字段创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `user_id` int NOT NULL COMMENT '申请证书的用户ID',
  `state` int unsigned NOT NULL COMMENT '证书状态（1：待审核， 2： 审核通过， 3：审核未通过）',
  `public_key` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '公钥',
  `country` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '国家',
  `province` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '州市',
  `locality` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '地区',
  `organization` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '组织',
  `organization_unit_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '部门',
  `common_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '姓名',
  `email_address` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '邮箱',
  PRIMARY KEY (`id`),
  KEY `idx_ca_requests_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8;
```

```sql
CREATE TABLE `certificates` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `user_id` int NOT NULL COMMENT '证书拥有者ID',
  `state` int unsigned NOT NULL COMMENT '状态（1 代表在使用中，2代表已撤销或过期）',
  `request_id` int NOT NULL COMMENT '证书请求ID',
  `expire_time` bigint NOT NULL COMMENT '过期时间戳',
  PRIMARY KEY (`id`),
  KEY `idx_certificates_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8;
```

```sql
CREATE TABLE `crls` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `certificate_id` int NOT NULL COMMENT '证书ID',
  `input_time` bigint NOT NULL COMMENT '加入时间戳',
  PRIMARY KEY (`id`),
  KEY `idx_crls_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COMMENT='证书吊销列表';
```

```sql
CREATE TABLE `user_tokens` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  `user_id` int NOT NULL,
  `token` varchar(64) NOT NULL,
  `expire_time` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_tokens_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=81 DEFAULT CHARSET=utf8;
```

```sql
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `username` varchar(16) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `authority` int DEFAULT NULL COMMENT '权限，1表示系统管理员',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8;
```



## 背景知识

### PKI

### 加密算法

#### 对称加密：3DES

> 3DES（或称为Triple DES）是三重数据加密算法（TDEA，Triple Data Encryption Algorithm）块密码的通称。它相当于是对每个数据块应用三次DES加密算法。由于计算机运算能力的增强，原版DES密码的密钥长度变得容易被暴力破解；3DES即是设计用来提供一种相对简单的方法，即通过增加DES的密钥长度来避免类似的攻击，而不是设计一种全新的块密码算法。

* [DES](http://39.106.168.39/#/detail/30)

#### 非对称加密 RSA

* [RSA](http://39.106.168.39/#/detail/39)

### 哈希算法

#### MD5

#### SHA256

### X.509

> X.509是密码学里公钥证书的格式标准。X509证书己应用在包括TLS/SSL在内的众多
>
> Internet协议里,同时它也用在很多非在线应用场景里,比如电子签名服务。X509证书里含有公钥、身份信息(比如网络主机名,组织的名称或个体名称等)和签名信息(可以是证书签发机构CA的签名,也可以是自签名)。对于一份经由可信的证书签发机构签名或者可以通过其它方式验证的证书，证书的拥有者就可以用证书及相应的私钥来创建安全的通信，对文档进行数字签名.另外除了证书本身功能,Ⅹ509还附带了证书吊销列表和用于从最终对证书进行签名的证书签发机构直到最终可信点为止的证书合法性验证算法。
>
> 在X.509里，组织机构通过发起证书签名请求(CSR)来得到一份签名的证书。首先需要生成一对钥匙对，然后用其中的私钥对CSR进行签名，并安全地保存私钥。CSR进而包含有请求发起者的身份信息、用来对此请求进行验真的的公钥以及所请求证书专有名称。CSR里还可能带有CA要求的其它有关身份证明的信息。然后CA对这个专有名称发布一份证书，并绑定一个公钥. 组织机构可以把受信的根证书分发给所有的成员，这样就可以使用公司的PKI系统了。

#### 证书格式

> 版本号(Version Number）：规范的版本号，目前为版本3，值为0x2；
>
> 序列号（Serial Number）：由CA维护的为它所发的每个证书分配的一的列号，用来追踪和撤销证书。只要拥有签发者信息和序列号，就可以唯一标识一个证书，最大不能过20个字节；
>
> 签名算法（Signature Algorithm）：数字签名所采用的算法，如：sha256-with-RSA-Encryptionccdsa-with-SHA2S6；
>
> 颁发者（Issuer）：发证书单位的标识信息，如 ” C=CN，ST=Beijing, L=Beijing, O=org.example.com，CN=ca.org。example.com ”；
>
> 有效期(Validity): 证书的有效期很，包括起止时间。
>
> 主体(Subject) : 证书拥有者的标识信息（Distinguished Name），如：" C=CN，ST=Beijing, L=Beijing,CN=person.org.example.com”；
>
> 主体的公钥信息(SubJect Public Key Info）：所保护的公钥相关的信息：
>
> 公钥算法 (Public Key Algorithm）公钥采用的算法；
>
> 主体公钥（Subject Unique Identifier）：公钥的内容。
>
> 颁发者唯一号（Issuer Unique Identifier）：代表颁发者的唯一信息，仅2、3版本支持，可选；
>
> 主体唯一号（Subject Unique Identifier）：代表拥有证书实体的唯一信息，仅2，3版本支持，可选：
>
> 扩展（Extensions，可选）: 可选的一些扩展。中可能包括：Subject Key Identifier：实体的秘钥标识符，区分实体的多对秘钥；
>
> Basic Constraints：一指明是否属于CA;Authority Key Identifier：证书颁发者的公钥标识符；
>
> CRL Distribution Points: 撤销文件的颁发地址；
>
> Key Usage：证书的用途或功能信息。

#### 证书吊销

在公共密钥基础设施（PKI）中， CSR 是证书申请者发送给 CA 机构用于为其签发数字证书的文件，其常见格式为 [PKCS](https://en.wikipedia.org/wiki/PKCS)＃10 和 SPKAC , 一般情况下，他会包含一个用户公钥和必要的用户信息，X.509 标准中，列举了如下列：

![](https://cdn.jsdelivr.net/gh/520MianXiangDuiXiang520/cdn/img/1606724475751-1606724475704.png)

申请者在生成证书之前，应该先生成一个密钥对，自己妥善保管私钥，对于 CSR 中的信息，应该使用私钥签名以防止被篡改。因此，一个标准 CSR 包含下面三部分信息：

1. 认证请求信息：包含主要的申请人信息和公钥
2. 签名算法标识：
3. CSR 签名：申请人应该使用自己的私钥对 CSR 信息进行签名，以确保其不被篡改。

本系统允许申请人通过填写表单快速提交个人信息和公钥，同时也允许申请人直接上传 CSR 文件。

本地通过 openSSL 生成 CSR 的方法：

```sh
# 生成随机私钥
openssl genrsa -des3 -out server.key 2048
# 使用私钥签发 CSR
openssl req -new -nodes -keyout server.key -out server.csr
```




## 模块划分和人员分工

| 模块                       | 内容                                            |      |
| -------------- | ------------------------------------------------------------ | ------ |
| 证书请求                   | 1. CSR文件上传<br />2. CSR表单填写              |      |
| 公钥提交<br />用户证书查询 |                                                 |      |
| 证书生成                   | 1. 待审核列表<br />2. 颁发证书<br />3. 驳回请求 |      |
| 登录注册鉴权               | 1. 登录<br />2. 注册<br />3. 登出<br /> |  |
| 证书吊销 | 1. 证书吊销请求<br />2. CRL 分发               |   |
| 代码签名   |                                            |  |