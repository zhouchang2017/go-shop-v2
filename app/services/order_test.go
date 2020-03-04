package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/tests"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestOrderService_Create(t *testing.T) {

	// preview
	mongodb.TestConnect()
	defer mongodb.Close()
	// todo: add item data into db to provide test
	// init
	authUser := tests.GenerateUser()
	// test case
	situations := []struct {
		name string
		param *OrderCreateOption
		wantErr bool
	}{
		{
			"normal case",
			&OrderCreateOption{
				UserAddress:  orderUserAddress{
					Id:           "123",
					ContactName:  "张三",
					ContactPhone: "13800138000",
					Province:     "广东省",
					City:         "深圳市",
					Areas:        "南山区",
					Addr:         "科苑天桥下",
				},
				TakeGoodType: models.OrderTakeGoodTypeOnline,
				OrderItems:   []*orderItem{
					{
						ItemId: "5e51e253ecbe820cbd5f6d77",
						Count:  1,
						Price:  109000,
					},
					{
						ItemId: "5e51e253ecbe820cbd5f6d80",
						Count:  2,
						Price:  109000,
					},
				},
				OrderAmount:  327000,
				ActualAmount: 327000,
			},
			false,
		},
		// todo: add more test case
	}
	for _, situation := range situations {
		t.Run(situation.name, func(t *testing.T) {
			// run
			srv := MakeOrderService()
			gotOrder, err := srv.Create(context.Background(), authUser, situation.param)
			// valid case
			if (err != nil) != situation.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, situation.wantErr)
				return
			}
			if !situation.wantErr {
				assert.NotEqual(t, "", gotOrder.ID)
				assert.NotEqual(t, "", gotOrder.OrderNo)
				assert.Equal(t, models.OrderStatusPrePay, gotOrder.Status)
				assert.Equal(t, len(situation.param.OrderItems), len(gotOrder.OrderItems))
				assert.Equal(t, situation.param.UserAddress.Id, gotOrder.UserAddress.Id)
				assert.Equal(t, situation.param.UserAddress.ContactName, gotOrder.UserAddress.ContactName)
				assert.Equal(t, situation.param.UserAddress.ContactPhone, gotOrder.UserAddress.ContactPhone)
				assert.Equal(t, situation.param.TakeGoodType, gotOrder.TakeGoodType)
				assert.Equal(t, situation.param.OrderAmount, gotOrder.OrderAmount)
				assert.Equal(t, situation.param.ActualAmount, gotOrder.ActualAmount)
			}
		})
	}
}