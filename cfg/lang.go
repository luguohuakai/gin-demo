package cfg

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	esTranslations "github.com/go-playground/validator/v10/translations/es"
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	zhTwTranslations "github.com/go-playground/validator/v10/translations/zh_tw"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

// Trans 定义一个全局翻译器T
var Trans ut.Translator

func InitLang() {
	locale := viper.GetString("app.lang")
	if locale == "" {
		locale = "zh"
	}
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		esT := es.New()
		jaT := ja.New()
		zhTw := zh_Hant_TW.New()

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT, esT, jaT, zhTw)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			fmt.Println(fmt.Sprintf("uni.GetTranslator(%s) failed", locale))
		} else {
			//fmt.Println(fmt.Sprintf("uni.GetTranslator(%s) succeed", locale))
		}

		// 添加额外翻译
		_ = v.RegisterTranslation("required_with", Trans, func(ut ut.Translator) error {
			return ut.Add("required_with", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_with", fe.Field())
			return t
		})

		// 注册翻译器
		var err error
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		case "es":
			err = esTranslations.RegisterDefaultTranslations(v, Trans)
		case "ja":
			err = jaTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh_tw":
			err = zhTwTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("switch locale error: %s", err))
		}
	}
}
