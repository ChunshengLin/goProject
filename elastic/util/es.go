package util

import (
	"elastic/logger"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
)


var defaultTimeOut = 20*time.Second

// NewEsClient 新建一个Es客户端实例
func NewEsClient() (*elastic.Client, error) {
	tag := "NewEsClient"

	configs := getEsConfig()
	url := fmt.Sprintf("http://%s:%s", configs["host"], configs["port"])
	timeout := cast.ToDuration(configs["timeout"])
	if timeout == 0 {
		timeout = defaultTimeOut
	}
	httpclient := &http.Client{
		Timeout: timeout,
	}
	options := []elastic.ClientOptionFunc {
		elastic.SetURL(url),
		elastic.SetHttpClient(httpclient),
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		logger.E(tag, "New EsClient failed, error:%+v", err)
		return nil, err
	}

	logger.I(tag, "create EsClient success")

	return client, nil
} 

// getEsConfig 获取es配置
func getEsConfig()  map[string]string {
	tag := "getEsConfig"

	_, curFile, _, _ := runtime.Caller(0)
	configPath := path.Join(path.Dir(path.Dir(curFile)), "./config/conf.ini")	
	cfg, err := ini.Load(configPath)
	if err != nil {
		logger.E(tag, "load conf.ini failed, err:%+v", err)
		fmt.Println("load conf.ini failed")
	}

	configs := make(map[string]string)
	configs["host"] = cfg.Section("elastic").Key("host").String()
	configs["port"] = cfg.Section("elastic").Key("port").String()
	configs["timeout"] = cfg.Section("elastic").Key("timeout").String()

	return configs
}

// GetEsFailedErr 构建批量操作失败的原因
func GetEsFailedErr(resp *elastic.BulkResponse) error {
	for _, item := range resp.Failed() {
		if item.Error == nil {
			continue
		}
		err := &elastic.Error{
			Status: item.Status,
			Details: item.Error,
		}
		return err
	}

	return nil
}