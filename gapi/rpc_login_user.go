package gapi

import (
	"context"
	db "github.com/pule1234/simple_bank/db/sqlc"
	"github.com/pule1234/simple_bank/pb"
	"github.com/pule1234/simple_bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 完成
func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	//密码校验
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//创建accesstoken
	accesstoken, accesspayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	// 存储到session中
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := &pb.LoginUserResponse{
		User:                  converter(user),
		SessionId:             session.ID.String(),
		AccessToken:           accesstoken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accesspayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(accesspayload.ExpiredAt),
	}
	return res, nil
}
