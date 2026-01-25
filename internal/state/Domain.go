package state

import "time"

type Domain struct {
	Name       string
	Exists     bool
	ExpiryDate time.Time
}
