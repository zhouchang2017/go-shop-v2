package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/utils"
	"strconv"
	"time"
	"unicode/utf8"
)

const (
	RefundStatusApply     = iota // 订单申请退款
	RefundStatusAgreed           // 同意退款
	RefundStatusReject           // 拒绝退款
	RefundStatusRefunding        // 退款中
	RefundStatusDone             // 退款已完成
	RefundStatusClosed           // 退款关闭
)

var RefundStatus []struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Class  string `json:"class"`
	Active bool   `json:"active"`
} = []struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Class  string `json:"class"`
	Active bool   `json:"active"`
}{
	{Name: "申请中", Value: RefundStatusApply, Class: "bg-red-400", Active: true},
	{Name: "同意退款", Value: RefundStatusAgreed, Class: "bg-yellow-400", Active: true},
	{Name: "拒绝退款", Value: RefundStatusReject, Class: "bg-red-200",},
	{Name: "退款中", Value: RefundStatusRefunding, Class: "bg-green-300", Active: true},
	{Name: "退款完成", Value: RefundStatusDone, Class: "bg-green-400",},
	{Name: "退款关闭", Value: RefundStatusClosed, Class: "bg-gray-200",},
}

type RefundCanceler struct {
	Type   string      `json:"type"`
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Avatar qiniu.Image `json:"avatar"`
}

