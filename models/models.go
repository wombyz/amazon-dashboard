package models

import "time"

// Product is the model for Amazon products by ASIN
type Product struct {
	ASIN      int
	Name      string
	Link      string
	ImgURL    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// StockReading is the model for stock Quantity readings of a particular product as at Date
type StockReading struct {
	ProductASIN int
	Date        time.Time
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User is the admin dashboard user model. Only users can log in and access the application
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}
