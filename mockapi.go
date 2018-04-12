package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/gomodule/redigo/redis"
	"time"
	"os"
)

var (
	addr     = flag.String("addr", ":8082", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
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

	//read redis
	pool := newPool("192.168.35.171","31321","")
	defer pool.Close()
	s := pool.Get()
	/*_,err := s.Do("mset","poem","东临碣石，以观沧海。水何澹澹，山岛竦峙。树木丛生，百草丰茂。秋风萧瑟，洪波涌起。日月之行，若出其中。星汉灿烂，若出其里。幸甚至哉，歌以咏志。","url","http://www.google.com")
	if err!=nil{
		fmt.Println("err is",err.Error())
	}*/
	cc,err :=s.Do("mget","poem")
	if err != nil{
		fmt.Println("read redis'err is",err.Error())
	}
	//fmt.Println("poem from redis is",cc)
	fmt.Fprintf(ctx, "poem is %s", cc)
	//c := pool.Get()
	//mset mget
	//fmt.Printf("redis's status is %s",pool.Stats())
	//fmt.Printf("ActiveCount:%d IdleCount:%d\r\n",pool.Stats().ActiveCount,pool.Stats().IdleCount)
	/*_,setErr := c.Do("mset","name","biaoge","url","http://xxbandy.github.io")
	errCheck("setErr",setErr)
	if r,mgetErr := redis.Strings(c.Do("mget","name","url")); mgetErr == nil {
		for _,v := range r {
			fmt.Println("mget ",v)
		}
	}*/
}

//构造一个链接函数，如果没有密码，passwd为空字符串
func redisConn(ip,port,passwd string) (redis.Conn, error) {
	c,err := redis.Dial("tcp",
		ip+":"+port,
		redis.DialConnectTimeout(5*time.Second),
		redis.DialReadTimeout(1*time.Second),
		redis.DialWriteTimeout(1*time.Second),
		redis.DialPassword(passwd),
		redis.DialKeepAlive(1*time.Second),
	)
	return c,err
}

//构造一个错误检查函数
func errCheck(tp string,err error) {
	if err != nil {
		fmt.Printf("sorry,has some error for %s.\r\n",tp,err)
		os.Exit(-1)
	}
}

//构造一个连接池
//url为包装了redis的连接参数ip,port,passwd
func newPool(ip,port,passwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:            5,    //定义redis连接池中最大的空闲链接为3
		MaxActive:          18,    //在给定时间已分配的最大连接数(限制并发数)
		IdleTimeout:        240 * time.Second,
		MaxConnLifetime:    300 * time.Second,
		Dial:               func() (redis.Conn,error) { return redisConn(ip,port,passwd) },
	}
}
