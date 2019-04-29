##爬虫系统

#### 分布式版 流程

* 先要有个爬中引擎 
#### ConcurrentEngine 爬虫引擎
```bash
#启动 elasticsearch 
docker run -d -p 9200:9200 elasticsearch

cd /Volumes/Data/Work/golang/src/gopcp.v2/chapter7
#启动存储服务
go run crawler_distributed/persist/server/itemserver.go -port 1234

# 启动 worker 服务
go run crawler_distributed/worker/server/worker.go -port 9000
go run crawler_distributed/worker/server/worker.go -port 9001
go run crawler_distributed/worker/server/worker.go -port 9002


#启动服务
go run crawler_distributed/main.go -itemsaver_host=":1234" -worker_hosts=":9000,:9001,:9002"

```
    
   