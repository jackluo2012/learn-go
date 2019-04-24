##爬虫系统

#### 并发版 流程

* 先要有个爬中引擎 
#### ConcurrentEngine 爬虫引擎
   > 1. 调度器 Scheduler
   > 2. 工作的个数 WorkerCount
   > 3. 执行的方法 Run

- Scheduler 调试器是个接口
    > 1. 有个直接向里面接提交的方法 Submit
    >>>   a. Sumit 里也是单独开的 goroutine 来处理
    
    > 2. 配置一个主 master ConfigureMasterWorkerChan
    
- WorkerCount       
    
    > 1. worker 的个数

- Run (request) 执行的要方法
    > 1.  初始化 in 和 out 通道
    > 2.  初始化Master Worker
    > 3.  新建 10 个 Worker ,传入in 和 out 两个通道
    > 4.  将Run请求的 Request 地址放入通道中
    > 5.  新建一个 for 死循环,接收的 ParseResult
    > 6.  将 收到的 ParseResult 中的新的 Request

- createWorker 创建新的 worker 进行工作,传入 Request 的通道和 ParseResult 通道
  > 1. 单独开一个 goroutine 来进行处理
  > 2. 一直 接收 Request 通道 
  > 3. 将通道的结果 ParseResult 传给 out 通道

----
#### 队列版 流程   
#### ConcurrentEngine 爬虫引擎
   > 1. 调度器 Scheduler
   > 2. 工作的个数 WorkerCount
   > 3. 执行的方法 Run
   
- ConcurrentEngine -> Scheduler 调试器是个接口
    > 1. 往调度器里提交任务的方法 Submit
    >>>   a. Sumit 里也是单独开的 goroutine 来处理
    > 2. 是否准备好工作 WorkerReady
    >>> b.直接执行的
    > 3. 有个启动调度器的方法 Run      

- Scheduler -> Submit
    > 1. 直接 将 Request 放入 requestChan 通道中
- Scheduler -> WorkerReady
    > 1. 将 worker 通道 放入 workerChan 的通道中    
- Scheduler -> Run  启动队列
    > 1. 初始化队列的 requestChan 和 workerChan 通道
    > 2. 申明 两个队列 一个放 requestQ 队列 一个放 workerQ 队列
    > 3. for 循环 不停的接收 调度里的 requestChan 和 workerChan (请求和 worker)
   
- ConcurrentEngine -> WorkerCount           
    > 1. worker 的个数
    
- ConcurrentEngine -> Run  (request) 执行的要方法
    > 1. 初始化解析器
    > 2. 启动 调度器
    > 3. 根据配置 调用 createWorker 创建 worker 的个数
    
 - createWorker(out chan ParseResult, s Scheduler) 传入的是调 解析器的通道和 调度器
    > 1. 初始化解析器
    
    