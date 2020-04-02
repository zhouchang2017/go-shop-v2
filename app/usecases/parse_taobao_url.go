package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/utils"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ParseUrlOption struct {
	Id           string                `json:"id"`
	Images       []qiniu.Image         `json:"images"`
	Description  string                `json:"description"`
	OptionValues []*models.OptionValue `json:"option_values" form:"option_values"`
}

type ParseUrlResponse struct {
	Images       []qiniu.Image         `json:"images"`
	Description  string                `json:"description"`
	OptionValues []*models.OptionValue `json:"option_values"`
}

// 正则匹配富文本img src地址
var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)

const parseTaobaoUrlCacheKey = "parser-taobao-url"

func getParseTaobaoCacheKey(id string) string {
	return fmt.Sprintf("%s-%s", parseTaobaoUrlCacheKey, id)
}
func ParseTaobaoUrl(ctx context.Context, form *ParseUrlOption) (res ParseUrlResponse) {
	if redis.GetConFn() != nil {
		result, err := redis.GetConFn().Get(getParseTaobaoCacheKey(form.Id)).Result()
		if err == nil {
			if result != "" {
				res = ParseUrlResponse{}
				if err := json.Unmarshal([]byte(result), &res); err == nil {
					return res
				}
			}
		}
	}
	images := fillImages(ctx, form.Images)
	description := fillDescription(ctx, form.Description)
	optionValues := fillOptionValues(ctx, form.OptionValues)

	res = ParseUrlResponse{
		Images:       images,
		Description:  description,
		OptionValues: optionValues,
	}

	if redis.GetConFn() != nil {
		if marshal, err := json.Marshal(res); err == nil {
			redis.GetConFn().Set(getParseTaobaoCacheKey(form.Id), marshal, time.Hour*24).Result()
		}
	}
	return res
}

func fillImages(ctx context.Context, images []qiniu.Image) (imgs []qiniu.Image) {
	var g errgroup.Group
	var mu sync.Mutex
	qiniuService := qiniu.GetQiniu()
	imgs = make([]qiniu.Image, len(images))
	sem := make(chan struct{}, 10)
	for index, img := range images {
		index, img := index, img
		sem <- struct{}{}
		g.Go(func() error {
			res, err := qiniuService.PutByUrl(ctx, img.Src(), utils.RandomString(32))
			if err == nil {
				mu.Lock()
				imgs[index] = res
				mu.Unlock()
			}
			<-sem
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return imgs
	}
	return imgs

}

func fillOptionValues(ctx context.Context, optionValues []*models.OptionValue) (response []*models.OptionValue) {
	var g errgroup.Group
	var mu sync.Mutex
	qiniuService := qiniu.GetQiniu()
	response = make([]*models.OptionValue, 0)
	sem := make(chan struct{}, 10)
	for _, value := range optionValues {
		value := value
		sem <- struct{}{}
		g.Go(func() error {
			if value.Image.IsUrl() {
				res, err := qiniuService.PutByUrl(ctx, value.Image.Src(), utils.RandomString(32))
				if err == nil {
					mu.Lock()
					response = append(response, &models.OptionValue{
						Id:    value.Id,
						Name:  value.Name,
						Image: &res,
					})
					mu.Unlock()
				}
				return nil
			}
			<-sem
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return response
	}

	return response
}

func fillDescription(ctx context.Context, description string) string {
	submatch := imgRE.FindAllStringSubmatch(description, -1)
	var g errgroup.Group
	qiniuService := qiniu.GetQiniu()
	sem := make(chan struct{}, 10)
	var syncMap sync.Map
	for _, match := range submatch {
		match := match
		sem <- struct{}{}
		g.Go(func() error {
			if len(match) >= 2 {
				res, err := qiniuService.PutByUrl(ctx, match[1], utils.RandomString(32))
				spew.Dump(res, err)
				if err == nil {
					syncMap.Store(match[1], res.Src())
				}
			}
			<-sem
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return description
	}
	syncMap.Range(func(key, value interface{}) bool {
		description = strings.ReplaceAll(description, key.(string), value.(string))
		return true
	})
	return description
}
