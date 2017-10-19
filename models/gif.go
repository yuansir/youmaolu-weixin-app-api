package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"youmaolu-wx-app-api/libs"
	config "github.com/spf13/viper"
	"math"
	"log"
	"fmt"
)

type Gif struct {
	gorm.Model
	Meta     Meta
	MetaId   uint
	ImageUrl string
	QiniuUrl string `gorm:"default:null"`
	Path     string `gorm:"default:null"`
	SpiderAt time.Time
}

func (m *Gif) List(page int) ([]Gif, int, int) {
	var data = []Gif{}
	var totalCount int
	config.SetConfigName("app")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	limit := config.GetInt("pagination.pageSize")
	offset := (page - 1) * limit
	libs.DB.Find(&data).Count(&totalCount)
	err := libs.DB.Order("id desc").Preload("Meta").Offset(offset).Limit(limit).Find(&data).Error
	if err != nil {
		log.Fatalln(err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return data, totalCount, totalPages
}

func (m *Gif) AddGif(item map[string]interface{}) uint {
	meta := Meta{}
	metaId := meta.AddMeta(item)
	var gif Gif

	libs.DB.Where("meta_id = ?", metaId).First(&gif)
	if gif.ID != 0 {
		return gif.ID
	}

	gif.MetaId = metaId
	gif.ImageUrl = item["image_url"].(string)
	gif.QiniuUrl = ""
	gif.Path = ""
	gif.SpiderAt = time.Now()
	err := libs.DB.Save(&gif).Error
	if err != nil {
		log.Fatal(err)
	}
	gif.uploadFile()
	return gif.ID

}

func (m *Gif) uploadFile()  {
	localFile, err := libs.DownloadFormUrl(m.ImageUrl, "./attachments/gif")
	if err != nil {
		fmt.Println(err)
	}
	file, err := libs.UploadToQiniu(localFile)
	if err != nil{
		fmt.Println(err)
	}
	m.Path = file
	libs.DB.Save(m)
}
