package handle

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (impl *Implement) Get(ctx *gin.Context) {
	resourceId := ctx.Param("resourceId")

	ok, bucket, err := impl.DB.FileIsExistById(resourceId)
	if !ok {
		if err != nil {
			logrus.Errorf("db check file is exist by id failed: %+v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		logrus.Info("the required file is not exist")
		ctx.Status(http.StatusBadRequest)
		return
	}

	reader, err := impl.Tool.Take(resourceId, bucket)
	if err != nil {
		logrus.Errorf("cloud get failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	logrus.Info("get successfully")

	if _, err := io.Copy(ctx.Writer, reader); err != nil {
		logrus.Errorf("copy object failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
