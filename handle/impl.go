package handle

import (
	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/id"
	"github.com/Squirrel-Qiu/image-bed/store"
)

type Implement struct {
	DB dbb.DBApi
	Generator id.Generator
	Cred *store.Credential
}
