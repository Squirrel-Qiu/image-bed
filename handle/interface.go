package handle

import (
	"github.com/gin-gonic/gin"

	"github.com/Squirrel-Qiu/image-bed/client"
	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/id"
)

type Api interface {
	Get(ctx *gin.Context)
	Upload(ctx *gin.Context)
}

func New(dbInstance dbb.DBApi, generator id.Generator, tool client.Tool) Api {
	return &Implement{DB: dbInstance, Generator: generator, Tool: tool}
}
