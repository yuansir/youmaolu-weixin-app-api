package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"youmaolu-wx-app-api/libs"
	config "github.com/spf13/viper"
	"math"
	"log"
	"strconv"
	"fmt"
)

type Product struct {
	gorm.Model
	Meta       Meta
	MetaId     uint
	ProductId  int64
	ProductUrl string
	ImageUrl   string
	QiniuUrl   string `gorm:"default:null"`
	Path       string `gorm:"default:null"`
	SpiderAt   time.Time
}

func (m *Product) List(page int) ([]Product, int, int) {
	var data = []Product{}
	var totalCount int
	config.SetConfigName("app")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	limit := config.GetInt("pagination.pageSize")
	offset := (page - 1) * limit
	libs.DB.Find(&data).Count(&totalCount)
	err := libs.DB.Preload("Meta").Offset(offset).Limit(limit).Find(&data).Error
	if err != nil {
		log.Fatalln(err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return data, totalCount, totalPages
}

func (m *Product) AddProduct(item map[string]interface{}) uint {
	meta := Meta{}
	metaId := meta.AddMeta(item)

	var product Product

	libs.DB.Where("meta_id = ?", metaId).First(&product)
	if product.ID != 0 {
		return product.ID
	}

	var s string = item["product_id"].(string)
	productId, _ := strconv.ParseInt(s, 10, 64)

	product.MetaId = metaId
	product.ProductId = productId
	product.ProductUrl = item["product_url"].(string)
	product.ImageUrl = item["image_url"].(string)
	product.QiniuUrl = ""
	product.Path = ""
	product.SpiderAt = time.Now()
	err := libs.DB.Save(&product).Error
	if err != nil {
		log.Print(err)
	}

	product.uploadFile()
	return product.ID
}

func (m *Product) uploadFile() {
	localFile, err := libs.DownloadFormUrl(m.ImageUrl, "./attachments/product")
	if err != nil {
		fmt.Println(err)
	}
	file, err := libs.UploadToQiniu(localFile)
	if err != nil {
		fmt.Println(err)
	}
	m.Path = file
	libs.DB.Save(m)
}
