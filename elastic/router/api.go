package router

import (
	"elastic/controller"
	"elastic/logger"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

// RegisterRouter
func RegisterRouter(engine *gin.Engine, esClient *elastic.Client) error {
	tag := "registerRouter"

	root := engine.Group("/order")

	orderController, err := controller.NewOrderController(esClient)
	if err != nil {
		logger.E(tag, "newOrderController failed, err:%+v", err)
		return err
	}
	root.POST("/add", orderController.Insert)
	root.POST("/update", orderController.Update)
	root.POST("/delete", orderController.Delete)
	root.GET("/info", orderController.MGet)
	root.POST("search", orderController.Search)

	return nil
}