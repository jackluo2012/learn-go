

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

#### V12 测试
```shell
    //创建钱包
    go run main.go createWallet
    // 加入到区块链中
    go run main.go createblockchain -address 1MiSmJgomsukJ6oPSLeDYgsFvSxcWnmwe1
    //查询余额
    go run main.go getbalance -address 1MiSmJgomsukJ6oPSLeDYgsFvSxcWnmwe1
    // 创建钱包
    go run main.go createWallet
    // 转帐
    go run main.go send -from 1MiSmJgomsukJ6oPSLeDYgsFvSxcWnmwe1 -to 1GmCgeaY7MaMGqCSb9ZQJARPkHGWPV81Pt -amount 4
    //查询 余额
    go run main.go getbalance -address 1MiSmJgomsukJ6oPSLeDYgsFvSxcWnmwe1
```

### 场景
- 1.中心节点创建blockchain
- 2.其他钱包节点连接到中心节点，并下载blockchain
- 3.一个或多个矿工节点连接到中心节点,并下载blockchain
- 4.钱包节点创建一个交易
- 5.矿工节点收到此交易并保存在mempool中
- 6.当mempool有足够交易时，矿工开始挖矿
- 7.当新的block生成时,其将被发送到中心节点
- 8.钱包节点与中心节点进行同步
- 9.钱包节点的拥有者将检查支付是否成功

### version 消息


### 场景进行演练
```shell
// 设置终端 端口
export NODE_ID=3000

go run main.go createWallet

go run main.go createblockchain -address 1zhpLAhUbQ7MTJ72iyfoVLh62BnYn3Njz

1JUXiiCPZsdcJLVWp6fpCwgi8bXPJLT5HK

124Qioe4NHSbzJHVZFTD3NFPMKCFC2dXAE

1DQcuNCCRBDfwxH3MDkoGyrjNkUbhKEXSS

go run main.go send -from 1zhpLAhUbQ7MTJ72iyfoVLh62BnYn3Njz -to 1JUXiiCPZsdcJLVWp6fpCwgi8bXPJLT5HK -amount 10 -mine

go run main.go send -from 1zhpLAhUbQ7MTJ72iyfoVLh62BnYn3Njz -to 124Qioe4NHSbzJHVZFTD3NFPMKCFC2dXAE -amount 10 -mine


go run main.go getbalance -address 1zhpLAhUbQ7MTJ72iyfoVLh62BnYn3Njz
go run main.go getbalance -address 1JUXiiCPZsdcJLVWp6fpCwgi8bXPJLT5HK
go run main.go getbalance -address 124Qioe4NHSbzJHVZFTD3NFPMKCFC2dXAE


go run main.go send -from 1JUXiiCPZsdcJLVWp6fpCwgi8bXPJLT5HK -to 124Qioe4NHSbzJHVZFTD3NFPMKCFC2dXAE -amount 1

go run main.go send -from 124Qioe4NHSbzJHVZFTD3NFPMKCFC2dXAE -to 1JUXiiCPZsdcJLVWp6fpCwgi8bXPJLT5HK -amount 1



```