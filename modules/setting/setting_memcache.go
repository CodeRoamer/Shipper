package setting

import (
	_ "github.com/gogits/cache/memcache"
)

func init() {
	EnableMemcache = true
}
