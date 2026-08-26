# 消息队列系统

纯 Go 标准库实现的消息队列管理平台，支持主题、消息、生产者、消费者、消费组、订阅关系、重试记录与死信队列的完整生命周期管理。

## 运行方式

在 `origin/` 目录下执行：

```bash
go run ./cmd/server
```

服务默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

HTTP 服务访问入口：`http://localhost:8080/`

## 认证

所有 API 请求需在 Header 中携带 Bearer Token：

```
Authorization: Bearer mq-demo-token
```

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/topics | 创建主题 |
| GET | /api/topics | 主题列表（支持 keyword/status 筛选） |
| GET | /api/topics/{id} | 主题详情 |
| PUT | /api/topics/{id} | 更新主题 |
| DELETE | /api/topics/{id} | 删除主题 |
| POST | /api/messages | 发布消息 |
| GET | /api/messages | 消息列表（支持 topic_id/status/partition/keyword 筛选） |
| GET | /api/messages/{id} | 消息详情 |
| PUT | /api/messages/{id} | 更新消息 |
| DELETE | /api/messages/{id} | 删除消息 |
| POST | /api/messages/{id}/deliver | 投递消息 |
| POST | /api/messages/{id}/ack | 确认消息（可选 group_id 推进 offset） |
| POST | /api/messages/{id}/fail | 标记投递失败（触发重试/死信） |
| POST | /api/messages/batch | 批量发布消息 |
| GET | /api/topics/{topic_id}/messages | 按主题查询消息 |
| POST | /api/producers | 创建生产者 |
| GET | /api/producers | 生产者列表 |
| GET | /api/producers/{id} | 生产者详情 |
| PUT | /api/producers/{id} | 更新生产者 |
| DELETE | /api/producers/{id} | 删除生产者 |
| POST | /api/consumers | 创建消费者 |
| GET | /api/consumers | 消费者列表 |
| GET | /api/consumers/{id} | 消费者详情 |
| PUT | /api/consumers/{id} | 更新消费者 |
| DELETE | /api/consumers/{id} | 删除消费者 |
| POST | /api/consumer-groups | 创建消费组 |
| GET | /api/consumer-groups | 消费组列表 |
| GET | /api/consumer-groups/{id} | 消费组详情 |
| PUT | /api/consumer-groups/{id} | 更新消费组 |
| DELETE | /api/consumer-groups/{id} | 删除消费组 |
| POST | /api/consumer-groups/{id}/advance | 推进 offset（delta） |
| POST | /api/subscriptions | 创建订阅关系 |
| GET | /api/subscriptions | 订阅列表 |
| GET | /api/subscriptions/{id} | 订阅详情 |
| PUT | /api/subscriptions/{id} | 更新订阅 |
| DELETE | /api/subscriptions/{id} | 删除订阅 |
| POST | /api/retries | 创建重试记录 |
| GET | /api/retries | 重试记录列表 |
| GET | /api/retries/{id} | 重试详情 |
| PUT | /api/retries/{id} | 更新重试记录 |
| DELETE | /api/retries/{id} | 删除重试记录 |
| GET | /api/messages/{message_id}/retries | 按消息查重试记录 |
| POST | /api/dead-letters | 创建死信 |
| GET | /api/dead-letters | 死信列表 |
| GET | /api/dead-letters/{id} | 死信详情 |
| PUT | /api/dead-letters/{id} | 更新死信 |
| DELETE | /api/dead-letters/{id} | 删除死信 |
| GET | /api/topics/{topic_id}/dead-letters | 按主题查死信 |
| GET | /api/stats/global | 全局统计 |
| GET | /api/stats/topics | 每主题消息统计 |
| GET | /api/stats/rates | 每主题消费速率/ack 率 |
| GET | /api/stats/dead-letter-topn?n=5 | 死信 TOP N |
| GET | /api/topics/{topic_id}/export | JSON 导出某主题全部数据 |

## 消息状态机

pending → delivered → acknowledged
pending → delivered → dead（重试超限）

非法流转返回 409 Conflict。

## 中间件

- 鉴权中间件（Bearer token）
- 限流中间件（令牌桶，默认 100 req/s）
- 日志中间件
- Recovery 中间件
