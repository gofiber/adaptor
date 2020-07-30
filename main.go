// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package adaptor

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gofiber/fiber"
	"github.com/gofiber/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// HTTPHandlerFunc wraps net/http handler func to fiber handler
func HTTPHandlerFunc(h http.HandlerFunc) func(*fiber.Ctx) {
	return HTTPHandler(h)
}

// HTTPHandler wraps net/http handler to fiber handler
func HTTPHandler(h http.Handler) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		handler := fasthttpadaptor.NewFastHTTPHandler(h)
		handler(c.Fasthttp)
	}
}

// FiberHandler wraps fiber handler to net/http handler
func FiberHandler(h func(*fiber.Ctx)) http.Handler {
	return FiberHandlerFunc(h)
}

// FiberHandlerFunc wraps fiber handler to net/http handler func
func FiberHandlerFunc(h fiber.Handler) http.HandlerFunc {
	return handlerFunc(fiber.New(), h)
}

// FiberApp wraps fiber app to net/http handler func
func FiberApp(app *fiber.App) http.HandlerFunc {
	return handlerFunc(app)
}

func handlerFunc(app *fiber.App, h ...fiber.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// New fasthttp request
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		// Convert net/http -> fasthttp request
		if r.Body != nil {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, utils.StatusMessage(fiber.StatusInternalServerError), fiber.StatusInternalServerError)
				return
			}
			req.Header.SetContentLength(len(body))
			_, _ = req.BodyWriter().Write(body)
		}
		req.Header.SetMethod(r.Method)
		req.SetRequestURI(r.RequestURI)
		req.SetHost(r.Host)
		for key, val := range r.Header {
			for _, v := range val {
				req.Header.Add(key, v)
			}
		}
		remoteAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		if err != nil {
			http.Error(w, utils.StatusMessage(fiber.StatusInternalServerError), fiber.StatusInternalServerError)
			return
		}

		// New fasthttp Ctx
		var fctx fasthttp.RequestCtx
		fctx.Init(req, remoteAddr, nil)
		if len(h) > 0 {
			// New fiber Ctx
			ctx := app.AcquireCtx(&fctx)
			defer app.ReleaseCtx(ctx)
			// Execute fiber Ctx
			h[0](ctx)
		} else {
			// Execute fasthttp Ctx though app.Handler
			app.Handler()(&fctx)
		}

		// Convert fasthttp Ctx > net/http
		fctx.Response.Header.VisitAll(func(k, v []byte) {
			w.Header().Set(string(k), string(v))
		})
		w.WriteHeader(fctx.Response.StatusCode())
		_, _ = w.Write(fctx.Response.Body())
	}
}
