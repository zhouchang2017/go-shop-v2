# go-shop-v2


#### Docker Compose

单独启动rabbitmq
```bash
docker-compose up -d rabbitmq
```

MongoDB 

keyfile
```shell script
$ openssl rand -base64 756 > mongod-keyfile
$ chmod 600 mongod-keyfile

# 在ubuntu 需要设置为docker用户所有
You'll just need to chown 999:999 keyfile before you run the docker container (you probably need sudo in there).

```

本地host设置
```
$ sudo vim /etc/hosts
+127.0.0.1 mongodb-primary
+127.0.0.1 mognodb-secondary
+127.0.0.1 mongodb-arbiter

```

副本集初始化
```bash
docker-compose exec mongodb-primary mongo --port 30000 -uroot -p12345678 /root/000_init_replSet.js
```


- 列表页 index
- 专题页 topic
    + 商品列表
- 文章页 article
    + relation some product link

列表页
```
/api/v1/index
```


产品删除
  - 删除items
  - 删除购物车中相关产品


#### ChangeLog 
- 2020-02-16
    - Category 模型不在包含 `[]*OptionValue`,每个产品都管理自身的 `[]*OptionValue`
    - `OptionValue`下的`Value`模型，新添加`Image string`字段

- 2020-02-21
    - 七牛(qiniu) pkg 添加 Image 类型
    
- 2020-02-29
    - mongoRep，上包裹一层 redisCache 在基础方法中实现缓存
        + `FindById`
        + `Create`
        + `Save`
        + `Update`
        + `Delete`
        + `DeleteMany`
        + `Restore`