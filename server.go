package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"youmaolu-wx-app-api/controllers"
	config "github.com/spf13/viper"
	"log"
	"youmaolu-wx-app-api/libs"
)

func init() {
	config.AddConfigPath("./config")
	config.SetConfigName("mysql")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	dbConfig := libs.DbConfig{
		config.GetString("default.host"),
		config.GetString("default.port"),
		config.GetString("default.database"),
		config.GetString("default.user"),
		config.GetString("default.password"),
	}

	libs.DB = dbConfig.InitDB()

	if config.GetBool("default.sql_log") {
		libs.DB.LogMode(true)
	}
}

func main() {
	app := iris.New()

	app.Get("/", func(ctx context.Context) {
		ctx.JSON(iris.Map{
			"status_code": iris.StatusOK,
			"message":     "数据不存在",
			"data":        iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusNotFound, func(ctx context.Context) {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"status_code": iris.StatusNotFound,
			"message":     "数据不存在",
			"data":        iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx context.Context) {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"status_code": iris.StatusInternalServerError,
			"message":     "服务器端错误",
			"data":        iris.Map{},
		})
	})

	app.Controller("/images", new(controllers.ImageController))
	app.Controller("/gifs", new(controllers.GifController))
	app.Controller("/products", new(controllers.ProductController))
	app.Controller("/videos", new(controllers.VideoController))

	app.Run(iris.Addr(":1121"), iris.WithConfiguration(iris.Configuration{
		TimeFormat: "2006-05-02 15:04:05",
	}))
}
