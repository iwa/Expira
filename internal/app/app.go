package app

import (
	"fmt"
	"time"

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
func (app *App) Start() error {
	// Initial domain update and notification
	utils.UpdateDomains(app.Store)
	utils.ReportStatusInConsole(app.Store)
	utils.Notify(app.Store, app.Config)

	// Start HTTP server
	server := api.NewServer("0.0.0.0:8080", app.Store)
	server.Start()

	// Start cron job
	cron.StartCronLoop(app.Store, app.Config)

	// Block the main go routine with error handling
	err := <-server.Errors()
	if shutdownErr := server.Shutdown(5 * time.Second); shutdownErr != nil {
		return fmt.Errorf("server error: %w, shutdown error: %v", err, shutdownErr)
	}
	return err
}
