package util

import (
	"elastic/model"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSnakeAndCamel(t *testing.T) {
	Convey("Snake and camel data should be success", t, func() {
		Convey("Snake data success", func(){
			data := model.Order {OrderId: 1}
			So(snakeEqual(SnakeData(data)), ShouldBeTrue)
		})
		Convey("Camel data success", func(){
			data := make(map[string]interface{})
			data["order_id"] = 1
			So(camelEqual(data), ShouldBeTrue)
		})
	})
}

func snakeEqual(actual map[string]interface{}) bool {
	actualJsonBytes, _ := json.Marshal(actual)

	tmpData := make(map[string]interface{})
	tmpData["order_id"] = 1
	expectedJsonBytes, _ := json.Marshal(tmpData)

	actualData := string(actualJsonBytes)
	expectedData := string(expectedJsonBytes)

	return actualData == expectedData
}

func camelEqual(actual map[string]interface{}) bool {
	var cameledData struct {OrderId int}
	_ = json.Unmarshal(CamelData(actual), &cameledData)
	actualJsonBytes, _ := json.Marshal(cameledData)
	actualData := string(actualJsonBytes)

	cameledData.OrderId = 1
	expectedJsonBytes, _ := json.Marshal(cameledData)
	expectedData := string(expectedJsonBytes)
	
	return actualData == expectedData
}
