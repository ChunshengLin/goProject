package main

import (
	"elastic/router"
	"elastic/util"

	"github.com/gin-gonic/gin"
)


func main() {
	esClient, err := util.NewEsClient()
	if err != nil {
		panic(err)
	}

	engine := gin.New()
	gin.SetMode(gin.ReleaseMode)

	if err := router.RegisterRouter(engine, esClient); err != nil {
		panic(err)
	}

	engine.Run()
}