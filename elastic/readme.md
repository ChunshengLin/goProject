### 目录结构

├── router                #路由目录
│   └── api.go            #路由注册实现
├── controller            #路由控制器
│   └── order.go        #订单路由控制实现
├── model                #model目录
│   └── order.go        #model实现
├── dao                    #dao目录
│   └── order.go         #dao数据操作实现
├── service                #服务目录
│   └── api.go             #es增删改查服务实现
├── util                     #工具目录
│   └── es.go              #es相关
│   └── json.go           #Json相关
├── bin                     #可执行文件目录
│   └── elastic             #可执行文件
├── cmd                    #命令代码目录
│   └── main.go          #main.go
├── conf                    #配置文件目录
│   └── conf.ini            #配置文件
├── go.mod
├── Makefile
├── readme.md