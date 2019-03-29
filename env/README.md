# Env
The `core/env` package provides functionality for ensuring we retrieve an environment variable

## Example

```go
    dbURL := env.MustGet("DATABASE_URL")
```