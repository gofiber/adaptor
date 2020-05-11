package adaptor

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gofiber/fiber"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// FiberHandlerFunc wraps net/http handler func to fiber handler
func HTTPHandlerFunc(h http.HandlerFunc) func(*fiber.Ctx) {
	return HTTPHandler(h)
}

// FiberHandler wraps net/http handler to fiber handler
func HTTPHandler(h http.Handler) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		handler := fasthttpadaptor.NewFastHTTPHandler(h)
		handler(c.Fasthttp)
	}
}

// HTTPHandler wraps fiber handler to net/http handler
func FiberHandler(h func(*fiber.Ctx)) http.Handler {
	return FiberHandlerFunc(h)
}

// HTTPHandlerFunc wraps fiber handler to net/http handler func
func FiberHandlerFunc(handler func(*fiber.Ctx)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx fasthttp.RequestCtx
		var req fasthttp.Request

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		req.Header.SetMethod(r.Method)
		req.SetRequestURI(r.RequestURI)
		req.Header.SetContentLength(len(body))
		req.SetHost(r.Host)
		for key, val := range r.Header {
			for _, v := range val {
				req.Header.Add(key, v)
			}
		}
		req.BodyWriter().Write(body)
		remoteAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		if err != nil {
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		ctx.Init(&req, remoteAddr, nil)

		resp := &ctx.Response
		w.WriteHeader(resp.StatusCode())
		resp.Header.VisitAll(func(k, v []byte) {
			sk := string(k)
			sv := string(v)
			w.Header().Add(sk, sv)
		})
		w.Write(resp.Body())
	})
}
