package yweet

import (
	"yatter-backend-go/app/domain/object/user"
	"yatter-backend-go/app/domain/object/yweet"
	"yatter-backend-go/app/usecase/query"
)

type UserResponse struct {
	ID             uint64 `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	CreatedAt      string `json:"created_at"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
	Note           string `json:"note"`
	Avatar         string `json:"avatar"`
	Header         string `json:"header"`
}

type AttachmentResponse struct {
	ID          uint64 `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type PostYweetsResponse struct {
	ID               uint64               `json:"id"`
	User             UserResponse         `json:"user"`
	Content          string               `json:"content"`
	CreatedAt        string               `json:"created_at"`
	ImageAttachments []AttachmentResponse `json:"image_attachments"`
}

func toUserResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:             u.ID(),
		Username:       u.Username(),
		CreatedAt:      u.CreatedAt().Format("2006-01-02T15:04:05.000Z"),
		DisplayName:    "",
		FollowersCount: 0,
		FollowingCount: 0,
		Note:           "",
		Avatar:         "",
		Header:         "",
	}
}

func toPostYweetsResponse(postedYweet *yweet.Yweet, author *user.User) PostYweetsResponse {
	return PostYweetsResponse{
		ID:               postedYweet.ID(),
		User:             toUserResponse(author),
		Content:          postedYweet.Content(),
		CreatedAt:        postedYweet.CreatedAt().Format("2006-01-02T15:04:05.000Z"),
		ImageAttachments: []AttachmentResponse{},
	}
}

func toGetYweetResponse(dto *query.FindYweetByIDDTO) PostYweetsResponse {
	return PostYweetsResponse{
		ID:        dto.YweetID,
		Content:   dto.YweetContent,
		CreatedAt: dto.YweetCreatedAt.Format("2006-01-02T15:04:05.000Z"),
		User: UserResponse{
			ID:        dto.UserID,
			Username:  dto.UserUsername,
			CreatedAt: dto.UserCreatedAt.Format("2006-01-02T15:04:05.000Z"),
		},
		ImageAttachments: []AttachmentResponse{},
	}
}
