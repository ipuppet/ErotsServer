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
	Logger = utils.Logger("user")
}

type Role struct {
	RoleId      int    `db:"role_id" form:"role_id" json:"role_id" binding:"number"`
	Name        string `db:"name" form:"name" json:"name" binding:"required"`
	Description string `db:"description" form:"description" json:"description" binding:"-"`
}

type Permission struct {
	RoleId       int    `db:"role_id" json:"role_id" binding:"-"`
	PermissionId int    `db:"permission_id" form:"permission_id" json:"permission_id" binding:"number"`
	Module       string `db:"module" form:"module" json:"module" binding:"required"`
	Description  string `db:"description" form:"description" json:"description" binding:"-"`
}
