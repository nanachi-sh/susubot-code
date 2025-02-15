package mysql

import (
	"context"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/sql/ldap"
	twoonone_model "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/model/twoonone"
	database_type "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types/database"
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type db_handler struct{}

type (
	db_action_inc_coin             struct{ float64 }
	db_action_dec_coin             struct{ float64 }
	db_action_inc_wincount         struct{ uint }
	db_action_dec_wincount         struct{ uint }
	db_action_inc_losecount        struct{ uint }
	db_action_dec_losecount        struct{ uint }
	db_action_update_getdaily_time struct{ time.Time }
)

func NewHandler() database_type.DatabaseModel {
	return &db_handler{}
}

func IncCoin(n float64) database_type.Action {
	return &db_action_inc_coin{n}
}

func DecCoin(n float64) database_type.Action {
	return &db_action_dec_coin{n}
}

func IncWinCount(n ...uint) database_type.Action {
	c := func() uint {
		if len(n) > 0 {
			return n[0]
		} else {
			return 1
		}
	}()
	return &db_action_inc_wincount{c}
}

func DecWinCount(n ...uint) database_type.Action {
	c := func() uint {
		if len(n) > 0 {
			return n[0]
		} else {
			return 1
		}
	}()
	return &db_action_dec_wincount{c}
}

func IncLoseCount(n ...uint) database_type.Action {
	c := func() uint {
		if len(n) > 0 {
			return n[0]
		} else {
			return 1
		}
	}()
	return &db_action_inc_losecount{c}
}

func DecLoseCount(n ...uint) database_type.Action {
	c := func() uint {
		if len(n) > 0 {
			return n[0]
		} else {
			return 1
		}
	}()
	return &db_action_dec_losecount{c}
}

func UpdateGetDailyTime(n ...time.Time) database_type.Action {
	t := func() time.Time {
		if len(n) > 0 {
			t := n[0]
			if t.IsZero() {
				t = time.Now()
			}
			return t
		} else {
			return time.Now()
		}
	}()
	return &db_action_update_getdaily_time{t}
}

func (s *db_action_inc_coin) Merge(logger logx.Logger, u *twoonone_model.UserTwoonone) {
	if s.float64 < 0 {
		s.float64 = 0
	}
	u.Coin += s.float64
}

func (s *db_action_dec_coin) Merge(logger logx.Logger, u *twoonone_model.UserTwoonone) {
	if s.float64 < 0 {
		s.float64 = 0
	}
	u.Coin -= s.float64
}

func (s *db_action_inc_wincount) Merge(logger logx.Logger, u *twoonone_model.UserTwoonone) {
	u.Wincount += int64(s.uint)
}

func (s *db_action_dec_wincount) Merge(logger logx.Logger, u *twoonone_model.UserTwoonone) {
	u.Wincount -= int64(s.uint)
}

func (s *db_action_inc_losecount) Merge(logger logx.Logger, u *twoonone_model.UserTwoonone) {
	u.Losecount += int64(s.uint)
}

func (s *db_action_dec_losecount) Merge(logger logx.Logger, u *twoonone_model.UserTwoonone) {
	u.Losecount -= int64(s.uint)
}

func (s *db_action_update_getdaily_time) Merge(logger logx.Logger, u *twoonone_model.UserTwoonone) {
	u.LastGetdaliyTime = s.Time
}

func (dbh *db_handler) GetUser(logger logx.Logger, id string) (database_type.User, *pkg_types.AppError) {
	if id == "" {
		logger.Error("invalid argument")
		return database_type.User{}, pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	u, err := getUser(logger, ctx, id)
	if err != nil {
		return database_type.User{}, err
	}
	return u, nil
}

func (dbh *db_handler) CreateUser(logger logx.Logger, id string) *pkg_types.AppError {
	if id == "" {
		logger.Error("invalid argument")
		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if serr := createUser(logger, ctx, twoonone_model.UserTwoonone{
		Id:               id,
		Wincount:         0,
		Losecount:        0,
		LastGetdaliyTime: time.Time{},
		Coin:             0,
	}); serr != nil {
		return serr
	}
	return nil
}

func (dbh *db_handler) UpdateUser(logger logx.Logger, id string, actions ...database_type.Action) *pkg_types.AppError {
	if id == "" {
		logger.Error("invalid argument")
		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	u := func() twoonone_model.UserTwoonone {
		u := new(twoonone_model.UserTwoonone)
		for _, v := range actions {
			v.Merge(logger, u)
		}
		return *u
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if serr := updateUser(logger, ctx, u); serr != nil {
		return serr
	}
	return nil
}

func (dbh *db_handler) DeleteUser(logger logx.Logger, id string) *pkg_types.AppError {
	if id == "" {
		logger.Error("invalid argument")
		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if serr := deleteUser(logger, ctx, id); serr != nil {
		return serr
	}
	return nil
}

func getUser(logger logx.Logger, ctx context.Context, id string) (database_type.User, *pkg_types.AppError) {
	ut, err := func() (*twoonone_model.UserTwoonone, *pkg_types.AppError) {
		u, err := configs.Model_TwoOnOne.FindOne(ctx, id)
		if err != nil {
			if myerr, ok := err.(*mysql.MySQLError); ok {
				switch myerr.Number {
				default:
					logger.Errorf("未处理错误：%v", myerr)
				}
			} else {
				switch err {
				default:
					logger.Errorf("未处理错误：%v", err)
				case twoonone_model.ErrNotFound:
					return nil, pkg_types.NewError(twoonone_pb.Error_ERROR_USER_NO_EXIST, "")
				}
			}
			return nil, pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
		}
		return u, nil
	}()
	if err != nil {
		return database_type.User{}, err
	}
	up, err := func() (*database_type.UserPublic, *pkg_types.AppError) {
		u, err := ldap.FindUser(logger, id)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}()
	if err != nil {
		return database_type.User{}, err
	}
	return database_type.User{
		UserTwoonone: *ut,
		UserPublic:   *up,
	}, nil
}

func updateUser(logger logx.Logger, ctx context.Context, u twoonone_model.UserTwoonone) *pkg_types.AppError {
	if err := configs.Model_TwoOnOne.Update(ctx, &u); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%v", myerr)
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误：%v", err)
			case twoonone_model.ErrNotFound:
				return pkg_types.NewError(twoonone_pb.Error_ERROR_USER_NO_EXIST, "")
			}
		}
		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	return nil
}

func deleteUser(logger logx.Logger, ctx context.Context, id string) *pkg_types.AppError {
	if err := configs.Model_TwoOnOne.Delete(ctx, id); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%v", myerr)
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误：%v", err)
			case twoonone_model.ErrNotFound:
				return pkg_types.NewError(twoonone_pb.Error_ERROR_USER_NO_EXIST, "")
			}
		}
		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	return nil
}

func createUser(logger logx.Logger, ctx context.Context, u twoonone_model.UserTwoonone) *pkg_types.AppError {
	if _, err := configs.Model_TwoOnOne.Insert(ctx, &u); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%v", myerr)
			case mysqlerr.ER_DUP_ENTRY:
				return pkg_types.NewError(twoonone_pb.Error_ERROR_USER_EXISTED, "")
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误：%v", err)
			}
		}
		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	return nil
}
