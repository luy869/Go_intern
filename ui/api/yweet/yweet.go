package yweet

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	ui_errors "yatter-backend-go/app/ui/api/pkg/errors"
	"yatter-backend-go/app/usecase/query"
	"yatter-backend-go/app/usecase/yweet"
	"yatter-backend-go/pkg/errors"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
}

func NewYweetHandler(createYweetUseCase yweet.CreateYweetUseCase, queryYweetService query.YweetQueryService) Handler {
	return &yweetHandlerImpl{
		createYweetUseCase: createYweetUseCase,
		queryYweetService:  queryYweetService,
	}
}

var _ Handler = (*yweetHandlerImpl)(nil)

type yweetHandlerImpl struct {
	createYweetUseCase yweet.CreateYweetUseCase
	queryYweetService  query.YweetQueryService
}

func (h *yweetHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.Header.Get("Authentication"), " ")
	if len(parts) != 2 || parts[0] != "username" {
		ui_errors.Handle(w, errors.ErrUnauthorized)
		return
	}
	username := parts[1]

	var req PostYweetsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ui_errors.Handle(w, errors.ErrBadRequest)
		return
	}

	ctx := r.Context()
	postedYweet, author, err := h.createYweetUseCase.Create(ctx, username, req.Yweet)
	if err != nil {
		ui_errors.Handle(w, err)
		return
	}

	resp := toPostYweetsResponse(postedYweet, author)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		ui_errors.Handle(w, errors.ErrInternal.WithDevMessage("failed to encode response"))
		return
	}
}

func (h *yweetHandlerImpl) Find(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		ui_errors.Handle(w, errors.ErrBadRequest)
		return
	}
	ctx := r.Context()

	dto, err := h.queryYweetService.FindYweetByID(ctx, id)
	if err != nil {

		ui_errors.Handle(w, err)
		return
	}

	resp := toGetYweetResponse(dto)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ui_errors.Handle(w, errors.ErrInternal.WithDevMessage("failed to encode response"))
		return
	}
}
