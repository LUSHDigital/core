# Env
The `core/env` package provides functionality for ensuring we retrieve an environment variable

## Example

```go
import "github.com/LUSHDigital/core/env"

func main() {
    dbURL := env.MustGet("DATABASE_URL")
    ...
}
```