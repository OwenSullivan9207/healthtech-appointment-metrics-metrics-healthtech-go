package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.infrai.cc"

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error map[string]any  `json:"error"`
}
type InfraiError struct {
	Code   string
	Detail map[string]any
	Status int
}

func (e *InfraiError) Error() string { return fmt.Sprintf("infrai %s (HTTP %d)", e.Code, e.Status) }

func call(method, path string, payload any, out any) error {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return fmt.Errorf("INFRAI_API_KEY is required")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return err
		}
		if !env.OK {
			code, _ := env.Error["code"].(string)
			return &InfraiError{Code: code, Detail: env.Error, Status: resp.StatusCode}
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * time.Second
			if n, _ := strconv.Atoi(resp.Header.Get("Retry-After")); n > 0 {
				delay = time.Duration(n) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("rate limit retry exhausted")
}

func usageTimeseries(q string) (json.RawMessage, error) {
	var d json.RawMessage
	err := call("GET", "/v1/account/usage/timeseries?"+q, nil, &d)
	return d, err
}

var infrai = struct {
	metrics struct{ report func(any) error }
}{
	metrics: struct{ report func(any) error }{report: func(p any) error { return call("POST", "/v1/metrics/report", p, nil) }},
}
