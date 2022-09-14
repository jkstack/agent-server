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

修改命令行交互