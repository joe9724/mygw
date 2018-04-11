package main

import (
	"flag"
	"log"
	"github.com/valyala/fasthttp"
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
	log.Println(ctx, "Hello, world! Requested path is %q", string(ctx.Path()))
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

	req.SetRequestURI("http://widsboat-api.bitekun.xin/joe9724/data_manage/1.0.0//device/list?operator_id=2&page=0&size=12")
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")
	resp.SkipBody = false
	resp.AppendBody([]byte("abc"))

	// strip other unneeded headers

	// alter other response data if needed
	// resp.Header.Set("Access-Control-Allow-Origin", "*")
	// resp.Header.Set("Access-Control-Request-Method", "OPTIONS,HEAD,POST")
	// resp.Header.Set("Content-Type", "application/json; charset=utf-8")
}

func main() {
	port := flag.String("port", "8082", "listen port")
	targetAddr := flag.String("target", "widsboat-api.bitekun.xin", "your server domain")
	flag.Parse()

	proxyClient.Addr = *targetAddr

	log.Println("port:", *port)
	log.Println("target:", *targetAddr)

	//setupMiddlewares(ReverseProxyHandler)
	/*requestHandler := func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	}*/

	// 创建自定义服务器。
	s := &fasthttp.Server{
		Handler: ReverseProxyHandler,
        MaxConnsPerIP:10,
		// Every response will contain 'Server: My super server' header.
		Name: "My super server",

		// Other Server settings may be set here.
	}

	if err := s.ListenAndServe("localhost:"+*port); err != nil {
		log.Fatalf("error in fasthttp server: %s", err)
	}
	log.Println("start server...")
}
