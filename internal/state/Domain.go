package state

import "time"

type Domain struct {
	Name       string
	Exists     bool
	ExpiryDate time.Time
}

// Calculate how many days between now and Domain expiry
// Diff is calculated using time.Now()
//
// If the domain is already expired, returns -1
func (d *Domain) GetDaysUntilExpiry() int {
	currentTime := time.Now()
	daysLeft := int(d.ExpiryDate.Sub(currentTime).Hours()/24) + 1 // Add 1 to include the current day

	if daysLeft < 0 {
		return -1
	}

	return daysLeft
}
