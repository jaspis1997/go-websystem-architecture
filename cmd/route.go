package main

import (
	"net/http"
	"playground/internal/external/web"
)

const (
	PathRoot = "/"
)

func registerRoutes() error {
	engine.GET(PathRoot, func(c web.Context) {
		c.Status(http.StatusOK)
	})
	return nil
}
