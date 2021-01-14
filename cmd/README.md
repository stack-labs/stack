# CMD

CMD is the init entrance for Stack APP

Most of its jobs are preparing configurations, setup components like client, server, and init the server

## Config Load Order

- Config Option(Hard code)
    - stack.Config 
- Config File: 
    - stack.yml. If there is no stack.yml, config init will skip it. All options please see: [stack.yml](../cmd/stack.yml)
    - yml files stored in same dir with the **stack.yml**, and you should set included for them in **stack.yml**.
    - we now suggest setting configs in yml, even though we support most type of config files eg. json, toml.
    - other config files setting at cmd by *--config*. eg: **go run main.go --config=path/to/a.json**.
- Environment Variables
    - we now support native Stack env for apps only, which means we can't inject user custom env into their options. so people should read env by themselves.
- CMD 
    