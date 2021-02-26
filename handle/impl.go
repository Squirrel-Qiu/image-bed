package handle

import (
	"github.com/Squirrel-Qiu/image-bed/client"
	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/id"
)

type Implement struct {
	DB        dbb.DBApi
	Generator id.Generator
	Tool client.Tool
}
