# CHANGELOG

## 1.0.0

1. 实现基本功能，支持[metrics-agent](https://github.com/jkstack/metrics-agent)和[example-agent](https://github.com/jkstack/example-agent)
2. 支持连接的agent列表查询
3. 支持swagger文档预览，地址: /docs/index.html

## 1.0.1

1. 修正默认配置文件无法启动问题
2. 实现deb和rpm打包

## 1.0.2

1. 增加/metrics/dynamic/temps接口
2. /metrics/dynamic接口中增加传感器温度数据
3. 增加传感器温度上报逻辑
4. 修正未配置kafka地址时的崩溃问题

## 1.0.3

1. 修改命令行交互
2. 修正rpm包升级时卸载服务的问题

## 1.0.4

1. 修改/metrics/dynamic下的进程列表返回内容，忽略一些空值
2. 修正agent升级后agent版本号和go版本号的埋点数据问题
3. 修改打包程序增加epoch

## 1.1.0

1. API接口返回数据结构中增加extime字段
2. 实现exec和file模块相关功能
3. go版本升级到1.19.2
4. 升级第三方库版本

## 1.1.1

1. 实现script模块相关功能
2. 修正/exec/run启动失败时不报错的问题
3. 修改命令行交互方式
4. 增加manifest.yaml配置项描述文件
5. go版本升级到1.19.3
6. 升级第三方库版本

## 1.1.2

修改打包脚本输出文件名

## 1.1.3

1. 修改启动、停止、重启、注册、反注册系统服务失败时的状态码
2. 修正服务无法重启的问题
3. metrics数据上报支持json格式，修改上报数据格式定义
4. 修正代码中的golint问题

## 1.1.4

1. 修改manifest.yaml文件中id字段类型
2. /agents和/agents/info接口增加is_busy返回值
3. /info/server接口中增加id返回值

## 1.1.5

1. 修正metrics-agent上报kafka连接问题
2. 修正打包脚本中没有正确更新文档的问题
3. go版本升级到1.19.4

## 1.2.0

1. 支持arm架构
2. 修正文档中的描述问题
3. go版本升级到1.19.5

## 1.3.0

1. 新增/agents/{id}/logs和/agents/{id}/log/download接口用于下载Agent执行日志
2. 优化swagger文档内容
3. 升级第三方库版本
4. go版本升级到1.20

## 1.3.1

1. 增加rpa模块的gPRC接口定义
2. 修改manifest.yaml中metrics.kafka.format字段类型

## 1.3.2

1. metrics: 新增nameserver相关字段，用于获取主机上的dns配置
2. go版本升级到1.20.1

## 1.3.3

调整打包脚本，新增oss上传逻辑

## 1.3.4

1. metrics: `/metrics/{id}/dynamic/usage`接口新增CPU的load1、load5、load15数据
2. metrics: `/metrics/{id}/dynamic/usage`接口新增磁盘的read_per_second、write_per_second、iops_in_progress数据

## 1.4.0

1. rpa:适配精鲲自研RPA-Agent
2. 新增:agent断开连接时的处理逻辑
3. 去除:manifest.yaml描述文件
4. 其他:go版本升级到1.20.2

## 1.4.1

1. 新增:grpc_listen配置项字段检查逻辑
2. 新增:支持rpa-agent的元素拾取功能
3. 修正:rpa-agent断开时运行状态没有释放的问题
4. 其他:go版本升级到1.20.10

## TODO

1. 支持IPMI agent
2. 升级第三方库