package core

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

var dbPath string
var conn *sql.DB
var okxCoinMap = make(map[string]map[string]string)

//6f9d25e43a6b43e2bf500f5d4c7f7a63

func Init(initSql string) {
	dir, _ := os.Getwd()
	dbPath = dir + string(os.PathSeparator) + "coin.db"
	log.Println(dbPath)
	var err error
	conn, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(initSql)
	err = conn.Ping()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	} else {
		log.Println("db connect success")
	}
}

func SaveOne(coin CoinInfo) {
	conn.Exec("insert into t_coin(pk_serial,symbol,num) values(?,?,?)", coin.PkSerial, coin.Symbol, coin.Num)
}

func SelectAll(pkSerial string) []CoinInfo {
	var coins []CoinInfo
	res, _ := conn.Query("select * from t_coin where pk_serial=?", pkSerial)
	for res.Next() {
		var coin CoinInfo
		res.Scan(&coin.PkSerial, &coin.Symbol, &coin.Num)
		coins = append(coins, coin)
	}
	return coins
}

func deleteOne(coin CoinInfo) {
	conn.Exec("delete from t_coin where pk_serial=? and symbol=?", coin.PkSerial, coin.Symbol)
}

const okxUrl string = "wss://wspap.okx.com:8443/ws/v5/public"

func OpenOkxWs(serial string, sourceConn *websocket.Conn) []CoinInfo {
	exchangeRate := 7.2
	resp, err := http.Get("https://api.frankfurter.app/latest?from=USD&to=CNY")
	if err == nil {
		buffer, _ := io.ReadAll(resp.Body)
		var respMap map[string]interface{}
		json.Unmarshal(buffer, &respMap)
		exchangeRate = respMap["rates"].(map[string]interface{})["CNY"].(float64)
	}

	coinList := SelectAll(serial)
	if len(coinList) == 0 {
		log.Fatal("failed to find any coin")
		sourceConn.Close()
		return nil
	}
	//open okx ws
	targetConn, _, err := websocket.DefaultDialer.Dial(okxUrl, nil)
	if err != nil {
		log.Fatal("failed to dial okx:", err)
		sourceConn.Close()
		return nil
	}
	//send subscribe
	argMap := make(map[string]string)
	argMap["channel"] = "index-tickers"
	var argArr []interface{}
	argArr = append(argArr, argMap)
	coinMap := make(map[string]any)
	coinMap["op"] = "subscribe"
	coinMap["args"] = argArr
	for _, coin := range coinList {
		argMap["instId"] = coin.Symbol + "-USDT"
		buf, _ := json.Marshal(coinMap)
		targetConn.WriteMessage(1, buf)
	}
	//receive msg
	go func() {
		for {
			_, p, err2 := targetConn.ReadMessage()
			if err2 != nil {
				log.Println("okx websocket read err:", err2)
				return
			}
			var resMap map[string]interface{}
			json.Unmarshal(p, &resMap)
			dataArr, ok := resMap["data"].([]interface{})
			if !ok {
				continue
			}
			data := dataArr[0].(map[string]interface{})
			key := data["instId"].(string)
			value := data["idxPx"].(string)
			m := okxCoinMap[serial]
			if m == nil {
				m = make(map[string]string)
			}
			m[key] = value
			okxCoinMap[serial] = m
		}
	}()
	//close target
	go func() {
		for {
			_, _, err2 := sourceConn.ReadMessage()
			if err2 != nil {
				targetConn.Close()
				log.Println("client websocket read err:", err2)
				return
			}
		}
	}()
	//send msg to client
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			m2 := okxCoinMap[serial]
			var cnyTotalAmount float64 = 0
			var usdTotalAmount float64 = 0
			for i, coin := range coinList {
				price := m2[coin.Symbol+"-USDT"]
				priceF, _ := strconv.ParseFloat(price, 64)
				usdAmount := coin.Num * priceF
				cnyAmount := usdAmount * exchangeRate

				usdTotalAmount += usdAmount
				cnyTotalAmount += cnyAmount
				coinList[i].Price = priceF

				coinList[i].CnyAmount = math.Trunc(cnyAmount)
				coinList[i].UsdAmount = math.Trunc(usdAmount)
			}
			respMap := map[string]interface{}{}
			respMap["coinList"] = coinList
			respMap["usdTotalAmount"] = math.Trunc(usdTotalAmount)
			respMap["cnyTotalAmount"] = math.Trunc(cnyTotalAmount)
			respMap["exchangeRate"] = exchangeRate
			buffer, _ := json.Marshal(respMap)
			err2 := sourceConn.WriteMessage(1, buffer)
			if err2 != nil {
				targetConn.Close()
			}
		}
	}()
	return nil
}
