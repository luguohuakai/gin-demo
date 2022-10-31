package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var DBErr error

// Init 初始化mysql
func Init() error {
	var i = 0
	for {
		if DB, DBErr = gorm.Open(
			"mysql",
			fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=2s&readTimeout=2s",
				viper.GetString("mysql.user"),
				viper.GetString("mysql.pwd"),
				viper.GetString("mysql.ip"),
				viper.GetInt("mysql.port"),
				viper.GetString("mysql.dbname"),
			),
		); DBErr != nil {
			zap.L().Error(fmt.Sprintf("mysql error: %s", DBErr.Error()))
			if i > 1 {
				break
			}
			time.Sleep(time.Second * 2)
			i++
			continue
		}
		if DBErr = DB.DB().Ping(); DBErr != nil {
			// 需进行重连操作
			zap.L().Error(fmt.Sprintf("mysql error: %s", DBErr.Error()))
			if i > 1 {
				break
			}
			time.Sleep(time.Second * 2)
			i++
			continue
		}
		break
	}

	if err := DB.DB().Ping(); err == nil {
		DBErr = nil
	}

	// 表结构默认会带s形成复数结构，设置取消复数形式
	DB.SingularTable(true)
	//DB.SetLogger(log.StandardLogger())
	//DB.LogMode(true)
	DB.DB().SetConnMaxLifetime(time.Duration(viper.GetInt("max_life_time")) * time.Minute)
	DB.DB().SetMaxOpenConns(viper.GetInt("max_open"))
	DB.DB().SetMaxIdleConns(viper.GetInt("max_idle"))

	if DBErr == nil {
		zap.L().Info(fmt.Sprintf("Mysql【%s:%d】init finished....", viper.GetString("mysql.ip"), viper.GetInt("mysql.port")))
	}

	return DBErr
}

func GetDB() *gorm.DB {
	if DB == nil {
		_ = Init()
	} else if err := DB.DB().Ping(); err != nil {
		zap.L().Error(fmt.Sprintf("DB.DB().Ping() error: %s ", err.Error()))
		// 需进行重连操作
		_ = Init()
	}

	return DB
}

func Close() {
	_ = DB.Close()
}
