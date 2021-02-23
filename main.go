package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/Squirrel-Qiu/image-bed/conf"
	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/handle"
	"github.com/Squirrel-Qiu/image-bed/id"
)

func main() {
	run()
}

func run() {
	url, user, password, listenAddr, credential := conf.ReadConf()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/image_bed", user, password, url))
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%+v", err)
	}

	dbInstance := dbb.InitDB(db)

	router := gin.New()
	router.Use(limits.RequestSizeLimiter(32 * 1024 * 1024))

	api := handle.New(dbInstance, &id.Generate{}, credential)
	router.POST("upload", api.Upload)
	router.GET("/get/:resourceId", api.Get)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // Interrupt Signal = syscall.SIGINT

	httpServer := http.Server{Handler: router, Addr: listenAddr}

	shutdownChan := make(chan struct{})

	go func() {
		<-signalChan
		timeout, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelFunc()
		if err := httpServer.Shutdown(timeout); err != nil {
			log.Println(err)
		}
		close(shutdownChan)
	}()

	log.Println("start http server")

	err = httpServer.ListenAndServe()
	switch err {
	case http.ErrServerClosed:
		<-shutdownChan

	default:
		log.Println(err)
	}

	if err := db.Close(); err != nil {
		log.Printf("%+v", xerrors.Errorf("close db failed: %w", err))
	}
}
