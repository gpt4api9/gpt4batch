# gpt4batch

简介：为[https://gpt4api.shop](https://gpt4api.shop)提供批量调用gpt-4接口服务脚本。

# 项目文件说明

- `gpt4batch/cmd/authsvc`: 获取用户批量调用的access_token.
- `gpt4batch/cmd/batchsvc`: 批量调用gpt-4接口服务.
- `gpt4batch/test/general`: 生成测试文件数据脚本.

# 批量脚本数据格式。
存放于`gpt4batch/test/example`

# 批量脚本字段说明
- id: 批量脚本唯一标识。
  - 解释：用于记录和回溯下载文件使用
- asks: 问题列表。
  - id: 问题唯一标识。
    - 解释：用于记录和回溯下载文件使用，对应那组对话
  - content: 提问内容.
  - image: 图片地址.
  - files: 文件地址.
- answers: 答案列表。
- iErr: 错误信息 <如果为空，则以为回到成功，反之错误>
  - code: 状态码
  - message: 错误信息
- extra: 额外扩展字段存储其他信息

# 请求路径地址

#### 普通版URL

- 文件上传
  - https://beta.gpt4api.plus/standard/uploaded
- All-Tools
  - https://beta.gpt4api.plus/standard/all-tools
- Gizmo
  - https://beta.gpt4api.plus/standard/gizmos

#### 并发版URL

- 文件上传
  - https://beta.gpt4api.plus/concurrent/uploaded
- All-Tools
  - https://beta.gpt4api.plus/concurrent/all-tools
- Gizmo
  - https://beta.gpt4api.plus/concurrent/gizmos

# 获取认证Access_Token

```shell
gpt4batch authsvc --help
Retrieve the GPT4Batch token for authentication. and only one valid token will be retained. please use this command with caution.

Usage:
  gpt4batch authsvc [flags]

Flags:
  -e, --email string      输入账号邮箱.如果没注册可访问官网.https://gpt4api.shop.
  -h, --help              help for authsvc
  -p, --password string   输入账号密码,如果没注册可访问官网.https://gpt4api.shop.
  -t, --ttl int           设置AccessToken过期时间，默认是60天. (default 86400)
  -u, --url string        设置获取Ak调用服务地址. (default "https://beta.gpt4api.shop/console/access_token")
```

# 开启批量调用

```shell
gpt4batch batchsvc --help
please proceed with caution when enabling batch calls to GP Ts scripts.

Usage:
  gpt4batch batchsvc [flags]

Flags:
  -d, --download-dir string             下载文件夹名称.如果未设置会存在当前文件夹目录.
  -p, --download-prefix string          设置文件下载前缀，防止下载文件名冲突覆盖. (default "GPT4API")
  -e, --enable-download                 是否开启文件下载. (default true)
  -n, --enable_nsq                      是否开启NSQ消息队列.
  -f, --fix                             是否开启续跑模式.
  -z, --gizmo-id string                 设置GPTs gizmo id的名称.
  -g, --goroutine int                   设置最大协程数量. (default 60)
  -h, --help                            help for batchsvc
  -s, --history_and_training_disabled   是否开启历史对话历史记录，默认是关闭的. (default true)
  -i, --in string                       输入文件路径，数据格式按照规定格式定义. (default "example.jsonl")
  -m, --model string                    设置调用GPTs的模型. (default "gpt-4-gizmo")
  -o, --out string                      输出文件路径，GPTs数据跑完存储数据的文件路径. (default "out.jsonl")
  -q, --qps int                         设置QPS并发量. (default 1)
  -r, --rdb                             是否开启RDB文件缓存持久化策略. (default true)
  -v, --rdb_interval int                RDB缓存时间间隔，默认是60分钟 (default 60)
  -l, --upload_url string               设置批量调用服务地址.普通版：standard 并发版：concurrent (default "https://beta.gpt4api.plus/standard/uploaded")
  -u, --url string                      设置批量调用服务地址.普通版：standard 并发版：concurrent (default "https://beta.gpt4api.plus/standard/all-tools")
```
