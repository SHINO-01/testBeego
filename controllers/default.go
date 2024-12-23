package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type CatController struct {
	web.Controller
}

type CatImage struct {
	ID     string  `json:"id"`
	URL    string  `json:"url"`
	Breeds []Breed `json:"breeds"`
}

type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
	Temperament string `json:"temperament"`
	WikipediaURL string `json:"wikipedia_url"`
}

type Vote struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
	Value   int    `json:"value"` // 1 for up, -1 for down
}

type FavoritePayload struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
}

// Utility function to fetch configuration values with default fallback
func getConfig(key, fallback string) string {
	value, err := web.AppConfig.String(key)
	if err != nil || value == "" {
		return fallback
	}
	return value
}

// Utility function to make an API request using goroutines and channels
func makeAPIRequest(endpoint, method string, payload interface{}) chan []byte {
	resultChan := make(chan []byte)

	go func() {
		defer close(resultChan)

		apiKey := getConfig("cat_api_key", "")
		baseURL := getConfig("api_base_url", "")

		var req *httplib.BeegoHTTPRequest
		url := fmt.Sprintf("%s/%s", baseURL, endpoint)

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

// Handler for rendering the index page
func (c *CatController) Get() {
	c.TplName = "index.tpl"
}

// Fetch a random cat image
func (c *CatController) GetRandomCat() {
	resultChan := makeAPIRequest("images/search?limit=1", "GET", nil)

	select {
	case result := <-resultChan:
		c.Ctx.Output.Header("Content-Type", "application/json")
		c.Ctx.Output.Body(result)
	case <-time.After(10 * time.Second):
		c.Data["json"] = map[string]string{"error": "Request timeout"}
		c.ServeJSON()
	}
}

// Fetch all cat breeds
func (c *CatController) GetBreeds() {
	resultChan := makeAPIRequest("breeds", "GET", nil)

	select {
	case result := <-resultChan:
		c.Ctx.Output.Header("Content-Type", "application/json")
		c.Ctx.Output.Body(result)
	case <-time.After(10 * time.Second):
		c.Data["json"] = map[string]string{"error": "Request timeout"}
		c.ServeJSON()
	}
}
// Get Breeds Images//
func (c *CatController) GetBreedImages() {
	logs.Info("GetBreedImages Endpoint Hit")

	breedID := c.GetString("breed_id")
	limit, err := c.GetInt("limit", 8)
	if err != nil{
		logs.Error("Invalid Limit Parameter: %v", err)
		c.Data["json"] = map[string]string{"error": "Invalid limit parameter"}
        c.ServeJSON()
        return
	}

	apiEndpoint := fmt.Sprintf("images/search?breed_ids=%s&limit=%d", breedID, limit)

    // Fire off the request in a goroutine
    resultChan := makeAPIRequest(apiEndpoint, "GET", nil)

    select {
    case result := <-resultChan:
        // Return the JSON response directly
        c.Ctx.Output.Header("Content-Type", "application/json")
        c.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timeout"}
        c.ServeJSON()
    }
}

// Submit a vote for a cat image
func (c *CatController) Vote() {
    logs.Info("Vote endpoint hit")

	bodyBytes := c.Ctx.Input.CopyBody(1024 * 1024) // up to 1MB (adjust as needed)
    if len(bodyBytes) == 0 {
        logs.Warn("Request body is empty after CopyBody")
        c.Data["json"] = map[string]string{"error": "Empty request body"}
        c.ServeJSON()
        return
    }
    logs.Info("Raw request body: %s", string(bodyBytes))

    var vote Vote

    // Log the raw request body for debugging
    logs.Info("Raw request body: %s", string(c.Ctx.Input.RequestBody))

    // Attempt to parse the request body
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vote); err != nil {
        logs.Error("Failed to parse request body: %v", err)
        c.Data["json"] = map[string]string{"error": "Invalid request body"}
        c.ServeJSON()
        return
    }

    logs.Info("Parsed Vote Struct: %+v", vote)

    // Validate required fields
    if vote.ImageID == "" || (vote.Value != 1 && vote.Value != -1) {
        logs.Warn("Invalid request parameters: image_id=%s, value=%d", vote.ImageID, vote.Value)
        c.Data["json"] = map[string]string{"error": "Invalid request parameters"}
        c.ServeJSON()
        return
    }

    // Assign sub_id if missing
    if vote.SubID == "" {
        vote.SubID = getConfig("sub_id", "default-user-id")
    }

    logs.Info("Final Vote Payload: %+v", vote)

    // Make the API request
    resultChan := makeAPIRequest("votes", "POST", vote)

    select {
    case result := <-resultChan:
        c.Ctx.Output.Header("Content-Type", "application/json")
        c.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timeout"}
        c.ServeJSON()
    }
}

// Fetch favorite cat images
func (c *CatController) GetFavorites() {
	subID := getConfig("sub_id", "default-user-id")
    
    // e.g. favourites?limit=28&order=DESC&sub_id=user-123
    apiEndpoint := fmt.Sprintf("favourites?limit=28&order=Desc&sub_id=%s", subID)

    resultChan := makeAPIRequest(apiEndpoint, "GET", nil)

	select {
	case result := <-resultChan:
		c.Ctx.Output.Header("Content-Type", "application/json")
		c.Ctx.Output.Body(result)
	case <-time.After(10 * time.Second):
		c.Data["json"] = map[string]string{"error": "Request timeout"}
		c.ServeJSON()
	}
}

// Add a cat image to favorites
func (c *CatController) AddFavorite() {
	var favorite FavoritePayload

	bodyBytes := c.Ctx.Input.CopyBody(1024 * 1024) // up to 1MB (adjust as needed)
    if len(bodyBytes) == 0 {
        logs.Warn("Request body is empty after CopyBody")
        c.Data["json"] = map[string]string{"error": "Empty request body"}
        c.ServeJSON()
        return
    }
    logs.Info("Raw request body: %s", string(bodyBytes))
	// Parse request body
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &favorite); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	// Validate input
	if favorite.ImageID == "" {
		c.Data["json"] = map[string]string{"error": "Image ID is required"}
		c.ServeJSON()
		return
	}

	// Assign sub_id from configuration
	favorite.SubID = getConfig("sub_id", "default-user-id")

	// Make the API request
	resultChan := makeAPIRequest("favourites", "POST", favorite)

	select {
	case result := <-resultChan:
		c.Ctx.Output.Header("Content-Type", "application/json")
		c.Ctx.Output.Body(result)
	case <-time.After(10 * time.Second):
		c.Data["json"] = map[string]string{"error": "Request timeout"}
		c.ServeJSON()
	}
}

// Remove a cat image from favorites
func (c *CatController) RemoveFavorite() {
	favoriteID := c.Ctx.Input.Param(":id")

	// Make the API request
	resultChan := makeAPIRequest(fmt.Sprintf("favourites/%s", favoriteID), "DELETE", nil)

	select {
	case result := <-resultChan:
		c.Ctx.Output.Header("Content-Type", "application/json")
		c.Ctx.Output.Body(result)
	case <-time.After(10 * time.Second):
		c.Data["json"] = map[string]string{"error": "Request timeout"}
		c.ServeJSON()
	}
}

func (c *CatController) GetVoteHistory() {
    logs.Info("GetVoteHistory endpoint hit")

    // Use sub_id from config OR your own logic
    subID := getConfig("sub_id", "default-user-id")

    // Build the endpoint for The Cat API:
    // e.g. votes?sub_id=user-123
    apiEndpoint := fmt.Sprintf("votes?sub_id=%s&limit=28&order=Desc", subID)

    // Make the goroutine-based request
    resultChan := makeAPIRequest(apiEndpoint, "GET", nil)

    select {
    case result := <-resultChan:
        c.Ctx.Output.Header("Content-Type", "application/json")
        c.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timeout"}
        c.ServeJSON()
    }
}

