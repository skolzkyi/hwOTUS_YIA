package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/app"
	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "./configs", "Path to config.yaml")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	err := config.Init(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config: ", config)
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		fmt.Println(err)
	}

	var storage app.Storage
	ctxStor, cancelStore := context.WithTimeout(context.Background(), config.GetdbTimeOut())
	if config.workWithDBStorage {
		//ToDo
	} else {
		storage = memorystorage.New()
		err = storage.Init(ctxStor)
		if err != nil {
			cancelStore()
			logg.Fatal("fatal error of inintialization memory storage: " + err.Error())
			fmt.Println(err)
		}
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, &config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Fatal("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

// cd C:\REPO\Go\!OTUS\hwOTUS_YIA\hw12_13_14_15_calendar
