package user

import (
	"encoding/json"
	"net/http"

	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/user/:id"
)

type Handler struct {
	Logger      logger.Logger
	UserService Service
}

func (h *Handler) Register(router *httprouter.Router) {
}

func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("GET USER BY EMAIL AND PASSWORD")

	params := httprouter.ParamsFromContext(r.Context())
	email := params.ByName("email")
	password := params.ByName("password")

	if email == "" || password == "" {
		// Bad Request Error
	}
}

func (h *Handler) JSON(w http.ResponseWriter, code int, data interface{}) error {
	obj, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(obj)

	return nil
}
