package applications

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TwoononeModel = (*customTwoononeModel)(nil)

type (
	// TwoononeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTwoononeModel.
	TwoononeModel interface {
		twoononeModel
		withSession(session sqlx.Session) TwoononeModel
	}

	customTwoononeModel struct {
		*defaultTwoononeModel
	}
)

// NewTwoononeModel returns a model for the database table.
func NewTwoononeModel(conn sqlx.SqlConn) TwoononeModel {
	return &customTwoononeModel{
		defaultTwoononeModel: newTwoononeModel(conn),
	}
}

func (m *customTwoononeModel) withSession(session sqlx.Session) TwoononeModel {
	return NewTwoononeModel(sqlx.NewSqlConnFromSession(session))
}
