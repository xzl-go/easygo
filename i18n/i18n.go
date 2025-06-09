package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/xzl-go/easygo/core"
)

// I18n 国际化管理器
type I18n struct {
	translations map[string]map[string]string
	defaultLang  string
}

// New 创建新的国际化管理器
func New(defaultLang string) *I18n {
	return &I18n{
		translations: make(map[string]map[string]string),
		defaultLang:  defaultLang,
	}
}

// LoadTranslations 加载翻译文件
func (i *I18n) LoadTranslations(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		lang := strings.TrimSuffix(filepath.Base(path), ".json")
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			return err
		}

		i.translations[lang] = translations
		return nil
	})
}

// Translate 获取翻译
func (i *I18n) Translate(key, lang string) string {
	if translations, ok := i.translations[lang]; ok {
		if translation, ok := translations[key]; ok {
			return translation
		}
	}
	if translations, ok := i.translations[i.defaultLang]; ok {
		if translation, ok := translations[key]; ok {
			return translation
		}
	}
	return key
}

// Middleware 创建国际化中间件
func (i *I18n) Middleware() core.HandlerFunc {
	return func(c *core.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = i.defaultLang
		}
		c.Set("lang", lang)
		c.Next()
	}
}
