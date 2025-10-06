package configs

import (
	"os"
	"strconv"
	"time"
)

type RateLimitVariables struct {
	ApiKeyLimit      int
	ApiKeyExpiration time.Duration
	IpLimit          int
	IpExpiration     time.Duration
}

func GetRateLimitVariables() RateLimitVariables {
	apiKeyLimit, err := strconv.Atoi(os.Getenv("API_KEY_LIMIT"))
	if err != nil {
		apiKeyLimit = 100
	}
	ipLimit, err := strconv.Atoi(os.Getenv("IP_LIMIT"))
	if err != nil {
		ipLimit = 100
	}
	apiKeyExpiration, err := time.ParseDuration(os.Getenv("API_KEY_EXPIRATION"))
	if err != nil {
		apiKeyExpiration = time.Second
	}
	ipExpiration, err := time.ParseDuration(os.Getenv("IP_EXPIRATION"))
	if err != nil {
		ipExpiration = time.Second
	}

	return RateLimitVariables{
		ApiKeyLimit:      apiKeyLimit,
		ApiKeyExpiration: apiKeyExpiration,
		IpLimit:          ipLimit,
		IpExpiration:     ipExpiration,
	}
}