// 退款
type Refund struct {
	Id          string          `json:"id" bson:"id"`                     // 退款单号
	OrderNo     string          `json:"order_no" bson:"order_no"`         // 订单号
	PaymentNo   string          `json:"payment_no" bson:"payment_no"`     // 微信订单号
	ReturnCode  string          `json:"return_code" bson:"return_code"`   // 微信执行状态
	TotalAmount uint64          `json:"total_amount" bson:"total_amount"` // 退款金额
	RefundDesc  string          `json:"refund_desc" bson:"refund_desc"`   // 退款原因
	Status      int             `json:"status"`
	Items       []*RefundItem   `json:"items"` // todo: 这里要加bson标签吗？
	CreatedAt   time.Time       `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	RejectDesc  string          `json:"reject_desc,omitempty" bson:"reject_desc,omitempty"` // 拒绝原因
	Canceler    *RefundCanceler `json:"canceler"`                                           // 关闭订单操作者
	// todo: add refund type,退款退货，仅退款，仅退货
	// todo: add logistics no
}

func (r *Refund) FillItemsFromOrder(o *Order) {
	for _, item := range r.Items {
		findItem := o.FindItem(item.ItemId)
		if findItem != nil {
			item.Item = findItem.Item
		}
	}
}

func (o Refund) StatusText() string  {
	switch o.Status {
	case RefundStatusApply:
		return "退款申请中"
	case RefundStatusAgreed:
		return "同意退款"
	case RefundStatusReject:
		return "拒绝退款"
	case RefundStatusRefunding:
		return "退款中"
	case RefundStatusDone:
		return "退款已完成"
	default:
		return "退款关闭"
	}
}

// 获取订单支付金额
func (o Refund) GetActualAmount() string {
	amount := float64(o.TotalAmount) / 100
	float := strconv.FormatFloat(amount, 'f', 2, 64)
	return fmt.Sprintf("￥%s", float)
}

// 获取退款物品名称，用户小程序订阅消息
func (o Refund) GoodsName(limit int, order *Order) string {
	o.FillItemsFromOrder(order)
	if len(o.Items) > 1 {
		name := o.Items[0].Item.Product.Name
		if limit == -1 {
			return fmt.Sprintf("%s(等商品)", name)
		}
		if utf8.RuneCountInString(name) > limit-8 {
			subString := utils.SubString(name, 0, limit-8)
			return fmt.Sprintf("%s...(等商品)", subString)
		}
		return fmt.Sprintf("%s(等商品)", name)
	}
	name := o.Items[0].Item.Product.Name
	if limit == -1 {
		return fmt.Sprintf("%s", name)
	}
	if utf8.RuneCountInString(name) > limit-3 {
		subString := utils.SubString(name, 0, limit-3)
		return fmt.Sprintf("%s...", subString)
	}
	return fmt.Sprintf("%s", name)
}

// 合计数量
func (r Refund) ItemCount() (count uint64) {
	for _, item := range r.Items {
		count += item.Qty
	}
	return
}

// 检测是否能关闭退款
func (r Refund) CanCancel() bool {
	if r.Status == RefundStatusApply || r.Status == RefundStatusAgreed {
		return true
	}
	return false
}

// 失败的退款记录
type FailedRefund struct {
	model.MongoModel `inline`
	RefundOn         string `json:"refund_on" bson:"refund_on"`
	OrderNo          string `json:"order_no" bson:"order_no"`
	ErrCode          string `json:"err_code" bson:"err_code"`
	ErrCodeDes       string `json:"err_code_des" bson:"err_code_des"`
}

func newRefund(order *Order, items ...*RefundItem) *Refund {
	r := &Refund{
		Id:         utils.RandomRefundOrderNo(""),
		OrderNo:    order.OrderNo,
		PaymentNo:  order.Payment.PaymentNo,
		ReturnCode: "",
		Status:     RefundStatusApply,
		Items:      items,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, amount := r.TotalQtyAndAmount()
	r.TotalAmount = amount
	return r
}

func (r Refund) FindItem(itemId string) (item *RefundItem, err error) {
	for _, i := range r.Items {
		if i.ItemId == itemId {
			return i, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r Refund) TotalQtyAndAmount() (qty uint64, amount uint64) {
	for _, item := range r.Items {
		qty += item.Qty
		amount += item.TotalAmount()
	}
	return
}

type RefundItem struct {
	ItemId string          `json:"item_id" bson:"item_id" form:"item_id"` // 退款
	Qty    uint64          `json:"qty"`                                   // 数量
	Amount uint64          `json:"amount"`                                // 退款金额
	Item   *AssociatedItem `json:"-" bson:"-"`
}

func (r RefundItem) TotalAmount() (amount uint64) {
	return r.Amount * r.Qty
}

// 退款选项结构
type RefundOption struct {
	Items []*RefundItem `json:"items"`
	Desc  string        `json:"desc"`
}

func (r RefundOption) isValid(order *Order) error {
	if len(r.Items) == 0 {
		return err2.Err422.F("缺少退款商品")
	}
	for _, item := range r.Items {
		if item.ItemId == "" {
			return err2.Err422.F("缺少退款商品id")
		}
		if item.Qty == 0 {
			return err2.Err422.F("参数异常，item.qty 不能为0")
		}
		orderItem := order.FindItem(item.ItemId)
		if item.Amount > uint64(orderItem.Amount) {
			return err2.Err422.F("退款金额不能大于付款进")
		}
		refunds := Refunds(order.Refunds)
		filterRefunds := refunds.FilterByStatus(RefundStatusApply, RefundStatusAgreed, RefundStatusRefunding, RefundStatusDone)
		refundQty, refundAmount := filterRefunds.CountItemQtyAndAmount(item.ItemId)
		if uint64(orderItem.Count)-refundQty < item.Qty {
			return err2.Err422.F("退款商品[%s]数量溢出，总共付款商品数量%d,已申请退款商品数量%d,当前申请退款数量%d", item.ItemId, orderItem.Count, refundQty, item.Qty)
		}
		if uint64(orderItem.TotalAmount())-refundAmount < item.TotalAmount() {
			return err2.Err422.F("退款商品[%s]退款溢出，总共付款金额%d,已申请退款金额%d,当前申请退款金额%d", item.ItemId, orderItem.TotalAmount(), refundAmount, item.TotalAmount())
		}
	}
	return nil
}

type Refunds []*Refund

func (r Refunds) FilterByStatus(status ...int) (res Refunds) {
	res = make([]*Refund, 0)
	for _, refund := range r {
		for _, s := range status {
			if refund.Status == s {
				res = append(res, refund)
				continue
			}
		}

	}
	return res
}

func (r Refunds) CountItemQtyAndAmount(itemId string) (qty uint64, amount uint64) {
	for _, refund := range r {
		if item, err := refund.FindItem(itemId); err == nil {
			qty += item.Qty
			amount += item.TotalAmount()
		}
	}
	return
}

func (r Refunds) CountQtyAndAmount() (qty uint64, amount uint64) {
	for _, refund := range r {
		q, a := refund.TotalQtyAndAmount()
		qty += q
		amount += a
	}
	return
}

func (o *Order) FindRefundByItemId(itemId string) (res Refunds) {
	res = make([]*Refund, 0)
	for _, refund := range o.Refunds {
		for _, item := range refund.Items {
			if itemId == item.ItemId {
				res = append(res, refund)
				continue
			}
		}
	}
	return res
}

// 申请退款
func (o *Order) MakeRefund(opt RefundOption) (*Refund, error) {
	if err := opt.isValid(o); err != nil {
		return nil, err
	}
	refund := newRefund(o, opt.Items...)
	refund.RefundDesc = opt.Desc
	o.Refunds = append(o.Refunds, refund)
	return refund, nil
}

// 搜索退款单
func (o *Order) FindRefund(refundId string) (refund *Refund, err error) {
	for _, r := range o.Refunds {
		if r.Id == refundId {
			return r, nil
		}
	}
	return nil, err2.Err404.F("refund no[%s],not exist!!", refundId)
}

// 刷新订单状态
func (o *Order) RefreshStatus() {
	// 如果全部退款完成，标记为交易成功
	// 取消退款，不影响订单状态
	if len(o.Refunds) > 0 {
		refunds := Refunds(o.Refunds)
		status := refunds.FilterByStatus(RefundStatusDone)
		qty, _ := status.CountQtyAndAmount()
		if uint64(o.ItemCount) == qty {
			// 该笔订单全部退款
			o.Status = OrderStatusDone // 交易成功
			o.CommentChannel = false   // 关闭评论通道
		}
	}
}
