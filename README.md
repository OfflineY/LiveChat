
# LiveChat
Lightweight chat application implemented with golang.

---
⚠️ v3.x is still developing.

## New features（developing）

- New ui（tailwindcss+react）
- New database（mongodb）

## Setup

```shell
git clone https://github.com/OfflineY/LiveChat.git
cd LiveChat

go get ...(Install the golang package)

go run main.go
```

## Api
<< GET/POST

RETURN:
```
"data": "", <- main data
"msg": "",  <- err or nil
```

### Search Group
[GET] api/groups/**search**

[PARAMS] `?name=[group_name]&id=[groups_id]`

take at **least one** of the two key values.

return `data:[...]` and `msg:nil/err`

### Create Group
[POST] api/groups/**create**

[JSON] `{"name":[new group name]}`

return `group_id` and `group_name`

### Group Historical messages
[GET] api/groups/**[group_id]**/messages

return `data:[...]` and `msg:nil/err`

### Group Info
[GET] api/groups/**[group_id]**/info
