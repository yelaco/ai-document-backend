package server

import "net/http"

type Router interface {
	GetHandler() http.Handler
}
