package models

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/model"
	"time"
)

// 促销类型
const (
	UnitSale          = iota // 单品折扣/定额，详情页直接可见优惠价格
	RecombinationSale        // 复合优惠，购物车，结账页可见优惠
)

// 单品折扣/定额 -> item 是否叠加

// 订单满减 ->   可指定商品[]product
// 订单满包邮 -> 可指定商品[]product
// Type = 0 不产生互斥
// 促销
type Promotion struct {
	model.MongoModel `inline`
	Name             string           `json:"name"`                     // 促销活动名称
	Description      string           `json:"description"`              // 活动说明
	Items            []*PromotionItem `json:"items" bson:"-"`           // 活动商品
	Type             int8             `json:"type"`                     // 活动类型
	Mutex            bool             `json:"mutex"`                    // 是否可叠加
	Rule             *PromotionRule   `json:"rule"`                     // 促销规则
	Policy           *PromotionPolicy `json:"policy"`                   // 促销策略
	Enable           bool             `json:"enable"`                   // 是否开启
	BeginAt          time.Time        `json:"begin_at" bson:"begin_at"` // 开始时间
	EndedAt          time.Time        `json:"ended_at" bson:"ended_at"` // 结束时间
}

type AssociatedPromotion struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        int8             `json:"type"`
	Mutex       bool             `json:"mutex"`
	Enable      bool             `json:"enable"`                   // 是否开启
	Rule        *PromotionRule   `json:"rule,omitempty"`           // 促销规则
	Policy      *PromotionPolicy `json:"policy,omitempty"`         // 促销策略
	BeginAt     time.Time        `json:"begin_at" bson:"begin_at"` // 开始时间
	EndedAt     time.Time        `json:"ended_at" bson:"ended_at"` // 结束时间
}

func (p Promotion) ToAssociated() *AssociatedPromotion {
	var description string
	if p.Description == "" && p.Rule != nil && p.Policy != nil {
		description = p.String()
	} else {
		description = p.Description
	}
	return &AssociatedPromotion{
		Id:          p.GetID(),
		Name:        p.Name,
		Description: description,
		Type:        p.Type,
		Mutex:       p.Mutex,
		Enable:      p.Enable,
		BeginAt:     p.BeginAt,
		EndedAt:     p.EndedAt,
		Policy:      p.Policy,
		Rule:        p.Rule,
	}
}

// 优惠描述
func (p Promotion) String() string {
	return fmt.Sprintf("%s,%s", p.Rule.Description(), p.Policy.Description())
}

type promotionItem struct {
	Promotion struct {
		Id   string `json:"id"`
		Name string `json:"name"` // 优惠描述
	} `json:"promotion"`
	ItemId      string `json:"item_id" bson:"item_id"`         // itemId
	OriginPrice int64  `json:"price"`                          // item 定价
	Price       int64  `json:"price"`                          // item价格
	Count       int64  `json:"count"`                          // 购买数量
	UnitAmount  int64  `json:"unit_amount" bson:"unit_amount"` // 优惠后单价
	Amount      int64  `json:"amount"`                         // 实际总额
}

// 单品促销验证
func (p Promotion) ValidationItemPrice(itemId string, price int64) bool {
	for _, product := range p.Items {
		for _, item := range product.Units {
			if item.ItemId == itemId {
				return item.Price == price
			}
		}
	}
	return false
}

// 单品优惠后价格
func (p Promotion) getItemSaleAmount(itemId string, originPrice int64) int64 {
	for _, product := range p.Items {
		for _, item := range product.Units {
			if item.ItemId == itemId {
				return item.Price
			}
		}
	}
	return originPrice
}

// 促进价格计算
func (p Promotion) Calculate(order *Order) (total int64, data []*promotionItem, err error) {
	data = []*promotionItem{}
	// 验证促销时间是否有效

	// 过滤符合规则商品
	items := p.filterIncludeItems(order)
	spew.Dump(items)
	if p.Type == 0 {
		for _, item := range items {
			if p.ValidationItemPrice(item.Item.Id, item.Price) {
				data = append(data, &promotionItem{
					Promotion: struct {
						Id   string `json:"id"`
						Name string `json:"name"`
					}{
						Id:   p.GetID(),
						Name: p.Name,
					},
					ItemId:      item.Item.Id,
					OriginPrice: item.Item.Price,
					Price:       item.Price,
					Count:       item.Count,
					UnitAmount:  item.Price,
					Amount:      item.Price * item.Count,
				})
			} else {
				// 验证不通过，价格被篡改
				saleAmount := p.getItemSaleAmount(item.Item.Id, item.Item.Price)
				return 0, nil, fmt.Errorf("item.id[%s] item.code[%s],实际优惠价格%d,提交优惠价格%d\n", item.Item.Id, item.Item.Code, saleAmount/100, item.Price/100)
			}
		}
	}
	// 验证是否复合规则
	if p.Rule.verify(items) {
		// 计算价格
		total, data = p.Policy.calculate(items)
		return
	}
	return 0, nil, errors.New("不符合促销规则")
}

// 过滤不存在优惠活动中的商品
func (p Promotion) filterIncludeItems(order *Order) (items []*OrderItem) {
	items = []*OrderItem{}
	if p.Type == UnitSale {
		for _, item := range order.OrderItems {
			if p.itemExist(item.Item.Id) {
				items = append(items, item)
			}
		}
	} else {
		// 如果是复合类型优惠，units 为空的话，
		// 该产品所有sku都加入活动，
		// 如果units不为空，则在数组里的sku才生效

		for _, item := range order.OrderItems {
			if p.itemExistWithProductId(item.Item.Id, item.Item.Product.Id) {
				items = append(items, item)
			}
		}
	}

	return items
}

