package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/tamper000/hysteria-auth/internal/repository"
)

type Repository interface {
	CreateUser(ctx context.Context, user repository.User) error
	DeleteUser(ctx context.Context, id string) error
	UserExists(ctx context.Context, auth string) (bool, error)
	GetIDByAuth(ctx context.Context, auth string) (string, bool, error)
}

type Handler struct {
	repo Repository
}

func New(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Register(ctx context.Context, input *User) (*ResultOutput, error) {
	exists, err := h.repo.UserExists(ctx, input.Body.Auth)
	if err != nil {
		return nil, huma.Error500InternalServerError("User verification error", err)
	}

	result := &ResultOutput{}
	if exists {
		result.Status = http.StatusConflict
		return result, nil
	}

	err = h.repo.CreateUser(ctx, repository.User{
		ID:       input.Body.ID,
		Auth:     input.Body.Auth,
		Optional: input.Body.Optional,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("User creation error", err)
	}

	result.Status = http.StatusCreated
	return result, nil
}

func (h *Handler) Delete(ctx context.Context, input *UserID) (*ResultOutput, error) {
	err := h.repo.DeleteUser(ctx, input.Body.ID)
	if err != nil {
		return nil, huma.Error404NotFound("User not found", err)
	}

	result := &ResultOutput{}
	result.Status = http.StatusOK
	return result, nil
}

func (h *Handler) Auth(ctx context.Context, input *UserAuth) (*AuthOutput, error) {
	ID, exists, err := h.repo.GetIDByAuth(ctx, input.Body.Auth)
	if err != nil || !exists {
		return nil, huma.Error404NotFound("User verification error", err)
	}

	result := &AuthOutput{}
	result.Body.ID = ID
	result.Body.Ok = true
	return result, nil
}
