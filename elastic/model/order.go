package model

import (
	"elastic/logger"
	"strconv"

	"github.com/olivere/elastic/v7"
)


type EsSearchOrderReq struct {
	OrderId			int			`json:"orderId,omitempty"`
	UserName  		string		`json:"userName,omitempty"`
	LowPrice		float64		`json:"lowPrice,omitempty"`
	HighPrice		float64		`json:"highPrice,omitempty"`
	ProductId		int			`json:"productId,omitempty"`
	ProductName		string		`json:"productName,omitempty"`
	PageNum			int			`json:"pageNum"`
	Size			int 		`json:"size"`
}

// CreateQuerys
func (e *EsSearchOrderReq) CreateQuerys() *SearchConditions{
	tag := "createQuerys"

	var searchConditons SearchConditions

	if len(strconv.FormatInt(int64(e.OrderId), 10)) > 0  && e.OrderId != 0 {
		searchConditons.MustQuery = append(searchConditons.MustQuery, elastic.NewTermsQuery("order_id", e.OrderId))
	}
	if len(e.UserName) > 0 {
		searchConditons.ShouldQuery = append(searchConditons.ShouldQuery, elastic.NewMatchQuery("user_name", e.UserName))
	}
	if len(strconv.FormatFloat(e.HighPrice, 'f', 2, 64)) > 0 && len(strconv.FormatFloat(e.LowPrice, 'f', 2, 64)) > 0 {
		searchConditons.MustQuery = append(searchConditons.MustQuery, elastic.NewRangeQuery("price").Gte(e.LowPrice).Lte(e.HighPrice))
	}
	if len(strconv.FormatInt(int64(e.ProductId), 10)) > 0 && e.ProductId != 0 {
		searchConditons.MustQuery = append(searchConditons.MustQuery, elastic.NewTermsQuery("product_id", e.ProductId))
	}
	if len(e.ProductName) > 0 {
		searchConditons.ShouldQuery = append(searchConditons.ShouldQuery, elastic.NewMatchQuery("product_name", e.ProductName))
	}
	if searchConditons.Sorters == nil {
		searchConditons.Sorters = append(searchConditons.Sorters, elastic.NewFieldSort("create_time").Desc())
	}

	if e.PageNum == 0 {
		e.PageNum = 1
	}
	if e.Size == 0 {
		e.Size = 1
	}

	searchConditons.From = (e.PageNum - 1)*e.Size
	searchConditons.Size = e.Size

	logger.I(tag, "orderId:%+v, name:%+v, productId:%+v, productName:%+v", e.OrderId, e.UserName, e.ProductId, e.ProductName)
	
	return &searchConditons
}

type Orders struct {
	Data	[]Order		`json:"data,omitempty"`
}

type Order struct {
	OrderId 		int		`json:"orderId,omitempty"`
	UserName		string	`json:"userName,omitempty"`
	Price			float64	`json:"price,omitempty"`
	ProductId  		int		`json:"productId,omitempty"`
	ProductName		string	`json:"productName,omitempty"`
	CreateTime		string	`json:"createTime,omitempty"`
	UpdateTime		string	`json:"updateTime,omitempty"`
}

type SearchConditions struct {
	MustQuery		[]elastic.Query
	MustNotQuery	[]elastic.Query
	ShouldQuery		[]elastic.Query
	Filters 		[]elastic.Query
	Sorters			[]elastic.Sorter
	From			int
	Size			int
}