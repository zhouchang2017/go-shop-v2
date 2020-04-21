package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

type ProductService struct {
	rep          *repositories.ProductRep
	promotionRep *repositories.PromotionRep
	itemRep      *repositories.ItemRep
}

func NewProductService(rep *repositories.ProductRep, promotionRep *repositories.PromotionRep, itemRep *repositories.ItemRep) *ProductService {
	return &ProductService{
		rep:          rep,
		promotionRep: promotionRep,
		itemRep:      itemRep,
	}
}

// 修改销量
func (this *ProductService) UpdateSalesQty(ctx context.Context, itemId string, count int64) (err error) {
	item, err := this.itemRep.UpdateSalesQty(ctx, itemId, count)
	if err == nil {
		product, err := this.FindById(ctx, item.Product.Id)
		if err != nil {
			// 404
			return nil
		}
		qty := int64(product.TotalSalesQty) + count
		if qty < 0 {
			product.TotalSalesQty = 0
		} else {
			product.TotalSalesQty = uint64(qty)
		}

		<-this.rep.Save(ctx, product)
	}

	return nil
}

func (this *ProductService) FindItemById(ctx context.Context, id string) (item *models.Item, err error) {
	item, err = this.rep.FindItemById(ctx, id)
	if err != nil {
		return
	}
	promotionItem, err := this.promotionRep.FindActivePromotionUnitSaleByProductId(ctx, item.Product.Id)
	if err == nil {
		price := promotionItem.FindPriceByItemId(item.GetID())
		if price != -1 {
			item.PromotionPrice = uint64(price)
			return item, nil
		}
	}
	item.PromotionPrice = item.Price

	return item, nil
}

func (this *ProductService) FindByIdWithItems(ctx context.Context, id string) (product *models.Product, err error) {
	res := <-this.rep.FindById(ctx, id)
	if res.Error != nil {
		return nil, res.Error
	}
	product = res.Result.(*models.Product)
	promotionItem, err := this.promotionRep.FindActivePromotionUnitSaleByProductId(ctx, product.GetID())
	if err == nil {
		// 存在单品促销
		product.PromotionPrice = promotionItem.MinPrice()
	} else {
		// 不存在单品促销，促销价设置为自身价格
		product.PromotionPrice = product.Price
	}

	itemRes := <-this.rep.GetItemRep().FindByProductId(ctx, id)
	if itemRes.Error != nil {
		return nil, itemRes.Error
	}
	items := itemRes.Result.([]*models.Item)

	for _, item := range items {
		if promotionItem != nil {
			itemPromotionPrice := promotionItem.FindPriceByItemId(item.GetID())
			if itemPromotionPrice != -1 {
				item.PromotionPrice = uint64(itemPromotionPrice)
				continue
			}
		}
		item.PromotionPrice = item.Price
	}
	product.Items = items

	return product, nil
}

func (this *ProductService) List(ctx context.Context, req *request.IndexRequest) (products []contracts.RelationsOption, pagination response.Pagination, err error) {
	req.SetSearchField("code")
	req.AppendProjection("images", bson.M{"$slice": 1})
	req.AppendProjection("_id", 1)
	req.AppendProjection("name", 1)
	req.AppendProjection("code", 1)
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	for _, product := range results.Result.([]*models.Product) {
		var avatar string
		if len(product.Images) == 1 {
			avatar = product.Images[0].Src()
		}
		products = append(products, contracts.RelationsOption{
			Id:     product.GetID(),
			Name:   product.Code,
			Avatar: avatar,
		})
	}
	pagination = results.Pagination
	return
}

