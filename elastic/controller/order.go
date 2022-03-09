package controller

import (
	"elastic/dao"
	"elastic/logger"
	"elastic/model"
	"elastic/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)


type OrderController struct {
	service service.OrderService
}

// NewOrderController
func NewOrderController(esClient *elastic.Client) (*OrderController, error) {
	tag := "NewOrderController"

	orderDao, err := dao.NewOrder(esClient)
	if err != nil {
		logger.E(tag, "NewOrder dao failed, err:%+v", err)
	}

	service := service.NewOrderService(orderDao)
	orderController := &OrderController{
		service: *service,
	}
	
	return orderController, nil
}


// Insert
func (o *OrderController) Insert(c *gin.Context) {
	tag := "Insert"

	orders := model.Orders{}
	if err := c.ShouldBind(&orders); err != nil {
		logger.E(tag, "params bind failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "params bind failed", "error": err.Error()})
		return
	}

	if err := o.service.BatchInsert(c, &orders); err != nil {
		logger.E(tag, "batchInsert failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code": 1001, "msg": "insert failed", "error": err.Error()})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// Update
func (o *OrderController) Update(c *gin.Context) {
	tag := "update"


	orders := model.Orders{}
	if err := c.ShouldBind(&orders); err != nil {
		logger.E(tag, "params bind failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "param bind failed", "error": err.Error()})
		return
	}

	if err := o.service.BatchUpdate(c, &orders); err != nil {
		logger.E(tag, "batchUpdate failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code": 1001, "msg": "update failed", "error": err.Error()})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// Delete
func (o *OrderController) Delete(c *gin.Context) {
	tag := "delete"

	orderIds := c.Query("orderIds")
	ids := make([]int, 0)

	for _, item := range strings.Split(orderIds, ",") {
		id, _ := strconv.Atoi(item)
		ids = append(ids, id)
	}

	if err := o.service.BatchDel(c, ids); err != nil {
		logger.E(tag, "batchDel failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code":1001, "msg": "delete failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// MGet
func (o *OrderController) MGet(c *gin.Context) {
	tag := "MGet"

	orderIds := c.Query("orderIds")
	ids := make([]int, 0)

	for _, item := range strings.Split(orderIds, ",") {
		id, _ := strconv.Atoi(item)
		ids = append(ids, id)
	}

	res, err := o.service.MGet(c, ids)
	if err != nil {
		logger.E(tag, "MGet failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "MGet failed", "error": err.Error()})
		return
	}
	if res == nil {
		logger.E(tag, "search failed, res is nil")
		c.JSON(http.StatusOK, gin.H{"code":1001, "msg": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": res})
}

// Search
func (o *OrderController) Search(c *gin.Context) {
	tag := "search"

	req := model.EsSearchOrderReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.E(tag, "params bind failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code":10000, "msg": "param bind failed", "error": err.Error()})
		return
	}

	res, err := o.service.Search(c, &req)
	if err != nil {
		logger.E(tag, "search failed, err:%+v", err)
		c.JSON(http.StatusOK, gin.H{"code":1001, "msg": "search failed", "error": err.Error()})
		return
	}
	if res == nil {
		logger.E(tag, "search failed, res is nil")
		c.JSON(http.StatusOK, gin.H{"code":1001, "msg": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code":200, "msg": "success", "data": res})
}