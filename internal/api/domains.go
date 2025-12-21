package api

import (
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/iwa/Expira/internal/state"
)

const domainsPageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Expira</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: #fff;
            padding: 2rem;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            color: #333;
            margin-bottom: 2rem;
        }
        .stats {
            display: flex;
            gap: 1rem;
            margin-bottom: 2rem;
        }
        .stat {
            padding: 0.5rem 1rem;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        .stat strong {
            display: block;
            font-size: 1.5rem;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            background: #fff;
            border: 1px solid #ddd;
        }
        th {
            background: #f8f9fa;
            padding: 1rem;
            text-align: left;
            font-weight: 600;
            border-bottom: 2px solid #ddd;
        }
        td {
            padding: 1rem;
            border-bottom: 1px solid #eee;
        }
        tr:hover {
            background: #f8f9fa;
        }
        .status-badge {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            border-radius: 4px;
            font-size: 0.875rem;
            font-weight: 600;
        }
        .status-critical {
            background: #fee;
            color: #c33;
        }
        .status-warning {
            background: #ffeaa7;
            color: #d63031;
        }
        .status-good {
            background: #d4edda;
            color: #155724;
        }
        .empty-state {
            padding: 4rem 2rem;
            text-align: center;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        .empty-state h2 {
            color: #6c757d;
            margin-bottom: 1rem;
        }
        .empty-state p {
            color: #adb5bd;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Expira</h1>

        {{if .Domains}}
        <div class="stats">
            <div class="stat">
                <span>Total Domains</span>
                <strong>{{.TotalCount}}</strong>
            </div>
            <div class="stat">
                <span>Critical (&lt; 7 days)</span>
                <strong>{{.CriticalCount}}</strong>
            </div>
            <div class="stat">
                <span>Warning (&lt; 30 days)</span>
                <strong>{{.WarningCount}}</strong>
            </div>
        </div>

        <table>
            <thead>
                <tr>
                    <th>Domain</th>
                    <th>Expiry Date</th>
                    <th>Days Until Expiry</th>
                    <th>Status</th>
                </tr>
            </thead>
            <tbody>
                {{range .Domains}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.ExpiryDateFormatted}}</td>
                    <td>{{.DaysUntilExpiry}}</td>
                    <td><span class="status-badge {{.StatusClass}}">{{.StatusText}}</span></td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <div class="empty-state">
            <h2>No domains configured</h2>
            <p>Add domains to your configuration to start tracking their expiry dates.</p>
        </div>
        {{end}}
    </div>
</body>
</html>
`

type domainViewModel struct {
	Name                string
	ExpiryDateFormatted string
	DaysUntilExpiry     int
	StatusText          string
	StatusClass         string
}

type pageData struct {
	Domains       []domainViewModel
	TotalCount    int
	CriticalCount int
	WarningCount  int
}

func DomainsPageHandlerFactory(store *state.DomainStore) http.HandlerFunc {
	tmpl := template.Must(template.New("domains").Parse(domainsPageTemplate))

	return func(w http.ResponseWriter, r *http.Request) {
		domains := store.GetAllDomains()

		var viewModels []domainViewModel
		criticalCount := 0
		warningCount := 0

		// Convert domains to view models
		for _, domain := range domains {
			daysUntil := int(time.Until(domain.ExpiryDate).Hours() / 24)

			var statusText string
			var statusClass string

			if daysUntil < 0 {
				statusText = "Expired"
				statusClass = "status-critical"
				criticalCount++
			} else if daysUntil < 7 {
				statusText = "Critical"
				statusClass = "status-critical"
				criticalCount++
			} else if daysUntil < 30 {
				statusText = "Warning"
				statusClass = "status-warning"
				warningCount++
			} else {
				statusText = "Good"
				statusClass = "status-good"
			}

			viewModels = append(viewModels, domainViewModel{
				Name:                domain.Name,
				ExpiryDateFormatted: domain.ExpiryDate.Format("Jan 02, 2006"),
				DaysUntilExpiry:     daysUntil,
				StatusText:          statusText,
				StatusClass:         statusClass,
			})
		}

		// Sort by days until expiry (ascending - most urgent first)
		sort.Slice(viewModels, func(i, j int) bool {
			return viewModels[i].DaysUntilExpiry < viewModels[j].DaysUntilExpiry
		})

		data := pageData{
			Domains:       viewModels,
			TotalCount:    len(viewModels),
			CriticalCount: criticalCount,
			WarningCount:  warningCount,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}
