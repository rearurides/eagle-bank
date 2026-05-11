package domain

type Currency string

const (
	GBP Currency = "GBP"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

// getMinorUnit returns the minor unit multiplier for a given currency.
// For example, for GBP, USD, and EUR, the minor unit is 100 (i.e., 1 pound/dollar/euro = 100 pence/cents).
func getMinorUnit(currency Currency) int64 {
	switch currency {
	case GBP, USD, EUR:
		return 100
	default:
		return 100
	}
}
