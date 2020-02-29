# go-shop-v2


#### Docker Compose

单独启动rabbitmq
```bash
docker-compose up -d rabbitmq
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


#### ChangeLog 
- 2020-02-16
    - Category 模型不在包含 `[]*OptionValue`,每个产品都管理自身的 `[]*OptionValue`
    - `OptionValue`下的`Value`模型，新添加`Image string`字段

- 2020-02-21
    - 七牛(qiniu) pkg 添加 Image 类型
    
- 2020-02-29
    - mongoRep 在基础方法中实现缓存,前提是需要先调用mongoRep下的`SetCache`方法
        + `FindById`
        + `Create`
        + `Save`
        + `Update`
        + `Delete`
        + `DeleteMany`
        + `Restore`