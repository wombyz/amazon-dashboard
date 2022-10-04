package handlers

import (
	"fmt"
	"github.com/wombyz/amazon-dashboard/internal/config"
	"github.com/wombyz/amazon-dashboard/internal/driver"
	"github.com/wombyz/amazon-dashboard/internal/render"
	"github.com/wombyz/amazon-dashboard/internal/repository"
	"github.com/wombyz/amazon-dashboard/internal/repository/dbrepo"
	"github.com/wombyz/amazon-dashboard/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
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
	products := models.CreateSampleData()

	data := make(map[string]interface{})
	data["sampleProducts"] = products

	render.Template(w, r, "best-sellers.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// ShowProduct is the products breakdown page handler
func (m *Repository) ShowProduct(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.URL.RequestURI(), "/")
	id, err := strconv.Atoi(exploded[3])
	log.Println(id)
	if err != nil {
		log.Println(err)
		return
	}

	products := models.CreateSampleData()

	data := make(map[string]interface{})
	data["product"] = products[id]

	render.Template(w, r, "show-product.page.tmpl", &models.TemplateData{
		Data: data,
	})
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

// PostUploadNewProducts is the UploadNewProducts function handler
func (m *Repository) PostUploadNewProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		fmt.Println("Get request?")
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			log.Println(err)
		}
	}

	render.Template(w, r, "import-data.page.tmpl", &models.TemplateData{})
}
