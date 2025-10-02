package main

import (
	"log"
	"os"
	"path/filepath"

	"p-monitor/internal/logs"
	"p-monitor/pkg/config"
	"p-monitor/pkg/display"
	"p-monitor/pkg/monitor"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	// Initialize logging
	logDir := filepath.Join(os.Getenv("HOME"), ".p-monitor", "logs")
	if err := logs.Init(logDir); err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logs.Error("Failed to load configuration: %v", err)
		cfg = config.Default()
	}

	// Create Fyne application
	a := app.NewWithID("com.p-monitor.app")
	a.SetIcon(nil) // We'll set the system tray icon instead

	// Check if we're running on desktop
	if desk, ok := a.(desktop.App); ok {
		// Initialize monitor
		monitor := monitor.New(cfg)

		// Initialize display
		display := display.New(desk, monitor, cfg)

		// Start monitoring
		go monitor.Start()

		// Start display
		display.Start()
	} else {
		logs.Error("This application requires desktop environment")
		os.Exit(1)
	}

	// Run the application
	a.Run()
}
