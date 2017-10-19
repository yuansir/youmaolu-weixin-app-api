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

type Image struct {
	gorm.Model
	Meta     Meta
	MetaId   uint
	ImageUrl string
	QiniuUrl string `gorm:"default:null"`
	Path     string `gorm:"default:null"`
	SpiderAt time.Time
}

func (m *Image) List(page int) ([]Image, int, int) {
	var data = []Image{}
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

func (m *Image) AddImage(item map[string]interface{}) uint {
	meta := Meta{}
	metaId := meta.AddMeta(item)
	var image Image

	libs.DB.Where("meta_id = ?", metaId).First(&image)
	if image.ID != 0 {
		return image.ID
	}

	image.MetaId = metaId
	image.ImageUrl = item["image_url"].(string)
	image.QiniuUrl = ""
	image.Path = ""
	image.SpiderAt = time.Now()
	err := libs.DB.Save(&image).Error
	if err != nil {
		log.Print(err)
	}

	image.uploadFile()
	return image.ID
}


func (m *Image) uploadFile()  {
		localFile, err := libs.DownloadFormUrl(m.ImageUrl, "./attachments/image")
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
