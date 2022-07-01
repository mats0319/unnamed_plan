# 未命名计划

[![Go Reference](https://pkg.go.dev/badge/github.com/mats9693/unnamed_plan.svg)](https://pkg.go.dev/github.com/mats9693/unnamed_plan)
[![Go Report Card](https://goreportcard.com/badge/github.com/mats9693/unnamed_plan)](https://goreportcard.com/report/github.com/mats9693/unnamed_plan)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)

一个web项目，用来学习和实践相关知识技能

## 编程语言/技术框架

- go 1.17
- vue 2，html+ts+less, recommend: node 16

## 工具

- nginx
- postgresql

## 内容

> 尽量少使用其他人的代码库，学会一部分，就自己实现一遍

- [x] 一个web项目，包含前端界面与后端程序，通过http协议建立联系
- [x] 微服务架构，界面的请求统一发送到API网关、网关再转发到给各业务服务；每个业务服务独立运行，在一个非核心服务故障时，其他服务还能支撑项目整体的基本运转
  - [x] API网关
  - [x] 配置中心
  - [x] 服务注册中心
  - [ ] 调用链追踪
- [ ] 性能监控
- [ ] 熔断、限流、降级
