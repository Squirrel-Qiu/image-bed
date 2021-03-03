package id

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"sync"

	"golang.org/x/xerrors"

	"github.com/Squirrel-Qiu/image-bed/dbb"
)

type Generator interface {
	GenerateId(idType string) (id string, err error)
}

type Generate struct {
	DB dbb.DBApi
	IdList []string
}

func (g *Generate) GenerateId(idType string) (id string, err error) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if len(g.IdList) == 0 {
		idValue, err := g.DB.GetIdValue(idType)
		if err != nil {
			return "", xerrors.Errorf("get id value failed: %w", err)
		}

		for i := idValue - 10; i <= idValue; i++ {
			buff := make([]byte, 8)
			binary.BigEndian.PutUint64(buff, uint64(i))
			b := md5.Sum(buff)
			m := hex.EncodeToString(b[:])
			g.IdList = append(g.IdList, m[:10])
		}
	}

	id = g.IdList[len(g.IdList)-1]
	g.IdList = g.IdList[:len(g.IdList)-1]

	return id, nil
}
