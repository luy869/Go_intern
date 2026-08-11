package yweet_test

import (
	"testing"
	"yatter-backend-go/app/domain/object/yweet"
	yatter_errors "yatter-backend-go/pkg/errors"
)

func Test_PendingYweet_SetContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		wantErr error
	}{
		{
			name:    "異常系: 内容が空",
			content: "",
			wantErr: yatter_errors.ErrBadRequest,
		},
		{
			name:    "正常系: コンテンツが存在する",
			content: "Hello",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			u := &yweet.PendingYweet{}

			err := u.SetContent(tt.content)

			if !yatter_errors.Is(err, tt.wantErr) {
				t.Errorf("SetContent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if u.Content() != tt.content {
					t.Errorf("content() = %v, want %v", u.Content(), tt.content)
				}
			}
		})
	}
}

func Test_PendingYweet_Userid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userid  uint64
		wantErr error
	}{
		{
			name:    "異常系: IDの要素が1",
			userid:  0,
			wantErr: yatter_errors.ErrInternal,
		},
		{
			name:    "正常系: IDの要素が1以上",
			userid:  3,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			u := &yweet.PendingYweet{}

			err := u.SetUserID(tt.userid)

			if !yatter_errors.Is(err, tt.wantErr) {
				t.Errorf("userid() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if u.UserID() != tt.userid {
					t.Errorf("userid() = %v, want %v", u.UserID(), tt.userid)
				}
			}
		})
	}
}
