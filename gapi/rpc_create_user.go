package gapi

import (
	"context"
	"time"

	db "github.com/520wheat/simplebank/db/sqlc"
	"github.com/520wheat/simplebank/pb"
	"github.com/520wheat/simplebank/util"
	"github.com/520wheat/simplebank/worker"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		Role:           "depositor",
	}

	user, err := server.store.CreateUser(ctx, arg)
        if err != nil {
                return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
        }

	taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
	}
	err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
        if err != nil {
		log.Error().Err(err).Msg("failed to enqueue verify email task")
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}