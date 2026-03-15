package i18n

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"strings"
)

//go:embed locales/*.json
var localeFS embed.FS

type I18n struct {
	translations map[string]map[string]string
	DefaultLang  string
}

func New(defaultLang string) (*I18n, error) {
	i18n := &I18n{
		translations: map[string]map[string]string{},
		DefaultLang:  defaultLang,
	}

	entries, err := localeFS.ReadDir("locales")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := localeFS.ReadFile("locales/" + entry.Name())
		if err != nil {
			return nil, err
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			return nil, err
		}

		lang := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		i18n.translations[lang] = translations
	}

	return i18n, nil
}

func (tr *I18n) T(lang, key string) string {
	if langTranslations, ok := tr.translations[lang]; ok {
		if value, ok := langTranslations[key]; ok {
			return value
		}
	}

	if defaultTranslations, ok := tr.translations[tr.DefaultLang]; ok {
		if value, ok := defaultTranslations[key]; ok {
			return value
		}
	}

	return key
}
