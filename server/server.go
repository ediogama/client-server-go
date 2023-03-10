package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
	ID         int    `json:"id"`
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
	ctxInsertDb, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ctxRequest, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite3", "file:table.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	req, err := http.NewRequestWithContext(ctxRequest, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	bodyMarshal := body[:len(body)-1]

	var dolar USDBRL

	err = json.Unmarshal(bodyMarshal[10:], &dolar)
	if err != nil {
		panic(err)
	}

	log.Println(dolar)

	const create string = `
	CREATE TABLE IF NOT EXISTS dolar (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		code VARCHAR(80),
		codein VARCHAR(80),
		name VARCHAR(80),
		high VARCHAR(80),
		low VARCHAR(80),
		varbid VARCHAR(80),
		pctchange VARCHAR(80),
		bid VARCHAR(80),
		ask VARCHAR(80),
		timestamp VARCHAR(80),
		createdate VARCHAR(80)
	);
	`
	_, err = db.Exec(create)
	if err != nil {
		panic(err)
	}

	err = insertCotacao(ctxInsertDb, db, dolar)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(string(dolar.Bid)))
}

func insertCotacao(ctx context.Context, db *sql.DB, cotacao USDBRL) error {
	stmt, err := db.Prepare("insert into dolar(code, codein, name, high, low, varbid, pctchange, bid, ask, timestamp, createdate) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, cotacao.Code, cotacao.Codein, cotacao.Name, cotacao.High, cotacao.Low, cotacao.VarBid, cotacao.PctChange, cotacao.Bid, cotacao.Ask, cotacao.Timestamp, cotacao.CreateDate)
	if err != nil {
		return err
	}
	return nil
}
