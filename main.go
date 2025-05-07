package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"smddgz/core"
)

//go:embed templates/*
var templates embed.FS

func main() {
	r := gin.Default()
	fs, _ := template.ParseFS(templates, "templates/*")
	r.SetHTMLTemplate(fs)
	r.GET("/", home)
	r.GET("/vmess", vmess)
	r.Run(":80")
}

func home(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func vmess(c *gin.Context) {
	c.String(http.StatusOK, core.GetLink())
}
