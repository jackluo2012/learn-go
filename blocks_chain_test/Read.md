

#### 创建区块
```shell
go run main.go createblockchain -address jackluo

go run main.go createblockchain -address 1NWjmibWEh712akAMuHoBg2vAcphvkqAes
```
#### 获取余额
```shell
go run main.go getbalance -address jackluo

go run main.go getbalance -address 1NWjmibWEh712akAMuHoBg2vAcphvkqAes


```
#### 转帐
```shell
go run main.go send -from jackluo -to tom -amount 30

go run main.go send -from jackluo -to tom -amount 10

go run main.go send -from 1JeLZ9dn5k5FrkQFWHrk2GvCNbCY9M2VSZ -to 1JcENpkKKxLhVw2rfiZz9aqSBhv8j5mw7R -amount 6

go run main.go send -from 1NWjmibWEh712akAMuHoBg2vAcphvkqAes -to 1GekiJaqVM98KJMaZqcVMDBMLMjPqBPyd8 -amount 6
```

#### 创建钱包
```shell 
go run main.go createWallet

```
#### 加入区块链中
```shell
go run main.go createblockchain -address 1JeLZ9dn5k5FrkQFWHrk2GvCNbCY9M2VSZ

```

#### 打印出 区块链
```shell

go run main.go printchain

```