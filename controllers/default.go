// controllers/cat_controller.go
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
    ID     string `json:"id"`
    URL    string `json:"url"`
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
    Value   int    `json:"value"` // 1 for up, -1 for down
}

type Favorite struct {
    ID      int    `json:"id"`
    ImageID string `json:"image_id"`
    URL     string `json:"url"`
}

func (c *CatController) Get() {
    c.TplName = "index.tpl"
}

func makeAPIRequest(endpoint string) chan []byte {
    resultChan := make(chan []byte)
    
    go func() {
        apiKey, _ := web.AppConfig.String("cat_api_key")
        baseURL, _ := web.AppConfig.String("api_base_url")
        
        req := httplib.Get(fmt.Sprintf("%s/%s", baseURL, endpoint))
        req.Header("x-api-key", apiKey)
        req.Header("Content-Type", "application/json")
        resp, err := req.Bytes()
        if err != nil {
            resp = []byte(`{"error": "Failed to fetch data"}`)
        }
        
        resultChan <- resp
    }()
    
    return resultChan
}

func (c *CatController) GetRandomCat() {
    resultChan := makeAPIRequest("images/search?limit=1")
    
    select {
    case result := <-resultChan:
        c.Ctx.Output.Header("Content-Type", "application/json")
        c.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timeout"}
        c.ServeJSON()
    }
}

func (c *CatController) GetBreeds() {
    resultChan := makeAPIRequest("breeds")
    
    select {
    case result := <-resultChan:
        c.Ctx.Output.Header("Content-Type", "application/json")
        c.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timeout"}
        c.ServeJSON()
    }
}

func (c *CatController) Vote() {
    var vote struct {
        ImageID string `json:"image_id"`
        SubID   string `json:"sub_id,omitempty"`
        Value   int    `json:"value"`
    }

    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vote); err != nil {
        c.Data["json"] = map[string]string{"error": "Invalid request body"}
        c.ServeJSON()
        return
    }

    // Validate fields
    if vote.ImageID == "" || (vote.Value != 1 && vote.Value != -1) {
        c.Data["json"] = map[string]string{"error": "Invalid request parameters"}
        c.ServeJSON()
        return
    }

    // Set default sub_id if not provided
    if vote.SubID == "" {
        vote.SubID = "user-shino33"
    }

    apiKey, _ := web.AppConfig.String("cat_api_key")
    baseURL, _ := web.AppConfig.String("api_base_url")

    // Make the API request
    req := httplib.Post(fmt.Sprintf("%s/votes", baseURL))
    req.Header("x-api-key", apiKey)
    req.Header("Content-Type", "application/json")
    req.JSONBody(vote)

    var result map[string]interface{}
    err := req.ToJSON(&result)

    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to submit vote"}
    } else {
        c.Data["json"] = result
    }
    c.ServeJSON()
}

func (c *CatController) GetFavorites() {
    resultChan := makeAPIRequest("favourites")
    
    select {
    case result := <-resultChan:
        c.Ctx.Output.Header("Content-Type", "application/json")
        c.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timeout"}
        c.ServeJSON()
    }
}

func (c *CatController) AddFavorite() {
    var favorite struct {
        ImageID string `json:"image_id"`
        SubID   string `json:"sub_id,omitempty"`
    }

    // Parse request body
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &favorite); err != nil {
        logs.Error("Invalid request body: %v", err)
        c.Data["json"] = map[string]string{"error": "Invalid request body"}
        c.ServeJSON()
        return
    }

    // Validate input
    if favorite.ImageID == "" {
        logs.Warn("Missing Image ID in request")
        c.Data["json"] = map[string]string{"error": "Image ID is required"}
        c.ServeJSON()
        return
    }

    if favorite.SubID == "" {
        favorite.SubID = "user-shino33"
    }

    logs.Info("Payload to send: %+v", favorite)

    apiKey, _ := web.AppConfig.String("cat_api_key")
    baseURL, _ := web.AppConfig.String("api_base_url")

    req := httplib.Post(fmt.Sprintf("%s/favourites", baseURL))
    req.Header("x-api-key", apiKey)
    req.Header("Content-Type", "application/json")

    reqBody, err := json.Marshal(favorite)
    if err != nil {
        logs.Error("Failed to serialize request body: %v", err)
        c.Data["json"] = map[string]string{"error": "Failed to serialize request body"}
        c.ServeJSON()
        return
    }

    // Debugging the exact body being sent
    logs.Info("Serialized Request Body: %s", string(reqBody))

    req.Body(reqBody)

    var result map[string]interface{}
    if err := req.ToJSON(&result); err != nil {
        logs.Error("API error: %v", err)
        c.Data["json"] = map[string]string{"error": "Failed to add favorite"}
    } else {
        logs.Info("API Response: %+v", result)
        c.Data["json"] = result
    }
    c.ServeJSON()
}

func (c *CatController) RemoveFavorite() {
    favoriteID := c.Ctx.Input.Param(":id")
    
    apiKey, _ := web.AppConfig.String("cat_api_key")
    baseURL, _ := web.AppConfig.String("api_base_url")
    
    req := httplib.Delete(fmt.Sprintf("%s/favourites/%s", baseURL, favoriteID))
    req.Header("x-api-key", apiKey)
    
    var result interface{}
    err := req.ToJSON(&result)
    
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to remove favorite"}
    } else {
        c.Data["json"] = result
    }
    c.ServeJSON()
}