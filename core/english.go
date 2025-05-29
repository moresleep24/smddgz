package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var audioPath string

func init() {
	dir, _ := os.Getwd()
	path := dir + string(os.PathSeparator) + "word"
	_, err := os.Stat(path)
	if err != nil {
		os.Mkdir(path, os.ModePerm)
	}
	audioPath = path
}

func SaveWordAudio(word string) {
	resp, _ := http.Get(fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=en&client=tw-ob", word))
	buf, _ := io.ReadAll(resp.Body)
	os.WriteFile(audioPath+string(os.PathSeparator)+word+".mp3", buf, os.ModePerm)
}

func SelectWord(pageNum int64, pageSize int64) []Word {
	var wordList []Word
	result, _ := conn.Query("select * from t_word order by created_date desc limit ? offset ?", pageSize, (pageNum-1)*pageSize)
	for result.Next() {
		var word Word
		_ = result.Scan(&word.Word, &word.CreatedDate)
		word.UsUrl = fmt.Sprintf("http://%s/word/%s.mp3", Ip, word.Word)
		wordList = append(wordList, word)
	}
	return wordList
}

func SaveWord(wordList []Word) {
	for _, word := range wordList {
		SaveWordAudio(word.Word)
		conn.Exec("insert into t_word (word,created_date) values (?,?)", word.Word, word.CreatedDate)
	}
}

func QueryWord(day string) []Word {
	var wordList []Word
	result, _ := conn.Query("select * from t_word where created_date=?", day)
	for result.Next() {
		var word Word
		_ = result.Scan(&word.Word, &word.CreatedDate)
		word.UsUrl = fmt.Sprintf("http://%s/word/%s.mp3", Ip, word.Word)
		wordList = append(wordList, word)
	}
	return wordList
}

func SelectDay() []string {
	res, _ := conn.Query("select created_date from t_word group by created_date order by created_date desc")
	var list []string
	for res.Next() {
		var day string
		_ = res.Scan(&day)
		list = append(list, day)
	}
	return list
}

func SelectByDay(day string) {

}

func readWord() {

}

type Word struct {
	Word        string `json:"word,omitempty"`
	CreatedDate string `json:"createdDate,omitempty"`
	UsUrl       string `json:"usUrl,omitempty"`
	UkUrl       string `json:"ukUrl,omitempty"`
}
