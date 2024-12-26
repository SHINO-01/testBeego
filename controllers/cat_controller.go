//cat_controller.go

package controllers

import (
    "fmt"
    "time"

    "github.com/beego/beego/v2/core/logs"
    "github.com/beego/beego/v2/server/web"
    // "myapp/models" // If needed for certain logic
)

type CatController struct {
    web.Controller
}

// Renders the index page
func (c *CatController) Get() {
    c.TplName = "index.tpl"
}

// Fetches a random cat image
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

// Fetches all cat breeds
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

// Fetch images for a specific breed
func (c *CatController) GetBreedImages() {
    logs.Info("GetBreedImages Endpoint Hit")

    breedID := c.GetString("breed_id")
    limit, err := c.GetInt("limit", 8)
    if err != nil {
        logs.Error("Invalid Limit Parameter: %v", err)
        c.Data["json"] = map[string]string{"error": "Invalid limit parameter"}
        c.ServeJSON()
        return
    }

    apiEndpoint := fmt.Sprintf("images/search?breed_ids=%s&limit=%d", breedID, limit)
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
