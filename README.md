# Stack-RPC

Golang微服务开发框架，基于Go-Micro（1.18）修改。

## 开始使用

安装依赖，假设已经安装好Golang开发环境。

安装protoc环境与Stack编译插件

```bash
$ go get github.com/golang/protobuf/protoc-gen-go@v1.3.2
$ go get github.com/stack-labs/stack-rpc/util/protoc-gen-stack 
```

## 与Go-Micro的差异

- 取消Micro工具集，以插件形式集成到Stack-RPC中
  - Web控制台插件
  - 网关插件
  - 注册中心插件
  - 配置中心插件
- 多类型服务同时部署
  - 支持RPC与HTTP同时暴露
- 增强日志特性
  - 支持动态修改日志级别
  - 支持日志自定义目录存储
  - 支持按级别存储不同文件
  - 支持按周期压缩日志文件
  - 支持按大小压缩日志文件
- 增强配置特性
  - 增加默认配置命名空间
  - 定义默认配置文件存储目录
  - 支持Apollo配置中心
- 删除不要的组件
  - Cloudflare

## 维护者团队

国内一线产商多年经验研发人员

成员列表(todo)：

