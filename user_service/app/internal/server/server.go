package server

import (
	"net/http"

	"github.com/juicyluv/sueta/user_service/app/config"
	"github.com/juicyluv/sueta/user_service/app/internal/handler"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		server: &http.Server{
			Handler: handler.NewHandler().Router(),
		},
	}
}
