package sqlc

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/520wheat/simplebank/util"
)

var testQueries *Queries
var testConnPool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		panic("cannot load config: " + err.Error())
	}

	testConnPool, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}

	testQueries = New(testConnPool)
	os.Exit(m.Run())
}