package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
)

type manualInventoryActionStatus int8

func (m manualInventoryActionStatus) Make(status int8) (res manualInventoryActionStatus) {
	switch status {
	case int8(MANUAL_FINISHED):
		res = MANUAL_FINISHED
	case int8(MANUAL_CANCELED):
		res = MANUAL_CANCELED
	default:
		res = MANUAL_SAVED
	}
	return res
}

const (
	MANUAL_SAVED    manualInventoryActionStatus = iota // 未提交
	MANUAL_FINISHED                                    // 完成
	MANUAL_CANCELED                                    // 取消
)

type manualInventoryActionType int8

func (m manualInventoryActionType) Make(t int8) (res manualInventoryActionType, err error) {
	switch t {
	case int8(MANUAL_TYPE_PUT):
		return MANUAL_TYPE_PUT, nil
	case int8(MANUAL_TYPE_TAKE):
		return MANUAL_TYPE_TAKE, nil
	}
	return res, fmt.Errorf("type %d not allow in [%d,%d]", t, MANUAL_TYPE_PUT, MANUAL_TYPE_TAKE)
}

const (
	MANUAL_TYPE_PUT  manualInventoryActionType = iota // 入库
	MANUAL_TYPE_TAKE                                  // 出库
)

// 标准库存操作
type ManualInventoryAction struct {
	model.MongoModel `inline`
	Type             manualInventoryActionType    `json:"type"`
	Shop             *AssociatedShop              `json:"shop" bson:"shop"` // 门店
	Items            []*ManualInventoryActionItem `json:"items"`
	User             *AssociatedAdmin             `json:"user"`
	Status           manualInventoryActionStatus  `json:"status"`
}

func (this *ManualInventoryAction) Types() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":  "入库",
			"value": MANUAL_TYPE_PUT,
		},
		{
			"name":   "出库",
			"value": MANUAL_TYPE_TAKE,
		},
	}
}

func (this *ManualInventoryAction) SetType(t int8) (error) {
	res, err := this.Type.Make(t)
	if err != nil {
		return err
	}
	this.Type = res
	return nil
}

func (this *ManualInventoryAction) SetStatus(status int8) *ManualInventoryAction {
	this.Status = this.Status.Make(status)
	return this
}

func (this *ManualInventoryAction) SetStatusToSaved() {
	this.SetStatus(0)
}

// 标准库存操作子项
type ManualInventoryActionItem struct {
	AssociatedItem `inline`
	Qty            int64           `json:"qty"`
	Status         InventoryStatus `json:"status"`
}

func (this *ManualInventoryActionItem) SetStatus(status int8) {
	this.Status = this.Status.Make(status)
}
