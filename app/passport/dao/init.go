package dao

import (
	"database/sql"
	"strings"

	"github.com/ipuppet/gtools/database"
	"github.com/ipuppet/gtools/utils"
)

var (
	Db         *sql.DB
	UserColumn string
)

func init() {
	Db = database.ConnectToMySQL("ums")

	columns, _ := utils.GetStructFieldNameToSnake(User{})
	UserColumn = "`" + strings.Join(columns, "`,`") + "`"
}
