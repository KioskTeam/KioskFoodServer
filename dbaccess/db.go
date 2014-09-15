package dbaccess

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))
