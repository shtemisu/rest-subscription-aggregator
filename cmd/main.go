package main

import (
	"subAggregator/internal/di"

	"go.uber.org/fx"
)

// @title           Subscription Aggregator API
// @version         1.0
// @description     REST-сервис агрегации онлайн-подписок
// @host            localhost:8080
// @BasePath        /
func main() {
	fx.New(di.Module()).Run()
}
