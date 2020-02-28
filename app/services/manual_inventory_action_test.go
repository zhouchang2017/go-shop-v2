package services

import (
	"context"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/utils"
	"testing"
)

// 批量随机添加库存
func TestManualInventoryActionService_Put(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	qiniu.NewQiniu(qiniu.Config{
		Drive:            "",
		QiniuDomain:      "http://q5q1efml2.bkt.clouddn.com",
		QiniuAccessKey:   "",
		QiniuSecretKey:   "",
		Bucket:           "",
		FileUploadAction: "",
	})

	adminService := MakeAdminService()
	admins, _, err2 := adminService.Pagination(context.Background(), &request.IndexRequest{})
	if err2 != nil {
		t.Fatal(err2)
	}

	user := admins[0]

	productService := MakeProductService()
	manualInventoryActionService := MakeManualInventoryActionService()

	var page int64 = 1
	hasNextPage := true
	for hasNextPage {
		indexRequest := &request.IndexRequest{}
		indexRequest.Page = page
		page++
		products, pagination, err := productService.Pagination(context.Background(), indexRequest)

		if err != nil {
			t.Fatal(err)
		}
		hasNextPage = pagination.HasNextPage

		for _, shop := range user.Shops {
			var inventoryActionPutOptionItems []*InventoryActionItemOption

			for _, product := range products {

				p, err := productService.FindByIdWithItems(context.Background(), product.GetID())
				if err != nil {
					t.Fatal(err)
				}
				for _, item := range p.Items {

					randomInt := utils.RandomInt(20)
					if randomInt > 0 {
						inventoryActionPutOptionItems = append(inventoryActionPutOptionItems, &InventoryActionItemOption{
							Id:     item.GetID(),
							Qty:    randomInt,
							Status: 0,
						})
					}

				}
			}

			if len(inventoryActionPutOptionItems) > 0 {
				option := &InventoryActionPutOption{
					ShopId: shop.Id,
					Items:  inventoryActionPutOptionItems,
				}

				puted, err := manualInventoryActionService.Put(context.Background(), option, user)
				if err != nil {
					t.Fatal(err)
				}

				_, err = manualInventoryActionService.StatusToFinished(context.Background(), puted.GetID())
				if err != nil {
					t.Fatal(err)
				}
			}

		}

	}
}
