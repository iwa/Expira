package main

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/iwa/Expira/internal/api"
	"github.com/iwa/Expira/internal/cron"
	"github.com/iwa/Expira/internal/utils"
)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(4).
	PaddingRight(4).
	MarginLeft(7).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#7D56F4"))

func main() {
	fmt.Println(titleStyle.Render("Domain Expiry Watcher"))

	// Load configuration and initialize domain store using dependency injection
	config, store := utils.LoadConfig()

	println("[INFO] Starting domain expiry watcher...")

	// Update domain expiry dates from WHOIS servers
	utils.UpdateDomains(store)

	// Display domain status in console
	utils.ReportStatusInConsole(store)

	// Send initial notifications if needed
	utils.Notify(store, config)

	// Setup HTTP API endpoints with dependency injection
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/status", api.StatusHandlerFactory(store))
	go http.ListenAndServe("0.0.0.0:8080", nil)

	// Start cron job for daily domain updates
	cron.StartCronLoop(store, config)

	select {} // Keep the main goroutine running
}
