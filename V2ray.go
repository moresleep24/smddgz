package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const Flag string = "WLAN"
const Path string = ""

//const Flag string = "ens3"
//const Path string = "/etc/v2ray/"

type ExportConfigInfo struct {
	Version    string `json:"v"`
	Ps         string `json:"ps"`
	Address    string `json:"add"`
	Port       int    `json:"port"`
	Aid        int    `json:"aid"`
	ConfigType string `json:"type"`
	Net        string `json:"net"`
	Path       string `json:"path"`
	Host       string `json:"host"`
	Id         string `json:"id"`
	Tls        string `json:"tls"`
}

type ImportConfigInfo struct {
	Inbounds []Inbounds `json:"inbounds"`
}

type Inbounds struct {
	Port           int            `json:"port"`
	Protocol       string         `json:"protocol"`
	Settings       Settings       `json:"settings"`
	StreamSettings StreamSettings `json:"streamSettings"`
}

type Settings struct {
	Clients []Clients `json:"clients"`
}

type Clients struct {
	Id      string `json:"id"`
	AlterId int    `json:"alterId"`
}

type StreamSettings struct {
	Network     string      `json:"network"`
	Security    string      `json:"security"`
	KcpSettings KcpSettings `json:"kcpSettings"`
}

type KcpSettings struct {
	Header Header `json:"header"`
}

type Header struct {
	HeaderType string `json:"type"`
}

func GetLink() string {
	var ic ImportConfigInfo
	ip := GetIp()
	file, _ := os.ReadFile(Path + "config.json")
	_ = json.Unmarshal(file, &ic)
	ib := ic.Inbounds[0]

	c := ExportConfigInfo{}
	c.Version = "2"
	c.Ps = ip + ":" + strconv.Itoa(ib.Port)
	c.Address = ip
	c.Port = ib.Port
	c.Aid = ib.Settings.Clients[0].AlterId
	c.ConfigType = ib.StreamSettings.KcpSettings.Header.HeaderType
	c.Net = ib.StreamSettings.Network
	c.Path = ""
	c.Host = ""
	c.Id = ib.Settings.Clients[0].Id
	c.Tls = ib.StreamSettings.Security

	marshal, _ := json.Marshal(c)
	log.Println(string(marshal))
	link := ib.Protocol + "://" + base64.StdEncoding.EncodeToString(marshal)
	followLink := base64.StdEncoding.EncodeToString([]byte(link))
	log.Println(followLink)
	return followLink
}

func GetIp() string {
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		if strings.EqualFold(Flag, i.Name) {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				if strings.Contains(addr.String(), ".") {
					ip, _, _ := net.ParseCIDR(addr.String())
					return ip.To4().String()
				}
			}
		}
	}
	return ""
}
