### Adaptor
Converter for net/http handlers to/from Fiber request handlers

### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/adaptor
```

### Signature
```go
adaptor.FiberHandler(h http.Handler) func(*fiber.Ctx)
adaptor.FiberHandlerFunc(h http.HandlerFunc) func(*fiber.Ctx)
adaptor.HTTPHandler(h func(*fiber.Ctx)) http.Handler
adaptor.HTTPHandlerFunc(handler func(*fiber.Ctx)) http.HandlerFunc
```

### Functions
| Name | Signature | Description
| :--- | :--- | :---
| FiberHandler | `FiberHandler(h http.Handler) func(*fiber.Ctx)` | net/http handler to Fiber handler wrapper
| FiberHandlerFunc | `FiberHandlerFunc(h http.HandlerFunc) func(*fiber.Ctx)` | net/http handler func to Fiber handler wrapper
| HTTPHandler | `HTTPHandler(h func(*fiber.Ctx)) http.Handler` | Fiber handler to net/http handler wrapper
| HTTPHandlerFunc | `HTTPHandlerFunc(h func(*fiber.Ctx)) http.HandlerFunc` | Fiber handler to net/http handler func wrapper

### net/http -> Fiber
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gofiber/adaptor"
	"github.com/gofiber/fiber"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	app := fiber.New()

	app.Get("/", adaptor.NewFiberHandlerFunc(greet))

	app.Listen(3000)
}

```

### Fiber -> net/http
```go
package main

import (
	"net/http"

	"github.com/gofiber/adaptor"
	"github.com/gofiber/fiber"
)

func main() {
	http.Handle("/", adaptor.NewHTTPHandlerFiber(func(c *fiber.Ctx) {
		c.Send("Hello World!")
	}))

	http.ListenAndServe(":3000", nil)
}

```
