package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/beego/beego/v2/server/web/context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "testBeego/mocks"
    "testBeego/models"
    "testBeego/controllers"
)

type ControllerTestSuite struct {
    suite.Suite
    mockClient *mocks.MockAPIClient
}

func (suite *ControllerTestSuite) SetupTest() {
    suite.mockClient = mocks.NewMockAPIClient()
}

func (suite *ControllerTestSuite) TearDownTest() {
    suite.mockClient.ClearCalls()
}

func TestControllerSuite(t *testing.T) {
    suite.Run(t, new(ControllerTestSuite))
}

// Helper function to setup controller test
func setupControllerTest(method, path string, body string) (*httptest.ResponseRecorder, *http.Request) {
    recorder := httptest.NewRecorder()
    req, _ := http.NewRequest(method, path, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    return recorder, req
}

func (suite *ControllerTestSuite) TestCatController_Get() {
    recorder, req := setupControllerTest("GET", "/", "")
    
    controller := &controllers.CatController{}
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.Get()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), "index.tpl", controller.TplName)
}

func (suite *ControllerTestSuite) TestCatController_GetRandomCat() {
    recorder, req := setupControllerTest("GET", "/random-cat", "")
    
    mockResponse := []byte(`[{"id":"test-id","url":"http://example.com/cat.jpg"}]`)
    suite.mockClient.SetResponse("images/search?limit=1", "GET", mockResponse)
    
    controller := &controllers.CatController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.GetRandomCat()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
    
    calls := suite.mockClient.GetCalls("images/search?limit=1", "GET")
    assert.Equal(suite.Suite.T(), 1, len(calls))
}

func (suite *ControllerTestSuite) TestVoteController_Vote() {
    vote := models.Vote{
        ImageID: "test-123",
        Value:   1,
        SubID:   "test-user",
    }
    voteJSON, _ := json.Marshal(vote)
    recorder, req := setupControllerTest("POST", "/vote", string(voteJSON))
    
    mockResponse := []byte(`{"message":"Success"}`)
    suite.mockClient.SetResponse("votes", "POST", mockResponse)
    
    controller := &controllers.VoteController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.Vote()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
    
    calls := suite.mockClient.GetCalls("votes", "POST")
    assert.Equal(suite.Suite.T(), 1, len(calls))
    assert.Equal(suite.Suite.T(), vote, calls[0].Payload)
}

func (suite *ControllerTestSuite) TestVoteController_Vote_InvalidPayload() {
    recorder, req := setupControllerTest("POST", "/vote", "{invalid json}")
    
    controller := &controllers.VoteController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.Vote()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    
    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "Invalid request body", response["error"])
}

func (suite *ControllerTestSuite) TestFavoritesController_GetFavorites() {
    recorder, req := setupControllerTest("GET", "/favorites", "")
    
    mockResponse := []byte(`[{"id":"fav-1","image_id":"img-1"}]`)
    suite.mockClient.SetResponse("favourites?limit=28&order=Desc&sub_id=default-user-id", "GET", mockResponse)
    
    controller := &controllers.FavoritesController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.GetFavorites()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
}

func (suite *ControllerTestSuite) TestFavoritesController_AddFavorite() {
    favorite := models.FavoritePayload{
        ImageID: "test-123",
        SubID:   "test-user",
    }
    favJSON, _ := json.Marshal(favorite)
    recorder, req := setupControllerTest("POST", "/favorites", string(favJSON))
    
    mockResponse := []byte(`{"id":1,"message":"Success"}`)
    suite.mockClient.SetResponse("favourites", "POST", mockResponse)
    
    controller := &controllers.FavoritesController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.AddFavorite()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
    
    calls := suite.mockClient.GetCalls("favourites", "POST")
    assert.Equal(suite.Suite.T(), 1, len(calls))
    assert.Equal(suite.Suite.T(), favorite, calls[0].Payload)
}

func (suite *ControllerTestSuite) TestFavoritesController_RemoveFavorite() {
    recorder, req := setupControllerTest("DELETE", "/favorites/123", "")
    
    mockResponse := []byte(`{"status":"success"}`)
    suite.mockClient.SetResponse("favourites/123", "DELETE", mockResponse)
    
    controller := &controllers.FavoritesController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    // Set up the path parameters
    ctx.Input.SetParam(":id", "123")
    
    controller.RemoveFavorite()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
    
    calls := suite.mockClient.GetCalls("favourites/123", "DELETE")
    assert.Equal(suite.Suite.T(), 1, len(calls))
}

// Additional test cases for better coverage

func (suite *ControllerTestSuite) TestCatController_GetRandomCat_Timeout() {
    recorder, req := setupControllerTest("GET", "/random-cat", "")

    // Simulate timeout in MockAPIClient
    ch := make(chan []byte)
    go func() {
        time.Sleep(11 * time.Second) // Simulates a delay greater than 10 seconds
        ch <- []byte(`{"error": "Request timeout"}`)
    }()
    suite.mockClient.SetResponse("images/search?limit=1", "GET", <-ch)

    controller := &controllers.CatController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.GetRandomCat()

    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "Request timeout", response["error"])
}

func (suite *ControllerTestSuite) TestVoteController_Vote_InvalidValue() {
    vote := models.Vote{
        ImageID: "test-123",
        Value:   2, // Invalid value, should be 1 or -1
        SubID:   "test-user",
    }
    voteJSON, _ := json.Marshal(vote)
    recorder, req := setupControllerTest("POST", "/vote", string(voteJSON))
    
    controller := &controllers.VoteController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.Vote()
    
    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "Invalid request parameters", response["error"])
}

