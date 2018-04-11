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

	postprocessResponse(resp)
}

func prepareRequest(req *fasthttp.Request) {
	// do not proxy "Connection" header.
	req.Header.Del("Connection")
	// strip other unneeded headers.

	// alter other request params before sending them to upstream host
	req.Header.SetHost(proxyAddr)

	req.SetRequestURI("http://www.baidu.com")
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")

	// strip other unneeded headers

	// alter other response data if needed
	// resp.Header.Set("Access-Control-Allow-Origin", "*")
	// resp.Header.Set("Access-Control-Request-Method", "OPTIONS,HEAD,POST")
	// resp.Header.Set("Content-Type", "application/json; charset=utf-8")
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
	fmt.Println("start server...")
}