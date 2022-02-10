package internal

import (
	"net/http"

	_ "github.com/juicyluv/sueta/user_service/app/docs"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	docsPath = "/docs/*any"
)

func InitSwagger(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, docsPath, httpSwagger.WrapHandler)
}
