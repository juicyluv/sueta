package user

import (
	"encoding/json"
	"net/http"

	"github.com/juicyluv/sueta/user_service/app/internal/handler"
	"github.com/juicyluv/sueta/user_service/app/internal/user/apperror"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/user/:uuid"
)

// Handler handles requests specified to user service.
type Handler struct {
	Logger      logger.Logger
	UserService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, userService Service) handler.Handling {
	return &Handler{
		Logger:      logger,
		UserService: userService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, h.GetUser)
}

// GetUser parses uuid from URL parameters, then, using user service,
// returns a user instance from database with given uuid or an error
// if there's no user with such uuid or something went wrong.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("GET USER")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	user, err := h.UserService.GetById(r.Context(), uuid)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, nil)
		return
	}

	h.JSON(w, http.StatusOK, user)
}

// JSON encodes to JSON format given data and sends a response
// to the client with a given http code and encoded data.
func (h *Handler) JSON(w http.ResponseWriter, code int, data interface{}) {
	obj, err := json.Marshal(data)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(obj)
}

// Error is a wrapper around JSON function. It responses with specified error and
// http code.
func (h *Handler) Error(w http.ResponseWriter, code int, err *apperror.AppError) {
	h.JSON(w, code, err)
}
