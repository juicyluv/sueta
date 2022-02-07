package handler

import "github.com/julienschmidt/httprouter"

type Handler struct {
	router *httprouter.Router
}

func NewHandler() *Handler {
	return &Handler{
		router: httprouter.New(),
	}
}

func (h *Handler) Router() *httprouter.Router {
	return h.router
}
