package repositories

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestPromotionRep_Create(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	// 从数据库检索一个产品
	itemRep := NewItemRep(NewBasicMongoRepositoryByDefault(&models.Item{}, mongodb.GetConFn()))
	productRep := NewProductRep(NewBasicMongoRepositoryByDefault(&models.Product{}, mongodb.GetConFn()), itemRep)

	startAt, err := time.Parse(`"2006-01-02 15:04:05"`, `"2020-03-01 00:00:00"`)
	if err != nil {
		t.Fatal(err)
	}

	// 新建促销计划
	promotion := &models.Promotion{
		Name:    "新春特惠",
		Type:    0,     // 单品优惠
		Mutex:   false, // 单品优惠不存在互斥
		Rule:    nil,   // 单品优惠不存在规则
		Policy:  nil,   // 单品优惠不存在
		Enable:  true,
		BeginAt: startAt,
		EndedAt: startAt.Add(time.Hour * 24 * 30),
	}

	var promotionItems []*models.PromotionItem

	results := <-productRep.Pagination(context.Background(), &request.IndexRequest{Page: 2})
	if results.Error != nil {
		t.Fatal(err)
	}
	products := results.Result.([]*models.Product)

	for _, product := range products {
		productId := product.GetID()
		product, err := productRep.WithItems(context.Background(), productId)
		if err != nil {
			t.Fatal(err)
		}
		// 促销计划产品
		p := models.NewPromotionItem(product.GetID())

		// 随机价格
		for index, item := range product.Items {
			var price int64
			if index%2 == 0 {
				// 偶数 7折
				price = item.Price * 70 / 100
			} else {
				// 奇数 8折
				price = item.Price * 80 / 100
			}
			err := p.AddUnit(item, price)
			if err != nil {
				t.Fatal(err)
			}
		}
		promotionItems = append(promotionItems, p)
	}

	promotion.Items = promotionItems

	promotionItemRep := NewPromotionItemRep(NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	rep := NewPromotionRep(NewBasicMongoRepositoryByDefault(&models.Promotion{}, mongodb.GetConFn()), promotionItemRep)

	created := <-rep.Create(context.Background(), promotion)

	if created.Error != nil {
		t.Fatal(created.Error)
	}

	spew.Dump(created.Result)
}

func TestPromotion_Calc(t *testing.T) {

	mongodb.TestConnect()
	defer mongodb.Close()

	// 从数据库检索一个产品
	itemRep := NewItemRep(NewBasicMongoRepositoryByDefault(&models.Item{}, mongodb.GetConFn()))
	productRep := NewProductRep(NewBasicMongoRepositoryByDefault(&models.Product{}, mongodb.GetConFn()), itemRep)

	startAt, err := time.Parse(`"2006-01-02 15:04:05"`, `"2020-03-01 00:00:00"`)
	if err != nil {
		t.Fatal(err)
	}

	// 新建促销计划
	promotion := &models.Promotion{
		Name:  "38女王节",                  // 满899 减50
		Type:  models.RecombinationSale, // 复合优惠
		Mutex: false,                    // 互斥
		Rule: &models.PromotionRule{
			Type:  models.AmountGreaterThanRule,
			Value: 89900,
		}, // 规则
		Policy: &models.PromotionPolicy{
			Type:  models.SalePolicy,
			Value: 5000,
		}, // 优惠策略
		Enable:  true,
		BeginAt: startAt,
		EndedAt: startAt.Add(time.Hour * 24 * 30),
	}
	promotion.ID = primitive.NewObjectID()
	promotion.CreatedAt = startAt
	promotion.UpdatedAt = time.Now()

	var promotionItems []*models.PromotionItem

	results := <-productRep.Pagination(context.Background(), &request.IndexRequest{Page: 2})
	if results.Error != nil {
		t.Fatal(err)
	}
	products := results.Result.([]*models.Product)

	for _, product := range products {
		productId := product.GetID()
		product, err := productRep.WithItems(context.Background(), productId)
		if err != nil {
			t.Fatal(err)
		}
		// 促销计划产品
		p := models.NewPromotionItem(product.GetID())
		promotionItems = append(promotionItems, p)
	}

	promotion.Items = promotionItems

	promotionItemRep := NewPromotionItemRep(NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	rep := NewPromotionRep(NewBasicMongoRepositoryByDefault(&models.Promotion{}, mongodb.GetConFn()), promotionItemRep)

	created := <-rep.Create(context.Background(), promotion)

	if created.Error != nil {
		t.Fatal(created.Error)
	}

	spew.Dump(created.Result)
}

func TestPromotionRep_resolveItemIds(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	promotionItemRep := NewPromotionItemRep(NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	rep := NewPromotionRep(NewBasicMongoRepositoryByDefault(&models.Promotion{}, mongodb.GetConFn()), promotionItemRep)

	ids := rep.resolveItemIds(context.Background(), "5e61b41ff6d70d76025b9f2c")

	spew.Dump(ids)

}

func TestPromotionRep_FindActivePromotionByProductId(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	promotionItemRep := NewPromotionItemRep(NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	rep := NewPromotionRep(NewBasicMongoRepositoryByDefault(&models.Promotion{}, mongodb.GetConFn()), promotionItemRep)

	items := rep.FindActivePromotionByProductId(context.Background(), "5e577e370d3f4744961cfcfd")

	spew.Dump(items)
}

func TestPromotionRep_FindProductPrice(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	promotionItemRep := NewPromotionItemRep(NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	rep := NewPromotionRep(NewBasicMongoRepositoryByDefault(&models.Promotion{}, mongodb.GetConFn()), promotionItemRep)

	price := rep.FindProductsPrice(context.Background(), "5e577e370d3f4744961cfcfd", "5e577fb20d3f4744961cfd2d","5e577fb20d3f4744961cfd2f")

	spew.Dump(price)
}

func TestPromotionRep_Delete(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	promotionItemRep := NewPromotionItemRep(NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	rep := NewPromotionRep(NewBasicMongoRepositoryByDefault(&models.Promotion{}, mongodb.GetConFn()), promotionItemRep)

	error := <- rep.Delete(context.Background(), "5e61c2447859f8230a004403")
	if error!=nil {
		t.Fatal(error)
	}
}