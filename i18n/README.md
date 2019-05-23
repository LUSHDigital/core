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
