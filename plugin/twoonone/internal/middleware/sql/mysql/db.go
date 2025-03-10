package mysql

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/configs"
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

func (s *db_action_inc_coin) Merge(logger logx.Logger, u *twoonone_model.Twoonone) {
	if s.float64 < 0 {
		s.float64 = 0
	}
	u.Coin += s.float64
}

func (s *db_action_dec_coin) Merge(logger logx.Logger, u *twoonone_model.Twoonone) {
	if s.float64 < 0 {
		s.float64 = 0
	}
	u.Coin -= s.float64
}

func (s *db_action_inc_wincount) Merge(logger logx.Logger, u *twoonone_model.Twoonone) {
	u.Wincount += int64(s.uint)
}

func (s *db_action_dec_wincount) Merge(logger logx.Logger, u *twoonone_model.Twoonone) {
	u.Wincount -= int64(s.uint)
}

func (s *db_action_inc_losecount) Merge(logger logx.Logger, u *twoonone_model.Twoonone) {
	u.Losecount += int64(s.uint)
}

func (s *db_action_dec_losecount) Merge(logger logx.Logger, u *twoonone_model.Twoonone) {
	u.Losecount -= int64(s.uint)
}

func (s *db_action_update_getdaily_time) Merge(logger logx.Logger, u *twoonone_model.Twoonone) {
	u.LastGetdaliyTime = s.Time.Unix()
}

func (dbh *db_handler) GetUser(logger logx.Logger, id string) (twoonone_model.Twoonone, error) {
	if id == "" {
		logger.Error("invalid argument")
		return twoonone_model.Twoonone{}, pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	u, err := getUser(logger, ctx, id)
	if err != nil {
		return twoonone_model.Twoonone{}, err
	}
	return u, nil
}

func (dbh *db_handler) CreateUser(logger logx.Logger, id string) error {
	return pkg_types.NewError(twoonone_pb.Error_ERROR_UNKNOWN, "不支持操作")
}

func (dbh *db_handler) UpdateUser(logger logx.Logger, id string, actions ...database_type.Action) error {
	if id == "" {
		logger.Error("invalid argument")
		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	u, err := func() (twoonone_model.Twoonone, error) {
		u, err := getUser(logger, ctx, id)
		if err != nil {
			return twoonone_model.Twoonone{}, err
		}
		for _, v := range actions {
			v.Merge(logger, &u)
		}
		return u, nil
	}()
	if err != nil {
		return err
	}
	u.Id = id
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if serr := updateUser(logger, ctx, u); serr != nil {
		return serr
	}
	return nil
}

func (dbh *db_handler) DeleteUser(logger logx.Logger, id string) error {
	return pkg_types.NewError(twoonone_pb.Error_ERROR_UNKNOWN, "不支持操作")
}

func getUser(logger logx.Logger, ctx context.Context, id string) (twoonone_model.Twoonone, error) {
	ut, err := func() (*twoonone_model.Twoonone, error) {
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
		return twoonone_model.Twoonone{}, err
	}
	return *ut, nil
}

func updateUser(logger logx.Logger, ctx context.Context, u twoonone_model.Twoonone) error {
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

// func deleteUser(logger logx.Logger, ctx context.Context, id string) error {
// 	if err := configs.Model_TwoOnOne.Delete(ctx, id); err != nil {
// 		if myerr, ok := err.(*mysql.MySQLError); ok {
// 			switch myerr.Number {
// 			default:
// 				logger.Errorf("未处理错误：%v", myerr)
// 			}
// 		} else {
// 			switch err {
// 			default:
// 				logger.Errorf("未处理错误：%v", err)
// 			case twoonone_model.ErrNotFound:
// 				return pkg_types.NewError(twoonone_pb.Error_ERROR_USER_NO_EXIST, "")
// 			}
// 		}
// 		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
// 	}
// 	return nil
// }

// func createUser(logger logx.Logger, ctx context.Context, u twoonone_model.Twoonone) error {
// 	if _, err := configs.Model_TwoOnOne.Insert(ctx, &u); err != nil {
// 		if myerr, ok := err.(*mysql.MySQLError); ok {
// 			switch myerr.Number {
// 			default:
// 				logger.Errorf("未处理错误：%v", myerr)
// 			case mysqlerr.ER_DUP_ENTRY:
// 				return pkg_types.NewError(twoonone_pb.Error_ERROR_USER_EXISTED, "")
// 			}
// 		} else {
// 			switch err {
// 			default:
// 				logger.Errorf("未处理错误：%v", err)
// 			}
// 		}
// 		return pkg_types.NewError(twoonone_pb.Error_ERROR_UNDEFINED, "")
// 	}
// 	return nil
// }
