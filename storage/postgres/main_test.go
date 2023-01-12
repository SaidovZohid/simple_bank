package postgres_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/SaidovZohid/simple_bank/config"
	"gitlab.com/SaidovZohid/simple_bank/storage"
)

var (
	dbManager storage.StorageI
)

func TestMain(m *testing.M) {
	cfg := config.Load("./../..")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	dbManager = storage.NewStorage(psqlConn)

	os.Exit(m.Run())
}
