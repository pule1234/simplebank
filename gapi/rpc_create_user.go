package gapi

import (
	"context"
	"github.com/lib/pq"
	db "github.com/pule1234/simple_bank/db/sqlc"
	"github.com/pule1234/simple_bank/pb"
	"github.com/pule1234/simple_bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 可直接进行调用
func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	// todo 对参数进行校验
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %s", pqErr.Detail)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	//传地址
	res := &pb.CreateUserResponse{
		User: converter(user),
	}
	return res, nil
}
