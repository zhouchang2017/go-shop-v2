package models

// 促销计划计算处理

type promotions []*Promotion

func (p promotions) findById(id string) *Promotion {
	for _, promotion := range p {
		if promotion.GetID() == id {
			return promotion
		}
	}
	return nil
}

// 全排列
func (p promotions) permute() []promotions {
	return permute(p)
}

func permute(input promotions) []promotions {
	return subNumberSlice(input)
}

func subNumberSlice(nums promotions) []promotions {
	if len(nums) == 0 {
		return nil
	}
	if len(nums) == 1 {
		return []promotions{{nums[0]}}
	}
	if len(nums) == 2 {
		return []promotions{{nums[0], nums[1]}, {nums[1], nums[0]}}
	}

	result := []promotions{}
	for index, value := range nums {
		var numsCopy = make(promotions, len(nums))
		copy(numsCopy, nums)
		numsSubOne := append(numsCopy[:index], numsCopy[index+1:]...)
		valueSlice := promotions{value}
		newSubSlice := subNumberSlice(numsSubOne)
		for _, newValue := range newSubSlice {
			result = append(result, append(valueSlice, newValue...))
		}
	}
	return result
}

type PromotionOrderItem struct {
	ItemId              string
	ProductId           string
	Qty                 uint64
	UnMutexPromotions   []*PromotionItem
	MutexPromotion      *PromotionItem
	promotions          []*PromotionItem
	MutexPromotionId    *string
	UnMutexPromotionIds []string
	Price               uint64
}

type PromotionInfo struct {
	Promotion *Promotion
	SalePrice uint64
}

type PromotionResult struct {
	SalePrices uint64
	Infos      []*PromotionInfo
}

type PromotionOverView struct {
	SalePrices uint64                    `json:"sale_prices"` // 优惠总额
	Infos      []*PromotionOverViewItem `json:"infos"`       // 优惠项目
}

type PromotionOverViewItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SalePrice   uint64  `json:"sale_price" bson:"sale_price"`
}

// 获取优惠总览
func (p *PromotionResult) Overview() *PromotionOverView {
	if p.Infos == nil {
		return nil
	}
	res := &PromotionOverView{}
	res.SalePrices = p.SalePrices
	infos := make([]*PromotionOverViewItem, 0)
	for _, item := range p.Infos {
		infos = append(infos, &PromotionOverViewItem{
			Id:          item.Promotion.GetID(),
			Name:        item.Promotion.Name,
			Description: item.Promotion.GetDescription(),
			SalePrice:   item.SalePrice,
		})
	}
	res.Infos = infos
	return res
}

// 获取优惠明细
func (p *PromotionResult) Detail() ItemPromotionInfos {
	if p.Infos == nil {
		return nil
	}
	res := ItemPromotionInfos{}
	for _, info := range p.Infos {
		for _, item := range info.Promotion.Items {
			for _, i := range item.items {
				find := res.FindByItemId(i.itemId)
				newItem := &ItemPromotion{
					Id:            info.Promotion.GetID(),
					Name:          info.Promotion.Name,
					Description:   info.Promotion.GetDescription(),
					SalePrice:     i.salePrice,
					UnitSalePrice: i.unitSalePrice,
				}
				if find != nil {
					find.Infos = append(find.Infos, newItem)
					find.SalePrices += i.salePrice
					find.UnitSalePrices += i.unitSalePrice
				} else {
					res = append(res, &ItemPromotionInfo{
						ItemId:         i.itemId,
						Infos:          []*ItemPromotion{newItem},
						SalePrices:     i.salePrice,
						UnitSalePrices: i.unitSalePrice,
					})
				}

			}
		}
	}
	return res
}

type ItemPromotionInfo struct {
	ItemId         string           `json:"item_id" bson:"-"`
	Infos          []*ItemPromotion `json:"infos"`                                    // 优惠明细
	SalePrices     uint64            `json:"sale_prices" bson:"sale_prices"`           // 总优惠金额
	UnitSalePrices uint64            `json:"unit_sale_prices" bson:"unit_sale_prices"` // 单件优惠金额
}

type ItemPromotionInfos []*ItemPromotionInfo

func (i ItemPromotionInfos) FindByItemId(id string) *ItemPromotionInfo {
	for _, item := range i {
		if item.ItemId == id {
			return item
		}
	}
	return nil
}

// 产品优惠明细
type ItemPromotion struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	SalePrice     uint64  `json:"sale_price" bson:"sale_price"`
	UnitSalePrice uint64  `json:"unit_sale_price" bson:"unit_sale_price"`
}

func promotionMax(arr []*PromotionResult) *PromotionResult {
	if len(arr) == 0 {
		return nil
	}
	max := arr[0]
	for _, item := range arr {
		if item.SalePrices > max.SalePrices {
			max = item
		}
	}
	return max
}

func PromotionCalculate(items ...*PromotionOrderItem) *PromotionResult {
	promotions := make(promotions, 0)
	for _, item := range items {
		if item.MutexPromotion != nil {
			item.promotions = append(item.promotions, item.MutexPromotion)
		}
		item.promotions = append(item.promotions, item.UnMutexPromotions...)
		for _, promotion := range item.promotions {
			find := promotions.findById(promotion.Promotion.Id)
			if find != nil {
				find.addItem(promotion, item.ItemId, item.Price, item.Qty)
			} else {
				newPromotion := promotion.Promotion.makePromotion()
				newPromotion.addItem(promotion, item.ItemId, item.Price, item.Qty)
				promotions = append(promotions, newPromotion)
			}
		}
	}
	var promotionResults []*PromotionResult
	if len(promotions) == 0 {
		return &PromotionResult{
			SalePrices: 0,
			Infos:      nil,
		}
	}
	permute := promotions.permute()
	for _, items := range permute {
		var salePrices uint64 = 0   //优惠总额
		var infos []*PromotionInfo // 优惠信息
		for _, item := range items {
			if promotionItem, salePrice, ok := item.calculate(salePrices); ok {
				salePrices += salePrice
				infos = append(infos, &PromotionInfo{
					Promotion: promotionItem,
					SalePrice: salePrice,
				})
			}
		}
		promotionResults = append(promotionResults, &PromotionResult{
			SalePrices: salePrices,
			Infos:      infos,
		})
	}
	return promotionMax(promotionResults)
}