// 验证item是否存在于活动中
func (p Promotion) itemExist(id string) bool {
	for _, item := range p.Items {
		if item.itemExist(id) {
			return true
		}
	}
	return false
}

func (p Promotion) itemExistWithProductId(id string, productId string) bool {
	for _, item := range p.Items {
		if item.itemExistWithProductId(id, productId) {
			return true
		}
	}
	return false
}

// 促销计划产品
type PromotionItem struct {
	model.MongoModel `inline`
	Promotion        *AssociatedPromotion `json:"promotion" bson:"promotion"`
	ProductId        string               `json:"product_id" bson:"product_id" form:"product_id"`
	Product          *AssociatedProduct   `json:"product,omitempty" bson:"-"` // 添加数据结构，方便前端展示
	Units            []*PromotionItemUnit `json:"units"`
}

func NewPromotionItem(productId string) *PromotionItem {
	return &PromotionItem{
		ProductId: productId,
		Units:     []*PromotionItemUnit{},
	}
}

func (p *PromotionItem) AddUnit(item *Item, price int64) error {
	if item.Price < price {
		return fmt.Errorf("价格设置异常,产品定价[%d],促销价格[%d]!!\n", item.Price/100, price/100)
	}
	if p.ProductId != item.Product.Id {
		return fmt.Errorf("item[%s] 不属于 product[%s]\n", item.GetID(), p.ProductId)
	}
	for _, unit := range p.Units {
		if unit.ItemId == item.GetID() {
			return fmt.Errorf("item[%s]以存在\n", unit.ItemId)
		}
	}
	p.Units = append(p.Units, &PromotionItemUnit{
		ItemId: item.GetID(),
		Item:   item.ToAssociated(),
		Price:  price,
	})
	return nil
}

// 验证item是否存在于集合中
func (p PromotionItem) itemExist(id string) bool {
	for _, item := range p.Units {
		if item.ItemId == id {
			return true
		}
	}
	return false
}

func (p PromotionItem) itemExistWithProductId(id string, productId string) bool {
	if p.ProductId == productId {
		if len(p.Units) == 0 {
			return true
		}

		for _, item := range p.Units {
			if item.ItemId == id {
				return true
			}
		}
	}
	return false
}

type PromotionItemUnit struct {
	ItemId string          `json:"item_id" bson:"item_id" form:"item_id"`
	Item   *AssociatedItem `json:"item" bson:"-"`      // 冗余数据，不存库
	Price  int64           `json:"price" bson:"price"` // 活动价
}

// Order Items ["123","33"]

// 促销流程
// 产品-> 可以对应多个促销计划
// 前端产品促销计划列表

const (
	UnLimitedRule         = iota // 不限规则
	AmountGreaterThanRule        // 金额大于
	QtyGreaterThanRule           // 数量大于
)

// 复合优惠，规则
type PromotionRule struct {
	Type  int8   `json:"type"`
	Value uint64 `json:"value"`
}

// 验证是否复合规则
func (p PromotionRule) verify(items []*OrderItem) bool {
	switch p.Type {
	case 0:
		return true
	case 1:
		var amount int64
		for _, item := range items {
			amount += item.Price
		}
		return uint64(amount) >= p.Value
	case 2:
		var countQty int64
		for _, item := range items {
			countQty += item.Count
		}
		return countQty > int64(p.Value)
	}
	return false
}

func (this PromotionRule) Description() string {
	switch this.Type {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("订单金额大于%d", this.Value/100)
	case 2:
		return fmt.Sprintf("单笔订单商品数量大于%d", this.Value)
	}
	return ""
}

// Verify() bool // 效验是否符合
// 价格大于
type PromotionRuleAmountGreaterThan struct {
	Value int64 // 最低价格
}

const (
	DiscountPolicy     = iota + 1 // 打折
	SalePolicy                    // 直减
	FreeShipmentPolicy            // 免邮
)

// 策略
type PromotionPolicy struct {
	Type  int8  `json:"type"`
	Value int64 `json:"value"`
}

func (p PromotionPolicy) Description() string {
	switch p.Type {
	case 1:
		return fmt.Sprintf("%d折", p.Value/10)
	case 2:
		return fmt.Sprintf("直减%d", p.Value/100)
	case 3:
		return "免运费"
	}
	return ""
}

func (p PromotionPolicy) calculate(items []*OrderItem) (total int64, data []*promotionItem) {
	data = []*promotionItem{}
	switch p.Type {
	case 1:
		// 打折
		for _, item := range items {
			amount := item.Price * p.Value / 100 // 折后价
			itemTotalAmount := amount * item.Count
			data = append(data, &promotionItem{
				ItemId:      item.Item.Id,    // skuId
				OriginPrice: item.Item.Price, // sku定价
				Price:       item.Price,      // 单品价格
				Count:       item.Count,      // 数量
				UnitAmount:  amount,          // 优惠后单价
				Amount:      itemTotalAmount, // item优惠后总金额
			})
			total += itemTotalAmount
		}
		return total, data

	case 2:
		// 直减，订单总额减
		// 总金额
		var totalAmount int64
		for _, item := range items {
			totalAmount += item.Amount * item.Count
		}
		// 优惠单价
		avgSale := p.Value / totalAmount
		for _, item := range items {
			data = append(data, &promotionItem{
				ItemId:      item.Item.Id,                      // skuId
				OriginPrice: item.Item.Price,                   // sku定价
				Price:       item.Price,                        // 单品价格
				Count:       item.Count,                        // 数量
				UnitAmount:  item.Price * avgSale,              // 优惠后单价
				Amount:      item.Price * avgSale * item.Count, // item优惠后总金额
			})
			total += item.Price * avgSale * item.Count
		}
		return total, data

	}
	return 0, nil
}
