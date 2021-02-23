package handle

import (
	"github.com/gin-gonic/gin"

	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/id"
	"github.com/Squirrel-Qiu/image-bed/store"
)

type Api interface {
	Get(ctx *gin.Context)
	Upload(ctx *gin.Context)
}

func New(dbInstance dbb.DBApi, generator id.Generator, credential *store.Credential) Api {
	return &Implement{DB: dbInstance, Generator: generator, Cred: credential}
}
