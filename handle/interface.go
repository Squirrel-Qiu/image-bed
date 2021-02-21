package handle

import (
	"github.com/gin-gonic/gin"

	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/handle/internal"
	"github.com/Squirrel-Qiu/image-bed/store"
)

type Api interface {
	Get(ctx *gin.Context)
	Upload(ctx *gin.Context)
}

func New(dbInstance dbb.DBApi, credential *store.Credential) Api {
	return &internal.Implement{DB: dbInstance, Cred: credential}
}
