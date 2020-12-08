# Stack-RPC

Stack-RPC旨在为中国开发者提供通用的分布式服务微服务开发库（比如配置管理、服务发现、熔断降级、路由、服务代理、安全、主从选举等）。基于Stack，开发者可以快速投入自身的业务开发中，只需要极少的学习成本。Stack适用于中小规模的开发场景，她可以轻易在桌面电脑、服务器、容器集群中搭建分布式服务。

## 开发手册

[开发文档](https://stacklabs.cn/docs/stack-rpc/introduce-cn)

[示例](https://github.com/stack-labs/stack-rpc-tutorials)

[插件库](https://github.com/stack-labs/stack-rpc-plugins)

## 交流

<div style="float:left">
<table width="60%">
    <tr>
        <td>公众号</td>
        <td>讨论群</td>
    </tr>
    <tr>
        <td><img alt="微信搜索公众号：StackHQ" src="https://github.com/stack-labs/Notice/raw/master/donation/wx_qrcode.jpg"> </td>
        <td><img alt="微信搜索公众号：MicroHQ，备注来源：“github”" src="https://github.com/stack-labs/Notice/raw/master/donation/wx_group_v1.jpg"> </td>
    </tr>
</table>
</div>

> 讨论群：微信搜索MicroHQ，备注来源：“github”

## 简单易用

启动一个微服务只需要如下代码

```
func main() {
  service := stack.NewService(stack.Name("stack.rpc.greeter"))
  service.Init()
  service.Run()
}
```

我们封装了微服务内在的复杂度，比如服务注册与发现、配置管理等。用户只需要花极小的成本学习如何暴露接口，如何启动服务，剩下的精力完全投放在业务需求的开发上。

## 特性

Stack-RPC既提供轻量的开发库，同时也提供对应高级别的扩展库，为大家带来开箱即用的开发体验。

支持的特性主要有：

- 分布式配置
- 服务注册与发现
- 服务路由
- 远程服务调用
- 负载均衡
- 链路中断与降级
- 分布式锁[todo]
- 主从选举[todo]
- 分布式广播

## 开始使用

我们为一直为大家准备持续开发、更新、愈加丰富的文档与资料：[StackLabs](https://stacklabs.cn/docs/stack-rpc/introduce-cn)

## 与Go-Micro的差异

- 取消Micro工具集，以服务插件形式集成到Stack-RPC中
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

## 鸣谢

- 感谢Go-Micro库，提供优秀的扩展性极强的原始框架，Stack-RPC作为衍生版本，受益颇多，同时Go-Micro的肄业也给Stack-RPC创造了生命
- 感谢Spring-Cloud，作为使用最广泛的开源分布式开发库，我们参考了她许多优秀的设计与文档
- 感谢各位Go-Micro的历史提交者，他们的代码永远运行在大家的内存中
- 感谢各位支持**StackLabs**中国发展的贡献者们
