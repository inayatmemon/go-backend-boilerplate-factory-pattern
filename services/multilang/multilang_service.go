package multilang_service

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	lang_constants "go_boilerplate_project/constants/lang"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func (s *service) GetMessage(lang, key string, params map[string]string) string {
	if s.Input.Bundle == nil {
		return key
	}
	if lang == "" {
		lang = s.Input.DefaultLang
	}
	lang = strings.ToLower(strings.TrimSpace(lang))
	if len(lang) > 2 {
		lang = lang[:2]
	}
	accept := lang
	localizer := i18n.NewLocalizer(s.Input.Bundle, lang, accept)
	templateData := make(map[string]interface{})
	for k, v := range params {
		switch k {
		case "field":
			templateData["Field"] = v
		case "param":
			templateData["Param"] = v
		case "detail":
			templateData["Detail"] = v
		default:
			templateData[k] = v
		}
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:   key,
		TemplateData: templateData,
	})
	if err != nil {
		s.Input.Logger.Debugw("Multilang message not found, using key as fallback",
			"key", key,
			"lang", lang,
			"error", err,
		)
		return key
	}
	return msg
}

func (s *service) GetLanguage(c *gin.Context) string {
	if c == nil {
		return s.Input.DefaultLang
	}
	lang := c.GetHeader(lang_constants.HeaderLanguage)
	if lang == "" {
		lang = c.GetHeader("Accept-Language")
		if idx := strings.Index(lang, ","); idx > 0 {
			lang = strings.TrimSpace(lang[:idx])
		}
		if idx := strings.Index(lang, "-"); idx > 0 {
			lang = lang[:idx]
		}
		if idx := strings.Index(lang, ";"); idx > 0 {
			lang = strings.TrimSpace(lang[:idx])
		}
	}
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return s.Input.DefaultLang
	}
	if len(lang) > 2 {
		lang = lang[:2]
	}
	return lang
}

// InitBundle creates and loads the i18n bundle from the languages directory.
func InitBundle(languagesPath string, defaultLang string) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.Make(defaultLang))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	languages := []string{"en", "es", "fr"}
	for _, lang := range languages {
		path := filepath.Join(languagesPath, fmt.Sprintf("%s.json", lang))
		if _, err := bundle.LoadMessageFile(path); err != nil {
			return nil, fmt.Errorf("load language file %s: %w", path, err)
		}
	}
	return bundle, nil
}
