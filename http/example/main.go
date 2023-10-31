package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

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

	if err := http.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}
