package display

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"p-monitor/internal/logs"
	"p-monitor/pkg/config"
	"p-monitor/pkg/monitor"
	"p-monitor/pkg/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
)

// Display handles the system tray display
type Display struct {
	app       desktop.App
	monitor   *monitor.Monitor
	config    *config.Config
	menu      *fyne.Menu
	menuItems map[string]*fyne.MenuItem
}

// New creates a new display instance
func New(app desktop.App, monitor *monitor.Monitor, cfg *config.Config) *Display {
	return &Display{
		app:       app,
		monitor:   monitor,
		config:    cfg,
		menuItems: make(map[string]*fyne.MenuItem),
	}
}

// Start starts the display system
func (d *Display) Start() {
	logs.Info("Starting display system")

	// Set system tray icon
	d.setSystemTrayIcon()

	// Create and set system tray menu
	d.createSystemTrayMenu()

	// Start update loop
	go d.updateLoop()
}

// setSystemTrayIcon sets the system tray icon
func (d *Display) setSystemTrayIcon() {
	iconPath := filepath.Join("assets", "icon.png")
	iconURI := storage.NewFileURI(iconPath)
	icon, err := storage.LoadResourceFromURI(iconURI)
	if err != nil {
		logs.Error("Failed to load system tray icon: %v", err)
		return
	}

	d.app.SetSystemTrayIcon(icon)
}

// createSystemTrayMenu creates the system tray menu
func (d *Display) createSystemTrayMenu() {
	d.menu = fyne.NewMenu("p-monitor")

	// Create initial menu items
	d.createInitialMenuItems()

	d.app.SetSystemTrayMenu(d.menu)
}

// createInitialMenuItems creates the initial menu structure
func (d *Display) createInitialMenuItems() {
	// Create menu items for metrics
	d.menuItems["disk"] = fyne.NewMenuItem("HDD: Loading...", nil)
	d.menuItems["memory"] = fyne.NewMenuItem("RAM: Loading...", nil)
	d.menuItems["cpu"] = fyne.NewMenuItem("CPU: Loading...", nil)

	// Add initial menu items
	d.menu.Items = append(d.menu.Items, d.menuItems["disk"])
	d.menu.Items = append(d.menu.Items, d.menuItems["memory"])
	d.menu.Items = append(d.menu.Items, d.menuItems["cpu"])

	// Add separator
	d.menu.Items = append(d.menu.Items, fyne.NewMenuItemSeparator())

	// Add configuration menu items
	d.addConfigMenuItems()

	// Set the menu to trigger initial update
	d.app.SetSystemTrayMenu(d.menu)
}

// updateMenu updates the system tray menu with current metrics
func (d *Display) updateMenu() {
	metrics := d.monitor.GetMetrics()
	if metrics == nil {
		return
	}

	// Recreate menu items with updated data
	d.recreateMenuItems(metrics)
}

// recreateMenuItems recreates the menu with updated metrics
func (d *Display) recreateMenuItems(metrics *types.SystemMetrics) {
	// Clear existing menu items
	d.menu.Items = nil
	d.menuItems = make(map[string]*fyne.MenuItem)

	// Add disk metrics
	d.menuItems["disk"] = d.createDiskMenuItem(metrics.Disk)
	d.menu.Items = append(d.menu.Items, d.menuItems["disk"])

	// Add memory metrics
	d.menuItems["memory"] = d.createMemoryMenuItem(metrics.Memory)
	d.menu.Items = append(d.menu.Items, d.menuItems["memory"])

	// Add CPU metrics
	d.menuItems["cpu"] = d.createCPUMenuItem(metrics.CPU)
	d.menu.Items = append(d.menu.Items, d.menuItems["cpu"])

	// Add GPU metrics
	for i, gpu := range metrics.GPUs {
		key := fmt.Sprintf("gpu%d", i)
		d.menuItems[key] = d.createGPUMenuItem(gpu, i)
		d.menu.Items = append(d.menu.Items, d.menuItems[key])
	}

	// Add separator
	d.menu.Items = append(d.menu.Items, fyne.NewMenuItemSeparator())

	// Add configuration menu items
	d.addConfigMenuItems()

	// Update the menu
	d.app.SetSystemTrayMenu(d.menu)
}

// createDiskMenuItem creates a disk metrics menu item
func (d *Display) createDiskMenuItem(disk *types.DiskMetrics) *fyne.MenuItem {
	var text string
	var icon fyne.Resource

	if disk.Error != "" {
		text = "HDD: n/a"
		icon = d.loadIcon("error-icon.png")
	} else {
		totalGB := float64(disk.Total) / (1024 * 1024 * 1024)
		text = fmt.Sprintf("HDD: %.1fGB (%.1f%%)", totalGB, disk.UsedPercent)
		icon = d.loadIcon("drive-icon.png")
	}

	item := fyne.NewMenuItem(text, nil)
	item.Icon = icon
	return item
}

// createMemoryMenuItem creates a memory metrics menu item
func (d *Display) createMemoryMenuItem(memory *types.MemoryMetrics) *fyne.MenuItem {
	var text string
	var icon fyne.Resource

	if memory.Error != "" {
		text = "RAM: n/a"
		icon = d.loadIcon("error-icon.png")
	} else {
		totalGB := float64(memory.Total) / (1024 * 1024 * 1024)
		text = fmt.Sprintf("RAM: %.1fGB (%.1f%%)", totalGB, memory.UsedPercent)
		icon = d.loadIcon("drive-icon.png") // Using drive icon for memory for now
	}

	item := fyne.NewMenuItem(text, nil)
	item.Icon = icon
	return item
}

