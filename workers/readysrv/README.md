# Ready Server
The package `core/workers/readysrv` is used to provide readiness checks for a service.

## Configuration
The readiness server can be configured through the environment to match setup in the infrastructure.

- `READINESS_INTERFACE` default: `:3674`
- `READINESS_PATH` default: `/ready`

## Examples

```go
srv := readysrv.New(readysrv.Checks{
    "google": readysrv.CheckerFunc(func() ([]string, bool) {
        if _, err := http.Get("https://google.com"); err != nil {
            return []string{err.Error()}, false
        }
        return []string{"google can be accessed"}, true
    }),
})
srv.Run(ctx, ioutil.Discard)
```