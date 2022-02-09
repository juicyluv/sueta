package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/juicyluv/sueta/user_service/app/internal/handler"
	"github.com/juicyluv/sueta/user_service/app/internal/user/apperror"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
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
	router.HandlerFunc(http.MethodGet, userURL, h.GetUser)
	router.HandlerFunc(http.MethodGet, usersURL, h.GetUserByEmailAndPassword)
	router.HandlerFunc(http.MethodPost, usersURL, h.CreateUser)
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
		h.InternalError(w, err.Error(), "")
		return
	}

	h.JSON(w, http.StatusOK, user)
}

// GetUser parses request body, then, using user service,
// returns a user instance from database with given uuid or an error
// if input is invalid, given email already taken or something went wrong.
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("CREATE USER")

	var input CreateUserDTO
	if err := h.readJSON(w, r, &input); err != nil {
		h.BadRequest(w, err.Error(), "invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.BadRequest(w, err.Error(), "input validation failed. please, provide valid values")
		return
	}

	userId, err := h.UserService.Create(r.Context(), &input)
	if err != nil {
		if err := apperror.ErrEmailTaken; err != nil {
			h.BadRequest(w, err.Error(), "")
			return
		}
		h.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}

	h.JSON(w, http.StatusCreated, map[string]string{"id": userId})
}

// GetUserByEmailAndPassword parses email and password from URL query,
// then validates it and finds a user with this parameters.
// If validation failed, returns a 400 Bad request response.
// If user not found, returns a 404 Not Found response.
// If user found, returns a user information.
func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("GET USER BY EMAIL AND PASSWORD")

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email == "" || password == "" {
		h.BadRequest(w, "empty email or password", "email and password must be provided")
		return
	}

	user, err := h.UserService.GetByEmailAndPassword(r.Context(), email, password)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			h.NotFound(w)
			return
		}
		h.BadRequest(w, err.Error(), "")
		return
	}

	h.JSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateUserPartially(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("UPDATE USER PARTIALLY")

	// params := httprouter.ParamsFromContext(r.Context())
	// uuid := params.ByName("uuid")

	var input UpdateUserDTO
	if err := h.readJSON(w, r, &input); err != nil {
		h.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.BadRequest(w, err.Error(), "you have provided invalid values")
		return
	}

	// TODO: update the user
}

// DeleteUser parses user uuid from URL parameters, then
// tries to delete the user with this uuid. If user with provided
// uuid doesn't exists, 404 Not Found response will be sent to the client.
// If user has been deleted, 200 OK status will be sent.
// If something went wrong on the server side, 500 Internal Server Error
// response will be sent.
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("DELETE USER")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	err := h.UserService.Delete(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			h.NotFound(w)
			return
		}
		h.InternalError(w, err.Error(), "something went wrong on the server side")
		return
	}

	w.WriteHeader(http.StatusOK)
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

// readJSON decodes request body to the given destination(usually model struct).
// Returns an error on failure.
func (h *Handler) readJSON(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	// Create a new decoder and check for unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dest)
	if err != nil {
		// If an error occurred, send an error mapped to JSON decoding error
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Syntax error
		case errors.As(err, &syntaxError):
			return fmt.Errorf(
				"request body contains badly-formatted JSON (at character %d)",
				syntaxError.Offset,
			)
		// Type error
		case errors.As(err, &unmarshalTypeError):
			// If there's an info for struct field, show what field contains an error
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"request body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}

			return fmt.Errorf(
				"request body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset,
			)
		// Unmarshall error
		case errors.As(err, &invalidUnmarshalError):
			// We are panicing here because this is unexpected error
			panic(err)
		// Empty JSON error
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		// Unknown field error
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown key %s", fieldName)

		// Return error as-is
		default:
			return err
		}
	}

	// Decode one more time to check wheter here is another JSON object
	if err = dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must only contain single JSON value")
	}

	return nil
}

// Error is a wrapper around JSON method.
// It responses with specified error and http code.
func (h *Handler) Error(w http.ResponseWriter, code int, message, developerMessage string) {
	appError := apperror.NewAppError(code, message, developerMessage)
	h.JSON(w, code, appError)
}

// BadRequest is a wrapper around Error method.
// Responses with 400 Bad Request status code and specified error message.
func (h *Handler) BadRequest(w http.ResponseWriter, message, developerMessage string) {
	h.Error(w, http.StatusBadRequest, message, developerMessage)
}

// Not Found is a wrapper around JSON method.
// Responses with 404 Not Found status code and specified error message.
func (h *Handler) NotFound(w http.ResponseWriter) {
	h.JSON(w, http.StatusNotFound, apperror.ErrNotFound)
}

// InternalError is a wrapper around Error method.
// Responses with 500 Internal Server Error status code and specified error message.
func (h *Handler) InternalError(w http.ResponseWriter, message, developerMessage string) {
	h.Error(w, http.StatusInternalServerError, message, developerMessage)
}
