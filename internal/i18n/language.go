package i18n

type Language struct {
	Code string
	Name string
}

var SupportedLanguages = []Language{
	{Code: "uk", Name: "Українська"},
	{Code: "en", Name: "English"},
	{Code: "pl", Name: "Polski"},
}

func IsSupported(code string) bool {
	for _, lang := range SupportedLanguages {
		if lang.Code == code {
			return true
		}
	}
	return false
}
