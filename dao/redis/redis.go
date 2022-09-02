package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Rds *redis.Client

// Init 初始化Redis
func Init() (err error) {
	Rds = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.ip"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.pwd"),
		DB:       viper.GetInt("redis.index"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	_, err = Rds.Ping().Result()

	if err != nil {
		zap.L().Error(fmt.Sprintf("redis init error: %s", err.Error()))
	} else {
		fmt.Println("Init redis succeed....")
		zap.L().Info(fmt.Sprintf("Redis【%s:%d】init finished....", viper.GetString("redis.ip"), viper.GetInt("redis.port")))
	}

	return
}

func Close() {
	_ = Rds.Close()
}

func GetRds() *redis.Client {
	if _, err := Rds.Ping().Result(); err != nil {
		fmt.Println(fmt.Sprintf("%s -> Reconnecting redis....", err))
		_ = Init()
		//defer Close()

		if _, err := Rds.Ping().Result(); err != nil {
			zap.L().Error(fmt.Sprintf("Redis closed error: %s", err))
		}
	}

	return Rds
}
