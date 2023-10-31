package http

import "github.com/marlaone/shepard/collections/hashmap"


type Router struct {
	routes hashmap.HashMap[Method, []Route]

	notFoundHandler Handler
}

func NewRouter() *Router {
	return &Router{
		routes: hashmap.New[Method, []Route](),
	}
}
