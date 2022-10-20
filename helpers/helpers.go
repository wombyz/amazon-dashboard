package helpers

import (
	"github.com/wombyz/amazon-dashboard/models"
	"math"
	"time"
)

func ConvertToDatetime(date string, t string) (time.Time, error) {
	var formattedDate time.Time

	layout := "1/2/2006 3:04:05 PM"
	formattedDate, err := time.Parse(layout, (date + " " + t))
	if err != nil {
		return formattedDate, err
	}

	return formattedDate, nil
}

func CalculateSalesSinceLastReading(r models.StockReading) float64 {
	otherSellersMultiplier := float64(r.NumOtherSellers) + 1.0
	sales := float64(r.QtyChangeSinceLast) * otherSellersMultiplier * r.Price

	return sales
}

func CalculateQuantityChangesSlice(q []int) []int {
	prev := 0
	var changeInStockSlice []int

	for i, qty := range q {
		if i == 0 {
			changeInStockSlice = append(changeInStockSlice, 0)
			prev = qty
			continue
		}
		change := prev - qty
		changeInStockSlice = append(changeInStockSlice, change)
		prev = qty
	}
	return changeInStockSlice
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
