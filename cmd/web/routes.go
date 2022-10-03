package main

import (
	"github.com/go-chi/chi"
	"github.com/wombyz/amazon-dashboard/internal/config"
	"github.com/wombyz/amazon-dashboard/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	//mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	// mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(handlers.Repo.Login))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		//mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.Dashboard)
		mux.Get("/search-products", handlers.Repo.SearchProducts)
		mux.Get("/best-sellers", handlers.Repo.BestSellers)
		mux.Get("/watchlist", handlers.Repo.Watchlist)
		mux.Get("/add-products", handlers.Repo.AddProducts)
		mux.Get("/import-data", handlers.Repo.ImportData)
		mux.Get("/logout", handlers.Repo.Logout)
	})

	return mux
}
