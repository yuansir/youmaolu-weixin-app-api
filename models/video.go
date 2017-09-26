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

type Video struct {
	gorm.Model
	Meta     Meta
	MetaId   uint
	VideoUrl string
	QiniuUrl string `gorm:"default:null"`
	Path     string `gorm:"default:null"`
	SpiderAt time.Time
}

func (m *Video) List(page int) ([]Video, int, int) {
	var data = []Video{}
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

func (m *Video) AddVideo(item map[string]interface{}) uint {
	meta := Meta{}
	metaId := meta.AddMeta(item)

	var video Video

	libs.DB.Where("meta_id = ?", metaId).First(&video)
	if video.ID != 0 {
		return video.ID
	}

	video.MetaId = metaId
	video.VideoUrl = item["video_url"].(string)
	video.QiniuUrl = ""
	video.Path = ""
	video.SpiderAt = time.Now()
	err := libs.DB.Save(&video).Error
	if err != nil {
		log.Print(err)
	}

	video.uploadFile()

	return video.ID
}

func (m *Video) uploadFile() {
	localFile, err := libs.DownloadFormUrl(m.VideoUrl, "./attachments/video")
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
