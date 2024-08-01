package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func failed(err error, msg string) bool {
	if err != nil {
		log.Println("[ERRO - SERVIDOR] ", msg, " > ", err)
	}
	return err != nil
}

const (
	QUOTE_URL    = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	HTTP_TIMEOUT = 200 * time.Millisecond
	DB_TIMEOUT   = 10 * time.Millisecond
)

type Quote struct {
	Code        string `gorm:"primaryKey"` //": "USD",
	Codein      string `gorm:"primaryKey"` //": "BRL",
	Name        string //": "Dólar Americano/Real Brasileiro",
	High        string `gorm:"type:float"`                                           //": "5.6717",
	Low         string `gorm:"type:float"`                                           //": "5.6185",
	VarBid      string `gorm:"type:float"`                                           //": "0.0118",
	PctChange   string `gorm:"type:float"`                                           //": "0.21",
	Bid         string `gorm:"type:float"`                                           //": "5.6557",
	Ask         string `gorm:"type:float"`                                           //": "5.6564",
	Timestamp   string `gorm:"type:int"`                                             //": "1722027598",
	Create_date string `gorm:"primaryKey;column:create_date;type:datetime;not null"` //": "2024-07-26 17:59:58"
}
type APIQuote struct {
	USDBRL Quote
}
type CotacaoHandler struct {
	db *gorm.DB
}

func Server() {

	// init database
	db, err := gorm.Open(sqlite.Open("usd.db"), &gorm.Config{})
	if failed(err, "Falha ao conectar com o banco de dados") {
		os.Exit(1)
	}
	db.AutoMigrate(&Quote{})

	// set up mux for the handler
	mux := http.NewServeMux()
	mux.Handle("/cotacao", CotacaoHandler{db})

	// start listening
	println("Server iniciado na porta 8080")

	http.ListenAndServe(":8080", mux)
}

func (h CotacaoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// get USD quote
	quote := consultaCotacao()

	// failed
	if quote == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	println("Resposta:")
	fmt.Println("    ", quote)
	println()

	// save to db
	err := persisteCotacao(h.db, quote)

	if err != nil {

		var msg string
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			msg = "Cotação não persistida. Registro já existente"
		}
		if strings.Contains(err.Error(), "deadline exceeded") {
			msg = "Estouro de tempo limite para pesistencia de dados"
		}
		failed(err, msg)

	} else {
		log.Println("Nova cotação persistida no banco 👍")
	}
	// answer to the client
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(quote.Bid))
}

func consultaCotacao() *Quote {
	// prepare context
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_TIMEOUT)
	defer cancel()

	// prepare request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, QUOTE_URL, nil)
	if failed(err, "Falha ao preparar request http") {
		return nil
	}

	// run the request
	resp, err := http.DefaultClient.Do(req)
	if failed(err, "Falha ao executar consulta remota") {
		return nil
	}
	defer resp.Body.Close()

	// load quote struct
	var quote APIQuote
	quoteText, err := io.ReadAll(resp.Body)
	if failed(err, "Falha ao ler o body da resposta") {
		return nil
	}
	json.Unmarshal(quoteText, &quote)

	return &quote.USDBRL
}

func persisteCotacao(db *gorm.DB, quote *Quote) error {
	// save response to local database
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	// try to persist
	tx := db.WithContext(ctx).Create(quote)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
