package database

import (
	"database/sql"
	"log"

	userPkg "ErotsServer/app/user/pkg"

	"github.com/ipuppet/gtools/utils"
)

var (
	Db     *sql.DB
	Logger *log.Logger
)

func init() {
	Db = userPkg.Db
	Logger = utils.Logger("admin")
}
