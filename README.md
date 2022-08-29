<div align="center">
  
# 📫OnlineChat
✨ 使用 golang 实现的简单聊天示例 ✨
  
</div>


## 🎈 使用

下载 version 里面的 exe 运行即可。

> [其他方法] 克隆此仓库，安装对应的库即可，然后：
>
> ```shell
> go run main.go
> ```

如果没有配置文件会自动初始化，初始化完成之后重启应用即可。

具体设置配置文件内写的很清楚。

## 🚁 服务端命令

> 未完工

- Help/help 帮助
- UserList/userList 获取用户在线列表

## 📝 TODO

- [ ] 【未来】服务端命令式控制（doing）
- [x] 【未来】增加版本更新检测 （done）
- [ ] 【未来】用户在线列表（doing）
- [ ] 【未来】使用gin和react构建web端（doing）

## 🚀 从 0.1 到 1.0

- 没有任何影响，可以直接替换 exe
- 增加了版本检测
- 优化了代码排版与注释
- 修复：回复出现乱码

## 🎡 技术栈

- Golang v1.18
- github.com/gorilla/websocket v1.5.0
- gopkg.in/ini.v1 v1.67.0
