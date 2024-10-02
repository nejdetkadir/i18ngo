package i18ngo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"
)

type (
	Context struct {
		availableLocales []string
		locales          []Locale
		currentLocale    Locale
		options          Options
		onLocaleChange   func(locale string)
		mutex            sync.RWMutex
	}
	Locale struct {
		locale string
		data   map[string]interface{}
	}
	Options struct {
		DefaultLocale string
		Locales       []LocaleOptions
		Debug         bool
		Separator     string
	}
	LocaleOptions struct {
		Locale string
		File   []byte
	}
	I18nGo interface {
		T(path string, options ...interface{}) string
		ChangeLocale(locale string) error
		CurrentLocale() string
	}
)

func New(options Options) (I18nGo, error) {
	var Locales []Locale
	var currentLocale Locale
	var availableLocales []string

	for _, locale := range options.Locales {
		var localeData map[string]interface{}
		if err := json.Unmarshal(locale.File, &localeData); err != nil {
			if options.Debug {
				log.Println(fmt.Sprintf("[DEBUG] [New] error unmarshalling JSON for locale '%s': %v", locale.Locale, err))
			}

			return nil, errors.New(fmt.Sprintf("error unmarshalling JSON for locale '%s': %v", locale.Locale, err))
		}

		Locales = append(Locales, Locale{
			locale: locale.Locale,
			data:   localeData,
		})

		if slices.Contains(availableLocales, locale.Locale) {
			if options.Debug {
				log.Println(fmt.Sprintf("[DEBUG] [New] locale '%s' is defined more than once", locale.Locale))
			}

			return nil, errors.New(fmt.Sprintf("locale '%s' is defined more than once", locale.Locale))
		}

		if locale.Locale == options.DefaultLocale {
			currentLocale = Locale{
				locale: locale.Locale,
				data:   localeData,
			}
		}

		availableLocales = append(availableLocales, locale.Locale)
	}

	if !slices.Contains(availableLocales, options.DefaultLocale) {
		if options.Debug {
			log.Println(fmt.Sprintf("[DEBUG] [New] default locale '%s' is not defined", options.DefaultLocale))
		}

		return nil, errors.New(fmt.Sprintf("default locale '%s' is not defined", options.DefaultLocale))
	}

	if options.Debug {
		log.Println(fmt.Sprintf("[DEBUG] [New] default locale is '%s'", options.DefaultLocale))
		log.Println(fmt.Sprintf("[DEBUG] [New] available locales are '%v'", availableLocales))
	}

	if options.Separator == "" {
		options.Separator = "."
	}

	ctx := Context{
		availableLocales: availableLocales,
		locales:          Locales,
		options:          options,
		currentLocale:    currentLocale,
	}

	return &ctx, nil
}

func (ctx *Context) T(path string, options ...interface{}) string {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()

	keys := strings.Split(path, ctx.options.Separator)
	requestedLocale := ctx.currentLocale.locale

	if len(options) > 0 {
		if optionsMap, ok := options[0].(*map[string]interface{}); ok {
			if scope, ok := (*optionsMap)["scope"].(string); ok {
				if scope != "" {
					scopeKeys := strings.Split(scope, ctx.options.Separator)
					keys = append(scopeKeys, keys...)
				}
			}

			if locale, ok := (*optionsMap)["locale"].(string); ok {
				requestedLocale = locale
			}
		}
	}

	currentKeyPath := strings.Join(keys, ctx.options.Separator)

	if len(keys) == 0 {
		if ctx.options.Debug {
			log.Println(fmt.Sprintf("[DEBUG] [T] invalid path '%s.%s'", requestedLocale, currentKeyPath))
		}

		return fmt.Sprintf("invalid path '%s'", currentKeyPath)
	}

	var current interface{} = ctx.currentLocale.data

	if requestedLocale != "" {
		for _, l := range ctx.locales {
			if l.locale == requestedLocale {
				current = l.data
				break
			}
		}
	}

	for _, key := range keys {
		if value, ok := current.(map[string]interface{})[key]; ok {
			current = value
		} else {
			if ctx.options.Debug {
				log.Println(fmt.Sprintf("[DEBUG] [T] missing translation for path '%s.%s'", requestedLocale, currentKeyPath))
			}

			return fmt.Sprintf("missing translation for path '%s.%s'", requestedLocale, currentKeyPath)
		}
	}

	val := fmt.Sprintf("%v", current)

	if len(options) > 0 {
		if optionsMap, ok := options[0].(*map[string]interface{}); ok {
			for key, value := range *optionsMap {
				val = strings.ReplaceAll(val, fmt.Sprintf("{{%s}}", key), fmt.Sprintf("%v", value))
			}
		}
	}

	return val
}

func (ctx *Context) ChangeLocale(locale string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	if !slices.Contains(ctx.availableLocales, locale) {
		if ctx.options.Debug {
			log.Println(fmt.Sprintf("[DEBUG] [changeLocale] locale '%s' is not available", locale))
		}

		return errors.New(fmt.Sprintf("locale '%s' is not available", locale))
	}

	if ctx.currentLocale.locale == locale {
		return nil
	}

	for _, l := range ctx.locales {
		if l.locale == locale {
			ctx.currentLocale = l

			if ctx.onLocaleChange != nil {
				ctx.onLocaleChange(locale)
			}

			if ctx.options.Debug {
				log.Println(fmt.Sprintf("[DEBUG] [changeLocale] locale changed to '%s'", locale))
			}

			return nil
		}
	}

	if ctx.options.Debug {
		log.Println(fmt.Sprintf("[DEBUG] [changeLocale] locale '%s' is not defined", locale))
	}

	return errors.New(fmt.Sprintf("locale '%s' is not defined", locale))
}

func (ctx *Context) CurrentLocale() string {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()

	return ctx.currentLocale.locale
}
