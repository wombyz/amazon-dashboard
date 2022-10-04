package models

import "time"

var sampleProducts []Product

func CreateSampleData() []Product {
	p1 := Product{
		ASIN:        "B07D15V7T3",
		Name:        "Farberware Professional Heat Resistant Nylon Meat and Potato Masher, Safe for Non-Stick Cookware, 10-Inch, Black",
		ListingURL:  "https://www.amazon.com.au/Farberware-5211438-Professional-Resistant-Masher-Safe/dp/B07D15V7T3",
		ImgURL:      "https://m.media-amazon.com/images/I/518Siu2n+AL._AC_UY218_.jpg",
		Rating:      4.8,
		ReviewCount: 36076,
		WeeklySales: 2340.32,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	p2 := Product{
		ASIN:        "B019HR91W4",
		Name:        "Ecolution Patented Microwave Micro-Pop Popcorn Popper, Borosilicate Glass, 3-in-1 Lid, Dishwasher Safe, BPA Free, 1.5 Quart Snack Size, Red.",
		ListingURL:  "https://www.amazon.com.au/Ecolution-Micro-Pop-Microwave-Popcorn-Popper/dp/B019HR91W4",
		ImgURL:      "https://m.media-amazon.com/images/I/81aTLT1IstL._AC_UY218_.jpg",
		Rating:      4.4,
		ReviewCount: 37556,
		WeeklySales: 3240.1,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	p3 := Product{
		ASIN:        "B00004OCIP",
		Name:        "OXO Good Grips Swivel Peeler",
		ListingURL:  "https://www.amazon.com.au/OXO-Good-Grips-Swivel-Peeler/dp/B00004OCIP",
		ImgURL:      "https://m.media-amazon.com/images/I/71UBG06NKFL._AC_UY218_.jpg",
		Rating:      4.8,
		ReviewCount: 44981,
		WeeklySales: 668.12,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	sampleProducts = append(sampleProducts, p1, p2, p3)

	return sampleProducts
}

// Product is the model for Amazon products by ASIN
type Product struct {
	ASIN        string
	Name        string
	ListingURL  string
	ImgURL      string
	Rating      float32
	ReviewCount int
	WeeklySales float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
