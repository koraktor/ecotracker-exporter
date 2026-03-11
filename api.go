package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const PowerDataPath = "v1/json"

var apiLog = log.Sugar().Named("api")
var client = http.Client{
	Timeout: 30 * time.Second,
}
var httpLog = log.Sugar().Named("http")

func getPowerData(host string, port int) (*powerData, error) {
	apiLog.Info("Querying power information …")

	var result *powerData
	var url = fmt.Sprintf("http://%s:%d/%s", host, port, PowerDataPath)

	return get(url, result)
}

func get(url string, result *powerData) (*powerData, error) {
	httpLog.Debugf("-> %s", url)

	bodyReader := bytes.NewReader([]byte{})

	req, err := http.NewRequest(http.MethodGet, url, bodyReader)
	if err != nil {
		apiLog.Fatalf("Error creating HTTP request: %s\n", err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		httpLog.Errorf("Error sending HTTP request: %s", err)
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		httpLog.Errorf("API endpoint not found: %s", req.URL)
		return nil, fmt.Errorf("API endpoint not found: %s", req.URL)
	}

	httpLog.Debugf("<- HTTP %s (%d bytes)", res.Status, res.ContentLength)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		httpLog.Errorf("Error reading HTTP response body: %s", err)
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		apiLog.Errorf("Error unmarshalling JSON response body: %s", err)
		return nil, err
	}

	return result, nil
}
