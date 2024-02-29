package main

import (
	"MURI-GO/dao/mysql"
	"MURI-GO/logger"
	"MURI-GO/router"
	"MURI-GO/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//1、加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Println("init settings failed")
		return
	}

	//2、日志初始化
	if err := logger.Init(); err != nil {
		fmt.Println("init logger failed")
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success .....")

	//初始化mysql
	if err := mysql.Init(); err != nil {
		zap.L().Fatal("init mysql failed", zap.Error(err))
	}
	defer mysql.Close()

	//注册路由
	r := gin.Default()
	router.Router.InitApiRouter(r)
	//启动gin server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("", zap.Error(err))
		}
	}()
	//优雅关闭
	//声明一个系统信号的channel，如果没有信号就一直阻塞，如果有就执行
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//设置ctx超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//关闭gin
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("gin server 关闭异常：", zap.Error(err))

	}

	zap.L().Info("gin server 退出成功")

}
