# Response

The response package provides ways to create HTTP + JSON responses in a consistent format.

Below are usage example which demonstrate it's use.

## Respond without data

```go
package handlers

import (
	"net/http"
	
	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/response"
)

func someHandler(w http.ResponseWriter, r *http.Request) {
    resp := &response.Response{
        Code:    http.StatusOK,
        Message: "some helpful message",
    }
    resp.WriteTo(w)
}
```

*Output*

```json
{
  "code": 200,
  "message": "some helpful message"
}
```

## Respond with data
```go
package handlers

import (
	"net/http"
	
	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/response"
)

func someHandler(w http.ResponseWriter, r *http.Request) {
    someData := map[string]interface{}{
        "hello": "world",
    }
    
    resp := &response.Response{
        Code:    http.StatusOK,
        Message: "some helpful message",
        Data: &response.Data{
            Type:    "something",
            Content: someData,
        },
    }
    resp.WriteTo(w)
}
```

*Output*

```json
{
  "code": 200,
  "message": "some helpful message",
  "data": {
    "something": {
      "hello": "world"
    }
  }
}
```

## Respond with data and pagination
```go
package handlers

import (
	"net/http"
	
	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/response"
)

func someHandler(w http.ResponseWriter, r *http.Request) {
    someData := map[string]interface{}{
        "hello": "world",
    }
    
    preq := pagination.Request{
        PerPage: 10,
        Page:    1,
    }
    paginator := pagination.MakeResponse(preq, 100)
    resp := &response.Response{
        Code:    http.StatusOK,
        Message: "some helpful message",
        Data: &response.Data{
            Type:    "something",
            Content: someData,
        },
        Pagination: &paginator,
    }
    resp.WriteTo(w)
}
```

*Output*

```json
{
  "code": 200,
  "message": "some helpful message",
  "data": {
    "something": {
      "hello": "world"
    }
  },
  "pagination": {
    "per_page": 10,
    "page": 1,
    "offset": 0,
    "total": 100,
    "last_page": 10
  }
}
```