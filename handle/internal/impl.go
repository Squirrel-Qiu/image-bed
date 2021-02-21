package internal

import (
	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/store"
)

type Implement struct {
	DB dbb.DBApi
	Cred *store.Credential
}
