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

# 版本

- v1.0.0: [20231226]
  - 添加多模态文件上传，对话调用，文件下载。
  - 增加nsq本地存储会话。
  - 支持GPTs市场接口服务调用。