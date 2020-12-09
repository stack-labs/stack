# stackctl

## install

- release binary
    - [releases](https://github.com/stack-labs/stack-rpc/releases)
- go get
    - `go get -u github.com/stack-labs/stack-rpc/util/istioctl`

## example

- 创建`service`

```shell script
stackctl new github.com/stack-labs/example --alias example -type service --gopath
```

- 创建`api`
  
```shell script
stackctl new github.com/stack-labs/example-api --alias example -type api --gopath
```

- 修改 FIXME 内容
    - 替换`example-api`中`path/to/service/proto/example`为`github.com/stack-labs/example/proto/example`
- 本地环境用`go.mod`用`replace`添加`github.com/stack-labs/example`包

```shell script
module github.com/stack-labs/example-api

go 1.14

replace github.com/stack-labs/example v1.0.0 => ../example

require (
	github.com/stack-labs/example v1.0.0
)
```
