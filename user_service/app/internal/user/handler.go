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
	logger      logger.Logger
	userService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, userService Service) handler.Handling {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userURL, h.GetUser)
	router.HandlerFunc(http.MethodGet, usersURL, h.GetUserByEmailAndPassword)
	router.HandlerFunc(http.MethodPost, usersURL, h.CreateUser)
	router.HandlerFunc(http.MethodPatch, userURL, h.UpdateUserPartially)
	router.HandlerFunc(http.MethodDelete, userURL, h.DeleteUser)
}

// GetUser godoc
// @Summary Show user information
// @Description Get user by uuid.
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User id"
// @Success 200 {object} User
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{uuid} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET USER")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	user, err := h.userService.GetById(r.Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrNoRows):
			h.NotFound(w)
		case errors.Is(err, apperror.ErrInvalidUUID):
			h.BadRequest(w, err.Error(), "")
		default:
			h.InternalError(w, err.Error(), "")
		}
		return
	}

	h.JSON(w, http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create user
// @Description Register a new user.
// @Tags users
// @Accept json
// @Produce json
// @Param input body user.CreateUserDTO true "JSON input"
// @Success 201 {object} internal.CreateUserResponse
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE USER")

	var input CreateUserDTO
	if err := h.readJSON(w, r, &input); err != nil {
		h.BadRequest(w, err.Error(), "invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.BadRequest(w, err.Error(), "input validation failed. please, provide valid values")
		return
	}

	if input.Password != input.RepeatPassword {
		h.BadRequest(w, "passwords don't match", "provided passwords must to match")
		return
	}

	userId, err := h.userService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmailTaken) {
			h.BadRequest(w, err.Error(), "")
			return
		}
		h.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}

	h.JSON(w, http.StatusCreated, map[string]string{"id": userId})
}

// GetUserByEmailAndPassword godoc
// @Summary Get user by email and password from query parameters
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "user email"
// @Param password query string true "user raw password"
// @Success 200 {object} User
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users [get]
func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET USER BY EMAIL AND PASSWORD")

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email == "" || password == "" {
		h.BadRequest(w, "empty email or password", "email and password must be provided")
		return
	}

	user, err := h.userService.GetByEmailAndPassword(r.Context(), email, password)
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

// UpdateUserPartially godoc
// @Summary Update user
// @Description Partially update the user with provided current password.
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User id"
// @Param input body user.UpdateUserDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{uuid} [patch]
func (h *Handler) UpdateUserPartially(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE USER PARTIALLY")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	var input UpdateUserDTO
	if err := h.readJSON(w, r, &input); err != nil {
		h.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.BadRequest(w, err.Error(), "you have provided invalid values")
		return
	}

	input.UUID = uuid

	err := h.userService.UpdatePartially(r.Context(), &input)
	if err != nil {
		switch err {
		case apperror.ErrNoRows:
			h.NotFound(w)
		case apperror.ErrWrongPassword:
			h.BadRequest(w, err.Error(), "you entered wrong password")
		default:
			h.InternalError(w, err.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete the user by uuid.
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User id"
// @Success 200
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{uuid} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE USER")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	err := h.userService.Delete(r.Context(), uuid)
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
