# 微信小程序【有猫撸】服务端API

服务端API基于Go Iris+MySQL，小程序客户端[**youmaolu-weixin-app **](https://github.com/yuansir/youmaolu-weixin-app)


## 配置

* 导入数据库 `database/tables.sql`
* 修改`config/*.toml.example ` => `config/*.toml`并配置

## 抓取脚本

* commands/spider/spider.go 抓取数据并将附件上传到七牛

