package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
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

func (a *AssociatedPromotion) makePromotion() *Promotion {
	objId, _ := primitive.ObjectIDFromHex(a.Id)
	promotion := &Promotion{
		Name:        a.Name,
		Description: a.Description,
		Type:        a.Type,
		Mutex:       a.Mutex,
		Rule:        a.Rule,
		Policy:      a.Policy,
		Enable:      a.Enable,
	}
	promotion.ID = objId
	return promotion
}

func (p Promotion) GetDescription() string {
	var description string
	if p.Description == "" && p.Rule != nil && p.Policy != nil {
		description = p.String()
	} else {
		description = p.Description
	}
	return description
}

func (p Promotion) ToAssociated() *AssociatedPromotion {
	return &AssociatedPromotion{
		Id:          p.GetID(),
		Name:        p.Name,
		Description: p.GetDescription(),
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
	return fmt.Sprintf("%s，%s", p.Rule.Description(), p.Policy.Description())
}

// 通过产品id获取产品优惠item
func (p Promotion) findItemByProductId(productId string) *PromotionItem {
	for _, item := range p.Items {
		if item.ProductId == productId {
			return item
		}
	}
	return nil
}

func (p *Promotion) addItem(promotionItem *PromotionItem, itemId string, price uint64, qty uint64) {
	find := p.findItemByProductId(promotionItem.ProductId)
	if find != nil {
		find.addItem(itemId, promotionItem.ProductId, price, qty)
	} else {
		promotionItem.addItem(itemId, promotionItem.ProductId, price, qty)
		p.Items = append(p.Items, promotionItem)
	}
}

func (p *Promotion) totalItemsAmountAndQty() (amount uint64, qty uint64) {
	for _, item := range p.Items {
		for _, i := range item.items {
			amount += i.qty * i.price
			qty += i.qty
		}
	}
	return
}

func (p *Promotion) assignSalePrice(salePrice uint64) {
	amount, _ := p.totalItemsAmountAndQty()
	avg := float64(salePrice) / float64(amount)
	var used uint64
	for index, item := range p.Items {
		for ind, i := range item.items {
			if index+1 == len(p.Items) && ind+1 == len(item.items) {
				i.salePrice = salePrice - used
				i.unitSalePrice = uint64(math.Ceil(float64(i.salePrice) / float64(i.qty)))
			} else {
				i.unitSalePrice = uint64(math.Ceil(avg * float64(i.price)))
				i.salePrice = i.unitSalePrice * i.qty
				used += i.salePrice
			}
		}
	}
}

func (p *Promotion) calculate(salePrices uint64) (promotion *Promotion, salePrice uint64, accessRule bool) {
	amount, qty := p.totalItemsAmountAndQty()
	amount -= salePrices
	// 判断规则
	if accessRule = p.Rule.verify(amount, qty); accessRule {
		// 策略
		salePrice = p.Policy.calculate(amount)
		p.assignSalePrice(salePrice)
		return p, salePrice, true
	}
	return nil, 0, false
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
func (p Promotion) ValidationItemPrice(itemId string, price uint64) bool {
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
func (p Promotion) getItemSaleAmount(itemId string, originPrice uint64) uint64 {
	for _, product := range p.Items {
		for _, item := range product.Units {
			if item.ItemId == itemId {
				return item.Price
			}
		}
	}
	return originPrice
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

type promotionItemItem struct {
	itemId        string
	price         uint64
	qty           uint64
	salePrice     uint64
	unitSalePrice uint64
}

// 促销计划产品
type PromotionItem struct {
	model.MongoModel `inline`
	items            []*promotionItemItem // 用户计算，参加该活动的产品
	Promotion        *AssociatedPromotion `json:"promotion" bson:"promotion"`
	ProductId        string               `json:"product_id" bson:"product_id" form:"product_id"`
	Product          *AssociatedProduct   `json:"product,omitempty" bson:"-"` // 添加数据结构，方便前端展示
	Units            []*PromotionItemUnit `json:"units"`
}

// 添加受改活动影响商品
func (p *PromotionItem) addItem(itemId string, productId string, price uint64, qty uint64) {
	if productId != p.ProductId {
		// 不符合规则
		return
	}
	if p.Promotion != nil {
		// 验证是否在活动内
		if p.Promotion.Type == 0 {
			// 单品活动
			if p.itemExist(itemId) {
				var exist bool
				for _, item := range p.items {
					if item.itemId == itemId {
						exist = true
						// 已存在同样产品
						if item.price == price {
							item.qty += qty
						}
					}
				}
				if !exist {
					p.items = append(p.items, &promotionItemItem{
						itemId: itemId,
						price:  price,
						qty:    qty,
					})
				}
			}
		}
		if p.Promotion.Type == 1 {
			// 复合活动
			if p.itemExistWithProductId(itemId, productId) {
				var exist bool
				for _, item := range p.items {
					if item.itemId == itemId {
						exist = true
						// 已存在同样产品
						if item.price == price {
							item.qty += qty
						}
					}
				}
				if !exist {
					p.items = append(p.items, &promotionItemItem{
						itemId: itemId,
						price:  price,
						qty:    qty,
					})
				}
			}
		}
	}
}

func NewPromotionItem(productId string) *PromotionItem {
	return &PromotionItem{
		ProductId: productId,
		Units:     []*PromotionItemUnit{},
	}
}

// units中最便宜的价格
func (p *PromotionItem) MinPrice() uint64 {
	var prices []int64
	for _, item := range p.Units {
		prices = append(prices, int64(item.Price))
	}
	return uint64(utils.Min(prices...))
}

// units中最贵的价格
func (p *PromotionItem) MaxPrice() uint64 {
	var prices []int64
	for _, item := range p.Units {
		prices = append(prices, int64(item.Price))
	}
	return uint64(utils.Max(prices...))
}

// 获取item促销价格
// 不存在返回-1
func (p PromotionItem) FindPriceByItemId(id string) int64 {
	for _, unit := range p.Units {
		if unit.ItemId == id {
			return int64(unit.Price)
		}
	}
	return -1
}

func (p *PromotionItem) AddUnit(item *Item, price uint64) error {
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

type PromotionItems []*PromotionItem

func (p PromotionItems) IncludeItem(id string, productId string) PromotionItems {
	items := make(PromotionItems, 0)
	for _, item := range p {
		if item.Promotion.Type == UnitSale {
			// 单品活动
			if item.itemExist(id) {
				items = append(items, item)
				continue
			}
		}
		// 复合活动
		if item.itemExistWithProductId(id, productId) {
			items = append(items, item)
		}
	}
	return items
}

func (p PromotionItems) FindById(id string) *PromotionItem {
	for _, item := range p {
		if item.Promotion.Id == id {
			return item
		}
	}
	return nil
}

type PromotionItemUnit struct {
	ItemId      string          `json:"item_id" bson:"item_id" form:"item_id"`
	Item        *AssociatedItem `json:"item" bson:"-"`                    // 冗余数据，不存库
	Price       uint64           `json:"price" bson:"price"`               // 活动价
	OriginPrice uint64           `json:"origin_price" bson:"origin_price"` // 原始价格
}

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
func (p PromotionRule) verify(amount uint64, qty uint64) bool {
	switch p.Type {
	case 0:
		return true
	case 1:
		return uint64(amount) >= p.Value
	case 2:
		return qty >= uint64(p.Value)
	}
	return false
}

func (this PromotionRule) Description() string {
	switch this.Type {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("订单金额大于%s元", utils.ToMoneyString(this.Value))
	case 2:
		return fmt.Sprintf("单笔订单商品数量大于%d件", this.Value)
	}
	return ""
}

const (
	DiscountPolicy     = iota + 1 // 打折
	SalePolicy                    // 直减
	FreeShipmentPolicy            // 免邮
)

// 策略
type PromotionPolicy struct {
	Type  int8  `json:"type"`
	Value uint64 `json:"value"`
}

func (p PromotionPolicy) Description() string {
	switch p.Type {
	case 1:
		return fmt.Sprintf("享%d折", p.Value/10)
	case 2:
		return fmt.Sprintf("优惠%s元", utils.ToMoneyString(p.Value))
	case 3:
		return "免运费"
	}
	return ""
}

func (p PromotionPolicy) calculate(amount uint64) (salePrice uint64) {
	switch p.Type {
	case 1:
		salePrice = amount - (amount * p.Value / 100)

	case 2:
		// 直减，订单总额减
		salePrice = p.Value
	default:
		salePrice = 0
	}
	return salePrice
}
