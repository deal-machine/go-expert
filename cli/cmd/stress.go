package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

type RunEFunc func(cmd *cobra.Command, args []string)

func RunPing(client *http.Client) RunEFunc {
	return func(cmd *cobra.Command, args []string) {
		var total, totalSuccess, totalError int64 = 0, 0, 0

		start := time.Now()
		concurrencyChannel := make(chan struct{}, concurrency)
		var wg sync.WaitGroup

		for r := 0; r < requests; r++ {
			wg.Add(1)
			concurrencyChannel <- struct{}{}

			go func() {
				defer wg.Done()
				atomic.AddInt64(&total, 1)

				res, err := client.Get(url)
				if res.StatusCode != 200 || err != nil {
					atomic.AddInt64(&totalError, 1)
					<-concurrencyChannel
					return
				}
				defer res.Body.Close()

				atomic.AddInt64(&totalSuccess, 1)
				<-concurrencyChannel
			}()
		}
		wg.Wait()

		totalTime := time.Since(start)
		message := fmt.Sprintf("\n%s\nTempo total: %ds \nTotal de Requisições: %d \nConcorrência: %d \nURL: %s \nTotal de Requisições com sucesso: %d \nTota com erros: %d\n", time.Now().Format(time.DateTime), time.Duration(totalTime.Seconds()), total, concurrency, url, totalSuccess, totalError)
		writeLog(message)
		fmt.Println(message)
	}
}

const LOGS_FILE_PATH = "logs.txt"

func writeLog(message string) {
	fileToWrite, err := os.OpenFile(LOGS_FILE_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer fileToWrite.Close()
	fileToWrite.WriteString(message)
}

var short_description = "Make requests for http url with concurrency"
var long_description = "A longer description that spans multiple lines and likely contains examples and usage of using your command."

var url string
var requests int
var concurrency int

func NewPingCommand(client *http.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "stress",
		Short: short_description,
		Long:  long_description,
		Run:   RunPing(client),
		PreRun: func(cmd *cobra.Command, args []string) {
			// channel for results of requests
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// generate report
		},
	}
}

func init() {
	stressCmd := NewPingCommand(&http.Client{})
	rootCmd.AddCommand(stressCmd)

	stressCmd.Flags().StringVarP(&url, "url", "u", "", "Url to stress")
	stressCmd.Flags().IntVarP(&requests, "requests", "r", 1, "Number of reques to stress, default 1")
	stressCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrency requests to stress, default 1")
	// stressCmd.Flags().StringP("url", "u", "", "Url to stress")
	// stressCmd.Flags().IntP("requests", "r", 1, "Number of reques to stress, default 1")
	// stressCmd.Flags().IntP("concurrency", "c", 1, "Number of concurrency requests to stress, default 1")
	stressCmd.MarkFlagsRequiredTogether("url", "requests", "concurrency")
}
