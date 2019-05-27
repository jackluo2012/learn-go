

#### 创建区块
```shell
go run main.go createblockchain -address jackluo
```
#### 获取余额
```shell
go run main.go getbalance -address jackluo
```
#### 转帐
```shell
go run main.go send -from jackluo -to tom -amount 30
go run main.go send -from jackluo -to tom -amount 10
```