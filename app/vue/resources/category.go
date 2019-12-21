package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
	"net/http"
)

func init() {
	register(NewCategory)
}

type Category struct {
	vue.AbstractResource
	model *models.Category
	rep   *repositories.CategoryRep
}

func NewCategory(model *models.Category, rep *repositories.CategoryRep) *Category {
	return &Category{model: model, rep: rep}
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

func (c *Category) IndexQuery(ctx *gin.Context, request *request.IndexRequest) {

}

func (c *Category) Model() interface{} {
	return c.model
}

func (c *Category) Repository() repository.IRepository {
	return c.rep
}

func (c Category) Make(model interface{}) vue.Resource {
	return &Category{model: model.(*models.Category)}
}

func (c *Category) SetModel(model interface{}) {
	c.model = model.(*models.Category)
}

func (c Category) Title() string {
	return "产品类目"
}

func (Category) Group() string {
	return "Product"
}

func (Category) Icon() string {
	return "i-grid"
}
