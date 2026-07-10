package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

type Localizer = i18n.Localizer

type Service struct {
	bundle *i18n.Bundle
}

func NewService() *Service {
	bundle := i18n.NewBundle(language.Ukrainian)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	return new(Service{bundle: bundle})
}

func (s *Service) LoadTranslations() error {
	for _, path := range []string{"locales/uk.json", "locales/en.json", "locales/pl.json"} {
		if _, err := s.bundle.LoadMessageFileFS(localeFS, path); err != nil {
			return fmt.Errorf("i18n: loading %s: %w", path, err)
		}
	}
	return nil
}

func (s *Service) FromLanguage(lang string) *Localizer {
	return i18n.NewLocalizer(s.bundle, lang)
}

func T(l *Localizer, messageID string, data ...map[string]any) string {
	cfg := new(i18n.LocalizeConfig{MessageID: messageID})
	if len(data) > 0 {
		cfg.TemplateData = data[0]
	}

	msg, err := l.Localize(cfg)
	if err != nil {
		return messageID
	}
	return msg
}