func (this *ProductService) FindByIds(ctx context.Context, ids []string) (products []*models.Product, err error) {
	results := <-this.rep.FindByIds(ctx, ids)
	if results.Error != nil {
		err = results.Error
		return
	}
	products = results.Result.([]*models.Product)
	var g errgroup.Group
	res := []*models.Product{}
	sem := make(chan struct{}, 10)

	for _, product := range products {
		product := product // local variable
		sem <- struct{}{}
		// 获取产品单品促销活动信息
		g.Go(func() error {
			promotionItem, err := this.promotionRep.FindActivePromotionUnitSaleByProductId(ctx, product.GetID())
			if err == nil {
				// 存在单品促销
				product.PromotionPrice = promotionItem.MinPrice()
			} else {
				// 不存在单品促销，促销价设置为自身价格
				product.PromotionPrice = product.Price
			}

			res = append(res, product)
			<-sem
			return nil
		})

		if err := g.Wait(); err != nil {
			return res, err
		}
		products = res
	}

	return products, nil
}

func (this *ProductService) RelationResolveIds(ctx context.Context, ids []string) (products []contracts.RelationsOption, err error) {
	products2, err := this.FindByIds(ctx, ids)
	if err != nil {
		return
	}
	for _, product := range products2 {
		products = append(products, contracts.RelationsOption{
			Id:     product.GetID(),
			Name:   product.Code,
			Avatar: product.Avatar.Src(),
		})
	}
	return
}

// 列表
func (this *ProductService) Pagination(ctx context.Context, req *request.IndexRequest) (products []*models.Product, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	includes := req.Includes()
	products = results.Result.([]*models.Product)

	hasItem := false
	for _, with := range includes {
		if with == "item" {
			hasItem = true
			break
		}
	}

	var g errgroup.Group
	res := []*models.Product{}
	sem := make(chan struct{}, 10)

	for _, product := range products {
		product := product // local variable
		sem <- struct{}{}
		// 获取产品单品促销活动信息
		g.Go(func() error {
			promotionItem, err := this.promotionRep.FindActivePromotionUnitSaleByProductId(ctx, product.GetID())
			if err == nil {
				// 存在单品促销
				product.PromotionPrice = promotionItem.MinPrice()
			} else {
				// 不存在单品促销，促销价设置为自身价格
				product.PromotionPrice = product.Price
			}
			if hasItem {
				// 加载items
				result := <-this.itemRep.FindByProductId(ctx, product.GetID())
				if result.Error == nil {
					items := result.Result.([]*models.Item)
					for _, item := range items {
						if promotionItem != nil {
							itemPromotionPrice := promotionItem.FindPriceByItemId(item.GetID())
							if itemPromotionPrice != -1 {
								item.PromotionPrice = uint64(itemPromotionPrice)
								continue
							}
						}
						item.PromotionPrice = item.Price
					}
					product.Items = items
				}
			}

			res = append(res, product)
			<-sem
			return nil
		})

		if err := g.Wait(); err != nil {
			return res, pagination, err
		}
		products = res
	}
	return products, results.Pagination, nil
}

type ProductCreateOption struct {
	Name         string                       `json:"name" form:"name" binding:"required,max=255"`
	Code         string                       `json:"code" form:"code" binding:"required,max=255"`
	Brand        *models.AssociatedBrand      `json:"brand" form:"brand"`
	Category     *models.AssociatedCategory   `json:"category" form:"category"`
	Attributes   []*models.ProductAttribute   `json:"attributes" form:"attributes"`
	Options      []*productOptionCreateOption `json:"options" form:"options"`
	Items        []*productItemCreateOption   `json:"items"`
	Description  string                       `json:"description"`
	Price        uint64                       `json:"price"`
	FakeSalesQty uint64                       `json:"fake_sales_qty" form:"fake_sales_qty"`
	Images       []qiniu.Image                `json:"images" form:"images"`
	OnSale       bool                         `json:"on_sale" form:"on_sale"`
	Sort         int64                        `json:"sort" form:"sort"`
}

type productOptionCreateOption struct {
	Id     string                            `json:"id" form:"id"`
	Uid    int                               `json:"uid" form:"uid"`
	Name   string                            `json:"name" form:"name"`
	Image  bool                              `json:"image" form:"image"`
	Values []*productOptionValueCreateOption `json:"values" form:"values"`
}

