package dao

import (
	"context"
	"elastic/logger"
	"elastic/model"
	"elastic/util"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
)


const (
	golangDate = "2006-01-02 15:04:05"
	retryTime = 3
	project = "elastic"
	dataName = "order"
	bodySetting = `{
		"settings":{
			"index":{
				"number_of_shards":2,
				"number_of_replicas":1
			}
		},
		"mappings":{
			"properties":{
				"order_id":		{"type": "long"},
				"user_name":	{"type": "text"},
				"price":		{"type": "long"},
				"porduct_id":	{"type": "long"},
				"product_name":	{"type": "text"},
				"create_time":	{"type": "date", "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"},
				"update_time":	{"type": "date", "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"}
			}
		}
	}`
)

type Order struct {
	index 		string
	bodySetting string 
	client 		*elastic.Client
}

// NewOrder 
func NewOrder(client *elastic.Client) (*Order, error) {
	tag := "NewOrder"

	index := fmt.Sprintf("%s_%s", project, dataName)
	order := &Order{
		index: 			index,
		bodySetting:	bodySetting,
		client:			client,
	}

	if err := order.validIndex(); err != nil {
		logger.E(tag, "validIndex failed, err%+v", err)
		return nil, err
	}

	return order, nil
}

// validIndex 校验index是否存在，不存在则创建
func (o *Order) validIndex() error {
	tag := "validIndex"

	ctx := context.Background()
	exist, err := o.client.IndexExists(o.index).Do(ctx)
	if err != nil {
		logger.E(tag, "valid index exist failed, err:%+v", err)
		return err
	}
	
	if !exist {
		_, err := o.client.CreateIndex(o.index).Body(o.bodySetting).Do(ctx)
		if err != nil {
			logger.E(tag, "create index failed, err:%+v", err)
			return err
		}

		logger.I(tag, "create index successed, index:%+v", o.index)
		return nil
	}

	logger.I(tag, "index exist, index:%+v", o.index)
	return nil
}

// BatchInsert 
func (o *Order) BatchInsert(ctx context.Context, orders *model.Orders) error {
	var err error
	for i := 0; i < retryTime; i++ {
		if err := o.batchInsert(ctx, orders); err != nil {
			continue
		}
		return nil
	}
	return err
}

// batchInsert
func (o *Order) batchInsert(ctx context.Context, orders *model.Orders) error {
	req := o.client.Bulk().Index(o.index)
	for _, order := range orders.Data {
		now := time.Now()
		order.CreateTime = now.Format(golangDate)
		order.UpdateTime = now.Format(golangDate)

		snakedOrder := util.SnakeData(order)

		doc := elastic.NewBulkIndexRequest().Id(strconv.FormatInt(int64(order.OrderId), 10)).Doc(snakedOrder)
		req.Add(doc)
	}

	if req.NumberOfActions() <= 0 {
		return nil
	}

	resp, err := req.Do(ctx)
	if err != nil {
		return err
	}

	if resp.Errors {
		return util.GetEsFailedErr(resp)
	}

	return nil
}

// BatchUpdate
func (o *Order) BatchUpdate(ctx context.Context, orders *model.Orders) error {
	var err error
	for i := 0; i < retryTime; i++ {
		if err = o.batchUpdate(ctx, orders); err != nil {
			continue
		}
		return nil
	}

	return err
}

// batchUpdate
func (o *Order) batchUpdate(ctx context.Context, orders *model.Orders) error {
	req := o.client.Bulk().Index(o.index)
	for _, order := range orders.Data {
		order.UpdateTime = time.Now().Format(golangDate)

		snakedOrder := util.SnakeData(order)

		doc := elastic.NewBulkUpdateRequest().Id(strconv.FormatInt(int64(order.OrderId), 10)).Doc(snakedOrder)
		req.Add(doc)
	}

	if req.NumberOfActions() < 0 {
		return nil
	}

	resp, err := req.Do(ctx)
	if err != nil {
		return err
	}

	if resp.Errors {
		return util.GetEsFailedErr(resp)
	}

	return nil
}

// MGet
func (o *Order) MGet(ctx context.Context, ids []int) (*model.Orders, error) {
	resData := make([]model.Order, 0)
	orderIdStrs := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		orderIdStrs = append(orderIdStrs, strconv.FormatInt(int64(id), 10))
	}

	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewTermsQuery("order_id", orderIdStrs...))
	sorter := elastic.NewFieldSort("create_time").Desc()

	resp, err := o.client.Search(o.index).
		Query(query).
		SortBy(sorter).
		Size(len(ids)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if resp.TotalHits() == 0 {
		return nil, nil
	}

	snakedOrder := util.SnakeData(model.Order{})
	var cameledOrder model.Order

	for _, item := range resp.Each(reflect.TypeOf(snakedOrder)) {
		tmpData := util.CamelData(item)
		_ = json.Unmarshal(tmpData, &cameledOrder)
		resData = append(resData, cameledOrder)
	}

	res := &model.Orders{
		Data: resData,
	}

	return res, nil
}

// BatchDel
func (o *Order) BatchDel(ctx context.Context, ids []int) error {
	var err error 
	for i := 0; i < retryTime; i++ {
		if _, err := o.batchDel(ctx, ids); err != nil {
			continue
		}
		return nil
	}

	return err
}

// batchDel
func (o *Order) batchDel(ctx context.Context, ids []int) (int64, error) {
	orderIdStrs := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		orderIdStrs = append(orderIdStrs, strconv.FormatInt(int64(id), 10))
	}

	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewTermsQuery("order_id", orderIdStrs...))
	resp, err := o.client.DeleteByQuery(o.index).Query(query).Size(len(ids)).Refresh("true").Do(ctx)

	if err != nil {
		return -1, err
	}

	return resp.Deleted, nil
}

// Search
func (o *Order) Search(ctx context.Context, searchConditions *model.SearchConditions) (*model.Orders, error) {
	query := elastic.NewBoolQuery()
	query.Must(searchConditions.MustQuery...)
	query.MustNot(searchConditions.MustNotQuery...)
	query.Should(searchConditions.ShouldQuery...)
	query.Filter(searchConditions.Filters...)
	
	if len(searchConditions.MustQuery) == 0 && len(searchConditions.MustNotQuery) == 0 && len(searchConditions.ShouldQuery) > 0 {
		query.MinimumShouldMatch("1")
	}

	resp, err := o.client.Search().
		Index(o.index).
		Query(query).
		SortBy(searchConditions.Sorters...).
		From(searchConditions.From).
		Size(searchConditions.Size).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil 
	}

	orderRes := make([]model.Order, 0)

	snakedOrder := util.SnakeData(model.Order{})
	var cameledOrder model.Order

	for _, item := range resp.Each(reflect.TypeOf(snakedOrder)) {
		tmpData := util.CamelData(item)
		_ = json.Unmarshal(tmpData, &cameledOrder)
		orderRes = append(orderRes, cameledOrder)
	}
	
	res := &model.Orders{
		Data: orderRes,
	}

	return res, nil
}
