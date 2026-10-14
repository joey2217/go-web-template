// Package validator 负责参数校验器的初始化与校验错误的中文翻译。
package validator

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	once       sync.Once
	translator ut.Translator
)

// Init 将 Gin 内置校验器接入中文翻译。可重复调用，仅首次生效。
func Init() {
	once.Do(func() {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return
		}
		locale := zh.New()
		uni := ut.New(locale, locale)
		trans, found := uni.GetTranslator("zh")
		if !found {
			return
		}
		if err := zhtranslations.RegisterDefaultTranslations(v, trans); err != nil {
			return
		}
		translator = trans
	})
}

// Translate 将校验错误转换为可读的中文提示；非校验错误返回原始信息。
func Translate(err error) string {
	if err == nil {
		return ""
	}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		msgs := make([]string, 0, len(verrs))
		for _, fe := range verrs {
			if translator != nil {
				msgs = append(msgs, fe.Translate(translator))
				continue
			}
			msgs = append(msgs, fallbackMessage(fe))
		}
		return strings.Join(msgs, "; ")
	}
	return err.Error()
}

// fallbackMessage 在翻译器不可用时给出兜底提示。
func fallbackMessage(fe validator.FieldError) string {
	return fmt.Sprintf("字段 %s 校验失败（%s）", fe.Field(), fe.Tag())
}
