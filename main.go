package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func IndexAPI(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "WELCOME FROM API!\n")
}

func HelloAPI(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello rock star API, %s!\n", ctx.UserValue("name"))
}

type HostMap map[string]fasthttp.RequestHandler

func (hm HostMap) HandleHost(ctx *fasthttp.RequestCtx) {
	path := string(ctx.URI().Path())
	// host := string(ctx.Host())

	if path == "/" {
		ctx.URI().SetPath("/index.html") // Path rewrite - nice!
	}

	fmt.Println("Path", path)

	hkey := "/"
	arr := strings.SplitN(path, "/", 3)
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
	const port = "7000"

	// API router
	a := router.New()
	a.GET("/api", IndexAPI)
	a.GET("/api/hello/:name", HelloAPI)

	// Static Handler for SPA
	fs := fasthttp.FS{
		Root:               "./static",
		GenerateIndexPages: false,
		Compress:           false,
	}
	staticHdlr := fs.NewRequestHandler()

	hm := make(HostMap)
	hm["/"] = staticHdlr
	hm["api"] = a.Handler

	fmt.Println("Listening on port:", port)
	log.Fatal(fasthttp.ListenAndServe(":"+port, hm.HandleHost))
}
