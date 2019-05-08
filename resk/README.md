### 红包秒杀系统 

#### apis 用记接口层  (网络适配,交互的代码)
- web 存放 web 接口
- thrift RPC 接口
#### services 应用层 (存放 service 的接口,定义所有的 services 接口,作为唯一的交互入口,所有的业务都封装在 services 内部)

#### core 再来创建核心领域层(存放业务逻辑代码,存放业务核心领域)
- acounts
- envelopes
- users

#### infra 基础设施层 (数据库,缓存,队列,和业务无关的基础代码)

####doc  辅助的包和目录(存放项目文档,脚本相关的内容)

#### brun 存放 main 函数和编译后的二进制代码

#### public 用来存放 css js 模板 和静态文件