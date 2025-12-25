package util

import (
	"fmt"
	"strings"
	"time"

	"math/rand"

	decimal "github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for range n {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomUsername generates a random username
func RandomUsername() string {
	return RandomString(8)
}

// RandomFullName generates a random username
func RandomFullName() string {
	return fmt.Sprintf("%s %s", RandomString(5), RandomString(8))
}

// RandomTitle generates a random title
func RandomTitle() string {
	return RandomString(14)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

// RandomFutureDate generates a random date in the future within a year
func RandomFutureDate() time.Time {
	next := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	return next.AddDate(2025, 0, rand.Intn(365))
}

// RandomMonthDate generates a random month first and end date
func RandomMonthDate() (time.Time, time.Time) {
	startDate := time.Date(2025, time.Month(rand.Intn(12)), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, startDate.Month(), 28, 0, 0, 0, 0, time.UTC)

	return startDate, endDate
}

// RandomYearDate generates a random year first and end date
func RandomYearDate() (time.Time, time.Time) {
	startDate := time.Date(2025+rand.Intn(100), 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(startDate.Year(), 12, 31, 0, 0, 0, 0, time.UTC)

	return startDate, endDate
}

// RandomMoney generates a random amount between -500 and 500
func RandomMoney() decimal.Decimal {
	return decimal.NewFromFloat(float64(rand.Intn(100000)) / 100)
}

// RandomBool generates a random bool
func RandomBool() bool {
	return rand.Intn(2) == 1
}

// RandomCurrency generates a random currency
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	return currencies[rand.Intn(len(currencies))]
}

// RandomRecurrency generates a random currency
func RandomRecurrency() string {
	recurrencies := []string{WEEKLY, MONTHLY, ANNUAL}
	return recurrencies[rand.Intn(len(recurrencies))]
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}
