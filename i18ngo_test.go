package i18ngo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var enData = []byte(`{
	"pages": {
		"login": {
			"buttons": {
				"login": "Login"
			}
		}
	},
	"common": {
		"greetings": {
			"hello": "Hello"
		}
	}
}`)

var trData = []byte(`{
	"pages": {
		"login": {
			"buttons": {
				"login": "Giriş"
			}
		}
	},
	"common": {
		"greetings": {
			"hello": "Merhaba"
		}
	}
}`)

func TestNewI18nInitialization(t *testing.T) {
	i18n, err := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	assert.NoError(t, err, "Expected no error during initialization")

	ctx := i18n.(*Context)

	assert.Equal(t, "en", ctx.currentLocale.locale, "Expected default locale to be 'en'")
	assert.Len(t, ctx.locales, 2, "Expected two locales to be loaded")
	assert.Contains(t, ctx.availableLocales, "en", "Expected availableLocales to include 'en'")
	assert.Contains(t, ctx.availableLocales, "tr", "Expected availableLocales to include 'tr'")
}

func TestTranslate(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	assert.Equal(t, "Login", i18n.T("pages.login.buttons.login"), "Expected 'Login' in English locale")

	_ = i18n.ChangeLocale("tr")
	assert.Equal(t, "Giriş", i18n.T("pages.login.buttons.login"), "Expected 'Giriş' in Turkish locale")
}

func TestMissingTranslation(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         true,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
		},
	})

	missingTranslation := i18n.T("pages.login.buttons.logout")
	expectedMessage := "missing translation for path 'en.pages.login.buttons.logout'"
	assert.Equal(t, expectedMessage, missingTranslation, "Expected missing translation message")
}

func TestTranslateWithParams(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: []byte(`{
				"pages": {
					"login": {
						"welcome": "Welcome, {{name}}!"
					}
				}
			}`)},
		},
	})

	params := map[string]interface{}{
		"name": "John",
	}

	translated := i18n.T("pages.login.welcome", &params)
	assert.Equal(t, "Welcome, John!", translated, "Expected 'Welcome, John!' with parameter replacement")
}

func TestTranslateWithScope(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	params := map[string]interface{}{
		"scope": "common.greetings",
	}

	translated := i18n.T("hello", &params)
	assert.Equal(t, "Hello", translated, "Expected 'Hello' from the 'common.greetings' scope in English")

	_ = i18n.ChangeLocale("tr")
	translated = i18n.T("hello", &params)
	assert.Equal(t, "Merhaba", translated, "Expected 'Merhaba' from the 'common.greetings' scope in Turkish")
}

func TestTranslateWithLocaleParam(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	params := map[string]interface{}{
		"locale": "tr",
	}

	translated := i18n.T("pages.login.buttons.login", &params)
	assert.Equal(t, "Giriş", translated, "Expected 'Giriş' translation when locale param is 'tr'")

	params["locale"] = "en"
	translated = i18n.T("pages.login.buttons.login", &params)
	assert.Equal(t, "Login", translated, "Expected 'Login' translation when locale param is 'en'")
}

func TestTranslateWithScopeAndLocale(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	params := map[string]interface{}{
		"scope":  "common.greetings",
		"locale": "tr",
	}

	translated := i18n.T("hello", &params)
	assert.Equal(t, "Merhaba", translated, "Expected 'Merhaba' from the 'common.greetings' scope in Turkish")

	params["locale"] = "en"
	translated = i18n.T("hello", &params)
	assert.Equal(t, "Hello", translated, "Expected 'Hello' from the 'common.greetings' scope in English")
}

func TestChangeLocale(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	ctx := i18n.(*Context)

	assert.Equal(t, "en", ctx.currentLocale.locale, "Expected current locale to be 'en'")

	err := i18n.ChangeLocale("tr")
	assert.NoError(t, err, "Expected no error when changing locale to 'tr'")
	assert.Equal(t, "tr", ctx.currentLocale.locale, "Expected current locale to be 'tr'")
}

func TestInvalidLocaleChange(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	err := i18n.ChangeLocale("fr")
	assert.Error(t, err, "Expected error when switching to an undefined locale")
	assert.EqualError(t, err, "locale 'fr' is not available", "Expected specific error message for undefined locale")
}

func TestSeparatorHandling(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Separator:     "/",
		Locales: []LocaleOptions{
			{Locale: "en", File: []byte(`{
				"pages": {
					"login": {
						"buttons": {
							"login": "Login"
						}
					}
				}
			}`)},
		},
	})

	assert.Equal(t, "Login", i18n.T("pages/login/buttons/login"), "Expected translation with custom separator")
}

func TestCurrentLocale(t *testing.T) {
	i18n, _ := New(Options{
		DefaultLocale: "en",
		Debug:         false,
		Locales: []LocaleOptions{
			{Locale: "en", File: enData},
			{Locale: "tr", File: trData},
		},
	})

	ctx := i18n.(*Context)

	assert.Equal(t, "en", ctx.CurrentLocale(), "Expected initial locale to be 'en'")

	_ = i18n.ChangeLocale("tr")
	assert.Equal(t, "tr", ctx.CurrentLocale(), "Expected current locale to be 'tr'")
}
