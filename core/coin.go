package core

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const coinUrl string = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map"
const quoteUrl string = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"
const headerKey string = "X-CMC_PRO_API_KEY"
const headerValue string = "9f403c6a-4f68-4b7d-9e7f-bbd443545efa"

func SearchCoin(symbol string) []CoinInfo {
	c := &http.Client{}
	request, _ := http.NewRequest("GET", coinUrl+"?symbol="+symbol, nil)
	request.Header.Add(headerKey, headerValue)

	resp, err := c.Do(request)
	if err != nil {
		log.Fatal("fetch map info err: ", err)
	}
	body, _ := io.ReadAll(resp.Body)

	var respInfo RespInfo
	_ = json.Unmarshal(body, &respInfo)

	return respInfo.Data
}

func getQuote(coinList []CoinInfo, convert string) []CoinInfo {
	if len(coinList) == 0 {
		return coinList
	}
	var idArr []string
	for _, info := range coinList {
		idArr = append(idArr, strconv.Itoa(info.Id))
	}
	client := &http.Client{}
	request, _ := http.NewRequest("GET", quoteUrl+"?convert="+convert+"&id="+strings.Join(idArr, ","), nil)
	request.Header.Add(headerKey, headerValue)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("invoke quoteUrl error", err)
	}
	body, _ := io.ReadAll(resp.Body)

	var quoteInfo map[string]interface{}
	err = json.Unmarshal(body, &quoteInfo)
	if err != nil {
		log.Fatal("quoteBody json unmarshal error", err)
	}
	dataMap := quoteInfo["data"].(map[string]interface{})

	for i, coin := range coinList {
		itemMap := dataMap[strconv.Itoa(coin.Id)].(map[string]interface{})
		price := itemMap["quote"].(map[string]interface{})[convert].(map[string]interface{})["price"].(float64)
		coinList[i].Price = price
		coinList[i].Symbol = itemMap["symbol"].(string)
	}
	return coinList
}

func SaveCoin(info CoinInfo) {
	dir, _ := os.Getwd()
	path := dir + string(os.PathSeparator) + "coinInfo.json"
	coinArr := []CoinInfo{info}
	coins := getQuote(coinArr, "USD")
	coinList := SelectCoin()
	coinList = append(coinList, coins...)
	buf, _ := json.Marshal(coinList)
	os.WriteFile(path, buf, 0744)
}

func DeleteCoin(id int) {
	dir, _ := os.Getwd()
	path := dir + string(os.PathSeparator) + "coinInfo.json"
	coinList := SelectCoin()
	for i, coin := range coinList {
		if coin.Id == id {
			infos := append(coinList[:i], coinList[i+1:]...)
			marshal, _ := json.Marshal(infos)
			err := os.WriteFile(path, marshal, 0744)
			if err != nil {
				log.Println("delete map info err: ", err)
			}
		}
	}
}

func SelectCoin() []CoinInfo {
	dir, _ := os.Getwd()
	path := dir + string(os.PathSeparator) + "coinInfo.json"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
	buf, _ := os.ReadFile(path)
	var coinList []CoinInfo
	json.Unmarshal(buf, &coinList)
	return getQuote(coinList, "USD")
}

func GetCoinData() map[string]interface{} {
	dataMap := make(map[string]interface{})
	list := SelectCoin()
	rate := getExchangeRate()

	var total float64
	var zhTotal float64

	for i, coin := range list {
		list[i].Amount = coin.Num * coin.Price
		list[i].ZhAmount = coin.Num * coin.Price * rate
		total += list[i].Amount
		zhTotal += list[i].ZhAmount
		list[i].Amount = math.Round(list[i].Amount)
		list[i].ZhAmount = math.Round(list[i].ZhAmount)
		list[i].Price = RoundToTwoDecimal(list[i].Price)
	}

	dataMap["total"] = total
	dataMap["zhTotal"] = zhTotal
	dataMap["coinList"] = list
	dataMap["rate"] = rate

	return dataMap
}

func getExchangeRate() float64 {
	var coin CoinInfo
	coin.Id = 825
	coinArr := []CoinInfo{coin}
	coinList := getQuote(coinArr, "CNY")
	return coinList[0].Price
}

func RoundToTwoDecimal(val float64) float64 {
	return math.Round(val*100000) / 100000
}

type RespInfo struct {
	Status interface{} `json:"status,omitempty"`
	Data   []CoinInfo  `json:"data,omitempty"`
}

type CoinInfo struct {
	Id        int     `json:"id"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
	Num       float64 `json:"num"`
	Amount    float64 `json:"amount"`
	ZhAmount  float64 `json:"zhAmount"`
	PkSerial  string  `json:"pkSerial"`
	CnyAmount float64 `json:"cnyAmount"`
	UsdAmount float64 `json:"usdAmount"`
}

func main() {
	//coin := searchCoin("BTC")
	//for _, datum := range coin.Data {
	//	fmt.Println(datum.Id, datum.Symbol, datum.Name)
	//}

	//var c1 CoinInfo
	//c1.Symbol = "BTC"
	//c1.Id = 1
	//c1.Num = 1
	//var c2 CoinInfo
	//c2.Symbol = "LTC"
	//c2.Id = 2
	//c2.Num = 1
	//var coinList []CoinInfo
	//coinList = append(coinList, c1)
	//coinList = append(coinList, c2)
	//quote := getQuote(coinList)
	//fmt.Println(quote)

	//dir, _ := os.Getwd()
	//print(dir + string(os.PathSeparator) + "coinInfo.json")
	//
	//coin := selectCoin()
	//log.Println(coin)

	//var c3 CoinInfo
	//c3.Symbol = "BTC"
	//c3.Id = 1
	//c3.Num = 1
	//
	//SaveCoin(c3)
	//
	//coin := selectCoin()
	//log.Println(coin)
	//DeleteCoin(1)

	println(getExchangeRate())

	data := GetCoinData()
	log.Println(data)
}
