package main

import (
	"context"
	"os"
	"time"

	httpClient "github.com/deal-machine/go-expert/client-server-api/http"
	"github.com/deal-machine/go-expert/client-server-api/logger"
)

const PREFIX = "[CLIENT] "

type PriceResponse struct {
	Code       string `json:"-"`
	Codein     string `json:"-"`
	Name       string `json:"-"`
	High       string `json:"-"`
	Low        string `json:"-"`
	VarBid     string `json:"-"`
	PctChange  string `json:"-"`
	Bid        string `json:"bid"`
	Ask        string `json:"-"`
	Timestamp  string `json:"-"`
	CreateDate string `json:"-"`
}

var loggr = logger.GetLogger(PREFIX)

func main() {
	starded := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	apiRequest := httpClient.APIRequest{
		Method: "GET",
		Url:    "http://localhost:8080/cotacao",
	}

	var priceResponse *PriceResponse
	_, err := httpClient.MakeRequest(ctx, apiRequest, &priceResponse)

	if err != nil {
		loggr.Println("Error on MakeRequest", err)
		return
	}

	writeOnFile("cotacao.txt", priceResponse.Bid)
	loggr.Println("Total time", time.Since(starded))
}

func writeOnFile(filePath string, value string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		loggr.Println("Error on OpenFile", err)
	}
	defer file.Close()

	file.WriteString("Dólar: " + value + "\n")
	loggr.Println("Success on write file")
}
