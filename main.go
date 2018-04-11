package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"
	"fmt"
)

var (
	proxyAddr   string
	proxyClient = &fasthttp.HostClient{
		IsTLS: false,
		Addr:  "api.yourdomain.com",

		// set other options here if required - most notably timeouts.
		// ReadTimeout: 60, // 如果在生产环境启用会出现多次请求现象
	}
)

func ReverseProxyHandler(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	resp := &ctx.Response

	prepareRequest(req)

	if err := proxyClient.Do(req, resp); err != nil {
		ctx.Logger().Printf("error when proxying the request: %s", err)
	}

	postprocessResponse(ctx,resp)
}

func prepareRequest(req *fasthttp.Request) {
	// do not proxy "Connection" header.
	req.Header.Del("Connection")
	// strip other unneeded headers.

	// alter other request params before sending them to upstream host
	req.Header.SetHost(proxyAddr)

	req.SetRequestURI("http://www.baidu.com")
}

func postprocessResponse(ctx *fasthttp.RequestCtx,resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")

	// strip other unneeded headers

	// alter other response data if needed
	// resp.Header.Set("Access-Control-Allow-Origin", "*")
	// resp.Header.Set("Access-Control-Request-Method", "OPTIONS,HEAD,POST")
	// resp.Header.Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}

func main() {
	port := flag.String("port", "8082", "listen port")
	targetAddr := flag.String("target", "www.baidu.com", "your server domain")
	flag.Parse()

	proxyClient.Addr = *targetAddr

	log.Println("port:", *port)
	log.Println("target:", *targetAddr)

	if err := fasthttp.ListenAndServe("localhost:"+*port, ReverseProxyHandler); err != nil {
		log.Fatalf("error in fasthttp server: %s", err)
	}
}