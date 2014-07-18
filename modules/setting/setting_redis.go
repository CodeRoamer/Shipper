package setting

import (
	_ "github.com/gogits/cache/redis"
	_ "github.com/gogits/session/redis"
)

func init() {
	EnableRedis = true
}
