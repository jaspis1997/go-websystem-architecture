package main

import (
	"playground/internal/entities"
	"playground/internal/external/log"
	"playground/internal/external/web"
)

var (
	HostPort func() (host string, port uint16)
)

var (
	// データ永続化処理を行う
	repository *entities.Repository
	// HTTP　リクエストのルーティングを行う
	engine web.Engine
	// スケジューリングされた処理を行う
	batch entities.BatchService
)

func main() {
	if err := initialize(); err != nil {
		log.Fatal(err)
	}
	if batch != nil {
		if err := batch.Start(); err != nil {
			log.Fatal(err)
		}
	}
	if err := registerRoutes(); err != nil {
		log.Fatal(err)
	}
	if err := engine.Run(HostPort()); err != nil {
		log.Fatal(err)
	}
	if batch != nil {
		if err := batch.Stop(); err != nil {
			log.Fatal(err)
		}
	}
}
