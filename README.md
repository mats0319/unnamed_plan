# 未命名计划

[![Go Reference](https://pkg.go.dev/badge/github.com/mats9693/unnamed_plan.svg)](https://pkg.go.dev/github.com/mats9693/unnamed_plan)
[![Go Report Card](https://goreportcard.com/badge/github.com/mats9693/unnamed_plan)](https://goreportcard.com/report/github.com/mats9693/unnamed_plan)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)

一个web项目，用来学习和实践相关知识技能，主要是后端

声明：项目脚本主要运行环境为windows，如有linux环境运行需求，请参考对应脚本或命令，自行执行

## 编程语言/技术框架

- go 1.17
- vue 2，html+ts+less, recommend: node 16

## 工具

- nginx
- postgresql

## 内容

> web相关、微服务相关内容，不使用其他人的代码库（如gin、etcd等），学会一部分，就自己实现一遍  
> 其他内容不限制使用其他人的代码库（如日志库、rpc框架）

- [x] 一个web项目，包含前端界面与后端程序，通过http协议建立联系
- [x] 微服务架构：客户端的请求统一发送到API网关，由网关转发给各业务服务；每个业务服务独立运行，在一个非核心服务故障时，其他服务还能支撑项目整体的基本运转
    - [x] API网关
    - [x] 配置中心
    - [x] 服务注册中心
    - [ ] 调用链追踪
    - [x] 配套的用户界面
- [ ] 性能监控
- [ ] 熔断、限流、降级

## 项目图

> 感谢：[绘图工具](https://excalidraw.com)，[压缩工具](https://tinypng.com)

结构图：

![结构图](/doc/img/structure.png)

启动图：

![启动图](/doc/img/start.png)

请求执行过程：

![请求执行过程](/doc/img/execute_request.png)
