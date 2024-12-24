package controllers

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/beego/beego/v2/core/logs"
    "github.com/beego/beego/v2/server/web"
    "testBeego/models"
)

type VoteController struct {
    web.Controller
}

// Submit a vote for a cat image
func (v *VoteController) Vote() {
    logs.Info("Vote endpoint hit")

    bodyBytes := v.Ctx.Input.CopyBody(1024 * 1024)
    if len(bodyBytes) == 0 {
        logs.Warn("Request body is empty after CopyBody")
        v.Data["json"] = map[string]string{"error": "Empty request body"}
        v.ServeJSON()
        return
    }
    logs.Info("Raw request body: %s", string(bodyBytes))

    var vote models.Vote
    if err := json.Unmarshal(bodyBytes, &vote); err != nil {
        logs.Error("Failed to parse request body: %v", err)
        v.Data["json"] = map[string]string{"error": "Invalid request body"}
        v.ServeJSON()
        return
    }

    // Validate fields
    if vote.ImageID == "" || (vote.Value != 1 && vote.Value != -1) {
        logs.Warn("Invalid params: image_id=%s, value=%d", vote.ImageID, vote.Value)
        v.Data["json"] = map[string]string{"error": "Invalid request parameters"}
        v.ServeJSON()
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
        v.Ctx.Output.Header("Content-Type", "application/json")
        v.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        v.Data["json"] = map[string]string{"error": "Request timeout"}
        v.ServeJSON()
    }
}

// Fetch the user's vote history
func (v *VoteController) GetVoteHistory() {
    logs.Info("GetVoteHistory endpoint hit")

    subID := getConfig("sub_id", "default-user-id")
    apiEndpoint := fmt.Sprintf("votes?sub_id=%s&limit=28&order=Desc", subID)

    resultChan := makeAPIRequest(apiEndpoint, "GET", nil)

    select {
    case result := <-resultChan:
        v.Ctx.Output.Header("Content-Type", "application/json")
        v.Ctx.Output.Body(result)
    case <-time.After(10 * time.Second):
        v.Data["json"] = map[string]string{"error": "Request timeout"}
        v.ServeJSON()
    }
}
