package validator

import (
	"errors"
	"reflect"
	"sync"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	trans    ut.Translator
	once     sync.Once
)

// getIns 获取单例对象
func getIns() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			return field.Tag.Get("label")
		})
		zhCn := zh.New()
		uni := ut.New(zhCn)
		trans, _ = uni.GetTranslator("zh")
		_ = zhTrans.RegisterDefaultTranslations(validate, trans)
	})
	return validate
}

// Check 判断是否符合规则
func Check(obj interface{}) (bool, string) {
	err := getIns().Struct(obj)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if !errors.As(err, &invalidValidationError) {
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				for _, err2 := range errs {
					return false, err2.Translate(trans)
				}
			}
		}
	}
	return true, ""
}
