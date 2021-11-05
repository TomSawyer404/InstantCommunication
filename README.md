# Golang命令行即时通讯项目

- 参考视频：[Aceld-8小时转职Golang工程师](https://www.bilibili.com/video/BV1gf4y1r79E?p=37)

## 食用指南

- Server端编译：使用`make server`编译服务器程序；
- Client端编译：使用`make client`编译客户端程序；
- 使用`./target/server`运行服务端，你可以使用`nc 127.0.0.1 8888`来连接服务器，也可以使用编译好的客户端连接服务器；
- 如果使用`nc`连接服务器，使用`rename|[Your_name]`来修改你的名字，使用`to|some_one|[Your Message]`来给某人私聊；
- 我暂时还没搞明白为什么，如果你用`nc`连接服务器，你可以随意发消息；但用编译好的client，你发的消息不能有空格；
- 服务端设置了超时下线功能，大概是100s好像；
