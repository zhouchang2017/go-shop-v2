package vue

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	err2 "go-shop-v2/pkg/err"
	"net/http"
	"reflect"
)

func (this *ResourceWarp) resourceCreateHandle(router gin.IRouter) {
	if creatable, ok := this.resource.(ResourceHttpCreate); ok && creatable.ResourceHttpCreate() {
		router.POST(this.UriKey(), func(c *gin.Context) {
			// 验证权限
			if !this.AuthorizedToCreate(c) {
				c.AbortWithStatus(403)
				return
			}

			// 表单验证
			fields, _ := this.resolveCreationFields(c)
			rules := make(map[string][]string)
			messages := make(map[string][]string)
			for _, field := range fields {
				rules[field.GetAttribute()] = []string{}
				messages[field.GetAttribute()] = []string{}
				for _, rule := range field.GetRules() {
					rules[field.GetAttribute()] = append(rules[field.GetAttribute()], rule.Rule)
					if rule.Message != "" {
						messages[field.GetAttribute()] = append(messages[field.GetAttribute()], fmt.Sprintf("%s:%s", rule.Rule, rule.Message))
					}
				}
			}

			var req map[string]interface{}

			opts := govalidator.Options{
				Request:         c.Request, // request object
				Rules:           rules,     // rules map
				Messages:        messages,  // custom message map (Optional)
				Data:            &req,
				RequiredDefault: true, // all the field to be pass the rules
			}

			v := govalidator.New(opts)

			if err := v.ValidateJSON(); len(err) > 0 {
				spew.Dump(err)
				errMessage := err2.NewFromCode(422).Data(err)
				err2.ErrorEncoder(nil, errMessage, c.Writer)
				return
			}
			spew.Dump(req)
			// 注入值
			i := reflect.ValueOf(this.resource.Model()).Elem().Type()
			model := reflect.New(i).Interface()
			for _, field := range fields {
				field.Fill(c, req, model)
			}
			//spew.Dump(fields)
			panic("test!")
			// 资源处理表单
			entity, err := creatable.CreateFormParse(c)
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			results := <-this.resource.Repository().Create(c, entity)
			if results.Error != nil {
				err2.ErrorEncoder(nil, results.Error, c.Writer)
				return
			}
			// created hook
			go this.resource.Created(c, results.Result)

			c.JSON(http.StatusCreated, gin.H{"id": results.Id})
		})
	}

}
