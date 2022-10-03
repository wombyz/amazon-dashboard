package handlers

import (
	"github.com/wombyz/amazon-dashboard/internal/config"
	"github.com/wombyz/amazon-dashboard/internal/render"
	"github.com/wombyz/amazon-dashboard/models"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Login is the login page handler
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	// if logged in, redirect to dashboard

	// if not logged in, display login page
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
}

// Dashboard is the admin dashboard page handler
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// BestSellers is the best sellers page handler
func (m *Repository) BestSellers(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "best-sellers.page.tmpl", &models.TemplateData{})
}

// SearchProducts is the search products page handler
func (m *Repository) SearchProducts(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-products.page.tmpl", &models.TemplateData{})
}

// Watchlist is the watchlist page handler
func (m *Repository) Watchlist(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "watchlist.page.tmpl", &models.TemplateData{})
}

// AddProducts is the add products page handler
func (m *Repository) AddProducts(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "add-products.page.tmpl", &models.TemplateData{})
}

// ImportData is the import data page handler
func (m *Repository) ImportData(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "import-data.page.tmpl", &models.TemplateData{})
}

// Logout is the logout function handler
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
}
