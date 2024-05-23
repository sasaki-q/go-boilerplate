package handlers

import (
	"io"
	"net/http"

	"boilerplate/server/restapi"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Handler struct {
	*UserHandler
}

func validate[T any](b *T, r io.ReadCloser, v *validator.Validate) error {
	if err := render.DecodeJSON(r, b); err != nil {
		return err
	}

	return v.Struct(b)
}

func successResponse(w http.ResponseWriter, r *http.Request, d interface{}) {
	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, d)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, d string) {
	render.Status(r, status)
	render.JSON(w, r, restapi.ErrorResponse{Message: d})
}
