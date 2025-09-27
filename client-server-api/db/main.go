package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/deal-machine/go-expert/client-server-api/logger"

	_ "modernc.org/sqlite"
)

var loggr = logger.GetLogger("[DB] ")

type CurrencyModel struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"var_bid"`
	PctChange  string `json:"pct_change"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func Init() *sql.DB {
	database, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		loggr.Fatalln("Error on sqlOpen", err)
	}

	_, err = database.Exec(createTableSql())
	if err != nil {
		loggr.Fatalln("Error on createTableSql", err)
	}

	loggr.Println("Database connected and created table!")
	return database
}

func Insert(ctx context.Context, database *sql.DB, cm CurrencyModel) bool {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err := database.ExecContext(ctx, insertSql(),
		cm.Code,
		cm.Codein,
		cm.Name,
		cm.High,
		cm.Low,
		cm.VarBid,
		cm.PctChange,
		cm.Bid,
		cm.Ask,
		cm.Timestamp,
		cm.CreateDate,
	)
	if err != nil {
		select {
		case <-ctx.Done():
			loggr.Println("Error on insertSql timeout/cancelation", err)
			return false
		default:
			loggr.Println("Error on insertSql", err)
			return false
		}
	}
	loggr.Println("Success on insert")
	return true
}

func createTableSql() string {
	return `
		CREATE TABLE IF NOT EXISTS currencies (
			id TEXT PRIMARY KEY,
			code TEXT,
			codein TEXT,
			name TEXT,
			high TEXT,
			low TEXT,
			var_bid TEXT,
			pct_change TEXT,
			bid TEXT,
			ask TEXT,
			timestamp TEXT,
			create_date TEXT
		);
	`
}

func insertSql() string {
	return `
		INSERT INTO currencies (
			code,
			codein,
			name,
			high,
			low,
			var_bid,
			pct_change,
			bid,
			ask,
			timestamp,
			create_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
}
