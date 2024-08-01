package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func failed(err error, msg string) bool {
	if err != nil {
		log.Println("[ERRO - CLIENTE] ", msg, " > ", err)
	}
	return err != nil
}

func Client() {
	println("Cliente iniciado")

	// set context
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// call server
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if failed(err, "Falha ao criar request") {
		os.Exit(1)
	}

	resp, err := http.DefaultClient.Do(req)
	if failed(err, "Falha ao executar request") {
		os.Exit(1)
	}
	defer resp.Body.Close()

	cotacaoUSD, err := io.ReadAll(resp.Body)
	if failed(err, "Falha ao interpretar cotação") {
		os.Exit(1)
	}
	if resp.StatusCode != 200 || cotacaoUSD == nil {
		failed(errors.New("bad satus code or content value"), "Falha ao receber cotação")
		os.Exit(1)
	}

	fmt.Println("COTAÇÃO USD: ", string(cotacaoUSD))

	// persist to file
	file, err := os.Create("cotacao.txt")
	if failed(err, "Falha ao criar aquivo local") {
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: " + string(cotacaoUSD))
	if failed(err, "Falha ao escrever cotação USD no arquivo") {
		os.Exit(1)
	}
	println("Execução cliente terminada com sucesso")
}