type productOptionValueCreateOption struct {
	Id    string       `json:"id" form:"id"`
	Uid   int          `json:"uid" form:"uid"`
	Name  string       `json:"name" form:"name"`
	Image *qiniu.Image `json:"image" form:"image"`
}

type productItemCreateOption struct {
	Id           string                            `json:"id"`
	Code         string                            `json:"code"`
	Price        uint64                            `json:"price,omitempty"`
	OptionValues []*productOptionValueCreateOption `json:"option_values" form:"option_values" `
	OnSale       bool                              `json:"on_sale" form:"on_sale"` // 上/下架 受product影响
	Qty          uint64                            `json:"qty" `                   // 可售数量
}

// 创建
func (this *ProductService) Create(ctx context.Context, opt ProductCreateOption) (product *models.Product, err error) {
	var images []qiniu.Image
	for _, image := range opt.Images {
		images = append(images, qiniu.NewImage(image.Src()))
	}
	model := &models.Product{
		Name:         opt.Name,
		Code:         opt.Code,
		Brand:        opt.Brand,
		Category:     opt.Category,
		Attributes:   opt.Attributes,
		Description:  opt.Description,
		Price:        opt.Price,
		FakeSalesQty: opt.FakeSalesQty,
		Images:       images,
		OnSale:       opt.OnSale,
	}
	model.SetAvatar()
	// uid => option
	optionMaps := map[int]*models.ProductOption{}
	valuesMaps := map[int]*models.OptionValue{}
	for _, o := range opt.Options {
		option := models.NewProductOption(o.Name)
		option.Image = o.Image
		for _, value := range o.Values {
			optionValue := option.NewValue(value.Name)
			if value.Image != nil {
				optionValue.Image = value.Image
			}
			valuesMaps[value.Uid] = optionValue
			option.AddValues(optionValue)
		}
		optionMaps[o.Uid] = option
		model.Options = append(model.Options, option)
	}
	var items []*models.Item
	for _, item := range opt.Items {
		newItem := &models.Item{
			Code:   item.Code,
			Price:  item.Price,
			OnSale: item.OnSale,
			Qty:    item.Qty,
		}
		for _, value := range item.OptionValues {
			if optVal, ok := valuesMaps[value.Uid]; ok {
				newItem.OptionValues = append(newItem.OptionValues, optVal)
			}
		}
		newItem.OnSale = model.OnSale
		items = append(items, newItem)
	}

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return
	}
	if err = session.StartTransaction(); err != nil {
		return
	}

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		created := <-this.rep.Create(sessionContext, model)
		if created.Error != nil {
			spew.Dump("created Product Error", created.Error)
			session.AbortTransaction(sessionContext)
			return created.Error
		}
		product = created.Result.(*models.Product)

		for _, item := range items {
			item.Product = product.ToAssociated()
			item.SetAvatar()
		}
		// 创建items
		newItems, err := this.itemRep.CreateMany(sessionContext, items)
		if err != nil {
			spew.Dump("item create many error:", err)
			session.AbortTransaction(sessionContext)
			return err
		}

		product.Items = newItems
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return
}

