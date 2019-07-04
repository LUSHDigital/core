# Internationalisation

The `core/i18n` package functions for dealing with internationalisation of services.

## Examples

### Parse locale from

```go
locale, _ := i18n.ParseLocale("EN_gb")
fmt.Println(locale)
// Output: "en-GB"
```

### Put locale through context

Setting the locale in a context.

```go
ctx = i18n.ContextWithLocale(context.Background(), "sv")
```

Retreiving a locale from context.

```go
locale = i18n.LocaleFromContext(ctx)
```

### Change the default locale

You can set the default locale to be any value you'd like. By default it's set to `en`.

```go
import "github.com/LUSHDigital/core/i18n"

func main() {
    i18n.DefaultLocale = "es"
}
```
