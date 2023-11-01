package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/http"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pprof.StopCPUProfile()
		os.Exit(1)
	}()

	r := http.NewRouter()

	r.Use(func(req *http.Request[http.RequestBody], next http.Next) shepard.Result[http.Response[http.Body], error] {
		// remove trailing slash
		if req.URL.Path != "/" && req.URL.Path[len(req.URL.Path)-1] == '/' {
			req.URL.Path = req.URL.Path[:len(req.URL.Path)-1]
		}
		return next()
	})

	r.Use(func(req *http.Request[http.RequestBody], next http.Next) shepard.Result[http.Response[http.Body], error] {
		req.URL.Query().Set("hello", "middleware")
		return next()
	})

	r.Register(http.Post("/", func(req *http.Request[http.RequestBody]) http.Response[http.Body] {
		parseRes := req.ParseForm()
		if parseRes.IsErr() {
			errRes := http.NewResponseBuilder(http.NewHttpResponseBytes()).Status(http.StatusCodeInternalServerError).Body(http.NewBytesBody()).Unwrap()
			errRes.Body().Write() <- []byte(parseRes.Err().Unwrap().Error())
			errRes.Body().Finish()
			return errRes
		}

		greet := "world"
		if req.Params().Has("hello") {
			greet = *req.Params().Get("hello").First().Unwrap()
		}

		res := http.NewResponseBuilder(http.NewHttpResponseBytes()).Status(200).Body(http.NewBytesBody()).Unwrap()
		res.Body().Write() <- []byte("Hello " + greet + "!")
		res.Body().Finish()
		return res
	}))

	r.Register(http.Get("/hello", func(req *http.Request[http.RequestBody]) http.Response[http.Body] {
		res := http.NewResponseBuilder(http.NewHttpResponseBytes()).Status(200).Body(http.NewBytesBody()).Unwrap()

		greetDefault := new(string)
		*greetDefault = "world"
		greet := req.URL.Query().Get("hello").First().UnwrapOr(greetDefault)

		res.Body().Write() <- []byte("Hello " + *greet + "!")
		res.Body().Finish()
		return res
	}))

	r.Register(http.Get("/hello.json", func(req *http.Request[http.RequestBody]) http.Response[http.Body] {
		greetDefault := new(string)
		*greetDefault = "world"
		greet := req.URL.Query().Get("hello").First().UnwrapOr(greetDefault)
		res := http.JsonResponse(struct {
			Hello string `json:"hello"`
		}{
			Hello: *greet,
		})
		return res
	}))

	if err := http.Serve(":8080", r); err != nil {
		log.Fatal(err)
	}
}
