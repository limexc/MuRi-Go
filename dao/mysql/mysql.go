package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"                  //gorm库
	_ "github.com/jinzhu/gorm/dialects/mysql" //gorm对应的mysql驱动
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

//db的初始化函数，与数据库建立连接
func Init() (err error) {
	//判断是否已经初始化了
	if isInit {
		return
	}
	//组装连接配置
	//parseTime是查询结果是否自动解析为时间
	//loc是Mysql的时区设置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"))
	//与数据库建立连接，生成一个*gorm.DB类型的对象
	GORM, err = gorm.Open("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}

	//打印sql语句
	GORM.LogMode(viper.GetBool("mysql.LogMode"))

	//开启连接池
	// 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	GORM.DB().SetMaxIdleConns(viper.GetInt("mysql.MaxIdleConns"))
	// 设置了连接可复用的最大时间
	GORM.DB().SetMaxOpenConns(viper.GetInt("mysql.MaxIdleConns"))
	// 设置了连接可复用的最大时间
	GORM.DB().SetConnMaxLifetime(30)

	isInit = true
	//GORM.AutoMigrate(model.Chart{}, model.Event{})
	//GORM.AutoMigrate(model.Agent{})
	zap.L().Info("连接数据库成功")
	return
}

//db的关闭函数
func Close() error {
	//logger.Info("关闭数据库连接")
	zap.L().Info("关闭数据库连接")
	return GORM.Close()
}