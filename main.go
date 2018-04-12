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

	//log.Println(ctx, "Hello, world! Requested path is %q", string(ctx.Path()))
	req := &ctx.Request
	resp := &ctx.Response

	prepareRequest(req)

	defer resp.SetConnectionClose()
	if err := proxyClient.Do(req, resp); err != nil {
		ctx.Logger().Printf("error when proxying the request: %s", err)
	}

	postprocessResponse(resp)
}

func prepareRequest(req *fasthttp.Request) {
	// do not proxy "Connection" header.
	//req.Header.Del("Connection")
	// strip other unneeded headers.

	// alter other request params before sending them to upstream host
	req.Header.SetHost(proxyAddr)

	req.SetRequestURI("http://realtime.inj100.jstv.com/real_i/v1?api_key=ec5cc8b1fb8c83eb92c66e5338182be9&api_sig=582fb08c4dba31a4e4a3894b96f1db15&imei=95343ade91432af2bdf1f120f02a29870505f83e&line_code=30&line_id=14&referer=www.baidu.com&timestamp=1523510279&updown_type=0")
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	//resp.Header.Del("Connection")
	//resp.SkipBody = false
	//resp.AppendBody([]byte("abc"))

	// strip other unneeded headers

	// alter other response data if needed
	// resp.Header.Set("Access-Control-Allow-Origin", "*")
	// resp.Header.Set("Access-Control-Request-Method", "OPTIONS,HEAD,POST")
	// resp.Header.Set("Content-Type", "application/json; charset=utf-8")
}

func main() {
	port := flag.String("port", "80", "listen port")
	targetAddr := flag.String("target", "realtime.inj100.jstv.com", "your server domain")
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
        //MaxConnsPerIP:1,
		//MaxRequestsPerConn:1,
		// Every response will contain 'Server: My super server' header.
		Name: "My super server",
		//DisableKeepalive:true,


		// Other Server settings may be set here.
	}


	if err := s.ListenAndServe(":"+*port); err != nil {
		log.Fatalf("error in fasthttp server: %s", err)
	}
	log.Println("start server...")
}
