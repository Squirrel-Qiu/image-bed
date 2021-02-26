package handle

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Squirrel-Qiu/image-bed/dbb"
)

const (
	URL = "https://..."
	IdType = "resource_id"
)


func (impl *Implement) Upload(ctx *gin.Context) {
	content, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logrus.Errorf("read file failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	h := sha256.Sum256(content)
	hash := hex.EncodeToString(h[:])
	ok, resourceId, err := impl.DB.FileIsExistByHash(hash)
	if !ok {
		if err != nil {
			logrus.Errorf("db check file is exist by hash failed: %+v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		logrus.Info("The upload file is exist")
		if _, err := ctx.Writer.WriteString(URL + resourceId); err != nil {
			logrus.Errorf("put object failed: %+v", err)
			return
		}
		return
	}

	resourceId, err = impl.Generator.GenerateId(IdType)
	if err != nil {
		logrus.Errorf("generate id failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	resource := new(dbb.Resource)
	resource.Id = resourceId
	now := time.Now()
	resource.Bucket = now.Format("2006-01")
	resource.CreateTime = now.Format("2006-01-02 15:04")
	resource.Hash = hash
	resource.Size = uint32(len(content))

	err = impl.DB.Store(resource)
	if err != nil {
		logrus.Errorf("db store resource failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	//m := md5.Sum(content)
	//cm := base64.StdEncoding.EncodeToString(m[:])
	reader := bytes.NewReader(content)

	err = impl.Tool.Storage(resourceId, resource.Bucket, reader)
	if err != nil {
		logrus.Errorf("cloud store failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if _, err := ctx.Writer.WriteString(URL + resourceId); err != nil {
		logrus.Errorf("put object failed: %+v", err)
		return
	}
}
