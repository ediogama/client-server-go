package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	bodyMarshal := body[:len(body)-1]

	//log.Println(string(bodyMarshal[10:]))

	var dolar USDBRL

	err = json.Unmarshal(bodyMarshal[10:], &dolar)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(string(dolar.Bid)))
}
