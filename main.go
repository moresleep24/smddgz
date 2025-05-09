package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"smddgz/core"
	"strconv"
)

//go:embed templates/*
var templates embed.FS

func main() {
	r := gin.Default()
	fs, _ := template.ParseFS(templates, "templates/*")
	r.SetHTMLTemplate(fs)
	r.GET("/", home)
	r.GET("/vmess", vmess)
	r.GET("/9f1c46f0b78e4f8d9e10f65e180bb58d", coin)
	r.POST("/saveCoin", saveCoin)
	r.POST("/deleteCoin", deleteCoin)
	r.Run(":80")
}

func deleteCoin(c *gin.Context) {
	idStr, _ := c.GetPostForm("id")
	id, _ := strconv.ParseInt(idStr, 10, 32)
	core.DeleteCoin(int(id))
	c.Redirect(http.StatusMovedPermanently, "/coin")
}

func saveCoin(c *gin.Context) {
	var coinInfo core.CoinInfo
	idStr, _ := c.GetPostForm("id")
	id, _ := strconv.ParseInt(idStr, 10, 32)
	numStr, _ := c.GetPostForm("num")
	num, _ := strconv.ParseFloat(numStr, 64)
	coinInfo.Id = int(id)
	coinInfo.Num = num

	log.Println("coinInfo:", coinInfo)
	core.SaveCoin(coinInfo)
	c.Redirect(http.StatusMovedPermanently, "/coin")
}

func home(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func vmess(c *gin.Context) {
	c.String(http.StatusOK, core.GetLink())
}

func coin(c *gin.Context) {
	key, b := c.GetQuery("key")
	res := core.GetCoinData()
	res["key"] = key
	if b {
		coins := core.SearchCoin(key)
		res["coins"] = coins
	}
	c.HTML(http.StatusOK, "coin.html", res)
}
