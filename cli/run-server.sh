#!/usr/bin/env bash
set -e

docker run -it --rm -p 3000:3000 golang:1.23 sh -c '
cat <<EOF > main.go
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync/atomic"
)

func main() {
	var total int64 = 0
	http.HandleFunc("/stress-test", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&total, 1)
		statusCode := http.StatusOK
		if rand.Int()%2 != 0 {
			statusCode = http.StatusBadRequest
		}
		message := fmt.Sprintf("call number %d with statusCode %d\\n", total, statusCode)
		fmt.Print(message)
		w.WriteHeader(statusCode)
		w.Write([]byte(message))
	})
	fmt.Println("server running on port 3000")
	http.ListenAndServe(":3000", nil)
}
EOF
go run main.go
'