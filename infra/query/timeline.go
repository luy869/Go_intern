package query

import (
	"context"
	"yatter-backend-go/app/usecase/query"

	"github.com/jmoiron/sqlx"
)

// usecase/query/timeline.go の契約を、このファイルがMySQL向けに実装している。
var _ query.TimelineQueryService = (*TimelineQueryServiceImpl)(nil)

// 読み取り系はトランザクションを使わないので、構造体が直接 *sqlx.DB を持つ。
// （書き込み系の infra/yweet.go は ctx からトランザクションを取り出す作りで、ここが違う）
type TimelineQueryServiceImpl struct {
	db *sqlx.DB
}

func NewTimelineQueryService(db *sqlx.DB) *TimelineQueryServiceImpl {
	return &TimelineQueryServiceImpl{db: db}
}

// FindPublicTimeline: 公開タイムラインを新しい順に取得する。
func (s *TimelineQueryServiceImpl) FindPublicTimeline(ctx context.Context, offset, limit int, onlyImage bool) ([]*query.TimelineItemDTO, error) {
	// 画像添付機能は未実装のため、only_image=true で絞ると該当する yweet は存在しない。
	// nil ではなく空スライスを返すのがポイント（nil を返すとJSONで null になってしまう）。
	if onlyImage {
		return []*query.TimelineItemDTO{}, nil
	}

	items := []*query.TimelineItemDTO{}

	// SelectContext は「複数行」を取得してスライスに詰めるメソッド。
	// （1行だけ取る場合は GetContext を使う。infra/query/yweet.go の方がその例）
	//
	// AS で別名を付けているのは、JOIN すると yweet.id と user.id が
	// どちらも "id" になって区別できなくなるため。
	// この別名が TimelineItemDTO の `db:"..."` タグと一致することでマッピングされる。
	err := s.db.SelectContext(ctx, &items, `
              SELECT
                      yweet.id AS yweet_id,
                      yweet.content AS yweet_content,
                      yweet.created_at AS yweet_created_at,
                      user.id AS user_id,
                      user.username AS user_username,
                      user.created_at AS user_created_at
              FROM yweet
              JOIN user ON yweet.user_id = user.id
              ORDER BY yweet.created_at DESC
              LIMIT ? OFFSET ?
      `, limit, offset)
	if err != nil {
		return nil, err
	}

	return items, nil
}
