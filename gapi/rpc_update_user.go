package gapi

import (
	"context"
	"time"

	db "github.com/520wheat/simplebank/db/sqlc"
	"github.com/520wheat/simplebank/pb"
	"github.com/520wheat/simplebank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = pgtype.Text{String: hashedPassword, Valid: true}
		arg.PasswordChangedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	}

	if req.FullName != nil {
		arg.FullName = pgtype.Text{String: req.GetFullName(), Valid: true}
	}

	if req.Email != nil {
		arg.Email = pgtype.Text{String: req.GetEmail(), Valid: true}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
