package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testBeego/controllers"
	"testBeego/models"
)

// MockAPIClient mocks the API client for testing
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) MakeRequest(endpoint, method string, payload interface{}) ([]byte, error) {
	args := m.Called(endpoint, method, payload)
	return args.Get(0).([]byte), args.Error(1)
}

var mockClient *MockAPIClient

func TestMain(m *testing.M) {
	mockClient = new(MockAPIClient)
	controllers.SetAPIClient(mockClient) // Replace the real API client with the mock
	m.Run()
}

func TestGetRandomCat(t *testing.T) {
	mockResponse := `[{"id":"abc123","url":"http://example.com/cat.jpg"}]`

	mockClient.On("MakeRequest", "images/search?limit=1", "GET", nil).Return([]byte(mockResponse), nil)

	r, _ := http.NewRequest("GET", "/api/cats/random", nil)
	w := httptest.NewRecorder()

	controller := &controllers.CatController{}
	controller.GetRandomCat(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, mockResponse, w.Body.String())
	mockClient.AssertExpectations(t)
}

func TestVote(t *testing.T) {
	votePayload := models.Vote{
		ImageID: "test_image",
		SubID:   "test_user",
		Value:   1,
	}
	mockResponse := `{"message":"Vote recorded"}`

	mockClient.On("MakeRequest", "votes", "POST", votePayload).Return([]byte(mockResponse), nil)

	body, _ := json.Marshal(votePayload)
	r, _ := http.NewRequest("POST", "/api/vote", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controller := &controllers.VoteController{}
	controller.Vote(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, mockResponse, w.Body.String())
	mockClient.AssertExpectations(t)
}

func TestAddFavorite(t *testing.T) {
	favoritePayload := models.FavoritePayload{
		ImageID: "test_image",
		SubID:   "test_user",
	}
	mockResponse := `{"message":"Favorite added"}`

	mockClient.On("MakeRequest", "favourites", "POST", favoritePayload).Return([]byte(mockResponse), nil)

	body, _ := json.Marshal(favoritePayload)
	r, _ := http.NewRequest("POST", "/api/favorites", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controller := &controllers.FavoritesController{}
	controller.AddFavorite(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, mockResponse, w.Body.String())
	mockClient.AssertExpectations(t)
}

func TestGetFavorites(t *testing.T) {
	mockResponse := `[{"id":"fav123","image_id":"test_image","sub_id":"test_user"}]`

	mockClient.On("MakeRequest", "favourites?limit=28&order=Desc&sub_id=test_user", "GET", nil).Return([]byte(mockResponse), nil)

	r, _ := http.NewRequest("GET", "/api/favorites", nil)
	w := httptest.NewRecorder()

	controller := &controllers.FavoritesController{}
	controller.GetFavorites(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, mockResponse, w.Body.String())
	mockClient.AssertExpectations(t)
}
