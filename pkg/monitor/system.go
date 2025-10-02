package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// getDiskUsage gets disk usage for the root filesystem
func getDiskUsage(path string) (*disk.UsageStat, error) {
	return disk.Usage(path)
}

// getMemoryUsage gets system memory usage
func getMemoryUsage() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

// getCPUUsage gets CPU usage percentage
func getCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	if len(percentages) == 0 {
		return 0, fmt.Errorf("no CPU usage data available")
	}

	return percentages[0], nil
}

// getCPUTemperature gets CPU temperature from thermal sensors
func getCPUTemperature() (float64, error) {
	// Try to read from thermal sensors
	// Common paths for CPU temperature on Linux
	thermalPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",                     // Most common
		"/sys/class/thermal/thermal_zone1/temp",                     // Alternative
		"/sys/devices/platform/coretemp.0/hwmon/hwmon*/temp1_input", // Intel
		"/sys/devices/virtual/thermal/thermal_zone0/temp",           // Virtual thermal
	}

	for _, path := range thermalPaths {
		if temp, err := readThermalFile(path); err == nil {
			return temp, nil
		}
	}

	// Try to find thermal zones dynamically
	if temp, err := findThermalZone(); err == nil {
		return temp, nil
	}

	return 0, fmt.Errorf("CPU temperature not available")
}

// readThermalFile reads temperature from a thermal file
func readThermalFile(path string) (float64, error) {
	// Handle wildcard paths
	if strings.Contains(path, "*") {
		matches, err := filepath.Glob(path)
		if err != nil || len(matches) == 0 {
			return 0, fmt.Errorf("no thermal files found")
		}
		path = matches[0]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	// Temperature is usually in millidegrees Celsius
	tempStr := strings.TrimSpace(string(data))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	// Convert from millidegrees to degrees
	return temp / 1000.0, nil
}

// findThermalZone finds and reads from available thermal zones
func findThermalZone() (float64, error) {
	thermalDir := "/sys/class/thermal"
	entries, err := os.ReadDir(thermalDir)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "thermal_zone") {
			tempPath := filepath.Join(thermalDir, entry.Name(), "temp")
			if temp, err := readThermalFile(tempPath); err == nil {
				// Check if this is a CPU thermal zone by reading type
				typePath := filepath.Join(thermalDir, entry.Name(), "type")
				if typeData, err := os.ReadFile(typePath); err == nil {
					zoneType := strings.TrimSpace(string(typeData))
					if strings.Contains(strings.ToLower(zoneType), "cpu") ||
						strings.Contains(strings.ToLower(zoneType), "x86") ||
						strings.Contains(strings.ToLower(zoneType), "core") {
						return temp, nil
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("no CPU thermal zone found")
}

// Public functions for testing
func GetDiskUsage(path string) (*disk.UsageStat, error) {
	return getDiskUsage(path)
}

func GetMemoryUsage() (*mem.VirtualMemoryStat, error) {
	return getMemoryUsage()
}

func GetCPUUsage() (float64, error) {
	return getCPUUsage()
}

func GetCPUTemperature() (float64, error) {
	return getCPUTemperature()
}
