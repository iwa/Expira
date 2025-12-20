package api

import (
	"fmt"
	"net/http"

	"github.com/iwa/Expira/internal/state"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	fmt.Fprintf(w, "Service is running")
}

// StatusHandlerFactory creates a status handler with access to the domain store.
// This allows the handler to use dependency injection instead of global state.
func StatusHandlerFactory(store *state.DomainStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		var status string = "Domain Expiry Watcher Status:\n\n"

		domains := store.GetAllDomains()

		for _, domain := range domains {
			if domain.ExpiryDate.IsZero() {
				status += fmt.Sprintf("Domain %s: Expiry date not set\n", domain.Name)
			} else {
				status += fmt.Sprintf("Domain %s: Expires on %s\n", domain.Name, domain.ExpiryDate.Format("2006-01-02"))
			}
		}

		fmt.Fprintf(w, "%s", status)
	}
}
