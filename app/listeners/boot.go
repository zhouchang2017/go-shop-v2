package listeners

import (
	"context"
	"go-shop-v2/app/email"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var listeners []rabbitmq.Listener
var adminRep *repositories.AdminRep

func ListenerInit() {
	listeners = []rabbitmq.Listener{
		OnOrderCreatedListener{}.Make(),            // 订单创建事件处理
		OrderPaidNotifyToAdmin{}.Make(),            // 订单付款通知管理员
		OrderPaidNotifyToUser{}.Make(),             // 订单付款通知客户
		OrderCloseNotifyToAdmin{}.Make(),           // 买家关闭订单通知管理员
		OrderClosedNotifyToUser{}.Make(),           // 管理员关闭订单通知用户
		OnOrderTimeOutListener{}.Make(),            // 订单超时自动关闭
		OnOrderApplyRefundListener{}.Make(),        // 订单申请退款
		OnOrderRefundChangeListener{}.Make(),       // 退款单状态改变
		OnOrderRefundCancelByUserListener{}.Make(), // 用户关闭退款
		//OnOrderRefundDoneListener{}.Make(),         // 退款完成事件
	}
}

func Boot(mq *rabbitmq.RabbitMQ) {
	for _, listener := range listeners {
		mq.Register(listener)
	}
}

func GetAppNotifications() []map[string]string {

	res := make([]map[string]string, 0)
	for _, listener := range listeners {

		if appNotifier, ok := listener.(AppNotifier); ok {

			res = append(res, map[string]string{
				"key":  utils.StructToName(listener.Event()),
				"name": appNotifier.Name(),
			})
		}

	}
	return res
}

type AppNotifier interface {
	Name() string
}

func resolveNotifierByEvent(ctx context.Context, event rabbitmq.Event) ([]email.Receiver, error) {
	if adminRep == nil {
		adminRep = repositories.MakeAdminRep()
	}
	receivers := make([]email.Receiver, 0)
	notifies, err := adminRep.FindByNotifies(ctx, []string{utils.StructToName(event)}, options.Find().SetProjection(bson.M{"email": 1, "nickname": 1}))
	if err != nil {
		return nil, err
	}
	for _, notify := range notifies {
		receivers = append(receivers, notify)
	}
	return receivers, nil
}

func sendEmail(event rabbitmq.Event, notify email.Notify) error {
	admins, err := resolveNotifierByEvent(context.Background(), event)
	if err == nil {
		if len(admins) > 0 {
			return email.Sends(notify, admins...)
		}
	}
	return nil
}
