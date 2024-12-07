package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"playground/internal/entities/constant"
	"playground/internal/external/log"
	"playground/internal/logics"

	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
)

func SetOutput(output io.Writer) {
	gin.DefaultWriter = output
}

func SetDebugMode() {
	gin.SetMode(gin.DebugMode)
}

func SetReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}

type (
	HandlerFunc = gin.HandlerFunc
	Context     = *gin.Context
)

type Engine interface {
	gin.IRouter
	Run(host string, port uint16) error
}

type engine struct {
	*gin.Engine
}

func (e *engine) Run(host string, port uint16) error {
	addr := func() string {
		return fmt.Sprintf("%s:%d", host, port)
	}
	server := &http.Server{
		Addr:    addr(),
		Handler: e.Engine.Handler(),
	}
	var wg sync.WaitGroup
	go gracefulShutdown(server, &wg)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	wg.Wait()
	return nil
}

func New() Engine {
	return &engine{gin.New()}
}

func gracefulShutdown(server *http.Server, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	log.Infof(constant.FormatLogStartShutdown, constant.DefaultShutdownTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultShutdownTimeout)
	defer cancel()
	go forceShutdown(ctx)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf(constant.FormatLogFailedShutdown, err)
	}

	logics.Done()
	log.Info(constant.LogLogicsProcessingCompleted)

	log.Info(constant.LogShutdownCompleted)
}

func forceShutdown(ctx context.Context) {
	<-ctx.Done()
	log.Fatal("force shutdown")
}
