package app

import (
	"fmt"
	"time"

	"github.com/iwa/Expira/internal/api"
	"github.com/iwa/Expira/internal/cron"
	"github.com/iwa/Expira/internal/state"
	"github.com/iwa/Expira/internal/utils"
	cronlib "github.com/robfig/cron/v3"
)

// App holds the application's core dependencies
type App struct {
	Config *state.Config
	Store  *state.DomainStore
	Cron   *cronlib.Cron
}

// New creates and initializes a new App instance
func New() *App {
	config, store := utils.LoadConfig()
	return &App{
		Config: config,
		Store:  store,
		Cron:   cronlib.New(),
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
	_, err := app.Cron.AddFunc("0 0 * * *", func() { // everyday at midnight
		cron.RunDailyUpdate(app.Store, app.Config)
	})
	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %w", err)
	}
	app.Cron.Start()
	fmt.Println("[INFO] Cron scheduler started - daily domain updates scheduled at midnight")

	// Block the main go routine with error handling
	errServer := <-server.Errors()
	if shutdownErr := server.Shutdown(5 * time.Second); shutdownErr != nil {
		app.Cron.Stop()
		return fmt.Errorf("server error: %w, shutdown error: %v", err, shutdownErr)
	}
	return errServer
}
