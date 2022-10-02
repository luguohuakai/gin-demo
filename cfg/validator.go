package cfg

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"time"
)

// InitValidator 初始化validator自定义验证验证方法
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("timing", timing)
		_ = v.RegisterValidation("topicUrl", topicUrl)
		_ = v.RegisterValidation("inArray", inArray)
	}
}

// 存在db的用户信息中创建时间与更新时间都要大于某一时间验证
func timing(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func topicUrl(fl validator.FieldLevel) bool {
	if url, ok := fl.Field().Interface().(string); ok {
		if matched, _ := regexp.MatchString(`\w{4,10}`, url); matched {
			return true
		}
	}
	return false
}

// 验证一个数组是否包含另一个数组
func inArray(fl validator.FieldLevel) bool {
	arr := strings.Split(fl.Param(), " ")
	//fmt.Println(fmt.Sprintf("%#v", arr))
	input := fl.Field().Interface()
	//fmt.Println(fmt.Sprintf("%#v", input))
	if val, ok := input.([]string); ok {
		for _, v := range val {
			for k, value := range arr {
				if v == value {
					break
				}
				if k == len(arr)-1 {
					return false
				}
			}
		}
	} else {
		return false
	}
	return true
}

// todo: 验证字符串中是否包含非16进制的字符
// todo: 数字必须能被5整除且不能为0
// todo: 必须包含 :; 顺序一定要对
// todo: 格式必须为 命令:[参数(可选)];
// todo: 如果存在参数参数格式必须为 K=V,K=V,K=V