// 更新
func (this *ProductService) Update(ctx context.Context, model *models.Product, opt ProductCreateOption) (product *models.Product, err error) {
	model.Name = opt.Name
	model.Brand = opt.Brand
	model.Attributes = opt.Attributes
	model.Description = opt.Description
	model.Price = opt.Price
	model.FakeSalesQty = opt.FakeSalesQty
	model.Images = opt.Images
	model.OnSale = opt.OnSale

	model.Options = []*models.ProductOption{}
	// uid => option
	optionMaps := map[int]*models.ProductOption{}
	valuesMaps := map[int]*models.OptionValue{}
	for _, o := range opt.Options {
		var option *models.ProductOption
		if o.Id != "" {
			option = models.MakeProductOption(o.Id, o.Name)
		} else {
			option = models.NewProductOption(o.Name)
		}
		option.Image = o.Image
		for _, value := range o.Values {
			var optionValue *models.OptionValue
			if value.Id != "" {
				optionValue = option.MakeValue(value.Id, value.Name)
			} else {
				optionValue = option.NewValue(value.Name)
			}
			if value.Image != nil {
				optionValue.Image = value.Image
			}
			valuesMaps[value.Uid] = optionValue
			option.AddValues(optionValue)
		}
		optionMaps[o.Uid] = option
		model.Options = append(model.Options, option)
	}

	model.SetAvatar()
	var newItems models.Items

	res := <-this.itemRep.FindByProductId(ctx, model.GetID())
	if res.Error != nil {
		err = res.Error
		return
	}
	productItems := models.Items(res.Result.([]*models.Item))
	for _, item := range opt.Items {
		var newItem *models.Item
		if item.Id != "" {
			find := productItems.FindById(item.Id)
			if find != nil {
				find.OptionValues = []*models.OptionValue{}
				newItem = find
			} else {
				newItem = models.NewItem()
			}
		} else {
			newItem = models.NewItem()
		}

		newItem.Product = model.ToAssociated()
		newItem.Qty = item.Qty
		newItem.Price = item.Price
		newItem.Code = item.Code

		for _, value := range item.OptionValues {
			if optVal, ok := valuesMaps[value.Uid]; ok {
				newItem.OptionValues = append(newItem.OptionValues, optVal)
			}
		}
		newItem.SetAvatar()
		newItem.OnSale = model.OnSale
		newItems = append(newItems, newItem)
	}

	// 被移除的item
	var deleteItemIds []string
	for _, item := range productItems {
		if newItems.FindById(item.GetID()) == nil {
			deleteItemIds = append(deleteItemIds, item.GetID())
		}
	}

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return
	}
	if err = session.StartTransaction(); err != nil {
		return
	}

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {

		// 更新产品
		saved := <-this.rep.Save(sessionContext, model)
		if saved.Error != nil {
			session.AbortTransaction(sessionContext)
			return saved.Error
		}
		product = saved.Result.(*models.Product)
		product.Items = []*models.Item{}
		// 更新变体
		for _, item := range newItems {
			if item.ID.IsZero() {
				// 新增变体
				created := <-this.itemRep.Create(sessionContext, item)
				if created.Error != nil {
					log.Errorf("update product %s add item error:%s", product.GetID(), created.Error)
					session.AbortTransaction(sessionContext)
					return created.Error
				}
				product.Items = append(product.Items, created.Result.(*models.Item))
			} else {
				saved := <-this.itemRep.Save(sessionContext, item)
				if saved.Error != nil {
					log.Errorf("update product %s save item[%s] error:%s", product.GetID(), item.GetID(), saved.Error)
					session.AbortTransaction(sessionContext)
					return saved.Error
				}
				product.Items = append(product.Items, saved.Result.(*models.Item))
			}
		}

		// 删除变体
		if err = <-this.itemRep.DeleteMany(sessionContext, deleteItemIds...); err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		session.CommitTransaction(sessionContext)
		return nil
	})

	session.EndSession(ctx)
	return
}

// 详情
func (this *ProductService) FindById(ctx context.Context, id string) (product *models.Product, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	product = byId.Result.(*models.Product)

	promotionItem, err := this.promotionRep.FindActivePromotionUnitSaleByProductId(ctx, product.GetID())
	if err != nil {
		// log
		product.PromotionPrice = product.Price
		return product, nil
	}
	product.PromotionPrice = promotionItem.MinPrice()
	return product, nil
}

// 删除
func (this *ProductService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *ProductService) Restore(ctx context.Context, id string) (product *models.Product, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Product), nil
}

// 所有OptionName
func (this *ProductService) AvailableOptionNames(ctx context.Context) (names []string) {
	return this.rep.AvailableOptionNames(ctx)
}
