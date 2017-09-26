package models

import (
	"github.com/jinzhu/gorm"
	"encoding/json"
	"youmaolu-wx-app-api/libs"
	"log"
)

type Meta struct {
	gorm.Model
	CardId         int64  `gorm:"default:0"`
	CardCategoryId int64  `gorm:"default:0"`
	ActivityId     int64  `gorm:"default:0"`
	Title          string `gorm:"default:null"`
	Overview       string `gorm:"default:null"`
	WebUrl         string `gorm:"default:null"`
	Author         string `gorm:"default:null"`
}

func (Meta) TableName() string {
	return "metas"
}

func (m *Meta) AddMeta(item map[string]interface{}) uint {
	var meta Meta
	id, _ := item["id"].(json.Number).Int64()
	libs.DB.Where("card_id = ?", id).First(&meta)

	if meta.ID != 0 {
		return meta.ID
	}
	cardCategoryId, _ := item["card_category_id"].(json.Number).Int64()
	activityId, _ := item["activity_id"].(json.Number).Int64()
	meta.CardId = id
	meta.CardCategoryId = cardCategoryId
	meta.ActivityId = activityId
	meta.Title = item["title"].(string)
	meta.Overview = item["overview"].(string)
	meta.WebUrl = item["web_url"].(string)
	meta.Author = item["author"].(string)

	err := libs.DB.Save(&meta).Error
	if err != nil {
		log.Print(err)
	}
	return meta.ID
}

func (m *Meta) GetMetaCountByActivityId(activityId int) int {
	var meta Meta
	var count int
	libs.DB.Where("activity_id = ?", activityId).First(&meta).Count(&count)
	return count
}
