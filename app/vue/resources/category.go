package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
	"net/http"
)

func init() {
	register(NewCategoryResource)
}

type Category struct {
	core.AbstractResource
	model interface{}
	rep   *repositories.CategoryRep
}

func (c *Category) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	panic("implement me")
}

func (c *Category) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	result := <-c.rep.FindById(ctx, id)
	return result.Result, result.Error
}

func (c *Category) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	results := <-c.rep.Pagination(ctx, req)
	return results.Result, results.Pagination, results.Error
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


			panels.NewPanel("销售属性",
				fields.NewTable("销售属性", "Options", func() []contracts.Field {
					return []contracts.Field{
						fields.NewTextField("名称", "Name"),
						fields.NewTextField("权重", "Sort"),
						fields.NewTextField("属性值", "Values"),
					}
				}),
			).SetWithoutPending(true),
		}
	}
}

func NewCategoryResource(rep *repositories.CategoryRep) *Category {
	return &Category{model: &models.Category{}, rep: rep}
}

type categoryForm struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func (c *Category) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	form := &categoryForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
	}
	category := model.(*models.Category)
	category.Name = form.Name
	return category, nil
}

func (c *Category) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &categoryForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
	}
	return models.NewCategory(form.Name), nil
}

type OptionForm struct {
	Name   string             `json:"name" form:"name"`
	Sort   int64              `json:"sort" form:"sort"`
	Values []*optionValueForm `json:"values" form:"values"`
}

type optionValueForm struct {
	Code  string `json:"code" form:"code"`
	Value string `json:"value" form:"value"`
}

// 添加销售属性处理函数
func (this *Category) addOption(router gin.IRouter, uri string, singularLabel string) {
	router.POST(fmt.Sprintf("%s/:%s/options", uri, singularLabel), func(c *gin.Context) {
		form := OptionForm{}
		err := c.ShouldBind(&form)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		option := models.NewProductOption(form.Name)
		option.Sort = form.Sort
		var values []*models.OptionValue
		for _, value := range form.Values {
			values = append(values, option.NewValue(value.Value, value.Code))
		}
		option.AddValues(values...)
		err = this.rep.AddOption(c, c.Param(singularLabel), option)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": option})
	})
}

// 更新销售属性处理函数
func (this *Category) updateOption(router gin.IRouter, uri string, singularLabel string) {
	router.PUT(fmt.Sprintf("%s/:%s/options/:optionId", uri, singularLabel), func(c *gin.Context) {
		form := OptionForm{}
		err := c.ShouldBind(&form)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		option := models.MakeProductOption(c.Param("optionId"), form.Name, form.Sort)
		values := []*models.OptionValue{}
		for _, value := range form.Values {
			values = append(values, option.NewValue(value.Value, value.Code))
		}
		option.AddValues(values...)
		this.rep.UpdateOption(c, c.Param(singularLabel), option)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": option})
	})
}

// 删除销售属性处理函数
func (this *Category) deleteOption(router gin.IRouter, uri string, singularLabel string) {
	router.DELETE(fmt.Sprintf("%s/:%s/options/:optionId", uri, singularLabel), func(c *gin.Context) {
		err := this.rep.DeleteOption(c, c.Param(singularLabel), c.Param("optionId"))
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusNoContent, nil)
	})
}

func (this *Category) CustomHttpRouters(router gin.IRouter, uri string, singularLabel string) {
	// ADD OPTION
	this.addOption(router, uri, singularLabel)
	// UPDATE OPTION
	this.updateOption(router, uri, singularLabel)
	// DELETE OPTION
	this.deleteOption(router, uri, singularLabel)
}

func (c *Category) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
}

func (c *Category) Model() interface{} {
	return c.model
}

func (c *Category) Repository() repository.IRepository {
	return c.rep
}

func (c Category) Make(model interface{}) contracts.Resource {
	return &Category{
		rep:   c.rep,
		model: model,
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
