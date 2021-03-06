package post

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/juicyluv/sueta/post_service/app/internal/handler"
	"github.com/juicyluv/sueta/post_service/app/internal/post/apperror"
	"github.com/juicyluv/sueta/post_service/app/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

const (
	postURL  = "/api/post"
	postsURL = "/api/posts/:uuid"
)

type Handler struct {
	logger      logger.Logger
	postService Service
}

func NewHandler(logger logger.Logger, postService Service) handler.Handling {
	return &Handler{
		logger:      logger,
		postService: postService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, postURL, h.GetPost)
	router.HandlerFunc(http.MethodPost, postsURL, h.CreatePost)
	router.HandlerFunc(http.MethodPatch, postURL, h.UpdatePostPartially)
	router.HandlerFunc(http.MethodDelete, postURL, h.DeletePost)
}

// GetPost godoc
// @Summary Show post information
// @Description Get post by uuid.
// @Tags posts
// @Accept json
// @Produce json
// @Param uuid path string true "Post id"
// @Success 200 {object} Post
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /posts/{uuid} [get]
func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET POST")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	post, err := h.postService.GetById(r.Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrNoRows):
			h.NotFound(w)
		default:
			h.InternalError(w, err.Error(), "")
		}
		return
	}

	h.JSON(w, http.StatusOK, post)
}

// CreatePost godoc
// @Summary Create post
// @Description Register a new post.
// @Tags posts
// @Accept json
// @Produce json
// @Param input body post.CreatepostDTO true "JSON input"
// @Success 201 {object} internal.CreatePostResponse
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /posts [post]
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE POST")

	var input CreatePostDTO
	if err := h.readJSON(w, r, &input); err != nil {
		h.BadRequest(w, err.Error(), "invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.BadRequest(w, err.Error(), apperror.ErrValidationFailed.Error())
		return
	}

	postId, err := h.postService.Create(r.Context(), &input)
	if err != nil {
		h.InternalError(w, fmt.Sprintf("cannot create post: %v", err), "")
		return
	}

	h.JSON(w, http.StatusCreated, map[string]string{"id": postId})
}

// UpdatePostPartially godoc
// @Summary Update post
// @Description Partially update the post with provided current password.
// @Tags posts
// @Accept json
// @Produce json
// @Param uuid path string true "Post id"
// @Param input body post.UpdatePostDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /posts/{uuid} [patch]
func (h *Handler) UpdatePostPartially(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE POST PARTIALLY")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	var input UpdatePostDTO
	if err := h.readJSON(w, r, &input); err != nil {
		h.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.BadRequest(w, err.Error(), apperror.ErrValidationFailed.Error())
		return
	}

	input.UUID = uuid

	err := h.postService.UpdatePartially(r.Context(), &input)
	if err != nil {
		switch err {
		case apperror.ErrNoRows:
			h.NotFound(w)
		default:
			h.InternalError(w, err.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeletePost godoc
// @Summary Delete post
// @Description Delete the post by uuid.
// @Tags posts
// @Accept json
// @Produce json
// @Param uuid path string true "Post id"
// @Success 200
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /posts/{uuid} [delete]
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE POST")

	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")

	err := h.postService.Delete(r.Context(), uuid)
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
