package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "WELCOME!\n")
}

func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

func IndexAPI(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "WELCOME FROM API!\n")
}

func HelloAPI(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello rock star API, %s!\n", ctx.UserValue("name"))
}

type HostMap map[string]fasthttp.RequestHandler

func (hm HostMap) HandleHost(ctx *fasthttp.RequestCtx) {
	host := string(ctx.Host())
	path := string(ctx.URI().Path())
	fmt.Println("Host", host)
	fmt.Println("Path", path)

	hkey := "/"
	arr := strings.SplitN(path, "/", 3)
	fmt.Printf("arr %#v\n", arr)
	if len(arr) > 1 && arr[1] == "api" {
		hkey = arr[1]
	}
	if handler := hm[hkey]; handler != nil {
		handler(ctx)
	} else {
		ctx.Error("Forbidden", 403)
	}
}

func main() {
	// Standard router
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/:name", Hello)

	// API router
	a := router.New()
	a.GET("/api", IndexAPI)
	a.GET("/api/hello/:name", HelloAPI)

	hm := make(HostMap)
	hm["/"] = r.Handler
	hm["api"] = a.Handler

	port := "7000"
	// if len(os.Args) > 1 {
	// 	port = os.Args[1]
	// }

	fmt.Println("Listening on port:", port )
	log.Fatal(fasthttp.ListenAndServe(":" + port, hm.HandleHost))
}