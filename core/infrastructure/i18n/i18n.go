package i18n

import "fmt"

type I18NProvider = func(key string) string

var i18nProviders []I18NProvider

// Provide 绑定I18N提供者
func Provide(p I18NProvider) {
	if p != nil {
		i18nProviders = append(i18nProviders, p)
	}
}

// T 获取翻译后的文字
func T(text string) string {
	for _, f := range i18nProviders {
		text = f(text)
	}
	return text
}

// Errorf 返回错误
func Errorf(format string, v ...interface{}) error {
	return fmt.Errorf(T(format), v...)
}
