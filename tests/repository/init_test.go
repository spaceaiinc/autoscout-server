package repository_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/tests"
)

var (
	dbm    *database.DB
	cfg    config.Config
	helper *tests.Helper
)

func TestMain(m *testing.M) {
	time.Local = utility.Tokyo

	// テスト用のDBを設定
	os.Setenv("DB_NAME", "autoscout_test")

	// envファイルローディング
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	dbm = database.NewDB(cfg.DB, false)
	helper = tests.NewHelper(cfg)

	// Do tests
	code := m.Run()

	os.Exit(code)
}
