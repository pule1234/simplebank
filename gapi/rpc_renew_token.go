package gapi

import (
	"context"
	"github.com/pule1234/simple_bank/pb"
	"github.com/pule1234/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	violidations := validateReNewAccessToken(req)
	if violidations != nil {
		return nil, invalidArgumentError(violidations)
	}
	//验证refrehtoken是否有效
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		return nil, err
	}

	if session.IsBlocked {
		return nil, status.Error(codes.Unauthenticated, "session is blocked")
	}

	if session.Username != refreshPayload.Username {
		return nil, status.Error(codes.PermissionDenied, "session username does not match")
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, status.Error(codes.Unauthenticated, "refresh token invalid")
	}

	//验证当前refreshtoken是否过期
	if time.Now().After(refreshPayload.ExpiredAt) {
		return nil, status.Error(codes.PermissionDenied, "session expired")
	}

	//当前token有效生成accessToken
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "can not create access token")
	}

	rsp := &pb.RenewAccessTokenResponse{
		AccessToken:         accessToken,
		AccessTokenExpireAt: timestamppb.New(accessPayload.ExpiredAt),
	}
	return rsp, nil
}

func validateReNewAccessToken(req *pb.RenewAccessTokenRequest) (violidations []*errdetails.BadRequest_FieldViolation) {
	// 校验refreshToken是否有效
	if err := val.ValidateRefreshToken(req.GetRefreshToken()); err != nil {
		violidations = append(violidations, fieldViolation("refreshToken", err))
	}
	return violidations
}
