# Env
The `core/env` package provides functionality for ensuring we retrieve an environment variable

## Example

### Loading default environment config

```go
env.TryLoadDefault()
```

### Loading default environment config together with other source

```go
env.TryLoadDefault("package/other.env")
```

### Expect environment variable to be set

```go
dbURL := env.MustGet("DATABASE_URL")
```

