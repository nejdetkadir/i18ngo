package main

import (
	"fmt"
	"github.com/nejdetkadir/i18ngo"
	"log"
	"os"
)

func main() {
	pwd, err := os.Getwd()
	enJSONFile, err := os.ReadFile(fmt.Sprintf("%s/example/en.json", pwd))

	if err != nil {
		log.Println("something went wrong while reading en.json file")

		panic(err)
	}

	trJSONFile, err := os.ReadFile(fmt.Sprintf("%s/example/tr.json", pwd))

	if err != nil {
		log.Println("something went wrong while reading tr.json file")

		panic(err)
	}

	i18n, err := i18ngo.New(i18ngo.Options{
		DefaultLocale: "en",
		Locales: []i18ngo.LocaleOptions{
			{
				Locale: "en",
				File:   enJSONFile,
			},
			{
				Locale: "tr",
				File:   trJSONFile,
			},
		},
	})

	if err != nil {
		log.Println("something went wrong while creating i18n instance")

		panic(err)
	}

	// with default locale (en)
	log.Println(fmt.Sprintf("current locale: %s", i18n.CurrentLocale()))
	log.Println(i18n.T("pages.home.title"))

	// try to change locale to turkish (tr)
	err = i18n.ChangeLocale("tr")

	// if there is an error while changing locale
	if err != nil {
		log.Println("something went wrong while changing locale")

		panic(err)
	}

	// after changing locale to turkish (tr)
	log.Println(fmt.Sprintf("current locale: %s", i18n.CurrentLocale()))
	log.Println(i18n.T("pages.home.title"))

	// with variables
	log.Println(i18n.T("pages.home.welcome", &map[string]interface{}{"name": "John"}))

	// with scope
	log.Println(i18n.T("home.title", &map[string]interface{}{"scope": "pages"}))
	log.Println(i18n.T("title", &map[string]interface{}{"scope": "pages.home"}))

	// with specific locale
	log.Println(i18n.T("pages.home.title", &map[string]interface{}{"locale": "en"}))

	// with missing key
	log.Println(i18n.T("pages.home.missing"))
}
