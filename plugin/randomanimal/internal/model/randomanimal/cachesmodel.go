package randomanimal

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CachesModel = (*customCachesModel)(nil)

type (
	// CachesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCachesModel.
	CachesModel interface {
		cachesModel
		withSession(session sqlx.Session) CachesModel
	}

	customCachesModel struct {
		*defaultCachesModel
	}
)

// NewCachesModel returns a model for the database table.
func NewCachesModel(conn sqlx.SqlConn) CachesModel {
	return &customCachesModel{
		defaultCachesModel: newCachesModel(conn),
	}
}

func (m *customCachesModel) withSession(session sqlx.Session) CachesModel {
	return NewCachesModel(sqlx.NewSqlConnFromSession(session))
}
