package repository

import "github.com/wombyz/amazon-dashboard/models"

type DatabaseRepo interface {
	InsertProducts(p []models.Product) error
	GetProductByASIN(asin string) (models.Product, error)
	InsertStockReadings(p []models.StockReading) error
	CalculateWeeklySalesByASIN(asin string) (float64, error)
	CalculateWeeklySalesForAll() error
	GetAllProducts() ([]models.Product, error)
	GetAllStockReadingsByASIN(asin string) ([]models.StockReading, error)
	DeleteStockReadings(s []string) error
}
