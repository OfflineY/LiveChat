<div style="text-align: center">

# LiveChat
使用 golang 实现的轻量级聊天应用

</div>

> v3.x 目前频繁更新中，敬请期待

## New features

- New ui（tailwindcss+react）
- New database（mongodb+...）

## Use

```shell
git clone https://github.com/OfflineY/LiveChat.git
go run main.go
```

## Service Api

api 返回数据格式基本如下
```
{
  "data": "",   <- main data
  "msg": ""     <- err or nil
  ...
}
```

#### Search Group
>[GET] api/groups/**search**
>
>[PARAMS] `?name=[group_name]&id=[groups_id]`

take at **least one** of the two key values.

return `data:[...]` and `msg:nil/err`

#### Create Group
>[POST] api/groups/**create**
> 
>[JSON] `{"name":[new group name]}`

return `group_id` and `group_name`

#### Group Historical messages
>[GET] api/groups/**[group_id]**/messages

return `data:[...]` and `msg:nil/err`

#### Group Info
>[GET] api/groups/**[group_id]**/info
