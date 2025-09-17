package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/deal-machine/go-expert/challenge01/db"
	httpClient "github.com/deal-machine/go-expert/challenge01/http"
	"github.com/deal-machine/go-expert/challenge01/logger"

	"github.com/google/uuid"
)

type ctxKey string

const REQ_ID ctxKey = "req_id"
const PREFIX = "[SERVER] "
const API_KEY = "989df200befeffa646face07dae9cdda48090d2882a91cdc8b8daa08cde51545"
const CURRENCY_TYPE = "USDBRL"

var loggr = logger.GetLogger(PREFIX)
var database *sql.DB

type CurrencyResponse struct {
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
	database = db.Init()
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoController)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	loggr.Fatalln(server.ListenAndServe())
}

func cotacaoController(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), REQ_ID, uuid.NewString())
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	currencyResponse, err := getCurrency(ctx)
	if err != nil {
		http.Error(w, "Currency Response Error", http.StatusInternalServerError)
		loggr.Println("Currency Response Error ", err)
		return
	}

	currencyModel := db.CurrencyModel{
		Code:       currencyResponse.Code,
		Codein:     currencyResponse.Codein,
		Name:       currencyResponse.Name,
		High:       currencyResponse.High,
		Low:        currencyResponse.Low,
		VarBid:     currencyResponse.VarBid,
		PctChange:  currencyResponse.PctChange,
		Bid:        currencyResponse.Bid,
		Ask:        currencyResponse.Ask,
		Timestamp:  currencyResponse.Timestamp,
		CreateDate: currencyResponse.CreateDate,
	}
	inserted := db.Insert(ctx, database, currencyModel)
	if !inserted {
		http.Error(w, "Error on insert", http.StatusInternalServerError)
		loggr.Println("Error on insert")
		return
	}
	curencyJson, err := json.Marshal(currencyResponse) // json.NewEncoder(resp).Encode(c)
	if err != nil {
		http.Error(w, "Encoding Error", http.StatusInternalServerError)
		loggr.Println("Encoding Error ", err)
		return
	}

	if reqId, ok := ctx.Value(REQ_ID).(string); ok {
		loggr.Print("Success on " + r.Method + "/" + r.Host + r.RequestURI + " reqId: " + reqId)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(curencyJson)
	w.Write([]byte("\n"))
}

func getCurrency(ctx context.Context) (*CurrencyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	apiRequest := httpClient.APIRequest{
		Method: "GET",
		Url:    "https://economia.awesomeapi.com.br/json/last/USD-BRL?token=" + API_KEY,
	}
	var currency map[string]CurrencyResponse

	_, err := httpClient.MakeRequest(ctx, apiRequest, &currency)
	if err != nil {
		select {
		case <-ctx.Done():
			loggr.Println("Timeout/Cancel error", err)
			return nil, err
		default:
			loggr.Println("Error", err)
			return nil, err
		}
	}

	if c, ok := currency[CURRENCY_TYPE]; ok {
		return &c, nil
	}

	return nil, errors.New("currency" + CURRENCY_TYPE + "not exists")
}
