package main

import (
	"dy/config"
	"dy/models"
	"dy/router"
	"dy/sql"
	"dy/utils"
	"flag"
	"github.com/gin-gonic/gin"
)

var (
	// specified config file path
	configFile string
)

func main() {
	flag.Parse()
	//ctx := context.Background()
	configFile = "./config/config.yaml"
	appConfig := config.NewConfigFile(configFile)
	cfg := utils.InitConfig(appConfig)
	// set app mode
	gin.SetMode(cfg.App.Mode)

	utils.InitMySQL(cfg.Mysql)
	utils.InitRedis(cfg.Redis)

	sql.InitMysqlSchema()

	models.InitProc()
	r := router.Router()
	r.Run(":8082") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
