package main

import (
	"context"
	"notify/api"
	"notify/cmd/config"
	"notify/internal/storage"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.GetConfig(&config.Config{})       // конфиг для подлючения к postgresSQL через config.yaml
	connect := storage.PostgresConnect(cfg.Storage) // подключение к postgresSQL
	connectStorage := storage.NewDB(connect)        // работа с БД
	ctx, cancel := context.WithCancel(context.Background())
	httpConfig := config.GetLocal(&config.HTTPConfig{}) // конфиг для подлючения к серверу через local.yaml
	go func() {
		for {
			api.NewApi(ctx, httpConfig, connectStorage)
		}
	}()

	osExit := make(chan os.Signal, 1)
	exitOnSignal(osExit)
	<-osExit
	cancel()
	time.Sleep(time.Second)
}

func exitOnSignal(osExit chan os.Signal) {
	signal.Notify(osExit, os.Interrupt, os.Kill)
}