// createCPUMenuItem creates a CPU metrics menu item
func (d *Display) createCPUMenuItem(cpu *types.CPUMetrics) *fyne.MenuItem {
	var text string
	var icon fyne.Resource

	if cpu.Error != "" {
		text = "CPU: n/a"
		icon = d.loadIcon("error-icon.png")
	} else {
		if cpu.Temperature > 0 {
			text = fmt.Sprintf("CPU: %.1f%% (%.1f°C)", cpu.UsagePercent, cpu.Temperature)
		} else {
			text = fmt.Sprintf("CPU: %.1f%%", cpu.UsagePercent)
		}
		icon = d.loadIcon("cpu-icon.png")
	}

	item := fyne.NewMenuItem(text, nil)
	item.Icon = icon
	return item
}

// createGPUMenuItem creates a GPU menu item with simplified label
func (d *Display) createGPUMenuItem(gpu *types.GPUMetrics, index int) *fyne.MenuItem {
	var text string

	// Create simplified GPU label
	gpuLabel := d.getSimplifiedGPULabel(gpu, index)

	if gpu.Error != "" {
		text = fmt.Sprintf("%s: n/a", gpuLabel)
	} else {
		if gpu.Temperature > 0 {
			text = fmt.Sprintf("%s: %.1f%% %.1f°C", gpuLabel, gpu.UsagePercent, gpu.Temperature)
		} else {
			text = fmt.Sprintf("%s: %.1f%%", gpuLabel, gpu.UsagePercent)
		}
	}

	item := fyne.NewMenuItem(text, nil)
	if gpu.Error != "" {
		item.Icon = d.loadIcon("error-icon.png")
	} else {
		item.Icon = d.loadIcon("gpu-icon.png")
	}

	return item
}

// getSimplifiedGPULabel returns a simplified GPU label
func (d *Display) getSimplifiedGPULabel(gpu *types.GPUMetrics, index int) string {
	switch gpu.Type {
	case "nvidia":
		return fmt.Sprintf("NVIDIA GPU %d", index)
	case "amd":
		return fmt.Sprintf("AMD GPU %d", index)
	case "integrated":
		return fmt.Sprintf("iGPU %d", index)
	default:
		return fmt.Sprintf("GPU %d", index)
	}
}

// addConfigMenuItems adds configuration menu items
func (d *Display) addConfigMenuItems() {
	// Update interval
	intervalText := fmt.Sprintf("Update: %d %s", d.config.UpdateInterval, d.config.TimeUnit)
	intervalItem := fyne.NewMenuItem(intervalText, func() {
		d.showIntervalDialog()
	})
	d.menu.Items = append(d.menu.Items, intervalItem)

	// Temperature unit
	tempText := fmt.Sprintf("Temperature: %s", strings.Title(d.config.TemperatureUnit))
	tempItem := fyne.NewMenuItem(tempText, func() {
		d.toggleTemperatureUnit()
	})
	d.menu.Items = append(d.menu.Items, tempItem)
}

// loadIcon loads an icon from the assets directory
func (d *Display) loadIcon(filename string) fyne.Resource {
	iconPath := filepath.Join("assets", filename)
	iconURI := storage.NewFileURI(iconPath)
	icon, err := storage.LoadResourceFromURI(iconURI)
	if err != nil {
		logs.Error("Failed to load icon %s: %v", filename, err)
		return nil
	}
	return icon
}

// updateLoop runs the update loop
func (d *Display) updateLoop() {
	// Update menu periodically
	ticker := time.NewTicker(2 * time.Second) // Update display every 2 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.updateMenu()
		}
	}
}

// showIntervalDialog shows a dialog to change update interval
func (d *Display) showIntervalDialog() {
	// This would show a dialog to change the interval
	// For now, we'll just toggle between some common values
	if d.config.TimeUnit == "seconds" {
		if d.config.UpdateInterval >= 60 {
			d.config.UpdateInterval = 1
			d.config.TimeUnit = "minutes"
		} else {
			d.config.UpdateInterval += 5
		}
	} else {
		if d.config.UpdateInterval >= 30 {
			d.config.UpdateInterval = 5
			d.config.TimeUnit = "seconds"
		} else {
			d.config.UpdateInterval += 5
		}
	}

	d.config.Save()
	logs.Info("Updated interval to %d %s", d.config.UpdateInterval, d.config.TimeUnit)

	// Trigger menu update to show new configuration
	d.updateMenu()
}

// toggleTemperatureUnit toggles between Celsius and Fahrenheit
func (d *Display) toggleTemperatureUnit() {
	if d.config.TemperatureUnit == "celsius" {
		d.config.TemperatureUnit = "fahrenheit"
	} else {
		d.config.TemperatureUnit = "celsius"
	}

	d.config.Save()
	logs.Info("Updated temperature unit to %s", d.config.TemperatureUnit)

	// Trigger menu update to show new configuration
	d.updateMenu()
}
