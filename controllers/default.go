// controllers/cat_controller.go
package controllers

import (
    "encoding/json"
    "fmt"
    "time"
    "github.com/beego/beego/v2/client/httplib"
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
    var vote Vote
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vote); err != nil {
        c.Data["json"] = map[string]string{"error": "Invalid request body"}
        c.ServeJSON()
        return
    }

    apiKey, _ := web.AppConfig.String("cat_api_key")
    baseURL, _ := web.AppConfig.String("api_base_url")
    
    req := httplib.Post(fmt.Sprintf("%s/votes", baseURL))
    req.Header("x-api-key", apiKey)
    req.Header("Content-Type", "application/json")
    req.JSONBody(vote)
    
    var result interface{}
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
    
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &favorite); err != nil {
        c.Data["json"] = map[string]string{"error": "Invalid request body"}
        c.ServeJSON()
        return
    }

    // Optional: Set a default sub_id if not provided
    if favorite.SubID == "" {
        favorite.SubID = "user-shino33"
    }

    apiKey, _ := web.AppConfig.String("cat_api_key")
    baseURL, _ := web.AppConfig.String("api_base_url")
    
    req := httplib.Post(fmt.Sprintf("%s/favourites", baseURL))
    req.Header("x-api-key", apiKey)
    req.Header("Content-Type", "application/json")
    req.JSONBody(favorite)
    
    var result interface{}
    err := req.ToJSON(&result)
    
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to add favorite"}
    } else {
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