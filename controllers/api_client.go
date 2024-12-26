//api_client.go

package controllers

import (
    "fmt"
    "github.com/beego/beego/v2/client/httplib"
    "github.com/beego/beego/v2/core/logs"
    "github.com/beego/beego/v2/server/web"
)

// getConfig fetches a config value or returns fallback if not found
func getConfig(key, fallback string) string {
    value, err := web.AppConfig.String(key)
    if err != nil || value == "" {
        return fallback
    }
    return value
}

// makeAPIRequest spawns a goroutine to call The Cat API, returning a channel
func makeAPIRequest(endpoint, method string, payload interface{}) chan []byte {
    resultChan := make(chan []byte)

    go func() {
        defer close(resultChan)

        apiKey := getConfig("cat_api_key", "")
        baseURL := getConfig("api_base_url", "")
        url := fmt.Sprintf("%s/%s", baseURL, endpoint)

        var req *httplib.BeegoHTTPRequest

        switch method {
        case "GET":
            req = httplib.Get(url)
        case "POST":
            req = httplib.Post(url)
            if payload != nil {
                req.JSONBody(payload)
            }
        case "DELETE":
            req = httplib.Delete(url)
        default:
            logs.Error("Invalid HTTP method: %s", method)
            resultChan <- []byte(`{"error": "Invalid HTTP method"}`)
            return
        }

        req.Header("x-api-key", apiKey)
        req.Header("Content-Type", "application/json")

        resp, err := req.Bytes()
        if err != nil {
            logs.Error("Error in API request: %v", err)
            resp = []byte(`{"error": "Failed to fetch data"}`)
        }

        resultChan <- resp
    }()

    return resultChan
}
