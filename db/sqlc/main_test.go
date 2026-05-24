package sqlc

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/520wheat/simplebank/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		panic("cannot load config: " + err.Error())
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}

	testQueries = New(connPool)
	os.Exit(m.Run())
}