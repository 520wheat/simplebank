package gapi

import (
	"context"

	db "github.com/520wheat/simplebank/db/sqlc"
	"github.com/520wheat/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	result, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email: %s", err)
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	}
	return rsp, nil
}
