# go-shop-v2


#### Docker Compose

单独启动rabbitmq
```bash
docker-compose up -d rabbitmq
```


```
docker swarm join --advertise-addr 150.158.121.107 --data-path-addr 150.158.121.107 --token SWMTKN-1-5srsj49v8js9r2lwmbnbmomxs5mna3g68ar4bojdrc4ok0gubk-d92i3bq59h7oczmpjzrpakpsw 106.54.17.169:2377
```

grafana
```
    # sudo mkdir -p /srv/docker/grafana/data; chown 472:472 /srv/docker/grafana/data

docker run -d --name=grafana --network=overlay --restart=always -v $HOME/configs/grafana:/var/lib/grafana -p 13000:3000 grafana/grafana
```

influxdb
```
docker run -d --name=influxdb --network=influxdb -e INFLUXDB_HTTP_AUTH_ENABLED=true -e INFLUXDB_ADMIN_USER=root -e INFLUXDB_ADMIN_PASSWORD=12345678 --restart=always -v $PWD/influxdb:/var/lib/influxdb -p 8086:8086 influxdb
```
MongoDB 

```
docker service create --replicas 1 \
--network overlay --mount type=volume,source=rsdata1,target=/data/db \
--mount type=bind,source=$HOME/mongod-keyfile,target=/etc/mongod-keyfile,readonly \
--constraint 'node.labels.mongo.rs==1' -p 27017:27017 \
-e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=12345678 \
--name mongo_rs1 mongo:latest mongod --bind_ip_all --replSet rs0 --auth --keyFile /etc/mongod-keyfile

service create --replicas 1 --network overlay --constraint 'node.labels.mongo.rs==2' --name tool busybox

docker service create --replicas 1 \
--network overlay --mount type=volume,source=rsdata2,target=/data/db \
--mount type=bind,source=$HOME/mongod-keyfile,target=/etc/mongod-keyfile,readonly \
-e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=12345678 \
--constraint 'node.labels.mongo.rs==1' \
--name mongo_rs2 mongo:latest mongod --bind_ip_all --replSet rs0 --auth --keyFile /etc/mongod-keyfile

docker service create --replicas 1 \
--network overlay --mount type=volume,source=rsdata4,target=/data/db \
--mount type=bind,source=$HOME/mongod-keyfile,target=/etc/mongod-keyfile,readonly \
-e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=12345678 \
--constraint 'node.labels.mongo.rs==3' \
--name mongo_rs3 mongo:latest mongod --bind_ip_all --replSet rs0 --auth --keyFile /etc/mongod-keyfile
```

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

```
docker run -it --rm -p 5050:5050 --name test-app -e DB_HOST=mongo_rs1 -e DB_NAME=go-shop -e DB_USERNAME=root -e DB_PASSWORD=12345678 -e DB_REPLICA_SET=rs0 zhouchang2018/test-demo:v2
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

推送通知模板
- 订单号
- 商品名称
- 订单金额
- 时间

dashboard
- 当日新订单
- 当日付款金额
- 当日付款订单笔数
- 待发货订单数
- 待付款订单数
- 当日新增用户数

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