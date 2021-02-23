package dbb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Squirrel-Qiu/image-bed/handle"
)

type DBApi interface {
	GetIdValue(idType string) (idValue int64, err error)
	Store(resource *handle.Resource) (err error)
	FileIsExistByHash(hash string) (ok bool, resourceId string, err error)
	FileIsExistById(resourceId string) (ok bool, bucket string, err error)
}

func InitDB(db *sql.DB) (DB DBApi) {
	return &Impl{DB: db}
}
