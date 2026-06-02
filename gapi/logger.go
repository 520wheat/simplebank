package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	if err != nil {
		log.Error().Err(err).
			Str("method", info.FullMethod).
			Dur("duration", duration).
			Msg("gRPC request failed")
	} else {
		log.Info().
			Str("method", info.FullMethod).
			Dur("duration", duration).
			Msg("gRPC request success")
	}

	return result, err
}