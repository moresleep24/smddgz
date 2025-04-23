package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

func main() {
	r := gin.Default()
	fs, _ := template.ParseFS(templates, "templates/*")
	r.SetHTMLTemplate(fs)
	r.GET("/", func(c *gin.Context) {
		// 返回 JSON 格式的响应
		c.HTML(http.StatusOK, "home.html", gin.H{
			"msg": "Hello",
		})
	})
	r.GET("/vmess", func(c *gin.Context) {
		// 返回 JSON 格式的响应
		c.String(http.StatusOK, GetLink())
	})
	r.Run(":80")
}
