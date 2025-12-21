package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/iwa/Expira/internal/app"
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

	app := app.New()
	if err := app.Start(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
