package dao

import (
	"database/sql"
	"log"

	"github.com/ipuppet/gtools/database"
	"github.com/ipuppet/gtools/utils"
)

var (
	Db     *sql.DB
	Logger *log.Logger
)

func init() {
	Db = database.ConnectToMySQL("")
	Logger = utils.Logger("admin")
}
