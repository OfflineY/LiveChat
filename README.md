
# LiveChat
Lightweight chat application implemented with golang.

> **Warning**
> 
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
>
> return example:
> ```
> "data": "", <- main data
> "msg": "",  <- err or nil
> ```

### Search Group
[GET] api/groups/**search**

[PARAMS] `?name=[group_name]&id=[groups_id]`

> **Note**
> 
> take at **least one** of the two key values.
> 
> return `data:[...]` and `msg:nil/err`


### Create Group
[POST] api/groups/**create**

[JSON] `{"name":[new group name]}`

> **Note**
>
> return `group_id` and `group_name`

### Group Historical messages
[GET] api/groups/**[group_id]**/messages

> **Note**
>
> return `data:[...]` and `msg:nil/err`

### Group Info
[GET] api/groups/**[group_id]**/info
