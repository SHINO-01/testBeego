package routers

import (
    "github.com/beego/beego/v2/server/web"
    "testBeego/controllers"
)

func init() {
    // Cat routes
    web.Router("/", &controllers.CatController{}, "get:Get")
    web.Router("/api/cats/random", &controllers.CatController{}, "get:GetRandomCat")
    web.Router("/api/breeds", &controllers.CatController{}, "get:GetBreeds")
    web.Router("/api/breed-images", &controllers.CatController{}, "get:GetBreedImages")

    // Voting routes
    web.Router("/api/vote", &controllers.VoteController{}, "post:Vote")
    web.Router("/api/vote_history", &controllers.VoteController{}, "get:GetVoteHistory")

    // Favorites routes
    web.Router("/api/favorites", &controllers.FavoritesController{}, "get:GetFavorites;post:AddFavorite")
    web.Router("/api/favorites/:id", &controllers.FavoritesController{}, "delete:RemoveFavorite")
}
