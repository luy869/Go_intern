package timeline

import (
	"yatter-backend-go/app/ui/api/yweet"
	"yatter-backend-go/app/usecase/query"
)

func toTimelineResponse(items []*query.TimelineItemDTO) []yweet.PostYweetsResponse {
	result := []yweet.PostYweetsResponse{}

	for _, item := range items {
		result = append(result, yweet.PostYweetsResponse{
			ID:        item.YweetID,
			Content:   item.YweetContent,
			CreatedAt: item.YweetCreatedAt.Format("2006-01-02T15:04:05.000Z"),
			User: yweet.UserResponse{
				ID:        item.UserID,
				Username:  item.UserUsername,
				CreatedAt: item.UserCreatedAt.Format("2006-01-02T15:04:05.000Z"),
			},
			ImageAttachments: []yweet.AttachmentResponse{},
		})
	}

	return result
}
