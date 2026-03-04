package multilang_service

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type Repository interface {
	// GetMessage returns translated message for key. Params: "field", "param", etc.
	GetMessage(lang, key string, params map[string]string) string
	// GetLanguage extracts language from gin context (header) or returns default.
	GetLanguage(c *gin.Context) string
}

type service struct {
	Input Input
}

type Input struct {
	Logger      *zap.SugaredLogger
	DefaultLang string
	Bundle      *i18n.Bundle
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for multilang service: %v", err)
	}
	input.Logger.Infow("Multilang service initialized (go-i18n)",
		"defaultLang", input.DefaultLang,
		"bundleLoaded", input.Bundle != nil,
	)
	return &service{Input: input}
}

func (i *Input) validateInput() error {
	if i == nil {
		return errors.New("input is nil for multilang service")
	}
	if i.Logger == nil {
		return errors.New("logger is nil for multilang service")
	}
	if i.DefaultLang == "" {
		return errors.New("defaultLang is required for multilang service")
	}
	return nil
}
