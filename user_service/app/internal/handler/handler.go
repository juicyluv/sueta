package handler

import "github.com/julienschmidt/httprouter"

// Handling describes new routes registration.
type Handling interface {
	Register(router *httprouter.Router)
}
