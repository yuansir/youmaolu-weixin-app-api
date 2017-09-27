package main

import (
	config "github.com/spf13/viper"
	"log"
	"youmaolu-wx-app-api/libs"
	"os"
	"time"
	"io"
	"github.com/parnurzeal/gorequest"
	"github.com/bitly/go-simplejson"
	"math"
	"strconv"
	"youmaolu-wx-app-api/models"
)

const (
	BaseUri           = "https://api.yirimao.com"
	ActiveListApi     = "/activity/get-activity-list"
	NewestActivityApi = "/activity/get-newest-activity"
)

var (
	LogInfo  *log.Logger
	LogError *log.Logger
)

func init() {
	os.Mkdir("logs", 0755)
	logFile, err := os.OpenFile("./logs/spider_"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln("open log file failed", err)
	}

	LogInfo = log.New(io.MultiWriter(logFile), "【Info】:", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(io.MultiWriter(logFile), "【Error】:", log.Ldate|log.Ltime|log.Lshortfile)

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

	newestActivity := fetchNewestActivity()
	id, _ := newestActivity.Get("data").Get("id").Int()

	//第一次执行抓取时全量抓取
	metas := models.Meta{}
	hasMeta := metas.GetMetaCountByActivityId(id - 1)


	if hasMeta < 1 {
		var j float64
		for j = 0; j < math.Ceil(float64(id)/float64(10)); j++ {
			params := "pageSize=10&pageIndex=" + strconv.FormatFloat(j, 'f', 6, 64) + "&version=1.6.1&androidVersion=13"
			activeList := fetchActiveList(params)
			for i := 0; i < id; i++ {
				cards, _ := activeList.Get("data").GetIndex(i).Get("cards").Array()

				if len(cards) > 0 {
					saveToDb(cards)
				}
			}
		}

	}

	cards, _ := newestActivity.Get("data").Get("cards").Array()

	saveToDb(cards)

}

func fetchActiveList(params string) *simplejson.Json {
	request := gorequest.New().Timeout(3 * time.Second)
	response, body, errors := request.Post(BaseUri + ActiveListApi).
	//Set("Content-Type", "application/x-www-form-urlencoded").
		Send(params).
		End()

	if response == nil || response.StatusCode != 200 {
		LogError.Fatalln(response, body, errors)
	} else {
		LogInfo.Println(response, body, errors)
	}

	jsonBody := []byte(body)
	json, err := simplejson.NewJson(jsonBody)

	if err != nil {
		LogError.Println("json error")
	}

	status, err := json.Get("status").Int()
	if status != 2000 || err != nil {
		LogError.Println("status error")
	}
	return json
}

func fetchNewestActivity() *simplejson.Json {
	request := gorequest.New().Timeout(3 * time.Second)
	response, body, errors := request.Get(BaseUri + NewestActivityApi).
		Set("Content-Type", "application/x-www-form-urlencoded").
		End()

	if response == nil || response.StatusCode != 200 {
		LogError.Fatalln(response, body, errors)
	} else {
		LogInfo.Println(response, body, errors)
	}

	jsonBody := []byte(body)
	json, err := simplejson.NewJson(jsonBody)

	if err != nil {
		LogError.Println("json error")
	}

	status, err := json.Get("status").Int()
	if status != 2000 || err != nil {
		LogError.Println("status error")
	}
	return json
}

func saveToDb(cards []interface{}) {
	for key, card := range cards {
		switch key {
		case 0:
			//image
			image := models.Image{}
			image.AddImage(card.(map[string]interface{}))
		case 1:
			//video
			video := models.Video{}
			video.AddVideo(card.(map[string]interface{}))
		case 2:
			//product
			product := models.Product{}
			product.AddProduct(card.(map[string]interface{}))
		case 3:
			gif := models.Gif{}
			gif.AddGif(card.(map[string]interface{}))
		default:
		}
	}
}
