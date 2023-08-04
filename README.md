
# LiveChat
Lightweight chat application implemented with golang.

> **Warning**
> v3.x is still developing.

## New features

- New ui（tailwindcss+react）
- New database（mongodb）

## Getting started

```shell
git clone https://github.com/OfflineY/LiveChat.git
cd LiveChat
go mod download
go run main.go
```

## Api

> **Note**
>[main.go](main.go)
> return example:`{"data": "[main data]", "msg": "[err or nil]"}`

---

### Group 对群组的操作

#### Search Group
[GET] api/groups/**search**

[PARAMS] `?name=[group_name]&id=[groups_id]`

> **Note**
> take at **least one** of the two key values.
> return `data:[...]` and `msg:nil/err`


#### Create Group
[POST] api/groups/**create**

[JSON] `{"name":"[group_name]"}`

> **Note**
> return `"data":{"group_id": "...", "group_status": "<T/F>", "group_name": "..."}` and `msg:nil/err`

#### Group Historical messages
[GET] api/groups/**[group_id]**/messages

> **Note**
> return `"data":[{"_id":"...","group_id":"...","group_name":"...","msg":"...","msg_type":"[text/image]","send_time":"...","url":"...","user_name":"..."}...]` and `msg:nil/err`

#### Group Info
[GET] api/groups/**[group_id]**/info

> **Note**
> return `"data":"[...]"` and `"msg":nil/err`

---

### User 用户操作
#### User Login

[POST] api/user/login

[JSON] `{"user_name": "[user_name]", "password": "[password]"}`

#### User Register

[POST] api/user/register

[JSON] `{"user_name": "[user_name]", "password": "[password]"}`



