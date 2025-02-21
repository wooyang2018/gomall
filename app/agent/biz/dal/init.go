package dal

import (
	"github.com/cloudwego/biz-demo/gomall/app/agent/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/agent/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
