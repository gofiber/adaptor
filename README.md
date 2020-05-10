### Proxy
Converter for net/http handlers to/from Fiber request handlers

### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/adaptor
```

### Signature
```go
adaptor.NewFiberHandler(h http.Handler) func(*fiber.Ctx)
adaptor.NewHTTPHandlerFiber(handler func(*fiber.Ctx)) http.Handler
adaptor.NewHTTPHandlerFasthttp(handler fasthttp.RequestHandler) http.Handler
```

### Functions
| Name | Signature | Description
| :--- | :--- | :---
| NewFiberHandler | `NewFiberHandler(h http.Handler) func(*fiber.Ctx)` | net/http handler to Fiber handler wrapper
| NewHTTPHandlerFiber | `NewHTTPHandlerFiber(handler func(*fiber.Ctx)) http.Handler` | Fiber handler to net/http handler wrapper
| NewHTTPHandlerFasthttp | `NewHTTPHandlerFasthttp(handler fasthttp.RequestHandler) http.Handler` | fasthttp handler to net/http handler wrapper

### net/http to Fiber
```go
package main

import (
    "fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/adaptor"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
    app := fiber.New()

    app.Get("/", adaptor.NewFiberHandlerFunc(greet))

    app.Listen(8080)
}
```

### Fiber to net/http
```go
package main

import (
    "net/http"
	"github.com/gofiber/fiber"
	"github.com/gofiber/adaptor"
)

func main() {
	http.Handle("/", adaptor.NewHTTPHandlerFiber(func (c *fiber.Ctx) {
        c.SendString("Hello World!")
    }))
	http.ListenAndServe(":8080", nil)
}
```

### fasthttp to net/http
```go
package main

import (
    "net/http"
	"github.com/valyala/fasthttp"
	"github.com/gofiber/adaptor"
)

func main() {
    http.Handle("/", adaptor.NewHTTPHandlerFasthttp(func(ctx *RequestCtx) {
        ctx.SetBodyString("Hello World!")
    }))
	http.ListenAndServe(":8080", nil)
}
```