package controllers

import (
	"github.com/kataras/iris/mvc"
	"youmaolu-wx-app-api/models"
	"strconv"
	"github.com/kataras/iris"
	config "github.com/spf13/viper"
	"log"
)

type ProductController struct {
	mvc.Controller
}

func (c *ProductController) Get() {
	products := models.Product{}
	page, err := strconv.Atoi(c.Ctx.URLParam("page"))

	if err != nil || page < 1 {
		page = 1
	}

	list, total, totalPages := products.List(page)

	config.SetConfigName("qiniu")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	for key, item := range list {

		if item.Path == "" {
			list[key].QiniuUrl = item.ImageUrl
		} else {
			list[key].QiniuUrl = config.GetString("default.url") + "/" + item.Path
		}
	}

	c.Ctx.JSON(iris.Map{
		"status_code": iris.StatusOK,
		"message":     "success",
		"data": iris.Map{
			"list":        list,
			"total":       total,
			"totalPages":  totalPages,
			"currentPage": page,
		},
	})
}
