package gapi

import (
	"fmt"

	db "github.com/520wheat/simplebank/db/sqlc"
	"github.com/520wheat/simplebank/pb"
	"github.com/520wheat/simplebank/token"
	"github.com/520wheat/simplebank/util"
	"github.com/520wheat/simplebank/worker"
)

// Server 实现 gRPC 服务接口
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer 创建 gRPC 服务实例
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
