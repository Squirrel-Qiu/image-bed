package dbb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Squirrel-Qiu/image-bed/dbb/internal"
	"github.com/Squirrel-Qiu/image-bed/store"
)

type DBApi interface {
	GetIdValue(idType string) (idValue int64, err error)
	Store(resource *store.Resource) (err error)
	FileIsExistByHash(hash string) (resourceId string, ok bool)
	FileIsExistById(resourceId string) (bucket string, ok bool, err error)
}

func InitDB(db *sql.DB) (DB DBApi) {
	return &internal.Impl{DB: db}
}
