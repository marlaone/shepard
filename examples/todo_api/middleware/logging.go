// Package middleware provides HTTP middleware for the Todo API.
package middleware

import (
	"log"
	"time"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/http"
)

// RequestLogger logs incoming requests and their response status.
func RequestLogger(ctx *http.Context, next func() shepard.Result[http.Response, http.HttpError]) shepard.Result[http.Response, http.HttpError] {
	start := time.Now()
	log.Printf("--> %s %s", ctx.Method(), ctx.Path())

	result := next()

	duration := time.Since(start)
	if result.IsOk() {
		log.Printf("<-- %d %s (%v)", result.Unwrap().StatusCode, ctx.Path(), duration)
	} else {
		log.Printf("<-- ERR %s (%v)", ctx.Path(), duration)
	}

	return result
}
