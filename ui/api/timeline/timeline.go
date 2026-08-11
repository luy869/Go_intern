package timeline

import (
	"encoding/json"
	"net/http"
	"strconv"
	ui_errors "yatter-backend-go/app/ui/api/pkg/errors"
	"yatter-backend-go/app/usecase/query"
	"yatter-backend-go/pkg/errors"
)

type Handler interface {
	Public(w http.ResponseWriter, r *http.Request)
}

func NewTimelineHandler(timelineQueryService query.TimelineQueryService) Handler {
	return &timelineHandlerImpl{
		timelineQueryService: timelineQueryService,
	}
}

var _ Handler = (*timelineHandlerImpl)(nil)

type timelineHandlerImpl struct {
	timelineQueryService query.TimelineQueryService
}

func (h *timelineHandlerImpl) Public(w http.ResponseWriter, r *http.Request) {
	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			ui_errors.Handle(w, errors.ErrBadRequest)
			return
		}
		offset = parsed
	}

	limit := 40
	if v := r.URL.Query().Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			ui_errors.Handle(w, errors.ErrBadRequest)
			return
		}
		limit = parsed
	}
	if limit > 80 {
		limit = 80
	}

	onlyImage := r.URL.Query().Get("only_image") == "true"

	items, err := h.timelineQueryService.FindPublicTimeline(r.Context(), offset, limit, onlyImage)
	if err != nil {
		ui_errors.Handle(w, err)
		return
	}

	resp := toTimelineResponse(items)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ui_errors.Handle(w, errors.ErrInternal.WithDevMessage("failed to encode response"))
		return
	}
}
