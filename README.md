<div style="text-align: center">

# LiveChat
使用 golang 实现的轻量级聊天应用

</div>

> v3.x 目前频繁更新中，敬请期待

## New features

- New ui（tailwindcss）
- New database（mongodb）

## Use

```shell
git clone https://github.com/OfflineY/LiveChat.git
go run main.go
```

## Service Api

**前端服务api部分，返回数据格式基本如下**
```
{
  "data": "",   <- main data
  "msg": ""     <- err or nil
  ...
}
```

### Search Group
[GET] api/groups/**search**

[PARAMS]

| key  | value         |
|------|---------------|
| name | *groups name* |
| id   | *groups id*   |

> Take at **least one** of the two key values.


### Create Group
[POST] api/groups/**create**

[JSON] `{"name":"new group name"}`

### Group Historical messages
[GET] api/groups/**[group_id]**/messages

### Group Info
[GET] api/groups/**[group_id]**/info