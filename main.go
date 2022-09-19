package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"srun/cfg"
	"srun/dao/mysql"
	"srun/dao/redis"
	"srun/logger"
	"srun/routes"
	"syscall"
	"time"
)

func main() {
	if err := cfg.Init(); err != nil {
		fmt.Println(fmt.Sprintf("Init config error: %s", err.Error()))
	}

	if err := logger.InitLogger(); err != nil {
		fmt.Println(fmt.Sprintf("Init logger error: %s", err.Error()))
	} else {
		defer func(l *zap.Logger) {
			err := l.Sync()
			if err != nil {
				fmt.Println(fmt.Sprintf("close logger error: %s", err.Error()))
			}
		}(zap.L())
		zap.L().Info("Init zap logger successful....")
	}

	if err := mysql.Init(); err != nil {
		fmt.Println(fmt.Sprintf("Init mysql error: %s", err.Error()))
	} else {
		defer mysql.Close()
		zap.L().Info("Init mysql successful....")
	}

	if err := redis.Init(); err != nil {
		fmt.Println(fmt.Sprintf("Init redis error: %s", err.Error()))
	} else {
		defer redis.Close()
		zap.L().Info("Init redis successful....")
	}

	if err := cfg.InitWebAuthn(); err != nil {
		fmt.Println(fmt.Sprintf("Init webauthn error: %s", err.Error()))
		zap.L().Error(fmt.Sprintf("Init webauthn error: %s", err.Error()))
	}

	r := routes.Setup()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if viper.GetString("app.protocol") == "https" {
			if err := server.ListenAndServeTLS(viper.GetString("app.cert_file"), viper.GetString("app.key_file")); err != nil && err != http.ErrServerClosed {
				zap.L().Error(fmt.Sprintf("Listen error: %s", err.Error()))
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				zap.L().Error(fmt.Sprintf("Listen error: %s", err.Error()))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Warn("Shutdown Server....")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error(fmt.Sprintf("Server Shutdown error: %s", err.Error()))
	}

	zap.L().Warn("Server exited")
}
