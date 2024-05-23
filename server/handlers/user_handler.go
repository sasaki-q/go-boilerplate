package handlers

import (
	"net/http"

	"boilerplate/server/models"
	"boilerplate/server/repositories"
	"boilerplate/server/restapi"
)

type UserHandler struct {
	h *Helper
	r repositories.IUserRepository
}

func NewUserHandler(h *Helper, r *repositories.UserRepository) *UserHandler {
	return &UserHandler{h, r}
}

func (h *UserHandler) GetUserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	us, err := h.r.List(ctx)
	if err != nil {
		h.h.l.Sugar().Errorf("GetUserList h.r.List / %v", err)
		errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	resp := []*restapi.UserResponse{}
	for _, v := range us {
		resp = append(resp, &restapi.UserResponse{
			Id:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rb := &restapi.UserInput{}
	if err := validate(rb, r.Body, h.h.v); err != nil {
		h.h.l.Sugar().Errorf("CreateUser validate / %v", err)
		errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.r.Create(ctx, &models.User{Name: rb.Name})
	if err != nil {
		h.h.l.Sugar().Errorf("CreateUser h.r.Create / %v", err)
		errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	return
}