func (suite *ControllerTestSuite) TestVoteController_GetVoteHistory() {
    recorder, req := setupControllerTest("GET", "/vote-history", "")
    
    mockResponse := []byte(`[{"id":"vote-1","image_id":"img-1","value":1}]`)
    suite.mockClient.SetResponse("votes?sub_id=default-user-id&limit=28&order=Desc", "GET", mockResponse)
    
    controller := &controllers.VoteController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx
    
    controller.GetVoteHistory()
    
    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
}

func (suite *ControllerTestSuite) TestCatController_GetBreeds() {
    recorder, req := setupControllerTest("GET", "/breeds", "")

    mockResponse := []byte(`[{"id":"breed-1","name":"Breed Name"}]`)
    suite.mockClient.SetResponse("breeds", "GET", mockResponse)

    controller := &controllers.CatController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.GetBreeds()

    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
}

func (suite *ControllerTestSuite) TestVoteController_GetVoteHistory_NoVotes() {
    recorder, req := setupControllerTest("GET", "/vote-history", "")

    mockResponse := []byte(`[]`) // No votes
    suite.mockClient.SetResponse("votes?sub_id=default-user-id&limit=28&order=Desc", "GET", mockResponse)

    controller := &controllers.VoteController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.GetVoteHistory()

    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
}


func (suite *ControllerTestSuite) TestFavoritesController_AddFavorite_InvalidPayload() {
    recorder, req := setupControllerTest("POST", "/favorites", `{}`) // Empty payload

    controller := &controllers.FavoritesController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.AddFavorite()

    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "Image ID is required", response["error"])
}


func (suite *ControllerTestSuite) TestCatAPIClient_ErrorHandling() {
    client := &controllers.CatAPIClient{}

    // Simulate a failed request
    resultChan := client.MakeAPIRequest("invalid-endpoint", "INVALID", nil)

    result := <-resultChan
    var response map[string]string
    json.Unmarshal(result, &response)

    assert.Equal(suite.Suite.T(), "Invalid HTTP method", response["error"])
}

func BenchmarkGetRandomCat(b *testing.B) {
    mockClient := mocks.NewMockAPIClient()
    mockResponse := []byte(`[{"id":"test-id","url":"http://example.com/cat.jpg"}]`)
    mockClient.SetResponse("images/search?limit=1", "GET", mockResponse)

    controller := &controllers.CatController{
        APIClient: mockClient,
    }

    for i := 0; i < b.N; i++ {
        recorder, req := setupControllerTest("GET", "/random-cat", "")
        ctx := context.NewContext()
        ctx.Reset(recorder, req)
        controller.Init(ctx, "", "", nil)
        controller.Ctx = ctx

        controller.GetRandomCat()
    }
}

func (suite *ControllerTestSuite) TestCatController_GetBreedImages_ValidBreed() {
    recorder, req := setupControllerTest("GET", "/breed-images?breed_id=abc123&limit=5", "")

    mockResponse := []byte(`[{"id":"image-1","url":"http://example.com/image1.jpg"}]`)
    suite.mockClient.SetResponse("images/search?breed_ids=abc123&limit=5", "GET", mockResponse)

    controller := &controllers.CatController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.GetBreedImages()

    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
}

func (suite *ControllerTestSuite) TestCatController_GetBreedImages_InvalidLimit() {
    recorder, req := setupControllerTest("GET", "/breed-images?breed_id=abc123&limit=invalid", "")

    controller := &controllers.CatController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.GetBreedImages()

    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "Invalid limit parameter", response["error"])
}

func (suite *ControllerTestSuite) TestFavoritesController_GetFavorites_Empty() {
    recorder, req := setupControllerTest("GET", "/favorites", "")

    mockResponse := []byte(`[]`) // No favorites
    suite.mockClient.SetResponse("favourites?limit=28&order=Desc&sub_id=default-user-id", "GET", mockResponse)

    controller := &controllers.FavoritesController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.GetFavorites()

    assert.Equal(suite.Suite.T(), http.StatusOK, recorder.Code)
    assert.Equal(suite.Suite.T(), string(mockResponse), recorder.Body.String())
}

func (suite *ControllerTestSuite) TestFavoritesController_GetFavorites_APIError() {
    recorder, req := setupControllerTest("GET", "/favorites", "")

    mockResponse := []byte(`{"error":"API failure"}`)
    suite.mockClient.SetResponse("favourites?limit=28&order=Desc&sub_id=default-user-id", "GET", mockResponse)

    controller := &controllers.FavoritesController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.GetFavorites()

    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "API failure", response["error"])
}

func (suite *ControllerTestSuite) TestVoteController_Vote_EmptyBody() {
    recorder, req := setupControllerTest("POST", "/vote", "")

    controller := &controllers.VoteController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.Vote()

    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "Empty request body", response["error"])
}

func (suite *ControllerTestSuite) TestVoteController_Vote_MissingImageID() {
    vote := models.Vote{
        Value: 1,
        SubID: "test-user",
    }
    voteJSON, _ := json.Marshal(vote)
    recorder, req := setupControllerTest("POST", "/vote", string(voteJSON))

    controller := &controllers.VoteController{
        APIClient: suite.mockClient,
    }
    ctx := context.NewContext()
    ctx.Reset(recorder, req)
    controller.Init(ctx, "", "", nil)
    controller.Ctx = ctx

    controller.Vote()

    var response map[string]string
    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.NoError(suite.Suite.T(), err)
    assert.Equal(suite.Suite.T(), "Invalid request parameters", response["error"])
}

func (suite *ControllerTestSuite) TestCatAPIClient_MissingAPIKey() {
    client := &controllers.CatAPIClient{}
    resultChan := client.MakeAPIRequest("images/search", "GET", nil)

    result := <-resultChan
    var response map[string]string
    json.Unmarshal(result, &response)

    assert.Equal(suite.Suite.T(), "Failed to fetch data", response["error"])
}









