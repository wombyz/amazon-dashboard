package models

import (
	"time"
)

// Product is the model for Amazon products by ASIN
type Product struct {
	ASIN        string
	Category    string
	Name        string
	ListingURL  string
	ImgURL      string
	Rating      float64
	ReviewCount int
	WeeklySales float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// StockReading is the model for stock Quantity readings of a particular product as at Date
type StockReading struct {
	RID                string
	ProductASIN        string
	Variant            string
	Quantity           int
	RecordedAt         time.Time
	QtyChangeSinceLast int
	Price              float64
	NumOtherSellers    int
	ShipsFrom          string
	SoldBy             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type StockSlice []StockReading

func (p StockSlice) Len() int {
	return len(p)
}

func (p StockSlice) Less(i, j int) bool {
	return p[i].RecordedAt.Before(p[j].RecordedAt)
}

func (p StockSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// User is the admin dashboard user model. Only users can log in and access the application
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}
