package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

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