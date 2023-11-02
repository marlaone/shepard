package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
	. "github.com/marlaone/shepard/http"
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

	r := NewRouter().
		Use(func(req *Request[RequestBody], next Next) shepard.Result[Response[Body], error] {
			// remove trailing slash
			if req.URL.Path != "/" && req.URL.Path[len(req.URL.Path)-1] == '/' {
				req.URL.Path = req.URL.Path[:len(req.URL.Path)-1]
			}
			return next()
		}).
		Use(func(req *Request[RequestBody], next Next) shepard.Result[Response[Body], error] {
			req.URL.Query().Set("hello", "middleware")
			return next()
		}).
		Route(Post("/", func(req *Request[RequestBody]) Response[Body] {
			parseRes := req.ParseForm()
			if parseRes.IsErr() {
				errRes := NewResponseBuilder(NewHttpResponseBytes()).Status(StatusCodeInternalServerError).Body(NewBytesBody()).Unwrap()
				errRes.Body().Write(slice.Init[byte]([]byte(parseRes.Err().Unwrap().Error())...))
				errRes.Body().Close()

				return errRes
			}

			greet := "world"
			if req.Params().Has("hello") {
				greet = *req.Params().Get("hello").First().Unwrap()
			}

			res := NewResponseBuilder(NewHttpResponseBytes()).Status(200).Body(NewBytesBody()).Unwrap()
			res.Body().Write(slice.Init[byte]([]byte("Hello " + greet + "!")...))
			res.Body().Close()
			return res
		})).
		Route(Get("/hello", func(req *Request[RequestBody]) Response[Body] {
			res := NewResponseBuilder(NewHttpResponseBytes()).Status(200).Body(NewBytesBody()).Unwrap()

			greetDefault := new(string)
			*greetDefault = "world"
			greet := req.URL.Query().Get("hello").First().UnwrapOr(greetDefault)

			res.Body().Write(slice.Init[byte]([]byte("Hello " + *greet + "!")...))
			res.Body().Close()
			return res
		})).
		Route(Get("/hello.json", func(req *Request[RequestBody]) Response[Body] {
			greetDefault := new(string)
			*greetDefault = "world"
			greet := req.URL.Query().Get("hello").First().UnwrapOr(greetDefault)
			res := JsonResponse(struct {
				Hello string `json:"hello"`
			}{
				Hello: *greet,
			})
			return res
		}))

	if err := Serve(":8080", r); err != nil {
		log.Fatal(err)
	}
}
