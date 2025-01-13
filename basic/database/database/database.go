package database

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
	"github.com/nanachi-sh/susubot-code/basic/database/db"
	"github.com/nanachi-sh/susubot-code/basic/database/log"
	database_pb "github.com/nanachi-sh/susubot-code/basic/database/protos/database"
)

var (
	logger = log.Get()
)

func Uno_CreateUser(req *database_pb.Uno_CreateUserRequest) *database_pb.Uno_CreateUserResponse {
	if req.Id == "" || req.Name == "" || req.Password == "" {
		return &database_pb.Uno_CreateUserResponse{
			Body: &database_pb.Uno_CreateUserResponse_Err{
				Err: database_pb.Errors_ValueError,
			},
		}
	}
	if err := db.Uno_CreateUser(req.Id, req.Name, req.Password); err != nil {
		switch (err.(sqlite3.Error)).Code {
		case sqlite3.ErrConstraint:
			return &database_pb.Uno_CreateUserResponse{
				Body: &database_pb.Uno_CreateUserResponse_Err{
					Err: database_pb.Errors_UserExist,
				},
			}
		}
		logger.Println(err)
		return &database_pb.Uno_CreateUserResponse{Body: &database_pb.Uno_CreateUserResponse_Err{
			Err: database_pb.Errors_Undefined,
		}}
	}
	return &database_pb.Uno_CreateUserResponse{
		Body: &database_pb.Uno_CreateUserResponse_Ok{Ok: true},
	}
}

func Uno_GetUser(req *database_pb.Uno_GetUserRequest) *database_pb.Uno_GetUserResponse {
	if req.Id == "" {
		return &database_pb.Uno_GetUserResponse{
			Body: &database_pb.Uno_GetUserResponse_Err{
				Err: database_pb.Errors_ValueError,
			},
		}
	}
	ui, err := db.Uno_GetUser(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &database_pb.Uno_GetUserResponse{
				Body: &database_pb.Uno_GetUserResponse_Err{
					Err: database_pb.Errors_UserNoExist,
				},
			}
		}
	}
	return &database_pb.Uno_GetUserResponse{
		Body: &database_pb.Uno_GetUserResponse_Userinfo{
			Userinfo: &database_pb.Uno_UserInfo{
				Id:        ui.Id,
				Name:      ui.Name,
				WinCount:  int32(ui.WinCount),
				LoseCount: int32(ui.LoseCount),
			},
		},
	}
}

func Uno_UpdateUser(req *database_pb.Uno_UpdateUserRequest) *database_pb.Uno_UpdateUserResponse {
	if req.Userinfo == nil || (req.Userinfo.Id == "" || req.Userinfo.Name == "" || req.Userinfo.WinCount < 0 || req.Userinfo.LoseCount < 0) {
		return &database_pb.Uno_UpdateUserResponse{
			Body: &database_pb.Uno_UpdateUserResponse_Err{
				Err: database_pb.Errors_ValueError,
			},
		}
	}
	if err := db.Uno_UpdateUser(db.Uno_UserInfo{
		Id:        req.Userinfo.Id,
		Name:      req.Userinfo.Name,
		WinCount:  int(req.Userinfo.WinCount),
		LoseCount: int(req.Userinfo.LoseCount),
	}); err != nil {
		if err == sql.ErrNoRows {
			return &database_pb.Uno_UpdateUserResponse{
				Body: &database_pb.Uno_UpdateUserResponse_Err{
					Err: database_pb.Errors_UserNoExist,
				},
			}
		}
		logger.Println(err)
		return &database_pb.Uno_UpdateUserResponse{Body: &database_pb.Uno_UpdateUserResponse_Err{
			Err: database_pb.Errors_Undefined,
		}}
	}
	return &database_pb.Uno_UpdateUserResponse{
		Body: &database_pb.Uno_UpdateUserResponse_Ok{Ok: true},
	}
}

func Uno_ChangePassword(req *database_pb.Uno_ChangePasswordRequest) *database_pb.Uno_ChangePasswordResponse {
	if req.Id == "" || req.NewPassword == "" || req.OldPassword == "" {
		return &database_pb.Uno_ChangePasswordResponse{
			Body: &database_pb.Uno_ChangePasswordResponse_Err{Err: database_pb.Errors_ValueError},
		}
	}
	ok, err := db.Uno_ChangePassword(req.Id, req.NewPassword, req.OldPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return &database_pb.Uno_ChangePasswordResponse{
				Body: &database_pb.Uno_ChangePasswordResponse_Err{
					Err: database_pb.Errors_UserNoExist,
				},
			}
		}
		logger.Println(err)
		return &database_pb.Uno_ChangePasswordResponse{Body: &database_pb.Uno_ChangePasswordResponse_Err{
			Err: database_pb.Errors_Undefined,
		}}
	}
	if !ok {
		return &database_pb.Uno_ChangePasswordResponse{
			Body: &database_pb.Uno_ChangePasswordResponse_Err{
				Err: database_pb.Errors_UserPasswordWrong,
			},
		}
	}
	return &database_pb.Uno_ChangePasswordResponse{
		Body: &database_pb.Uno_ChangePasswordResponse_Ok{
			Ok: true,
		},
	}
}

func Uno_VerifyUser(req *database_pb.Uno_VerifyUserRequest) *database_pb.Uno_VerifyUserResponse {
	if req.Id == "" || req.Password == "" {
		return &database_pb.Uno_VerifyUserResponse{
			Body: &database_pb.Uno_VerifyUserResponse_Err{
				Err: database_pb.Errors_ValueError,
			},
		}
	}
	ok, err := db.Uno_VerifyUser(req.Id, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return &database_pb.Uno_VerifyUserResponse{
				Body: &database_pb.Uno_VerifyUserResponse_Err{
					Err: database_pb.Errors_UserNoExist,
				},
			}
		}
		logger.Println(err)
		return &database_pb.Uno_VerifyUserResponse{Body: &database_pb.Uno_VerifyUserResponse_Err{
			Err: database_pb.Errors_Undefined,
		}}
	}
	if !ok {
		return &database_pb.Uno_VerifyUserResponse{
			Body: &database_pb.Uno_VerifyUserResponse_Err{
				Err: database_pb.Errors_UserPasswordWrong,
			},
		}
	}
	return &database_pb.Uno_VerifyUserResponse{
		Body: &database_pb.Uno_VerifyUserResponse_Ok{
			Ok: true,
		},
	}
}
