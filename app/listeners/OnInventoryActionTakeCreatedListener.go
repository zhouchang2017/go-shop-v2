package listeners

import (
	"context"
	"encoding/json"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/message"
	"log"
)

// 超时关闭出货单
type TimeOutCloseInventoryAction struct {
	actionSrv *services.ManualInventoryActionService
}

func NewTimeOutCloseInventoryAction() *TimeOutCloseInventoryAction {
	return &TimeOutCloseInventoryAction{actionSrv: services.MakeManualInventoryActionService()}
}

// 对应的触发事件
func (TimeOutCloseInventoryAction) Event() message.Event {
	return events.InventoryActionTakeCreated{}
}

// 队列名称
func (TimeOutCloseInventoryAction) QueueName() string {
	return "TimeOutCloseInventoryAction"
}

// 错误处理
func (TimeOutCloseInventoryAction) OnError(err error) {
	log.Printf("TimeOutCloseInventoryAction Error:%s\n", err)
}

// 处理逻辑
func (this TimeOutCloseInventoryAction) Handler(data []byte) error {
	var action models.ManualInventoryAction
	err := json.Unmarshal(data, &action)
	if err != nil {
		return err
	}
	// 已取消、已完成操作不处理
	if action.IsCancel() || action.IsFinished() {
		return nil
	}
	return this.actionSrv.Cancel(context.Background(), action.GetID())
}
