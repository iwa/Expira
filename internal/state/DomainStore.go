package state

import (
	"sync"
)

// DomainStore provides thread-safe access to domain data.
// It uses a read-write mutex to allow concurrent reads while ensuring safe writes.
type DomainStore struct {
	mu      sync.RWMutex
	domains map[string]Domain
}

// NewDomainStore creates a new DomainStore instance
func NewDomainStore() *DomainStore {
	return &DomainStore{
		domains: make(map[string]Domain),
	}
}

// GetDomain retrieves a domain by name.
// Returns the domain and true if found, or an empty domain and false if not found.
func (ds *DomainStore) GetDomain(name string) (Domain, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	domain, ok := ds.domains[name]
	return domain, ok
}

// SetDomain stores or updates a domain.
// This method is thread-safe and can be called from multiple goroutines.
func (ds *DomainStore) SetDomain(name string, domain Domain) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.domains[name] = domain
}

// GetAllDomains returns a copy of all domains.
// This ensures the returned map cannot cause race conditions if modified by the caller.
func (ds *DomainStore) GetAllDomains() map[string]Domain {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Create a copy to avoid exposing internal map
	domainsCopy := make(map[string]Domain, len(ds.domains))
	for k, v := range ds.domains {
		domainsCopy[k] = v
	}

	return domainsCopy
}

// SetDomains replaces all domains with the provided map.
// This is useful for bulk initialization.
func (ds *DomainStore) SetDomains(domains map[string]Domain) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.domains = domains
}

// Count returns the number of domains in the store
func (ds *DomainStore) Count() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return len(ds.domains)
}
