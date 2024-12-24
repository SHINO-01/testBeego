package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"runtime"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	_ "testBeego/routers"   // Import routers to initialize routes

	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestVote(t *testing.T) {
	validPayload := `{"image_id":"test_image","value":1}`
	invalidPayload := `{"image_id":"","value":1}`

	Convey("Subject: Test Vote Endpoint\n", t, func() {
		Convey("Valid Payload Should Return Status 200", func() {
			r, _ := http.NewRequest("POST", "/api/vote", strings.NewReader(validPayload))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			So(w.Code, ShouldEqual, 200)
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})

		Convey("Invalid Payload Should Return Status 400", func() {
			r, _ := http.NewRequest("POST", "/api/vote", strings.NewReader(invalidPayload))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			So(w.Code, ShouldEqual, 400)
		})
	})
}

func TestAddFavorite(t *testing.T) {
	validPayload := `{"image_id":"test_image"}`
	invalidPayload := `{"invalid_field":"test_image"}`

	Convey("Subject: Test Add Favorite Endpoint\n", t, func() {
		Convey("Valid Payload Should Return Status 200", func() {
			r, _ := http.NewRequest("POST", "/api/favorites", strings.NewReader(validPayload))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			So(w.Code, ShouldEqual, 200)
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})

		Convey("Invalid Payload Should Return Status 400", func() {
			r, _ := http.NewRequest("POST", "/api/favorites", strings.NewReader(invalidPayload))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			So(w.Code, ShouldEqual, 400)
		})
	})
}

func TestGetRandomCat(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/cats/random", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Subject: Test Get Random Cat Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestGetVotingHistory(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/vote_history", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Subject: Test Get Voting History Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestInvalidVotePayload(t *testing.T) {
    invalidPayload := `{"invalid_field":"value"}`
    r, _ := http.NewRequest("POST", "/api/vote", strings.NewReader(invalidPayload))
    r.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    beego.BeeApp.Handlers.ServeHTTP(w, r)

    logs.Trace("testing", "TestInvalidVotePayload", "Code[%d]\n%s", w.Code, w.Body.String())

    Convey("Subject: Test Invalid Vote Payload\n", t, func() {
        Convey("Status Code Should Be 400", func() {
            So(w.Code, ShouldEqual, 400)
        })
    })
}

func TestMissingVoteHeaders(t *testing.T) {
    validPayload := `{"image_id":"test_image","value":1}`
    r, _ := http.NewRequest("POST", "/api/vote", strings.NewReader(validPayload))
    w := httptest.NewRecorder() // No Content-Type header

    beego.BeeApp.Handlers.ServeHTTP(w, r)
    logs.Trace("testing", "TestMissingVoteHeaders", "Code[%d]\n%s", w.Code, w.Body.String())

    Convey("Subject: Test Missing Vote Headers\n", t, func() {
        Convey("Status Code Should Be 400 or 415", func() {
            So(w.Code, ShouldBeIn, []int{400, 415})
        })
    })
}
