package handlers

import (
	"encoding/csv"
	"fmt"
	"github.com/wombyz/amazon-dashboard/helpers"
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
	"time"
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

	//data := make(map[string]interface{})
	//data["sampleProducts"] = products

	render.Template(w, r, "best-sellers.page.tmpl", &models.TemplateData{
		//Data: data,
	})
}

// ShowProduct is the products breakdown page handler
func (m *Repository) ShowProduct(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.URL.RequestURI(), "/")
	asin := exploded[3]

	product, err := m.DB.GetProductByASIN(asin)
	if err != nil {
		log.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["product"] = product

	render.Template(w, r, "product.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// SearchProducts is the search products page handler
func (m *Repository) SearchProducts(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.URL.RequestURI(), "=")

	if len(exploded) > 1 {
		asin := exploded[1]
		weeklySales, err := m.DB.CalculateWeeklySalesByASIN(asin)
		if err != nil {
			log.Println(err)
			// Raise error message
			render.Template(w, r, "search-products.page.tmpl", &models.TemplateData{})
		}
		log.Println("Weekly Sales:", weeklySales)

		//data := make(map[string]interface{})
		//data["product"] = product

		//render.Template(w, r, "product.page.tmpl", &models.TemplateData{
		//	Data: data,
		//})
	}

	render.Template(w, r, "search-products.page.tmpl", &models.TemplateData{})

}

// SearchProducts is the search products page handler
func (m *Repository) SearchProductByID(w http.ResponseWriter, r *http.Request) {

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

// PostAddProducts is the UploadNewProducts function handler
func (m *Repository) PostAddProducts(w http.ResponseWriter, r *http.Request) {
	var filePath string

	// Save csv data
	if r.Method == "GET" {
		fmt.Println("Get request?")
	} else {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Println(err)
			return
		}
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		fmt.Fprintf(w, "%v", handler.Header)
		filePath = "./cache/" + handler.Filename
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
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

	// Read csv data into product objects
	var productsToUpload []models.Product

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	productsData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	for _, record := range productsData {
		if record[0] == "ASIN" {
			continue
		}
		rt, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			rt = 0
			log.Println("error when parsing rating string to float:", err)
		}

		rc, err := strconv.Atoi(record[6])
		if err != nil {
			rc = 0
			log.Println("error when parsing review count string to int:", err)
		}

		r := models.Product{
			ASIN:        record[0],
			Category:    record[1],
			Name:        record[2],
			ListingURL:  record[3],
			ImgURL:      record[4],
			Rating:      rt,
			ReviewCount: rc,
			WeeklySales: 0.0,
		}

		productsToUpload = append(productsToUpload, r)
	}

	// Save to database
	err = m.DB.InsertProducts(productsToUpload)
	if err != nil {
		log.Println(err)
	}

	render.Template(w, r, "add-products.page.tmpl", &models.TemplateData{})
}

// PostUploadData is the upload new stock data function handler
func (m *Repository) PostUploadData(w http.ResponseWriter, r *http.Request) {
	var filePath string

	// Save csv data
	if r.Method == "GET" {
		fmt.Println("Get request?")
	} else {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Println(err)
			return
		}
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		fmt.Fprintf(w, "%v", handler.Header)
		filePath = "./cache/" + handler.Filename
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
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

	// Read csv data into StockReading objects
	var stockReadingsToUpload []models.StockReading

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	stockData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	for _, reading := range stockData {
		if reading[0] == "RID (reading id)" {
			continue
		}

		date, err := helpers.ConvertToDatetime(reading[3], reading[4])
		if err != nil {
			log.Println("error when converting datetime:", err)
			return
		}

		quantity, err := strconv.Atoi(reading[6])
		if err != nil {
			log.Println("error when converting quantity to int:", err)
			quantity = 1234
		}

		rawPrice := reading[7]
		if rawPrice == "" {
			continue
		}
		rawPrice = strings.Replace(rawPrice, ",", "", -1)
		rawPrice = strings.Replace(rawPrice, "$", "", -1)

		price, err := strconv.ParseFloat(rawPrice, 64)
		if err != nil {
			log.Println("error when converting price to float64:", err)
			return
		}

		numOtherSellers, err := strconv.Atoi(reading[9])
		if err != nil {
			log.Println("error when converting other seller # to int:", err)
			numOtherSellers = 0
		}

		r := models.StockReading{
			RID:             reading[0],
			ProductASIN:     reading[1],
			Variant:         reading[2],
			Quantity:        quantity,
			RecordedAt:      date,
			Price:           price,
			NumOtherSellers: numOtherSellers,
			ShipsFrom:       reading[10],
			SoldBy:          reading[11],
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}
		stockReadingsToUpload = append(stockReadingsToUpload, r)
	}

	// Save to database
	err = m.DB.InsertStockReadings(stockReadingsToUpload)
	if err != nil {
		log.Println(err)
	}

	render.Template(w, r, "import-data.page.tmpl", &models.TemplateData{})
}
