package dbrepo

import (
	"context"
	"fmt"
	"github.com/wombyz/amazon-dashboard/helpers"
	"github.com/wombyz/amazon-dashboard/models"
	"log"
	"sort"
	"strings"
	"time"
)

func (m *postgresDBRepo) InsertProducts(p []models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into products 
    			(asin, category, name, listing_url, img_url, rating, review_count, weekly_sales, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (asin) DO NOTHING;
			`

	for i, product := range p {
		log.Printf("Inserting product %d/%d into products", i, len(p))
		_, err := m.DB.ExecContext(ctx, stmt,
			product.ASIN,
			product.Category,
			product.Name,
			product.ListingURL,
			product.ImgURL,
			product.Rating,
			product.ReviewCount,
			product.WeeklySales,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	log.Printf("Successfully inserted %d products into the products table!", len(p))
	return nil
}

func (m *postgresDBRepo) InsertStockReadings(p []models.StockReading) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into stock_data (rid, product_asin, variant, recorded_at, quantity,
			price, num_of_other_sellers, ships_from, sold_by, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (rid) DO NOTHING;
			`

	for i, reading := range p {
		log.Printf("Inserting stock reading %d/%d into products", i, len(p))
		_, err := m.DB.ExecContext(ctx, stmt,
			reading.RID,
			reading.ProductASIN,
			reading.Variant,
			reading.RecordedAt,
			reading.Quantity,
			reading.Price,
			reading.NumOtherSellers,
			reading.ShipsFrom,
			reading.SoldBy,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	log.Printf("Successfully inserted %d stock readings into the stock_data table!", len(p))
	return nil
}

func (m *postgresDBRepo) GetProductByASIN(asin string) (models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var product models.Product

	query := `
			select
				asin, name, listing_url, img_url, rating, review_count, weekly_sales, created_at, updated_at
			from
			    products
			where
			    asin = $1
			`

	row := m.DB.QueryRowContext(ctx, query, asin)
	err := row.Scan(
		&product.ASIN,
		&product.Name,
		&product.ListingURL,
		&product.ImgURL,
		&product.Rating,
		&product.ReviewCount,
		&product.WeeklySales,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return product, err
	}

	return product, nil
}

func (m *postgresDBRepo) CalculateWeeklySalesByASIN(asin string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var weeklySales float64
	var stockReadings []models.StockReading

	currentDate := time.Now()
	startOfWeekDate := time.Now().Add(-time.Duration(7) * (time.Duration(24) * time.Hour))

	query := `
			select
				product_asin, variant, recorded_at, quantity, price, num_of_other_sellers, ships_from, sold_by
			from
			    stock_data
			where
			    product_asin = $1
			    and $2 < recorded_at and $3 > recorded_at;
			`

	rows, err := m.DB.QueryContext(ctx, query, asin, startOfWeekDate, currentDate)
	if err != nil {
		return weeklySales, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.StockReading
		err := rows.Scan(
			&i.ProductASIN,
			&i.Variant,
			&i.RecordedAt,
			&i.Quantity,
			&i.Price,
			&i.NumOtherSellers,
			&i.ShipsFrom,
			&i.SoldBy,
		)

		if err != nil {
			return weeklySales, err
		}

		stockReadings = append(stockReadings, i)
	}

	log.Println("Number of stock readings:", len(stockReadings))

	if err = rows.Err(); err != nil {
		return weeklySales, err
	}

	sellers := make(map[string]bool)
	var quantities []int

	for _, r := range stockReadings {
		log.Println(r.SoldBy)
		quantities = append(quantities, r.Quantity)
		sellers[strings.TrimSpace(r.SoldBy)] = true
	}

	if len(sellers) > 1 {
		log.Println("more than one seller this week. ")
		return weeklySales, nil
	}

	var stockDataMap = make(map[string]models.StockReading)
	for i, r := range stockReadings {
		stockDataMap[fmt.Sprintf("%d", i)] = r
	}

	log.Println(stockDataMap)
	//Sort the map by date
	sortedReadings := make(models.StockSlice, 0, len(stockDataMap))
	for _, d := range stockDataMap {
		sortedReadings = append(sortedReadings, d)
	}
	sort.Sort(sortedReadings)

	log.Println(sortedReadings)

	var qtySlice []int
	for _, row := range sortedReadings {
		qtySlice = append(qtySlice, row.Quantity)
	}
	log.Println(qtySlice)
	changeSlice := helpers.CalculateQuantityChangesSlice(qtySlice)

	log.Println(changeSlice)

	for i, row := range sortedReadings {
		row.QtyChangeSinceLast = changeSlice[i]

		if row.QtyChangeSinceLast > 0 {
			weeklySales += helpers.CalculateSalesSinceLastReading(row)
		}
	}

	return weeklySales, nil
}
