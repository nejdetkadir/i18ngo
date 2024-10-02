![Build and test](https://github.com/nejdetkadir/i18ngo/actions/workflows/main.yml/badge.svg?branch=main)
![Go Version](https://img.shields.io/badge/go_version-_1.23.1-007d9c.svg)


![cover](docs/cover.png)

# I18nGo

I18nGo is a simple internationalization library for Golang that enables translation and locale management. It provides utilities to change locales, handle translations with scoped paths, and supports dynamic translation replacement using templates.

## Features
- Multiple Locales: Load and manage multiple locales for translations.
- Change Locales Dynamically: Easily switch between different languages.
- Translation with Scoped Paths: Fetch translations using dot-separated paths.
- Dynamic Value Replacement: Replace placeholders in translations dynamically with provided values.
- Debugging: Enable debugging to log details about missing translations and locale management.
- Customizable Separator: Use a custom separator for translation paths.

## Installation
To install I18nGo, use the following:

```bash
go get github.com/nejdetkadir/i18ngo
```

## Usage
### Initialize I18nGo
First, create and initialize I18nGo with options like default locale and available locales.

```go
package main

import (
    "github.com/nejdetkadir/i18ngo"
    "os"
)

func main() {
    enData, _ := os.ReadFile("foo/bar/en.json")
    trData, _ := os.ReadFile("foo/bar/tr.json")
	
    /*
        json files should be like:
        {
            "pages": {
                "login": {
                    "buttons": {
                        "login": "Login"
                    }
                }
            }
        }
   */

    i18n, err := i18ngo.New(i18ngo.Options{
        DefaultLocale: "en",
        Debug:         true,
        Locales: []i18ngo.LocaleOptions{
            {Locale: "en", File: enData},
            {Locale: "tr", File: trData},
        },
    })

    if err != nil {
        panic(err)
    }

    // Use the i18ngo instance for translations
    translated := i18n.T("pages.login.buttons.login")
    println(translated) // Output: Login
}
```

### Changing Locale
Change the locale dynamically using ChangeLocale.

```go
i18n.ChangeLocale("tr")
translated := i18n.T("pages.login.buttons.login")
println(translated) // Output: Giriş
```

### Dynamic Value Replacement
You can pass dynamic values for translation using the T function.

```go
/*
    json files should be like:
    {
        "pages": {
            "welcome": "Welcome, {{name}}!"
        }
    }
*/

params := map[string]interface{}{"name": "John"}
translated := i18n.T("pages.welcome", &params)
println(translated) // Output: Welcome, John!
```

### Scoped Translations
Use a scope to define a context for your translation paths.

```go
/*
    json files should be like:
    {
        "common": {
            "greetings": {
                "hello": "Hello"
            }
        }
    }
*/

params := map[string]interface{}{"scope": "common.greetings"}
translated := i18n.T("hello", &params)
println(translated) // Output: Hello
```

## Example
Check out the [example](example/main.go) for a complete demonstration of I18nGo usage.

## Unit Testing
The package includes comprehensive unit tests to ensure correct behavior across various scenarios. To run the tests, use the following:

```b![cover.png](../../../Downloads/cover.png)ash
go test ./...
```

## Contributing
Bug reports and pull requests are welcome on GitHub at https://github.com/nejdetkadir/i18ngo. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [code of conduct](https://github.com/nejdetkadir/i18ngo/blob/main/CODE_OF_CONDUCT.md).

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.