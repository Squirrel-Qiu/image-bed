package store

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"time"

	"golang.org/x/xerrors"

	"github.com/Squirrel-Qiu/image-bed/dbb"
)

const ResourceId = "resource_id"

type Resource struct {
	Id         string
	Bucket     string
	Hash       string
	CreateTime string
	Size       uint32
}

type Generate struct {
	IdType string
	IdList []string
}

func Store(hash string, db dbb.DBApi) (bucket, resourceId string, err error) {
	g := Generate{IdType: ResourceId}

	// id/generate
	if len(g.IdList) == 0 {
		idValue, err := db.GetIdValue(g.IdType) // 原子操作
		if err != nil {
			return "", "", xerrors.Errorf("get id value failed: %w", err)
		}

		for i := idValue - 10; i <= idValue; i++ {
			buff := make([]byte, 8)
			binary.BigEndian.PutUint64(buff, uint64(i))
			m := fmt.Sprintf("%x", md5.Sum(buff))
			g.IdList = append(g.IdList, m[:10])
		}
	}

	var resource Resource
	resource.Id = g.IdList[len(g.IdList)-1]
	g.IdList = g.IdList[:len(g.IdList)]

	resource.Bucket = time.Now().Format("2006-01")
	resource.CreateTime = time.Now().Format("2006-01-02 15:04")
	resource.Hash = hash
	//resource.Size = uint32(len(data))

	if err = db.Store(&resource); err != nil {
		return "", "", xerrors.Errorf("db store failed: %w", err)
	}

	return resource.Bucket, resource.Id, nil
}
