package adaptor

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gofiber/fiber"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// NewFiberHandler wraps net/http handler to fiber handler
func NewFiberHandler(h http.Handler) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		handler := fasthttpadaptor.NewFastHTTPHandler(h)
		handler(c.Fasthttp)
	}
}

// NewHTTPHandlerFiber wraps fiber handler to net/http handler
func NewHTTPHandlerFiber(handler func(*fiber.Ctx)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initFasthttpContext(w, r)
		handler(&fiber.Ctx{
			Fasthttp: ctx,
		})
		writeFasthttpResponse(&ctx.Response, w)
	})
}

// NewHTTPHandlerFasthttp wraps fasthttp RequestHandler to net/http handler
func NewHTTPHandlerFasthttp(handler fasthttp.RequestHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initFasthttpContext(w, r)
		handler(ctx)
		writeFasthttpResponse(&ctx.Response, w)
	})
}

func initFasthttpContext(w http.ResponseWriter, r *http.Request) *fasthttp.RequestCtx {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return &ctx
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
		return &ctx
	}
	ctx.Init(&req, remoteAddr, nil)

	return &ctx
}

func writeFasthttpResponse(resp *fasthttp.Response, w http.ResponseWriter) {
	w.WriteHeader(resp.StatusCode())
	resp.Header.VisitAll(func(k, v []byte) {
		sk := string(k)
		sv := string(v)
		w.Header().Add(sk, sv)
	})
	w.Write(resp.Body())
}
