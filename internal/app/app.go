package app

import (
	"net/http"

	"github.com/iwa/Expira/internal/api"
	"github.com/iwa/Expira/internal/cron"
	"github.com/iwa/Expira/internal/state"
	"github.com/iwa/Expira/internal/utils"
)

// App holds the application's core dependencies
type App struct {
	Config *state.Config
	Store  *state.DomainStore
}

// New creates and initializes a new App instance
func New() *App {
	config, store := utils.LoadConfig()
	return &App{
		Config: config,
		Store:  store,
	}
}

// Start runs the application
func (app *App) Start() {
	// Initial domain update and notification
	utils.UpdateDomains(app.Store)
	utils.ReportStatusInConsole(app.Store)
	utils.Notify(app.Store, app.Config)

	// Setup HTTP server
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/status", api.StatusHandlerFactory(app.Store))
	go http.ListenAndServe("0.0.0.0:8080", nil)

	// Start cron job
	cron.StartCronLoop(app.Store, app.Config)
}
