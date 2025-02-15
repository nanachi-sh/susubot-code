package twoonone

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserTwoononeModel = (*customUserTwoononeModel)(nil)

type (
	// UserTwoononeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTwoononeModel.
	UserTwoononeModel interface {
		userTwoononeModel
		withSession(session sqlx.Session) UserTwoononeModel
	}

	customUserTwoononeModel struct {
		*defaultUserTwoononeModel
	}
)

// NewUserTwoononeModel returns a model for the database table.
func NewUserTwoononeModel(conn sqlx.SqlConn) UserTwoononeModel {
	return &customUserTwoononeModel{
		defaultUserTwoononeModel: newUserTwoononeModel(conn),
	}
}

func (m *customUserTwoononeModel) withSession(session sqlx.Session) UserTwoononeModel {
	return NewUserTwoononeModel(sqlx.NewSqlConnFromSession(session))
}
