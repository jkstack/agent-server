# agent-server

[![agent](https://github.com/jkstack/agent-server/actions/workflows/build.yml/badge.svg)](https://github.com/jkstack/agent-server/actions/workflows/build.yml)
[![go-mod](https://img.shields.io/github/go-mod/go-version/jkstack/agent-server)](https://github.com/jkstack/agent-server)
[![license](https://img.shields.io/github/license/jkstack/agent-server)](https://www.gnu.org/licenses/agpl-3.0.txt)

jkstack agent统一管理服务，与同类产品相比有以下优势：

1. 支持大规模集群
2. 支持agent配置的动态更新和统一下发
3. 支持agent监控数据的统一收集以及grafana的展示
4. 提供restful接口方便二次开发

## 功能描述

1. 支持已连接agent的列表和基本信息获取
2. 支持服务器端状态获取
3. 主机监控信息采集，需依赖[metrics-agent](https://github.com/jkstack/metrics-agent)
   - 采集任务的批量启动/停止
   - 采集节点的动态更改任务配置
4. 支持agent节点状态监控和grafana展示

## 快速部署

服务端程序推荐使用`linux`系统进行部署

1. 根据当前操作系统下载`deb`或`rpm`安装包，[下载地址](https://github.com/jkstack/agent-server/releases/latest)
2. 使用`rpm`或`dpkg`命令安装该软件包，程序将被安装到`/opt/agent-server`目录下
3. 按需修改配置文件，配置文件将被安装在`/opt/agent-server/conf/server.conf`目录下
4. 使用以下命令启动服务器端程序

    /opt/agent-server/bin/agent-server -action start
5. 检查当前服务启动状态

    curl http://<服务端IP>:<端口号(默认13081)>/api/info/server

## restful接口

* 在部署完成后可通过`http://<服务端IP>:<端口号(默认13081)>/docs/index.html`插件接口文档
* 也可下载`http://<服务端IP>:<端口号(默认13081)>/docs/doc.json`文件后导入到apifox或postman进行调试