package main

import (
	"go.uber.org/fx"
	"video-balancer/internal/app"
)

func main() {
	fx.New(app.NewApp()).Run()
}
