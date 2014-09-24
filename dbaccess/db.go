package dbaccess

import (
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))

// CacheIsRecent cheches for the given time to be within 5 seconds of time.Now()
func CacheIsRecent(cacheTime time.Time) bool {
	return time.Now().Sub(cacheTime) < time.Duration(5)*time.Second
}
