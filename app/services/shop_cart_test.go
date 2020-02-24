package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/qiniu"
	"testing"
)

func TestShopCartService_Add(t *testing.T) {
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

	itemService := MakeItemService()

	item, err := itemService.FindById(context.Background(), "5e51e253ecbe820cbd5f6d80")
	if err != nil {
		t.Fatal(err)
	}

	shopCartService := MakeShopCartService()

	addItem, err := shopCartService.Add(context.Background(), "2", item, 1, true)

	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(addItem)
}

func TestShopCartService_Update(t *testing.T) {
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

	shopCartService := MakeShopCartService()

	update, err := shopCartService.Update(context.Background(), "5e521bec6dd2f323e7379ae8", 5, false)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(update)
}

func TestShopCartService_Delete(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	shopCartService := MakeShopCartService()

	force := ctx.WithForce(context.Background(), true)
	err := shopCartService.Delete(force, "5e52205cebf9fce4371999ec", "5e523ca749a3a822faefca0f")

	if err!=nil {
		t.Fatal(err)
	}
}
