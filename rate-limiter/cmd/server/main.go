package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/deal-machine/go-expert/rate-limiter/configs"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/api/routers"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/persistence/redis"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/persistence/repository"
	"github.com/deal-machine/go-expert/rate-limiter/internal/strategy"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("cmd/server/.env"); err != nil {
		log.Println("Error on loading environment variables")
		// return
	}

	// conexão redis
	conn, err := redis.NewRedisClient(context.Background(), getDbConn())
	if err != nil {
		log.Fatalln("error on connect redis db", err)
		return
	}
	defer conn.Close()

	repo := repository.NewRedisRepository(conn)
	s := strategy.NewRateLimitStrategy(repo)
	rlv := configs.GetRateLimitVariables()

	r := routers.Routers(s, rlv)

	serverport := os.Getenv("SERVER_PORT")
	if serverport == "" {
		serverport = "8080"
	}
	log.Println("Server is running on port", serverport)
	http.ListenAndServe(":"+serverport, r)
}

func getDbConn() string {
	dbport := os.Getenv("DB_PORT")
	if dbport == "" {
		dbport = "6379"
	}
	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		dbhost = "localhost"
	}
	return dbhost + ":" + dbport // "localhost:6379"
}
