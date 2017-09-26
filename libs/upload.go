package libs

import (
	"net/url"
	"fmt"
	"os"
	"net/http"
	"io"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"context"
	"strings"
	"log"
	config "github.com/spf13/viper"
)

func DownloadFormUrl(urlString string, dir string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
	}

	filename := dir + u.Path
	output, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer output.Close()

	response, err := http.Get(urlString)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return "", err
	}

	log.Println(n)

	return filename, nil

}

func UploadToQiniu(localFile string) (string, error) {
	config.SetConfigName("qiniu")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	accessKey := config.GetString("default.accessKey")
	secretKey := config.GetString("default.secretKey")
	bucket := config.GetString("default.bucket")

	tokens := strings.Split(localFile, "attachments/")
	key := tokens[len(tokens)-1]

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{

	}
	//putExtra.NoCrc32Check = true
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return ret.Key, nil

}
