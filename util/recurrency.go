package util

const (
	WEEKLY = "WEEKLY"
	MONTHLY = "MONTHLY"
	ANNUAL = "ANNUAL"
)

// IsSupportedRecurrency returns true if the recurrency is supported
func IsSupportedRecurrency(recurrency string) bool {
	switch recurrency {
	case WEEKLY, MONTHLY, ANNUAL:
		return true
	}

	return false
}