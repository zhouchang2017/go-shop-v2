package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
)

// 表单验证
func Validator(ctx *gin.Context, fields []contracts.Field) (req map[string]interface{}, err error) {
	rules := make(map[string][]string)
	messages := make(map[string][]string)
	for _, field := range fields {
		rules[field.GetAttribute()] = []string{}
		messages[field.GetAttribute()] = []string{}
		for _, rule := range field.GetRules() {
			rules[field.GetAttribute()] = append(rules[field.GetAttribute()], rule.GetRule())
			if rule.GetMessage() != "" {
				messages[field.GetAttribute()] = append(messages[field.GetAttribute()], fmt.Sprintf("%s:%s", rule.GetRule(), rule.GetMessage()))
			}
		}
	}

	req = map[string]interface{}{}

	opts := govalidator.Options{
		Request:         ctx.Request, // request object
		Rules:           rules,       // rules map
		Messages:        messages,    // custom message map (Optional)
		Data:            &req,
		RequiredDefault: true, // all the field to be pass the rules
	}

	v := govalidator.New(opts)

	if err := v.ValidateJSON(); len(err) > 0 {
		return req, err2.NewFromCode(422).Data(err)
	}
	return req, nil
}

func ResourceUpdateValidator()  {

}