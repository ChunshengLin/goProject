package service

import (
	"context"
	"elastic/dao"
	"elastic/model"
)


type OrderService struct {
	order *dao.Order
}

// NewOrderService
func NewOrderService(order *dao.Order) *OrderService {
	orderService := &OrderService{
		order: order,
	}

	return orderService
}

// BatchInsert
func (o *OrderService) BatchInsert(ctx context.Context, orders *model.Orders) error {
	return o.order.BatchInsert(ctx, orders)
}

// BatchUpdate
func (o *OrderService) BatchUpdate(ctx context.Context, orders *model.Orders) error {
	return o.order.BatchUpdate(ctx, orders)
}

// BatchDel
func (o *OrderService) BatchDel(ctx context.Context, ids []int) error {
	return o.order.BatchDel(ctx, ids)
}

// MGet
func (o *OrderService) MGet(ctx context.Context, ids []int) (*model.Orders, error) {
	return o.order.MGet(ctx, ids)
}

// Search
func (o *OrderService) Search(ctx context.Context, req *model.EsSearchOrderReq) (*model.Orders, error) {
	return o.order.Search(ctx, req.CreateQuerys())
}