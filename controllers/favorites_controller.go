//favorites_controller.go

package controllers

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/beego/beego/v2/core/logs"
    "github.com/beego/beego/v2/server/web"
    "testBeego/models"
)

type FavoritesController struct {
    web.Controller
}

// Fetch the user's favorites (latest 28)
func (f *FavoritesController) GetFavorites() {
    subID := getConfig("sub_id", "default-user-id")
    apiEndpoint := fmt.Sprintf("favourites?limit=28&order=Desc&sub_id=%s", subID)

    resultChan := makeAPIRequest(apiEndpoint, "GET", nil)

    select {
    case result := <-resultChan:
        f.Ctx.Output.Header("Content-Type", "application/json")
        f.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        f.Data["json"] = map[string]string{"error": "Request timeout"}
        f.ServeJSON()
    }
}

// Add a cat image to favorites
func (f *FavoritesController) AddFavorite() {
    bodyBytes := f.Ctx.Input.CopyBody(1024 * 1024)
    if len(bodyBytes) == 0 {
        logs.Warn("Request body is empty after CopyBody")
        f.Data["json"] = map[string]string{"error": "Empty request body"}
        f.ServeJSON()
        return
    }
    logs.Info("Raw request body: %s", string(bodyBytes))

    var fav models.FavoritePayload
    if err := json.Unmarshal(bodyBytes, &fav); err != nil {
        f.Data["json"] = map[string]string{"error": "Invalid request body"}
        f.ServeJSON()
        return
    }

    if fav.ImageID == "" {
        f.Data["json"] = map[string]string{"error": "Image ID is required"}
        f.ServeJSON()
        return
    }

    // Assign sub_id if missing
    if fav.SubID == "" {
        fav.SubID = getConfig("sub_id", "default-user-id")
    }

    resultChan := makeAPIRequest("favourites", "POST", fav)

    select {
    case result := <-resultChan:
        f.Ctx.Output.Header("Content-Type", "application/json")
        f.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        f.Data["json"] = map[string]string{"error": "Request timeout"}
        f.ServeJSON()
    }
}

// Remove a cat image from favorites
func (f *FavoritesController) RemoveFavorite() {
    favoriteID := f.Ctx.Input.Param(":id")
    apiEndpoint := fmt.Sprintf("favourites/%s", favoriteID)

    resultChan := makeAPIRequest(apiEndpoint, "DELETE", nil)

    select {
    case result := <-resultChan:
        f.Ctx.Output.Header("Content-Type", "application/json")
        f.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        f.Data["json"] = map[string]string{"error": "Request timeout"}
        f.ServeJSON()
    }
}
