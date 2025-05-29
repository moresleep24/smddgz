package main

import (
	"bufio"
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"smddgz/core"
	"strconv"
	"strings"
	"time"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/init.sql
var initSql []byte

const coinUri string = "/9f1c46f0b78e4f8d9e10f65e180bb58d"
const okxCoinUrl string = "/6f9d25e43a6b43e2bf500f5d4c7f7a63"

func init() {
	core.Init(string(initSql))
}

func main() {
	r := gin.Default()
	tmpl := template.New("example").Delims("[[", "]]")
	fs, _ := tmpl.ParseFS(templates, "templates/*")
	r.Static("/word", "./word")
	r.SetHTMLTemplate(fs)
	r.GET("/", home)
	r.GET("/vmess", vmess)
	r.GET(coinUri, coin)
	r.GET(okxCoinUrl, okxCoin)
	r.GET("/ws", ws)
	r.POST("/saveCoin", saveCoin)
	r.POST("/deleteCoin", deleteCoin)
	r.POST("/uploadWord", uploadWord)
	r.GET("/selectWord", selectWord)
	r.GET("/word", word)
	r.GET("/selectDay", selectDay)
	r.GET("/queryWord", queryWord)
	r.Run(":80")
}

func queryWord(context *gin.Context) {
	day, _ := context.GetQuery("day")
	if !strings.EqualFold(day, "") {
		context.JSON(http.StatusOK, core.QueryWord(day))
	} else {
		context.JSON(http.StatusOK, nil)
	}
}

func selectDay(context *gin.Context) {
	context.JSON(http.StatusOK, core.SelectDay())
}

func word(context *gin.Context) {
	context.HTML(http.StatusOK, "english.html", gin.H{})
}

func selectWord(context *gin.Context) {
	pageNum, _ := context.GetQuery("pageNum")
	pageSize, _ := context.GetQuery("pageSize")
	pageNumInt, _ := strconv.ParseInt(pageNum, 10, 64)
	pageSizeInt, _ := strconv.ParseInt(pageSize, 10, 64)
	wordList := core.SelectWord(pageNumInt, pageSizeInt)
	context.JSON(http.StatusOK, wordList)
}

func uploadWord(context *gin.Context) {
	file, _ := context.FormFile("file")
	open, _ := file.Open()
	scanner := bufio.NewScanner(open)
	var wordList []core.Word
	for scanner.Scan() {
		text := scanner.Text()
		if strings.EqualFold(text, "") {
			break
		}
		now := time.Now()
		format := now.Format("2006-01-02")

		var word core.Word
		word.Word = text
		word.CreatedDate = format
		wordList = append(wordList, word)
	}
	if len(wordList) != 0 {
		core.SaveWord(wordList)
	}
	context.JSON(200, "ok")
}

func okxCoin(c *gin.Context) {
	c.HTML(http.StatusOK, "okx-coin.html", nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域，生产中要注意安全策略
	},
}

func ws(context *gin.Context) {
	pkSerial, _ := context.GetQuery("pkSerial")
	sourceConn, _ := upgrader.Upgrade(context.Writer, context.Request, nil)
	core.OpenOkxWs(pkSerial, sourceConn)
}

func deleteCoin(c *gin.Context) {
	idStr, _ := c.GetPostForm("id")
	id, _ := strconv.ParseInt(idStr, 10, 32)
	core.DeleteCoin(int(id))
	c.Redirect(http.StatusMovedPermanently, coinUri)
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
	c.Redirect(http.StatusMovedPermanently, coinUri)
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"msg":  "Hello World",
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
}

func vmess(c *gin.Context) {
	c.String(http.StatusOK, core.GetLink())
}

func coin(c *gin.Context) {
	key, b := c.GetQuery("key")
	res := core.GetCoinData()
	res["key"] = key
	res["coinUrl"] = coinUri
	if b {
		coins := core.SearchCoin(key)
		res["coins"] = coins
	}
	c.HTML(http.StatusOK, "coin.html", res)
}
