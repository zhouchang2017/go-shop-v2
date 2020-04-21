package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

type Category struct {
	core.AbstractResource
	model   interface{}
	service *services.CategoryService
}

// 创建方法
func (c *Category) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	option := services.CategoryCreateOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}
	category, err := c.service.Create(ctx, option)
	if err != nil {
		return "", err
	}

	return core.CreatedRedirect(c, category.GetID()), nil
}

// 更新方法
func (c *Category) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	option := services.CategoryCreateOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}
	category, err := c.service.Update(ctx, model.(*models.Category), option)
	if err != nil {
		return "", err
	}

	return core.UpdatedRedirect(c, category.GetID()), nil
}

// 删除方法
func (c *Category) Destroy(ctx *gin.Context, id string) (err error) {
	return c.service.Delete(ctx, id)
}

// 详情方法
func (c *Category) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return c.service.FindById(ctx, id)
}

// 列表页
func (c *Category) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return c.service.Pagination(ctx, req)
}

func (c *Category) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (c *Category) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (c *Category) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (c *Category) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (c *Category) Policy() interface{} {
	return nil
}

func (c *Category) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("名称", "Name"),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),
		}
	}
}

func NewCategoryResource() *Category {
	return &Category{model: &models.Category{}, service: services.MakeCategoryService()}
}

//
//// 添加销售属性处理函数
//func (this *Category) addOption(router gin.IRouter, uri string, singularLabel string) {
//	router.POST(fmt.Sprintf("%s/:%s/options", uri, singularLabel), func(c *gin.Context) {
//		form := OptionForm{}
//		err := c.ShouldBind(&form)
//		if err != nil {
//			err2.ErrorEncoder(nil, err, c.Writer)
//			return
//		}
//		option := models.NewProductOption(form.Name)
//		option.Sort = form.Sort
//		var values []*models.OptionValue
//		for _, value := range form.Values {
//			values = append(values, option.NewValue(value.Name, value.Code))
//		}
//		option.AddValues(values...)
//		err = this.rep.AddOption(c, c.Param(singularLabel), option)
//		if err != nil {
//			err2.ErrorEncoder(nil, err, c.Writer)
//			return
//		}
//		c.JSON(http.StatusCreated, gin.H{"data": option})
//	})
//}
//
//// 更新销售属性处理函数
//func (this *Category) updateOption(router gin.IRouter, uri string, singularLabel string) {
//	router.PUT(fmt.Sprintf("%s/:%s/options/:optionId", uri, singularLabel), func(c *gin.Context) {
//		form := OptionForm{}
//		err := c.ShouldBind(&form)
//		if err != nil {
//			err2.ErrorEncoder(nil, err, c.Writer)
//			return
//		}
//		option := models.MakeProductOption(c.Param("optionId"), form.Name, form.Sort)
//		values := []*models.OptionValue{}
//		for _, value := range form.Values {
//			values = append(values, option.NewValue(value.Name, value.Code))
//		}
//		option.AddValues(values...)
//		this.rep.UpdateOption(c, c.Param(singularLabel), option)
//		if err != nil {
//			err2.ErrorEncoder(nil, err, c.Writer)
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{"data": option})
//	})
//}
//
//// 删除销售属性处理函数
//func (this *Category) deleteOption(router gin.IRouter, uri string, singularLabel string) {
//	router.DELETE(fmt.Sprintf("%s/:%s/options/:optionId", uri, singularLabel), func(c *gin.Context) {
//		err := this.rep.DeleteOption(c, c.Param(singularLabel), c.Param("optionId"))
//		if err != nil {
//			err2.ErrorEncoder(nil, err, c.Writer)
//			return
//		}
//		c.JSON(http.StatusNoContent, nil)
//	})
//}

func (c *Category) Model() interface{} {
	return c.model
}

func (c Category) Make(model interface{}) contracts.Resource {
	return &Category{
		service: c.service,
		model:   model,
	}
}

func (c *Category) SetModel(model interface{}) {
	c.model = model
}

func (c Category) Title() string {
	return "产品类目"
}

// 左侧导航栏分组
func (Category) Group() string {
	return "Product"
}

func (Category) Icon() string {
	return "icons-grid"
}
