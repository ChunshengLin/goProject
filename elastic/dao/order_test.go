package dao

import (
	"context"
	"elastic/model"
	"elastic/util"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDataOperation(t *testing.T) {
	Convey("Data operation should be success", t, func() {
		Convey("Batch insert data success", func(){
			orders := model.Orders{
				Data: []model.Order{{
					OrderId: 3,
					UserName: "赵",
					Price: 99,
					ProductId: 2,
					ProductName: "北京-三年级-数学",},
					},
				}
		
			esClient, _ := util.NewEsClient()
		
			orderEs, _ := NewOrder(esClient)
			ctx := context.Background()
		
			err := orderEs.BatchInsert(ctx, &orders)
			So(err, ShouldBeNil)
		})
	})
}