package query

import (
	"context"
	"database/sql"
	stderr "errors"
	"yatter-backend-go/app/usecase/query"
	"yatter-backend-go/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// import に errors が2つある理由:
//   stderr "errors"                 ... 標準ライブラリ。errors.Is（エラーの種類判定）に使う
//   "yatter-backend-go/pkg/errors"  ... このプロジェクト独自。ErrNotFound などに使う
// 両方そのまま書くと名前が衝突するので、片方に stderr という別名を付けている。

// usecase/query/yweet.go の契約を、このファイルがMySQL向けに実装している。
var _ query.YweetQueryService = (*YweetQueryServiceImpl)(nil)

// 読み取り系なのでトランザクションを使わず、直接 *sqlx.DB を持つ。
type YweetQueryServiceImpl struct {
	db *sqlx.DB
}

func NewYweetQueryService(db *sqlx.DB) *YweetQueryServiceImpl {
	return &YweetQueryServiceImpl{db: db}
}

// FindYweetByID: 指定IDの yweet を1件、投稿者情報と一緒に取得する。
func (s *YweetQueryServiceImpl) FindYweetByID(ctx context.Context, yweetID uint64) (*query.FindYweetByIDDTO, error) {
	var dto query.FindYweetByIDDTO

	// GetContext は「1行だけ」取得するメソッド。
	// （複数行なら SelectContext。infra/query/timeline.go の方がその例）
	//
	// AS で別名を付けているのは、JOIN すると yweet.id と user.id が
	// どちらも "id" になって区別できなくなるため。
	// この別名が FindYweetByIDDTO の `db:"..."` タグと一致することでマッピングされる。
	err := s.db.GetContext(ctx, &dto, `
		SELECT
			yweet.id AS yweet_id,
			yweet.content AS yweet_content,
			yweet.created_at AS yweet_created_at,
			user.id AS user_id,
			user.username AS user_username,
			user.created_at AS user_created_at
		FROM yweet
		JOIN user ON yweet.user_id = user.id
		WHERE yweet.id = ?
	`, yweetID)
	if err != nil {
		// 該当する行が無い場合、sqlx は sql.ErrNoRows を返す。
		// これを ErrNotFound（404）に変換することで、
		// 最終的にクライアントに 404 が返る（ui/api/pkg/errors の Handle がcodeを見て変換する）。
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.WithDevMessage("yweet not found")
		}
		// それ以外（DB接続エラーなど）はそのまま返す。
		// ここで一律 ErrNotFound にしてしまうと、サーバー障害まで404になってしまい原因が分からなくなる。
		return nil, err
	}
	return &dto, nil
}
