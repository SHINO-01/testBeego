// routers/router.go
package routers

import (
    "testBeego/controllers"
    "github.com/beego/beego/v2/server/web"
)

func init() {
    web.Router("/", &controllers.CatController{})
    web.Router("/api/cats/random", &controllers.CatController{}, "get:GetRandomCat")
    
    web.Router("/api/breeds", &controllers.CatController{}, "get:GetBreeds")
    web.Router("/api/breed-images", &controllers.CatController{}, "get:GetBreedImages")

    web.Router("/api/vote", &controllers.CatController{}, "post:Vote")
    
    web.Router("/api/favorites", &controllers.CatController{}, "get:GetFavorites")
    web.Router("/api/favorites", &controllers.CatController{}, "post:AddFavorite")
    web.Router("/api/favorites/:id", &controllers.CatController{}, "delete:RemoveFavorite")

    web.Router("/api/vote_history", &controllers.CatController{}, "get:GetVoteHistory")
}