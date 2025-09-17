package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type APIRequest struct {
	Method string
	Url    string
}

func MakeRequest[T any](ctx context.Context, req APIRequest, res *T) (*T, error) {
	client := http.Client{}

	reqCtx, _ := http.NewRequestWithContext(ctx, req.Method, req.Url, nil)
	resp, err := client.Do(reqCtx)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, err
		default:
			return nil, err
		}

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response T
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
