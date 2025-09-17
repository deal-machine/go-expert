package logger

import "log"

func GetLogger(prefix string) *log.Logger {
	log.SetPrefix(prefix)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return log.Default()
}
